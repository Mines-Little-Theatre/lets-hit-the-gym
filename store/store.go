package store

import (
	"database/sql"
	"fmt"

	"github.com/Mines-Little-Theatre/lets-hit-the-gym/store/queries"

	_ "modernc.org/sqlite"
)

const (
	applicationID uint32 = 0x4c696654
	userVersion   uint32 = 3
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

func (s *Store) UpdateToken(token string) error {
	return putKV(s.db, "token", token)
}

func (s *Store) GetChannelID() (string, error) {
	return getKV[string](s.db, "channel_id")
}

func (s *Store) UpdateChannelID(id string) error {
	return putKV(s.db, "channel_id", id)
}

func (s *Store) GetLastScheduleMessageID() (string, error) {
	return getKV[string](s.db, "last_schedule_message_id")
}

func (s *Store) UpdateLastScheduleMessageID(id string) error {
	return putKV(s.db, "last_schedule_message_id", id)
}

func (s *Store) GetDay(name string) (*Day, error) {
	rows, err := s.db.Query(queries.Get("get_day"), name)
	if err != nil {
		return nil, err
	}
	if !rows.Next() {
		return nil, rows.Err()
	}

	day := new(Day)
	var (
		workoutTitle, workoutDescription *string
		workoutColor                     *int

		routineTitle, routineDescription *string
	)
	err = rows.Scan(&day.OpenHour, &day.CloseHour, &workoutTitle, &workoutDescription, &workoutColor, &routineTitle, &routineDescription)
	if err != nil {
		return nil, err
	}
	if workoutTitle != nil {
		day.Workout = &Workout{
			Title:       *workoutTitle,
			Description: *workoutDescription,
			Color:       *workoutColor,
		}
		if routineTitle != nil {
			day.Workout.Routines = append(day.Workout.Routines, Routine{
				Title:       *routineTitle,
				Description: *routineDescription,
			})
			for rows.Next() {
				var routine Routine
				err := rows.Scan(new(int), new(int), new(string), new(string), new(int), &routine.Title, &routine.Description)
				if err != nil {
					return nil, err
				}
				day.Workout.Routines = append(day.Workout.Routines, routine)
			}
			if rows.Err() != nil {
				return nil, rows.Err()
			}
			return day, nil
		}
	}

	err = rows.Close()
	if err != nil {
		return nil, err
	}

	return day, nil
}

func (s *Store) GetDayNames() ([]string, error) {
	result := make([]string, 0)
	rows, err := s.db.Query(queries.Get("get_day_names"))
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
