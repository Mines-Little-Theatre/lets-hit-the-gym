package store

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/Mines-Little-Theatre/lets-hit-the-gym/store/queries"

	_ "modernc.org/sqlite"
)

const (
	applicationID uint32 = 0x4c696654
	userVersion   uint32 = 4
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
		return nil, errors.Join(err, db.Close())
	}

	row := db.QueryRow("PRAGMA application_id;")
	var dbAppId uint32
	err = row.Scan(&dbAppId)
	if err != nil {
		return nil, errors.Join(err, db.Close())
	}

	if dbAppId != applicationID && dbAppId != 0 {
		return nil, errors.Join(
			fmt.Errorf("application_id mismatch: expected %d, but was %d", applicationID, dbAppId),
			db.Close(),
		)
	}

	row = db.QueryRow("PRAGMA user_version;")
	var dbUserVer uint32
	err = row.Scan(&dbUserVer)
	if err != nil {
		return nil, errors.Join(err, db.Close())
	}

	if dbUserVer > userVersion {
		return nil, errors.Join(
			fmt.Errorf("user_version is too high: expected %d or lower, but was %d", userVersion, dbUserVer),
			db.Close(),
		)
	} else if dbAppId == 0 && dbUserVer != 0 {
		return nil, errors.Join(
			fmt.Errorf("application id was zero but user version was nonzero (%d)", userVersion),
			db.Close(),
		)
	}

	for dbUserVer < userVersion {
		tx, err := db.Begin()
		if err != nil {
			return nil, errors.Join(err, db.Close())
		}
		_, err = tx.Exec(queries.GetMigration(dbUserVer + 1))
		if err != nil {
			return nil, errors.Join(err, tx.Rollback(), db.Close())
		}
		err = tx.Commit()
		if err != nil {
			return nil, errors.Join(err, db.Close())
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

func (s *Store) GetSignupEmbedIndex() (int, error) {
	return getKV[int](s.db, "signup_embed_index")
}

func (s *Store) UpdateLastScheduleMessage(id string, signupEmbedIndex int) error {
	return transact(s.db, func(tx *sql.Tx) error {
		err := putKV(tx, "last_schedule_message_id", id)
		if err != nil {
			return err
		}
		err = putKV(tx, "signup_embed_index", signupEmbedIndex)
		if err != nil {
			return err
		}
		_, err = tx.Exec(queries.Get("clear_arrivals"))
		if err != nil {
			return err
		}
		return nil
	})
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
					return nil, errors.Join(err, rows.Close())
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
			return nil, errors.Join(err, rows.Close())
		}
		result = append(result, s)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return result, nil
}

func (s *Store) GetAllArrivals() ([]HourArrivals, error) {
	result := make([]HourArrivals, 0)
	rows, err := s.db.Query(queries.Get("get_arrivals"))
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var (
			hour   int
			userID string
		)
		err = rows.Scan(&hour, &userID)
		if err != nil {
			return nil, errors.Join(err, rows.Close())
		}
		if len(result) == 0 || hour != result[len(result)-1].Hour {
			result = append(result, HourArrivals{Hour: hour})
		}
		users := &result[len(result)-1].ArrivingUsers
		*users = append(*users, userID)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return result, nil
}

func (s *Store) GetArrivingUsers(hour int) ([]string, error) {
	result := make([]string, 0)
	rows, err := s.db.Query(queries.Get("get_hour_arrivals"), hour)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var id string
		err := rows.Scan(&id)
		if err != nil {
			return nil, errors.Join(err, rows.Close())
		}
		result = append(result, id)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return result, nil
}

func (s *Store) SetUserArrivals(id string, hours []int) error {
	return transact(s.db, func(tx *sql.Tx) error {
		_, err := tx.Exec(queries.Get("remove_user_arrivals"), id)
		if err != nil {
			return err
		}
		for _, hour := range hours {
			_, err := tx.Exec(queries.Get("add_arrival"), id, hour)
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func (s *Store) ClearArrivals() error {
	_, err := s.db.Exec(queries.Get("clear_arrivals"))
	return err
}

func (s *Store) Close() error {
	return s.db.Close()
}
