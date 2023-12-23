package service

import (
	"context"

	myctx "github.com/antsrp/gdb_ex/internal/context"
	"github.com/antsrp/gdb_ex/internal/domain/models"
	"github.com/antsrp/gdb_ex/internal/interfaces/db"
	"github.com/antsrp/gdb_ex/internal/interfaces/logger"
)

type Service struct {
	dr     db.DepartmentsRepository
	er     db.EmployeesRepository
	pr     db.ProjectsRepository
	conn   db.Connection
	logger logger.Logger
	txKey  myctx.TxKey
}

func NewService(dr db.DepartmentsRepository, er db.EmployeesRepository, pr db.ProjectsRepository, conn db.Connection,
	logger logger.Logger, txKey myctx.TxKey) *Service {
	return &Service{
		dr:     dr,
		er:     er,
		pr:     pr,
		conn:   conn,
		logger: logger,
		txKey:  txKey,
	}
}

func (s Service) AssignEmployeeToDepartment(employee models.Employee, dep models.Department) error {
	// in transaction...
	tx, err := s.conn.CreateTransaction()
	if err != nil {
		return err
	}

	ctx := context.WithValue(context.Background(), s.txKey, tx)

	defer func() {
		if err := tx.Rollback(); err != nil {
			s.logger.Info(err.Error())
		}
	}()

	if d, err := s.dr.Create(ctx, dep); err != nil {
		return err
	} else {
		employee.Department.ID = d.ID
	}
	if err := s.er.Update(ctx, employee); err != nil {
		return err
	}

	return tx.Commit()
}
