package store

import (
	"database/sql"
	"errors"

	"github.com/Mines-Little-Theatre/lets-hit-the-gym/store/queries"
)

type querier interface {
	QueryRow(string, ...any) *sql.Row
	Exec(string, ...any) (sql.Result, error)
}

func transact(db *sql.DB, f func(*sql.Tx) error) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	err = f(tx)
	if err != nil {
		return errors.Join(err, tx.Rollback())
	}

	return tx.Commit()
}

// func transactQuery[T any](db *sql.DB, f func(*sql.Tx) (T, error)) (T, error) {
// 	tx, err := db.Begin()
// 	if err != nil {
// 		var zv T
// 		return zv, err
// 	}

// 	result, err := f(tx)
// 	if err != nil {
// 		return result, errors.Join(err, tx.Rollback())
// 	}

// 	return result, tx.Commit()
// }

func getKV[T any](q querier, key string) (T, error) {
	row := q.QueryRow(queries.Get("get_kv"), key)
	var value T
	err := row.Scan(&value)
	return value, err
}

func putKV(q querier, key string, value any) error {
	_, err := q.Exec(queries.Get("put_kv"), key, value)
	return err
}
