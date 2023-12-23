package service_test

import (
	"log"
	"testing"

	source "github.com/antsrp/gdb_ex"
	myctx "github.com/antsrp/gdb_ex/internal/context"
	"github.com/antsrp/gdb_ex/internal/domain/models"
	idb "github.com/antsrp/gdb_ex/internal/interfaces/db"
	"github.com/antsrp/gdb_ex/internal/usecases/injections"
	"github.com/stretchr/testify/require"

	"github.com/antsrp/gdb_ex/internal/usecases/service"
	"github.com/antsrp/gdb_ex/pkg/helpers"
)

var (
	srv *service.Service
)

func TestMain(m *testing.M) {

	logger, err := injections.BuildMockLogger()
	if err != nil {
		log.Fatal(err)
	}

	settings, err := injections.BuildConnectionSettings("DB", "test.env")
	if err != nil {
		logger.Fatal(err.Error())
	}

	var txKey myctx.TxKey = "tx"
	connection, err := injections.BuildConnectionPgx(settings, logger, txKey)
	if err != nil {
		logger.Fatal(err.Error())
	}
	defer helpers.HandleCloser(logger, "database connection", connection)

	srv, _, err = injections.BuildServicePgx(connection, logger, txKey)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	pgxMigrationTool, err := injections.BuildGooseMigrator(source.Migrations, settings.Type, connection)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	folder := "migrations/postgres"

	if err = pgxMigrationTool.Up(folder); err != nil {
		logger.Error(err.Error())
		return
	}

	defer func() {
		if err = pgxMigrationTool.DownAll(folder); err != nil {
			logger.Error(err.Error())
			return
		}
	}()

	m.Run()
}

func TestInsertDepartments(t *testing.T) {

	deps := []models.Department{
		{Name: "IT"},
		{Name: "Staff"},
		{Name: "Legal"},
	}

	for _, d := range deps {
		_, actual := srv.CreateDepartment(d)
		require.NoError(t, actual)
	}
}

func TestInsertEmployees(t *testing.T) {

	emps := []models.Employee{
		{Name: "James", Surname: "Bond", Age: 30, Salary: 10000, Department: models.Department{ID: 3}},
		{Name: "Mark", Surname: "Twain", Age: 58, Salary: 65986, Department: models.Department{ID: 2}},
		{Name: "Thierry", Surname: "Henry", Age: 49, Salary: 212435, Department: models.Department{ID: 2}},
		{Name: "Magnus", Surname: "Carlsen", Age: 33, Salary: 99999, Department: models.Department{ID: 1}},
		{Name: "Cristiano", Surname: "Ronaldo", Age: 37, Salary: 77777, Department: models.Department{ID: 1}},
		{Name: "Kylian", Surname: "Mbappe", Age: 24, Salary: 101212, Department: models.Department{ID: 3}},
		{Name: "Lionel", Surname: "Messi", Age: 35, Salary: 88881, Department: models.Department{ID: 1}},
	}

	for _, e := range emps {
		_, actual := srv.CreateEmployee(e)
		require.NoError(t, actual)
	}
}

func TestInsertProjects(t *testing.T) {

	projects := []models.Project{
		{Name: "Project1"},
		{Name: "Project2"},
		{Name: "Project3"},
		{Name: "Project4"},
	}

	for _, p := range projects {
		_, actual := srv.CreateProject(p)
		require.NoError(t, actual)
	}
}

func TestInsertProjectsEmployees(t *testing.T) {

	tests := []struct {
		p models.Project
		e models.Employee
	}{
		{models.Project{ID: 1}, models.Employee{ID: 3}},
		{models.Project{ID: 1}, models.Employee{ID: 6}},
		{models.Project{ID: 3}, models.Employee{ID: 6}},
		{models.Project{ID: 4}, models.Employee{ID: 5}},
		{models.Project{ID: 3}, models.Employee{ID: 4}},
		{models.Project{ID: 3}, models.Employee{ID: 2}},
		{models.Project{ID: 2}, models.Employee{ID: 2}},
		{models.Project{ID: 4}, models.Employee{ID: 1}},
	}

	for _, test := range tests {
		actual := srv.AddEmployeeToProject(test.e, test.p)
		require.NoError(t, actual)
	}
}

func TestSelectEmployees(t *testing.T) {
	expected := []models.Employee{
		{ID: 1, Name: "James", Surname: "Bond", Age: 30, Salary: 10000, Department: models.Department{ID: 3, Name: "Legal"}},
		{ID: 2, Name: "Mark", Surname: "Twain", Age: 58, Salary: 65986, Department: models.Department{ID: 2, Name: "Staff"}},
		{ID: 3, Name: "Thierry", Surname: "Henry", Age: 49, Salary: 212435, Department: models.Department{ID: 2, Name: "Staff"}},
		{ID: 4, Name: "Magnus", Surname: "Carlsen", Age: 33, Salary: 99999, Department: models.Department{ID: 1, Name: "IT"}},
		{ID: 5, Name: "Cristiano", Surname: "Ronaldo", Age: 37, Salary: 77777, Department: models.Department{ID: 1, Name: "IT"}},
		{ID: 6, Name: "Kylian", Surname: "Mbappe", Age: 24, Salary: 101212, Department: models.Department{ID: 3, Name: "Legal"}},
		{ID: 7, Name: "Lionel", Surname: "Messi", Age: 35, Salary: 88881, Department: models.Department{ID: 1, Name: "IT"}},
	}

	actual, err := srv.GetAllEmployees()

	require.NoError(t, err)

	require.Equal(t, len(expected), len(actual))

	n := len(expected)

	for i := 0; i < n; i++ {
		require.Equalf(t, expected[i], actual[i], "employee %d, expected: %v, actual: %v\n", i+1, expected[i], actual[i])
	}
}

func TestUpdateEmployee(t *testing.T) {
	expected := models.Employee{ID: 1, Name: "Robert"}

	err := srv.UpdateEmployee(expected)

	require.NoError(t, err)

	actual, err := srv.GetEmployee(models.Employee{ID: 1})

	require.NoError(t, err)

	require.Equalf(t, expected.Name, actual.Name, "employee name expected: %s, actual: %s\n", expected.Name, actual.Name)
}

func TestDeleteEmptyProject(t *testing.T) {
	projToDel := models.Project{ID: 16}

	err := srv.DeleteProject(projToDel)

	require.EqualError(t, err, idb.ErrNoRowsDelete.Error())
}

func TestDeleteProject(t *testing.T) {
	projToDel := models.Project{ID: 2}

	err := srv.DeleteProject(projToDel)

	require.NoError(t, err)

	expected := []models.Project{
		{ID: 1, Name: "Project1"},
		{ID: 3, Name: "Project3"},
		{ID: 4, Name: "Project4"},
	}

	actual, err := srv.GetAllProjects()
	require.NoError(t, err)

	require.Equal(t, expected, actual)
}

var (
	lastEmployeeId uint
)

func TestAssignEmployeeToNewDepartmentFail(t *testing.T) {
	employee := models.Employee{Name: "Sam", Surname: "Shankland", Age: 32, Salary: 323234, Department: models.Department{
		ID: 1,
	}}
	em, err := srv.CreateEmployee(employee)

	require.NoError(t, err)

	lastEmployeeId = em.ID

	dep := models.Department{Name: "Staff"}

	employee = models.Employee{ID: lastEmployeeId}

	err = srv.AssignEmployeeToDepartment(employee, dep)

	require.ErrorIs(t, err, idb.ErrExistingNameDepartment)

	actual, err := srv.GetEmployee(employee)

	require.NoError(t, err)

	require.Equal(t, uint(1), actual.Department.ID)
}

func TestAssignEmployeeToNewDepartment(t *testing.T) {
	dep := models.Department{Name: "Accounting"}

	employee := models.Employee{ID: lastEmployeeId}

	err := srv.AssignEmployeeToDepartment(employee, dep)

	require.NoError(t, err)

	actual, err := srv.GetEmployee(employee)

	require.NoError(t, err)

	require.NotEqual(t, 1, actual.Department.ID)
}
