package gorm

import "github.com/antsrp/gdb_ex/internal/domain/models"

type DepartmentModel struct {
	ID   uint   `gorm:"primaryKey"`
	Name string `gorm:"constraint:UNIQUE"`
}

func (DepartmentModel) TableName() string {
	return "departments"
}

type EmployeeModel struct {
	ID           uint `gorm:"primaryKey"`
	Name         string
	Surname      string
	Age          int
	Salary       int
	DepartmentID uint            `gorm:"column:departmentid"`
	Department   DepartmentModel `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

func (EmployeeModel) TableName() string {
	return "employees"
}

func createEmployeeModel(e models.Employee) EmployeeModel {
	return EmployeeModel{
		ID:           e.ID,
		Name:         e.Name,
		Surname:      e.Surname,
		Age:          e.Age,
		Salary:       e.Salary,
		DepartmentID: e.Department.ID,
		Department:   DepartmentModel(e.Department),
	}
}

type getEmployeeStruct struct {
	ID      uint
	Name    string
	Surname string
	Age     int
	Salary  int
	DepID   uint   `gorm:"column:depid"`
	DepName string `gorm:"column:depname"`
}

type ProjectModel struct {
	ID   uint `gorm:"primaryKey"`
	Name string
}

func (ProjectModel) TableName() string {
	return "projects"
}

type ProjectEmployeeRelationModel struct {
	ProjectID  uint          `gorm:"column:projectid"`
	Project    ProjectModel  `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	EmployeeID uint          `gorm:"column:employeeid"`
	Employee   EmployeeModel `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

func (ProjectEmployeeRelationModel) TableName() string {
	return "projects_employees"
}
