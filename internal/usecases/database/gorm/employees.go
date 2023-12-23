package gorm

import (
	"context"
	"errors"
	"fmt"

	"github.com/antsrp/gdb_ex/internal/domain/models"
	irepo "github.com/antsrp/gdb_ex/internal/interfaces/db"
	"gorm.io/gorm"
)

var _ irepo.EmployeesRepository = &EmployeeStorage{}

type EmployeeStorage struct {
	//conn *connection
	conn *Connection
}

// func NewEmployeeStorage(conn *connection) (*EmployeeStorage, error) {
func NewEmployeeStorage(conn *Connection) (*EmployeeStorage, error) {
	s := &EmployeeStorage{conn: conn}

	return s, nil
}

func (es EmployeeStorage) Create(ctx context.Context, employee models.Employee) (models.Employee, error) {
	operand := presentOperand(ctx, es.conn.logger, es.conn.gdb, es.conn.txKey)

	model := createEmployeeModel(employee)

	if result := operand.Model(&EmployeeModel{}).Create(&model); result.Error != nil {
		return models.Employee{}, fmt.Errorf("can't create employee: %w", result.Error)
	}
	employee.ID = model.ID

	return employee, nil
}

func (es EmployeeStorage) Update(ctx context.Context, employee models.Employee) error {
	operand := presentOperand(ctx, es.conn.logger, es.conn.gdb, es.conn.txKey)

	model := createEmployeeModel(employee)

	if result := operand.Model(&EmployeeModel{}).Where("id = ?", model.ID).Updates(&model); result.Error != nil {
		return fmt.Errorf("can't update employee: %w", result.Error)
	} else if result.RowsAffected == 0 {
		return irepo.ErrNoRowsUpdate
	}

	return nil
}

func (es EmployeeStorage) Delete(ctx context.Context, employee models.Employee) error {
	operand := presentOperand(ctx, es.conn.logger, es.conn.gdb, es.conn.txKey)

	model := createEmployeeModel(employee)

	if result := operand.Model(&EmployeeModel{}).Delete(&model); result.Error != nil {
		return fmt.Errorf("can't delete employee: %w", result.Error)
	} else if result.RowsAffected == 0 {
		return irepo.ErrNoRowsDelete
	}

	return nil
}

func (es EmployeeStorage) GetOne(ctx context.Context, employee models.Employee) (*models.Employee, error) {
	operand := presentOperand(ctx, es.conn.logger, es.conn.gdb, es.conn.txKey)

	var (
		result getEmployeeStruct
	)

	if err := operand.Table(fmt.Sprintf("%s AS e", EmployeeModel{}.TableName())).
		Select("e.id, e.name, e.surname, e.age, e.salary, d.id AS depid, d.name AS depname").
		Joins("JOIN Departments d ON e.departmentID = d.id").
		Where("e.id = ?", employee.ID).
		Scan(&result).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = irepo.ErrNoRowsSelect
		}
		return nil, fmt.Errorf("can't select employee: %w", err)
	}

	e := models.Employee{
		ID:      result.ID,
		Name:    result.Name,
		Surname: result.Surname,
		Age:     result.Age,
		Salary:  result.Salary,
		Department: models.Department{
			ID:   result.DepID,
			Name: result.DepName,
		},
	}

	return &e, nil
}

func (es EmployeeStorage) GetAll(ctx context.Context) ([]models.Employee, error) {
	operand := presentOperand(ctx, es.conn.logger, es.conn.gdb, es.conn.txKey)

	var (
		employees []models.Employee
	)

	var results []getEmployeeStruct

	if err := operand.Table(fmt.Sprintf("%s AS e", EmployeeModel{}.TableName())).
		Select("e.id, e.name, e.surname, e.age, e.salary, d.id AS depid, d.name AS depname").
		Joins("JOIN Departments d ON e.departmentID = d.id").
		Scan(&results).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = irepo.ErrNoRowsSelect
		}
		return nil, fmt.Errorf("can't select employee: %w", err)
	}

	for _, r := range results {
		employees = append(employees, models.Employee{
			ID:      r.ID,
			Name:    r.Name,
			Surname: r.Surname,
			Age:     r.Age,
			Salary:  r.Salary,
			Department: models.Department{
				ID:   r.DepID,
				Name: r.DepName,
			},
		})
	}

	return employees, nil
}

func (es EmployeeStorage) AddToProject(ctx context.Context, employee models.Employee, project models.Project) error {
	operand := presentOperand(ctx, es.conn.logger, es.conn.gdb, es.conn.txKey)

	model := ProjectEmployeeRelationModel{
		Employee: createEmployeeModel(employee),
		Project:  ProjectModel(project),
	}

	if result := operand.Model(&ProjectEmployeeRelationModel{}).Create(&model); result.Error != nil {
		return fmt.Errorf("can't add employee to project: %w", result.Error)
	}

	return nil
}
