package goose_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateWithPgx(t *testing.T) {
	actual := pgxMigrationTool.Up(folder)

	assert.NoError(t, actual)
}

func TestDropToWithPgx(t *testing.T) {
	actual := pgxMigrationTool.DownTo(folder, 20231220160554)

	assert.NoError(t, actual)
}

func TestUpToWithPgx(t *testing.T) {
	actual := pgxMigrationTool.UpTo(folder, 20231220160716)

	assert.NoError(t, actual)
}

func TestDropWithPgx(t *testing.T) {
	actual := pgxMigrationTool.DownAll(folder)

	assert.NoError(t, actual)
}
