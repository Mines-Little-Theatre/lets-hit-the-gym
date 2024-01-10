package queries

import (
	"embed"
	"fmt"
)

//go:embed *.sql migrations/*.sql
var fs embed.FS

// Get returns the text of the query with the given name, or panics if it does not exist.
func Get(name string) string {
	data, err := fs.ReadFile(name + ".sql")
	if err != nil {
		panic(err)
	}
	return string(data)
}

// GetMigration returns the text of the migration to the given user version from its predecessor, or panics if it does not exist.
func GetMigration(userVersion uint32) string {
	data, err := fs.ReadFile(fmt.Sprintf("migrations/%02d.sql", userVersion))
	if err != nil {
		panic(err)
	}
	return string(data)
}
