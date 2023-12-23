package db

import (
	"context"

	"github.com/antsrp/gdb_ex/internal/domain/models"
)

type EmployeesRepository interface {
	Create(context.Context, models.Employee) (models.Employee, error)
	Update(context.Context, models.Employee) error
	Delete(context.Context, models.Employee) error
	GetOne(context.Context, models.Employee) (*models.Employee, error)
	GetAll(context.Context) ([]models.Employee, error)
	AddToProject(context.Context, models.Employee, models.Project) error
}
