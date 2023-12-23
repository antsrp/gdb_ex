package db

import (
	"context"

	"github.com/antsrp/gdb_ex/internal/domain/models"
)

type DepartmentsRepository interface {
	Create(context.Context, models.Department) (models.Department, error)
	Update(context.Context, models.Department) error
	Delete(context.Context, models.Department) error
	GetOne(context.Context, models.Department) (*models.Department, error)
	GetAll(context.Context) ([]models.Department, error)
}
