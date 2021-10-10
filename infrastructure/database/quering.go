package database

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/artem-malko/auth-and-go/infrastructure/caller"
	"github.com/pkg/errors"
)

type QueryBuilder interface {
	ToSql() (string, []interface{}, error)
}

func GetRawFieldName(fieldName string) string {
	return strings.Split(fieldName, ".")[1]
}

func QueryRow(executor QueryExecutor, queryBuilder QueryBuilder) (*sql.Row, error) {
	callerName := "repository unknown caller"

	// Remove getCaller from stack
	if caller, ok := caller.GetCaller(2); ok == true {
		callerName = caller
	}

	return QueryRowWithCaller(callerName, executor, queryBuilder)
}

func QueryRowWithCaller(caller string, executor QueryExecutor, queryBuilder QueryBuilder) (*sql.Row, error) {
	if executor == nil {
		return nil, errors.New("no query executor for " + caller)
	}

	query, args, err := queryBuilder.ToSql()

	if err != nil {
		return nil, errors.Wrap(err, caller+" build query error")
	}

	var row *sql.Row

	row = executor.QueryRow(query, args...)

	return row, nil
}

func Exec(executor QueryExecutor, queryBuilder QueryBuilder) error {
	callerName := "repository unknown caller"

	// Remove getCaller from stack
	if caller, ok := caller.GetCaller(2); ok == true {
		callerName = caller
	}

	return ExecWithCaller(callerName, executor, queryBuilder)
}

func ExecWithCaller(caller string, executor QueryExecutor, queryBuilder QueryBuilder) error {
	if executor == nil {
		return errors.New("no query executor for " + caller)
	}

	query, args, err := queryBuilder.ToSql()

	if err != nil {
		return errors.Wrap(err, caller+" build query error")
	}

	var res sql.Result

	res, err = executor.Exec(query, args...)

	if err != nil {
		return errors.Wrap(err, caller+" exec query error")
	}

	rowsAffected, err := res.RowsAffected()

	if err != nil {
		return errors.Wrap(err, caller+" rowsAffected calc error")
	}

	if rowsAffected == 0 {
		return ErrRepositoryNoRowsAffected
	}

	return nil
}

func Query(executor QueryExecutor,
	queryBuilder QueryBuilder,
	scanFunc func(rows *sql.Rows) error,
) error {
	callerName := "repository unknown caller"

	// Remove getCaller from stack
	if caller, ok := caller.GetCaller(2); ok == true {
		callerName = caller
	}

	return QueryWithCaller(callerName, executor, queryBuilder, scanFunc)
}

func QueryWithCaller(
	caller string,
	executor QueryExecutor,
	queryBuilder QueryBuilder,
	scanFunc func(rows *sql.Rows) error,
) error {
	if executor == nil {
		return errors.New("no query executor for " + caller)
	}

	query, args, err := queryBuilder.ToSql()

	fmt.Println(query)

	if err != nil {
		return errors.Wrap(err, caller+" build query error")
	}

	var rows *sql.Rows

	rows, err = executor.Query(query, args...)

	if err != nil {
		return errors.Wrap(err, caller+" query error")
	}

	for rows.Next() {
		err := scanFunc(rows)

		if err != nil {
			return err
		}
	}

	err = rows.Err()

	if err != nil {
		return errors.Wrap(err, caller+" can`t get info")
	}

	return nil
}

func RunWithTransaction(db *sql.DB, callerName string, function func(tx *sql.Tx) error) error {
	tx, err := db.Begin()

	if err != nil {
		return errors.Wrap(err, "transaction opening error for "+callerName)
	}

	defer func() {
		// @TODO think about it
		//if p := recover(); p != nil {
		//	tx.Rollback()
		//	panic(p) // re-throw panic after Rollback
		//} else if err != nil {
		//	tx.Rollback() // err is non-nil; don't change it
		//} else {
		//	err = tx.Commit() // err is nil; if Commit returns error update err
		//}
		if err != nil {
			rollbackErr := tx.Rollback()

			if rollbackErr != nil {
				err = errors.Wrap(rollbackErr, "transaction rollbackErr"+callerName)
			}
		}
	}()

	err = function(tx)

	if err != nil {
		return err
	}

	err = tx.Commit()

	if err != nil {
		return errors.Wrap(err, "transaction commit error for "+callerName)
	}

	return nil
}
