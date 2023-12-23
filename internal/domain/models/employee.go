package models

type Employee struct {
	ID         uint
	Name       string
	Surname    string
	Age        int
	Salary     int
	Department Department
}
