package store

import (
	"database/sql"
	"fmt"

	"github.com/Mines-Little-Theatre/lets-hit-the-gym/store/queries"

	_ "modernc.org/sqlite"
)

const (
	applicationID uint32 = 0x4c696654
	userVersion   uint32 = 2
)

// much of this is copied from the lean bot, maybe I should make a library

type Store struct {
	db *sql.DB
}

func Open(dataSourceName string) (*Store, error) {
	db, err := sql.Open("sqlite", dataSourceName)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec("PRAGMA foreign_keys = ON;")
	if err != nil {
		db.Close()
		return nil, err
	}

	row := db.QueryRow("PRAGMA application_id;")
	var dbAppId uint32
	err = row.Scan(&dbAppId)
	if err != nil {
		db.Close()
		return nil, err
	}

	if dbAppId != applicationID && dbAppId != 0 {
		db.Close()
		return nil, fmt.Errorf("application_id mismatch: expected %d, but was %d", applicationID, dbAppId)
	}

	row = db.QueryRow("PRAGMA user_version;")
	var dbUserVer uint32
	err = row.Scan(&dbUserVer)
	if err != nil {
		db.Close()
		return nil, err
	}

	if dbUserVer > userVersion {
		db.Close()
		return nil, fmt.Errorf("user_version is too high: expected %d or lower, but was %d", userVersion, dbUserVer)
	} else if dbAppId == 0 && dbUserVer != 0 {
		db.Close()
		return nil, fmt.Errorf("application id was zero but user version was nonzero (%d)", userVersion)
	}

	for dbUserVer < userVersion {
		_, err := db.Exec(queries.GetMigration(dbUserVer + 1))
		if err != nil {
			db.Close()
			return nil, err
		}
		dbUserVer++
	}

	return &Store{db: db}, nil
}

func (s *Store) GetToken() (string, error) {
	return getKV[string](s.db, "token")
}

func (s *Store) GetChannelID() (string, error) {
	return getKV[string](s.db, "channel_id")
}

func (s *Store) GetLastScheduleMessageID() (string, error) {
	return getKV[string](s.db, "last_schedule_message_id")
}

func (s *Store) UpdateLastScheduleMessageID(id string) error {
	return putKV(s.db, "last_schedule_message_id", id)
}

func (s *Store) GetWorkout(name string) (Workout, error) {
	var result Workout

	row := s.db.QueryRow(queries.Get("get_workout"), name)
	err := row.Scan(&result.Title, &result.Description, &result.Color)
	if err != nil {
		return result, err
	}

	result.Routines = make([]Routine, 0)
	rows, err := s.db.Query(queries.Get("get_workout_routines"), name)
	if err != nil {
		return result, err
	}
	for rows.Next() {
		var r Routine
		err := rows.Scan(&r.Title, &r.Description)
		if err != nil {
			return result, err
		}
		result.Routines = append(result.Routines, r)
	}
	if err := rows.Err(); err != nil {
		return result, err
	}

	return result, nil
}

func (s *Store) GetWorkoutNames() ([]string, error) {
	result := make([]string, 0)
	rows, err := s.db.Query(queries.Get("get_workout_names"))
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var s string
		err := rows.Scan(&s)
		if err != nil {
			return nil, err
		}
		result = append(result, s)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return result, nil
}

func (s *Store) Close() error {
	return s.db.Close()
}
