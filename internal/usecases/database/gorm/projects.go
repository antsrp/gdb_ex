package gorm

import (
	"context"
	"errors"
	"fmt"

	"github.com/antsrp/gdb_ex/internal/domain/models"
	irepo "github.com/antsrp/gdb_ex/internal/interfaces/db"
	"gorm.io/gorm"
)

var _ irepo.ProjectsRepository = &ProjectStorage{}

type ProjectStorage struct {
	//conn *connection
	conn *Connection
}

// func NewProjectStorage(conn *connection) (*ProjectStorage, error) {
func NewProjectStorage(conn *Connection) (*ProjectStorage, error) {
	s := &ProjectStorage{conn: conn}

	return s, nil
}

func (ps ProjectStorage) Create(ctx context.Context, project models.Project) (models.Project, error) {
	operand := presentOperand(ctx, ps.conn.logger, ps.conn.gdb, ps.conn.txKey)

	if result := operand.Model(&models.Project{}).Create(&project); result.Error != nil {
		return models.Project{}, fmt.Errorf("can't create project: %w", result.Error)
	}

	return project, nil
}

func (ps ProjectStorage) Update(ctx context.Context, project models.Project) error {
	operand := presentOperand(ctx, ps.conn.logger, ps.conn.gdb, ps.conn.txKey)

	if result := operand.Model(&models.Project{}).Updates(&project); result.Error != nil {
		return fmt.Errorf("can't update project: %w", result.Error)
	} else if result.RowsAffected == 0 {
		return irepo.ErrNoRowsUpdate
	}

	return nil
}

func (ps ProjectStorage) Delete(ctx context.Context, project models.Project) error {
	operand := presentOperand(ctx, ps.conn.logger, ps.conn.gdb, ps.conn.txKey)

	if result := operand.Model(&models.Project{}).Delete(&project); result.Error != nil {
		return fmt.Errorf("can't delete project: %w", result.Error)
	} else if result.RowsAffected == 0 {
		return irepo.ErrNoRowsDelete
	}

	return nil
}

func (ps ProjectStorage) GetOne(ctx context.Context, project models.Project) (*models.Project, error) {
	operand := presentOperand(ctx, ps.conn.logger, ps.conn.gdb, ps.conn.txKey)

	var p models.Project

	if tx := operand.Model(&models.Project{}).Find(&p, &project); tx.Error != nil {
		var err error
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			err = irepo.ErrNoRowsSelect
		} else {
			err = tx.Error
		}
		return nil, fmt.Errorf("can't select project: %w", err)
	} else if tx.RowsAffected == 0 {
		return nil, irepo.ErrNoRowsSelect
	}

	return &p, nil
}

func (ps ProjectStorage) GetAll(ctx context.Context) ([]models.Project, error) {
	operand := presentOperand(ctx, ps.conn.logger, ps.conn.gdb, ps.conn.txKey)

	var projects []models.Project

	if tx := operand.Model(&models.Project{}).Find(&projects); tx.Error != nil {
		var err error
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			err = irepo.ErrNoRowsSelect
		} else {
			err = tx.Error
		}
		return nil, fmt.Errorf("can't select project: %w", err)
	} else if tx.RowsAffected == 0 {
		return nil, irepo.ErrNoRowsSelect
	}

	return projects, nil
}
