package goose_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateWithGorm(t *testing.T) {
	actual := gormMigrationTool.Up(folder)

	assert.NoError(t, actual)
}

func TestDropToWithGorm(t *testing.T) {
	actual := gormMigrationTool.DownTo(folder, 20231220160554)

	assert.NoError(t, actual)
}

func TestUpToWithGorm(t *testing.T) {
	actual := gormMigrationTool.UpTo(folder, 20231220160716)

	assert.NoError(t, actual)
}

func TestDropWithGorm(t *testing.T) {
	actual := gormMigrationTool.DownAll(folder)

	assert.NoError(t, actual)
}
