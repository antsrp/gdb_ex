package gorm

import (
	"context"
	"errors"
	"fmt"

	"github.com/antsrp/gdb_ex/internal/domain/models"
	irepo "github.com/antsrp/gdb_ex/internal/interfaces/db"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

var _ irepo.DepartmentsRepository = &DepartmentStorage{}

type DepartmentStorage struct {
	//conn *connection
	conn *Connection
}

// func NewDepartmentStorage(conn *connection) (*DepartmentStorage, error) {
func NewDepartmentStorage(conn *Connection) (*DepartmentStorage, error) {
	s := &DepartmentStorage{conn: conn}

	return s, nil
}

func (ds DepartmentStorage) Create(ctx context.Context, department models.Department) (models.Department, error) {
	operand := presentOperand(ctx, ds.conn.logger, ds.conn.gdb, ds.conn.txKey)

	if result := operand.Model(&models.Department{}).Create(&department); result.Error != nil {
		err := result.Error
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			err = irepo.ErrExistingNameDepartment
		}
		/*if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			err = irepo.ErrExistingNameDepartment
		} else {
			err = result.Error
		} */
		return models.Department{}, fmt.Errorf("can't create department with name %s: %w", department.Name, err)
	}

	return department, nil
}

func (ds DepartmentStorage) Update(ctx context.Context, department models.Department) error {
	operand := presentOperand(ctx, ds.conn.logger, ds.conn.gdb, ds.conn.txKey)

	if result := operand.Model(&models.Department{}).Updates(&department); result.Error != nil {
		return fmt.Errorf("can't update department: %w", result.Error)
	} else if result.RowsAffected == 0 {
		return irepo.ErrNoRowsUpdate
	}

	return nil
}

func (ds DepartmentStorage) Delete(ctx context.Context, department models.Department) error {
	operand := presentOperand(ctx, ds.conn.logger, ds.conn.gdb, ds.conn.txKey)

	if result := operand.Model(&models.Department{}).Delete(&department); result.Error != nil {
		return fmt.Errorf("can't delete department: %w", result.Error)
	} else if result.RowsAffected == 0 {
		return irepo.ErrNoRowsDelete
	}

	return nil
}

func (ds DepartmentStorage) GetOne(ctx context.Context, department models.Department) (*models.Department, error) {
	operand := presentOperand(ctx, ds.conn.logger, ds.conn.gdb, ds.conn.txKey)

	var p models.Department

	if tx := operand.Model(&models.Department{}).Find(&p, &department); tx.Error != nil {
		var err error
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			err = irepo.ErrNoRowsSelect
		} else {
			err = tx.Error
		}
		return nil, fmt.Errorf("can't select department: %w", err)
	} else if tx.RowsAffected == 0 {
		return nil, irepo.ErrNoRowsSelect
	}

	return &p, nil
}

func (ds DepartmentStorage) GetAll(ctx context.Context) ([]models.Department, error) {
	operand := presentOperand(ctx, ds.conn.logger, ds.conn.gdb, ds.conn.txKey)

	var departments []models.Department

	if tx := operand.Model(&models.Department{}).Find(&departments); tx.Error != nil {
		var err error
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			err = irepo.ErrNoRowsSelect
		} else {
			err = tx.Error
		}
		return nil, fmt.Errorf("can't select department: %w", err)
	} else if tx.RowsAffected == 0 {
		return nil, irepo.ErrNoRowsSelect
	}
	return departments, nil
}
