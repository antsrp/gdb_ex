package postgres

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/antsrp/gdb_ex/internal/domain/models"
	irepo "github.com/antsrp/gdb_ex/internal/interfaces/db"
	"github.com/jackc/pgx/v5"
)

var _ irepo.EmployeesRepository = &EmployeeStorage{}

const (
	createEmployeeStmtQuery       = "INSERT INTO Employees (name, surname, age, salary, departmentID) VALUES ($1, $2, $3, $4, $5) RETURNING id;"
	deleteEmployeeStmtQuery       = "DELETE FROM Employees WHERE id = $1;"
	getAllEmployeesStmtQuery      = `SELECT e.id, e.name, e.surname, e.age, e.salary, d.id, d.name FROM Employees e JOIN Departments d ON e.departmentID = d.id`
	getOneEmployeeStmtQuery       = getAllEmployeesStmtQuery + ` WHERE e.id = $1`
	addToProjectEmployeeStmtQuery = "INSERT INTO Projects_Employees (projectId, employeeId) VALUES ($1, $2);"
)

type EmployeeStorage struct {
	conn *Connection
}

func NewEmployeeStorage(conn *Connection) (*EmployeeStorage, error) {
	return &EmployeeStorage{conn: conn}, nil
}

func (es EmployeeStorage) Create(ctx context.Context, employee models.Employee) (models.Employee, error) {
	operand := presentOperand(ctx, es.conn.logger, es.conn.pc, es.conn.txKey)
	row := operand.QueryRow(es.conn.ctx, createEmployeeStmtQuery, employee.Name, employee.Surname, employee.Age, employee.Salary, employee.Department.ID)
	var id uint

	if err := row.Scan(&id); err != nil {
		return models.Employee{}, fmt.Errorf("can't create employee: %w", err)
	}

	return models.Employee{ID: id}, nil
}

func (es EmployeeStorage) Update(ctx context.Context, employee models.Employee) error {
	var sb strings.Builder
	var columnsToUpdate []string

	if employee.Name != "" {
		columnsToUpdate = append(columnsToUpdate, fmt.Sprintf("name = '%s'", employee.Name))
	}
	if employee.Surname != "" {
		columnsToUpdate = append(columnsToUpdate, fmt.Sprintf("surname = '%s'", employee.Surname))
	}
	if employee.Age != 0 {
		columnsToUpdate = append(columnsToUpdate, fmt.Sprintf("age = %d", employee.Age))
	}
	if employee.Salary != 0 {
		columnsToUpdate = append(columnsToUpdate, fmt.Sprintf("salary = %d", employee.Salary))
	}
	if employee.Department.ID != 0 {
		columnsToUpdate = append(columnsToUpdate, fmt.Sprintf("departmentId = %d", employee.Department.ID))
	}

	if len(columnsToUpdate) == 0 { // nothing to update
		return nil
	}

	sb.WriteString("UPDATE Employees SET ")
	sb.WriteString(strings.Join(columnsToUpdate, ","))
	sb.WriteString(" WHERE id = $1")

	operand := presentOperand(ctx, es.conn.logger, es.conn.pc, es.conn.txKey)

	if tag, err := operand.Exec(es.conn.ctx, sb.String(), employee.ID); err != nil {
		return fmt.Errorf("can't update employee: %w", err)
	} else if tag.RowsAffected() == 0 {
		return irepo.ErrNoRowsUpdate
	}

	return nil
}

func (es EmployeeStorage) Delete(ctx context.Context, employee models.Employee) error {
	operand := presentOperand(ctx, es.conn.logger, es.conn.pc, es.conn.txKey)
	_, err := operand.Exec(es.conn.ctx, deleteEmployeeStmtQuery, employee.ID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			err = irepo.ErrNoRowsDelete
		}
		return fmt.Errorf("can't delete employee: %w", err)
	}

	return nil
}

func (es EmployeeStorage) GetOne(ctx context.Context, employee models.Employee) (*models.Employee, error) {
	operand := presentOperand(ctx, es.conn.logger, es.conn.pc, es.conn.txKey)
	row := operand.QueryRow(es.conn.ctx, getOneEmployeeStmtQuery, employee.ID)

	var e models.Employee
	if err := row.Scan(&e.ID, &e.Name, &e.Surname, &e.Age, &e.Salary, &e.Department.ID, &e.Department.Name); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			err = irepo.ErrNoRowsSelect
		}
		return nil, fmt.Errorf("can't select employee: %w", err)
	}

	return &e, nil
}

func (es EmployeeStorage) GetAll(ctx context.Context) ([]models.Employee, error) {
	operand := presentOperand(ctx, es.conn.logger, es.conn.pc, es.conn.txKey)
	rows, err := operand.Query(es.conn.ctx, getAllEmployeesStmtQuery)
	if err != nil {
		return nil, fmt.Errorf("can't select all employees: %w", err)
	}
	defer rows.Close()

	/*employees, err := pgx.CollectRows(rows, pgx.RowToStructByPos[models.Employee])
	if err != nil {
		return nil, fmt.Errorf("can't select employees: %w", err)
	} */
	var employees []models.Employee

	for rows.Next() {
		var e models.Employee
		if err := rows.Scan(&e.ID, &e.Name, &e.Surname, &e.Age, &e.Salary, &e.Department.ID, &e.Department.Name); err != nil {
			return nil, fmt.Errorf("can't scan employee: %w", err)
		}
		employees = append(employees, e)
	}

	return employees, nil
}

func (es EmployeeStorage) AddToProject(ctx context.Context, employee models.Employee, project models.Project) error {
	operand := presentOperand(ctx, es.conn.logger, es.conn.pc, es.conn.txKey)
	_, err := operand.Exec(es.conn.ctx, addToProjectEmployeeStmtQuery, project.ID, employee.ID)
	if err != nil {
		return fmt.Errorf("can't add employee to project: %w", err)
	}

	return nil
}
