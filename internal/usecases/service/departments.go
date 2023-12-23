package service

import (
	"context"

	"github.com/antsrp/gdb_ex/internal/domain/models"
)

func (s Service) CreateDepartment(dep models.Department) (models.Department, error) {
	return s.dr.Create(context.TODO(), dep)
}

func (s Service) UpdateDepartment(dep models.Department) error {
	return s.dr.Update(context.TODO(), dep)
}

func (s Service) DeleteDepartment(dep models.Department) error {
	return s.dr.Delete(context.TODO(), dep)
}

func (s Service) GetDepartment(dep models.Department) (*models.Department, error) {
	return s.dr.GetOne(context.TODO(), dep)
}

func (s Service) GetAllDepartments() ([]models.Department, error) {
	return s.dr.GetAll(context.TODO())
}
