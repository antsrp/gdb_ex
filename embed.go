package source

import "embed"

//go:embed migrations/postgres/*
var Migrations embed.FS
