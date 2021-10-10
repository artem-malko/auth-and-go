package database

import (
	"context"
	"database/sql"
	"time"

	"github.com/apex/log"
	"github.com/dimiro1/health"

	// инициализация pgx как драйвера для базы данных
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/pkg/errors"
)

const MigrationsTableName = "migrations"

var (
	// ErrDBConnAttemptsFailed ошибка для случаев невозможности подключения к БД
	ErrDBConnAttemptsFailed = errors.New("All attempts to connect db failed")
	dbDriverName            = "pgx"
	connMaxLifetime         = time.Duration(5) * time.Minute
	maxConnectionAttempts   = 30
	heartbeatInterval       = time.Duration(5) * time.Second
)

type RowScanner interface {
	Scan(destination ...interface{}) error
}

type RowsScanner interface {
	Scan(destination ...interface{}) error
	Err() error
	Next() bool
}

type QueryExecutor interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
}

// Database - компонент для подключения к БД
type Database struct {
	db       *sql.DB
	Logger   log.Interface
	stopChan chan bool
	health   health.Health
}

// New инициализирует подключение к БД
func New(dsn string, maxOpenConns int, logger log.Interface) (*Database, error) {
	db, err := sql.Open(dbDriverName, dsn)

	if err != nil {
		return nil, errors.Wrap(err, "Can't  open database")
	}

	db.SetMaxOpenConns(maxOpenConns)
	db.SetConnMaxLifetime(connMaxLifetime)

	res := &Database{
		db:       db,
		Logger:   logger,
		stopChan: make(chan bool, 1),
	}

	res.health = health.NewHealth()
	res.health.Down()

	err = res.connect()

	if err != nil {
		return nil, errors.Wrap(err, "Can't connect to database")
	}

	return res, nil
}

// connect осуществляет попытку подключения к БД
func (d *Database) connect() error {
	var dbError error

	for attempt := 1; attempt <= maxConnectionAttempts; attempt++ {
		dbError = d.db.Ping()
		if dbError == nil {
			break
		}
		nextAttemptWait := time.Duration(attempt) * time.Second

		d.Logger.
			WithField("source", "database_connect").
			WithError(dbError).
			Errorf(
				"Attempt %v: Can't establish a connection with the database. Wait for %v.",
				attempt, nextAttemptWait,
			)
		time.Sleep(nextAttemptWait)
	}

	if dbError != nil {
		return ErrDBConnAttemptsFailed
	}

	d.health.Up()

	return nil

}

// Close закрывает соединение с БД
func (d *Database) Close() error {
	d.stopChan <- true

	err := d.db.Close()
	if err != nil {
		return errors.Wrap(err, "Can't close database")
	}
	return nil

}

// WatchConnection вешает watcher на соединение с базой и в случае ошибки с ним - рестартит ее
func (d *Database) WatchConnection(errChan chan<- error) {
	for {
		select {
		case <-d.stopChan:
			d.Logger.Info("DB heartbeat was stopped")
			return
		default:
			// Базовый подход с использованием метода Ping() не работает. При выключении БД "DB heartbeat is OK".
			// Причина такого поведения связата с тем, что соединение открывается только в случае необходимости,
			// а таковой нет. Похожая проблема описана в https://github.com/go-sql-driver/mysql/issues/82
			_, dbError := d.db.Exec("SELECT 1;")
			if dbError == nil {
				d.Logger.Debug("DB heartbeat is OK")
				time.Sleep(heartbeatInterval)
				continue
			}

			d.health.Down().AddInfo("error", dbError.Error())

			d.Logger.Warn("DB heartbeat has problem")
			err := d.connect()
			if err != nil {
				d.Logger.Warn("DB heartbeat retry was failed")
				errChan <- err
			}
			d.Logger.Warn("DB heartbeat retry was success")
		}
	}
}

// DB возвращает указатель на коннект к БД
func (d *Database) DB() *sql.DB {
	return d.db
}

// Check возвращает состояние подключения к базе
func (d *Database) Check() health.Health {
	return d.health
}
