package db

import (
	"fmt"

	"github.com/antsrp/gdb_ex/pkg/infrastructure/db"
	"github.com/antsrp/gdb_ex/pkg/parsers/env"
)

func InitSettings(prefix string, filenames ...string) (*db.Settings, error) {
	if err := env.Load(filenames...); err != nil {
		return nil, err
	}
	settings, err := env.Parse[db.Settings](prefix)
	if err != nil {
		return nil, fmt.Errorf("can't launch config of database: %w", err)
	}

	return settings, nil
}
