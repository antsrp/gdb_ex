package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/antsrp/gdb_ex/internal/domain/models"
	irepo "github.com/antsrp/gdb_ex/internal/interfaces/db"
	"github.com/jackc/pgx/v5"
)

var _ irepo.ProjectsRepository = &ProjectStorage{}

const (
	createProjectStmtQuery  = "INSERT INTO Projects (name) VALUES ($1) RETURNING id;"
	updateProjectStmtQuery  = "UPDATE Projects SET name = $1 WHERE id = $2;"
	deleteProjectStmtQuery  = "DELETE FROM Projects WHERE id = $1;"
	getAllProjectsStmtQuery = "SELECT id, name FROM Projects"
	getOneProjectStmtQuery  = getAllProjectsStmtQuery + " WHERE id = $1;"
)

type ProjectStorage struct {
	conn *Connection
}

func NewProjectStorage(conn *Connection) (*ProjectStorage, error) {
	return &ProjectStorage{conn: conn}, nil
}

func (ps ProjectStorage) Create(ctx context.Context, project models.Project) (models.Project, error) {
	operand := presentOperand(ctx, ps.conn.logger, ps.conn.pc, ps.conn.txKey)

	row := operand.QueryRow(ps.conn.ctx, createProjectStmtQuery, project.Name)

	var id uint

	if err := row.Scan(&id); err != nil {
		return models.Project{}, fmt.Errorf("can't create project: %w", err)
	}

	return models.Project{ID: id}, nil
}

func (ps ProjectStorage) Update(ctx context.Context, project models.Project) error {
	operand := presentOperand(ctx, ps.conn.logger, ps.conn.pc, ps.conn.txKey)
	tag, err := operand.Exec(ps.conn.ctx, updateProjectStmtQuery, project.Name, project.ID)
	if err != nil {
		return fmt.Errorf("can't update project: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return irepo.ErrNoRowsUpdate
	}

	return nil
}

func (ps ProjectStorage) Delete(ctx context.Context, project models.Project) error {
	operand := presentOperand(ctx, ps.conn.logger, ps.conn.pc, ps.conn.txKey)
	tag, err := operand.Exec(ps.conn.ctx, deleteProjectStmtQuery, project.ID)
	if err != nil {
		return fmt.Errorf("can't delete project: %w", err)
	} else if tag.RowsAffected() == 0 {
		return irepo.ErrNoRowsDelete
	}

	return nil
}

func (ps ProjectStorage) GetOne(ctx context.Context, project models.Project) (*models.Project, error) {
	operand := presentOperand(ctx, ps.conn.logger, ps.conn.pc, ps.conn.txKey)
	row := operand.QueryRow(ps.conn.ctx, getOneProjectStmtQuery, project.ID)

	var p models.Project
	if err := row.Scan(&p.ID, &p.Name); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			err = irepo.ErrNoRowsSelect
		}
		return nil, fmt.Errorf("can't select project: %w", err)
	}

	return &p, nil
}

func (ps ProjectStorage) GetAll(ctx context.Context) ([]models.Project, error) {
	operand := presentOperand(ctx, ps.conn.logger, ps.conn.pc, ps.conn.txKey)
	rows, err := operand.Query(ps.conn.ctx, getAllProjectsStmtQuery)
	if err != nil {
		return nil, fmt.Errorf("can't select all projects: %w", err)
	}
	defer rows.Close()

	var projects []models.Project

	for rows.Next() {
		var p models.Project
		if err := rows.Scan(&p.ID, &p.Name); err != nil {
			return nil, fmt.Errorf("can't scan project: %w", err)
		}
		projects = append(projects, p)
	}

	return projects, nil
}
