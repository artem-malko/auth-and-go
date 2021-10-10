package database

import "errors"

var (
	ErrRepositoryNoRowsAffected  = errors.New("no rows affected")
	ErrRepositoryNoEntitiesFound = errors.New("no entities found")
	ErrRepositoryEntityExists    = errors.New("entity exists")
)
