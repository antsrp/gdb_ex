package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/antsrp/gdb_ex/internal/domain/models"
	irepo "github.com/antsrp/gdb_ex/internal/interfaces/db"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

var _ irepo.DepartmentsRepository = &DepartmentStorage{}

const (
	createDepartmentStmtQuery  = "INSERT INTO Departments (name) VALUES ($1) RETURNING id;"
	updateDepartmentStmtQuery  = "UPDATE Departments SET name = $1 WHERE id = $2;"
	deleteDepartmentStmtQuery  = "DELETE FROM Departments WHERE id = $1;"
	getAllDepartmentsStmtQuery = "SELECT id, name FROM Departments"
	getOneDepartmentStmtQuery  = getAllDepartmentsStmtQuery + " WHERE id = $1;"
)

type DepartmentStorage struct {
	conn *Connection
}

func NewDepartmentStorage(conn *Connection) (*DepartmentStorage, error) {
	return &DepartmentStorage{conn: conn}, nil
}

func (ds DepartmentStorage) Create(ctx context.Context, department models.Department) (models.Department, error) {
	operand := presentOperand(ctx, ds.conn.logger, ds.conn.pc, ds.conn.txKey)

	row := operand.QueryRow(ds.conn.ctx, createDepartmentStmtQuery, department.Name)

	var id uint

	if err := row.Scan(&id); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			err = irepo.ErrExistingNameDepartment
		}
		return models.Department{}, fmt.Errorf("can't create department with name %s: %w", department.Name, err)
	}

	return models.Department{ID: id}, nil
}

func (ds DepartmentStorage) Update(ctx context.Context, department models.Department) error {
	operand := presentOperand(ctx, ds.conn.logger, ds.conn.pc, ds.conn.txKey)
	tag, err := operand.Exec(ds.conn.ctx, updateDepartmentStmtQuery, department.Name, department.ID)
	if err != nil {
		return fmt.Errorf("can't update department: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return irepo.ErrNoRowsUpdate
	}

	return nil
}

func (ds DepartmentStorage) Delete(ctx context.Context, department models.Department) error {
	operand := presentOperand(ctx, ds.conn.logger, ds.conn.pc, ds.conn.txKey)
	_, err := operand.Exec(ds.conn.ctx, deleteDepartmentStmtQuery, department.ID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			err = irepo.ErrNoRowsDelete
		}
		return fmt.Errorf("can't delete department: %w", err)
	}

	return nil
}

func (ds DepartmentStorage) GetOne(ctx context.Context, department models.Department) (*models.Department, error) {
	operand := presentOperand(ctx, ds.conn.logger, ds.conn.pc, ds.conn.txKey)
	row := operand.QueryRow(ds.conn.ctx, getOneDepartmentStmtQuery, department.ID)

	var d models.Department
	if err := row.Scan(&d.ID, &d.Name); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			err = irepo.ErrNoRowsSelect
		}
		return nil, fmt.Errorf("can't select department: %w", err)
	}

	return &d, nil
}

func (ds DepartmentStorage) GetAll(ctx context.Context) ([]models.Department, error) {
	operand := presentOperand(ctx, ds.conn.logger, ds.conn.pc, ds.conn.txKey)
	rows, err := operand.Query(ds.conn.ctx, getAllDepartmentsStmtQuery)
	if err != nil {
		return nil, fmt.Errorf("can't select all departments: %w", err)
	}
	defer rows.Close()

	var departments []models.Department

	for rows.Next() {
		var d models.Department
		if err := rows.Scan(&d.ID, &d.Name); err != nil {
			return nil, fmt.Errorf("can't scan department: %w", err)
		}
		departments = append(departments, d)
	}

	return departments, nil
}
