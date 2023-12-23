package postgres_test

import (
	"context"
	"log"
	"testing"

	source "github.com/antsrp/gdb_ex"
	myctx "github.com/antsrp/gdb_ex/internal/context"
	"github.com/antsrp/gdb_ex/internal/domain/models"
	idb "github.com/antsrp/gdb_ex/internal/interfaces/db"
	"github.com/antsrp/gdb_ex/internal/usecases/database/postgres"
	"github.com/antsrp/gdb_ex/internal/usecases/injections"
	"github.com/antsrp/gdb_ex/pkg/helpers"
	"github.com/stretchr/testify/require"
)

var (
	txKey myctx.TxKey = "tx"
	conn  idb.Connection
	ds    *postgres.DepartmentStorage
	es    *postgres.EmployeeStorage
	ps    *postgres.ProjectStorage
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

	connection, err := injections.BuildConnectionPgx(settings, logger, txKey)
	if err != nil {
		logger.Fatal(err.Error())
	}
	conn = connection
	defer helpers.HandleCloser(logger, "database connection", connection)

	storages, err := injections.BuildStoragesPgx(connection)
	if err != nil {
		logger.Error(err.Error())
	}
	ds, es, ps = storages.Ds, storages.Es, storages.Ps

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

	ctx := context.Background()

	for _, d := range deps {
		_, actual := ds.Create(ctx, d)
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

	ctx := context.Background()

	for _, e := range emps {
		_, actual := es.Create(ctx, e)
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

	ctx := context.Background()

	for _, p := range projects {
		_, actual := ps.Create(ctx, p)
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

	ctx := context.Background()

	for _, test := range tests {
		actual := es.AddToProject(ctx, test.e, test.p)
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

	actual, err := es.GetAll(context.Background())

	require.NoError(t, err)

	require.Equal(t, len(expected), len(actual))

	n := len(expected)

	for i := 0; i < n; i++ {
		require.Equalf(t, expected[i], actual[i], "employee %d, expected: %v, actual: %v\n", i+1, expected[i], actual[i])
	}
}

func TestUpdateEmployee(t *testing.T) {
	ctx := context.Background()
	expected := models.Employee{ID: 1, Name: "Robert"}

	err := es.Update(ctx, expected)

	require.NoError(t, err)

	actual, err := es.GetOne(ctx, models.Employee{ID: 1})

	require.NoError(t, err)

	require.Equalf(t, expected.Name, actual.Name, "employee name expected: %s, actual: %s\n", expected.Name, actual.Name)
}

func TestDeleteEmptyProject(t *testing.T) {
	ctx := context.Background()
	projToDel := models.Project{ID: 16}

	err := ps.Delete(ctx, projToDel)

	require.EqualError(t, err, idb.ErrNoRowsDelete.Error())
}

func TestDeleteProject(t *testing.T) {
	ctx := context.Background()
	projToDel := models.Project{ID: 2}

	err := ps.Delete(ctx, projToDel)

	require.NoError(t, err)

	expected := []models.Project{
		{ID: 1, Name: "Project1"},
		{ID: 3, Name: "Project3"},
		{ID: 4, Name: "Project4"},
	}

	actual, err := ps.GetAll(ctx)
	require.NoError(t, err)

	require.Equal(t, expected, actual)
}

func TestCreateBadDepInTransaction(t *testing.T) {
	tx, err := conn.CreateTransaction()

	require.NoError(t, err)

	defer tx.Rollback()

	ctxTx := context.WithValue(context.Background(), txKey, tx)

	project := models.Project{Name: "Some project"}
	p, err := ps.Create(ctxTx, project)

	require.NoError(t, err)

	id := p.ID

	dep := models.Department{Name: "Staff"}

	_, err = ds.Create(ctxTx, dep)

	require.ErrorIs(t, err, idb.ErrExistingNameDepartment)
	tx.Rollback()

	_, err = ps.GetOne(context.Background(), models.Project{ID: id})

	require.ErrorIs(t, err, idb.ErrNoRowsSelect)
}

func TestAssignEmployeeToNewDepartment(t *testing.T) {
	employee := models.Employee{Name: "Sam", Surname: "Shankland", Age: 32, Salary: 323234, Department: models.Department{
		ID: 1,
	}}
	em, err := es.Create(context.Background(), employee)
	require.NoError(t, err)

	tx, err := conn.CreateTransaction()
	defer tx.Rollback()

	require.NoError(t, err)

	dep := models.Department{Name: "Accounting"}

	ctxTx := context.WithValue(context.Background(), txKey, tx)
	d, err := ds.Create(ctxTx, dep)

	require.NoError(t, err)

	em.Department.ID = d.ID

	err = es.Update(ctxTx, em)

	require.NoError(t, err)

	tx.Commit()

	actual, err := es.GetOne(context.Background(), em)

	require.NoError(t, err)

	require.Equal(t, d.ID, actual.Department.ID)
}
