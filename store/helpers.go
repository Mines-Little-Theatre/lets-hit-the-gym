package store

import (
	"database/sql"

	"github.com/Mines-Little-Theatre/lets-hit-the-gym/store/queries"
)

func getKV[T any](db *sql.DB, key string) (T, error) {
	row := db.QueryRow(queries.Get("get_kv"), key)
	var value T
	err := row.Scan(&value)
	return value, err
}

func putKV(db *sql.DB, key string, value any) error {
	_, err := db.Exec(queries.Get("put_kv"), key, value)
	return err
}
