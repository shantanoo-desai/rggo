//go:build !inmemory && !containers
// +build !inmemory,!containers

package repository

import (
	"database/sql"
	"sync"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/shantanoo-desai/rggo/pomo/pomodoro"
)

const (
	createTableInterval string = `CREATE TABLE IF NOT EXISTS "interval" (
	"id" INTEGER,
	"start_time" DATETIME NOT NULL,
	"planned_duration" INTEGER DEFAULT 0,
	"actual_duration" INTEGER DEFAULT 0,
	"category" TEXT NOT NULL,
	"state" INTEGER DEFAULT 1,
	PRIMARY KEY("id")
	);`
)

type dbRepo struct {
	db *sql.DB
	sync.RWMutex
}

func NewSQLite3Repo(dbfile string) (*dbRepo, error) {
	db, err := sql.Open("sqlite3", dbfile)
	if err != nil {
		return nil, err
	}

	db.SetConnMaxLifetime(30 * time.Minute)
	db.SetMaxOpenConns(1)

	if err := db.Ping(); err != nil {
		return nil, err
	}

	if _, err := db.Exec(createTableInterval); err != nil {
		return nil, err
	}

	return &dbRepo{
		db: db,
	}, nil
}

func (r *dbRepo) Create(i pomodoro.Interval) (int64, error) {
	r.Lock()
	defer r.Unlock()

	insStmt, err := r.db.Prepare("INSERT INTO interval VALUES(NULL, ?,?,?,?,?)")
	if err != nil {
		return 0, err
	}

	defer insStmt.Close()

	res, err := insStmt.Exec(i.StartTime, i.PlannedDuration, i.ActualDuration, i.Category, i.State)
	if err != nil {
		return 0, err
	}

	var id int64
	if id, err = res.LastInsertId(); err != nil {
		return 0, err
	}

	return id, nil
}

func (r *dbRepo) Update(i pomodoro.Interval) error {
	r.Lock()
	defer r.Unlock()

	updStmt, err := r.db.Prepare(
		"UPDATE interval SET start_time=?,actual_duration=?,state=? WHERE id=?")

	if err != nil {
		return err
	}

	defer updStmt.Close()

	res, err := updStmt.Exec(i.StartTime, i.ActualDuration, i.State, i.ID)
	if err != nil {
		return err
	}

	_, err = res.RowsAffected()
	return err
}

func (r *dbRepo) ByID(id int64) (pomodoro.Interval, error) {
	r.RLock()
	defer r.RUnlock()

	row := r.db.QueryRow("SELECT * FROM interval WHERE id=?", id)

	i := pomodoro.Interval{}

	err := row.Scan(&i.ID, &i.StartTime, &i.PlannedDuration, &i.ActualDuration, &i.Category, &i.State)

	return i, err
}

func (r *dbRepo) Last() (pomodoro.Interval, error) {
	r.RLock()
	defer r.RUnlock()

	last := pomodoro.Interval{}

	err := r.db.QueryRow("SELECT * FROM interval ORDER BY id desc LIMIT 1").
		Scan(
			&last.ID, &last.StartTime, &last.PlannedDuration,
			&last.ActualDuration, &last.Category, &last.State,
		)
	if err == sql.ErrNoRows {
		return last, pomodoro.ErrNoIntervals
	}

	if err != nil {
		return last, err
	}

	return last, nil
}

func (r *dbRepo) Breaks(n int) ([]pomodoro.Interval, error) {
	r.RLock()
	defer r.RUnlock()

	stmt := `SELECT * FROM interval WHERE category LIKE '%Break'
	ORDER BY id DESC LIMIT ?`

	rows, err := r.db.Query(stmt, n)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	data := []pomodoro.Interval{}

	for rows.Next() {
		i := pomodoro.Interval{}
		err = rows.Scan(&i.ID, &i.StartTime, &i.PlannedDuration, &i.ActualDuration, &i.Category, &i.State)
		if err != nil {
			return nil, err
		}

		data = append(data, i)
	}

	err = rows.Err()

	if err != nil {
		return nil, err
	}

	return data, nil
}

func (r *dbRepo) CategorySummary(day time.Time, filter string) (time.Duration, error) {
	r.RLock()
	defer r.RUnlock()

	stmt := `SELECT sum(actual_duration) FROM interval
	WHERE category LIKE ? AND
	strftime('%Y-%m-%d', start_time, 'localtime')=
	strftime('%Y-%m-%d', ?, 'localtime')`

	var ds sql.NullInt64
	err := r.db.QueryRow(stmt, filter, day).Scan(&ds)

	var d time.Duration

	if ds.Valid {
		d = time.Duration(ds.Int64)
	}
	return d, err
}
