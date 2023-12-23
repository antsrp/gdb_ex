package models

type Project struct {
	ID   uint
	Name string
}

type ProjectEmployeeRelation struct {
	Project  Project
	Employee Employee
}
