package migrations

type Migrator interface {
	Up(string) error
	Down(string) error
	UpTo(string, int64) error
	DownTo(string, int64) error
	DownAll(string) error
}
