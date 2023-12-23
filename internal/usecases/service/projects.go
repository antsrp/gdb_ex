package service

import (
	"context"

	"github.com/antsrp/gdb_ex/internal/domain/models"
)

func (s Service) CreateProject(project models.Project) (models.Project, error) {
	return s.pr.Create(context.TODO(), project)
}

func (s Service) UpdateProject(project models.Project) error {
	return s.pr.Update(context.TODO(), project)
}

func (s Service) DeleteProject(project models.Project) error {
	return s.pr.Delete(context.TODO(), project)
}

func (s Service) GetProject(project models.Project) (*models.Project, error) {
	return s.pr.GetOne(context.TODO(), project)
}

func (s Service) GetAllProjects() ([]models.Project, error) {
	return s.pr.GetAll(context.TODO())
}
