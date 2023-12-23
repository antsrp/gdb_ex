package main

import (
	"fmt"
	"log"

	myctx "github.com/antsrp/gdb_ex/internal/context"
	"github.com/antsrp/gdb_ex/internal/usecases/injections"
	"github.com/antsrp/gdb_ex/pkg/helpers"
)

func handlePanic() {
	if err := recover(); err != nil {
		log.Println("program is ending with error: ", err)
	}
}

func main() {
	defer handlePanic()

	logger, err := injections.BuildZapLogger()
	if err != nil {
		log.Fatal(err)
	}

	settings, err := injections.BuildConnectionSettings("DB")
	if err != nil {
		logger.Fatal(err.Error())
	}

	var txKey myctx.TxKey = "tx"
	connection, err := injections.BuildConnectionPgx(settings, logger, txKey)
	if err != nil {
		logger.Fatal(err.Error())
	}
	defer helpers.HandleCloser(logger, "database connection", connection)

	srv, _, err := injections.BuildServicePgx(connection, logger, txKey)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	employees, err := srv.GetAllEmployees()
	if err != nil {
		logger.Error("can't get employees: %s", err.Error())
		return
	}

	fmt.Println("employees:")
	fmt.Printf("%-3s\t%-30s\t%-30s\t%-3s\t%-10s\t%-30s\n", "id", "name", "surname", "age", "salary", "department")
	for _, e := range employees {
		fmt.Printf("%-3d\t%-30s\t%-30s\t%-3d\t%-10d\t%-30s\n", e.ID, e.Name, e.Surname, e.Age, e.Salary, e.Department.Name)
	}
}
