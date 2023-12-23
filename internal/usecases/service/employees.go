package service

import (
	"context"

	"github.com/antsrp/gdb_ex/internal/domain/models"
)

func (s Service) CreateEmployee(employee models.Employee) (models.Employee, error) {
	return s.er.Create(context.TODO(), employee)
}

func (s Service) UpdateEmployee(employee models.Employee) error {
	return s.er.Update(context.TODO(), employee)
}

func (s Service) DeleteEmployee(employee models.Employee) error {
	return s.er.Delete(context.TODO(), employee)
}

func (s Service) GetEmployee(employee models.Employee) (*models.Employee, error) {
	return s.er.GetOne(context.TODO(), employee)
}

func (s Service) GetAllEmployees() ([]models.Employee, error) {
	return s.er.GetAll(context.TODO())
}

func (s Service) AddEmployeeToProject(employee models.Employee, project models.Project) error {
	return s.er.AddToProject(context.TODO(), employee, project)
}
