package db

import (
	"context"

	"github.com/antsrp/gdb_ex/internal/domain/models"
)

type ProjectsRepository interface {
	Create(context.Context, models.Project) (models.Project, error)
	Update(context.Context, models.Project) error
	Delete(context.Context, models.Project) error
	GetOne(context.Context, models.Project) (*models.Project, error)
	GetAll(context.Context) ([]models.Project, error)
}
