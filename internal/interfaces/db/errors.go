package db

import "errors"

var (
	ErrNoRowsDelete = errors.New("no object to delete")
	ErrNoRowsSelect = errors.New("no rows selected")
	ErrNoRowsUpdate = errors.New("no rows updated")

	ErrExistingNameDepartment = errors.New("department with this name is already exist")
)
