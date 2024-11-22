package storage

import (
	"database/sql"
	"fmt"
	"os"
	"path"
	"remi/pkg/model"
	"remi/pkg/remi_errors"
	"time"

	_ "github.com/mattn/go-sqlite3"

	appdir "github.com/emersion/go-appdir"
)

var sqlitePath string = appdir.New("remi").UserData()
var sqliteFile string = "remi.db"

var idFile string = "last_ids.txt" // Stores the associations for the last reminder listing.

var DB *sql.DB

var schema = `
CREATE TABLE IF NOT EXISTS reminders (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	timestamp INTEGER NOT NULL,
	message TEXT NOT NULL,
	has_time BOOLEAN NOT NULL,
	show_from INTEGER NOT NULL
);
`

// An sqlite compatible Reminder type.
type dbReminder struct {
	id        int
	timestamp int64
	message   string
	hasTime   bool
	showFrom  int64
}

func setIdFile(reminders []dbReminder) error {
	file, err := os.Create(path.Join(sqlitePath, idFile))
	if err != nil {
		return err
	}
	defer file.Close()

	file.Truncate(0)
	for i, reminder := range reminders {
		_, err = file.WriteString(fmt.Sprintf("%d:%d\n", i, reminder.id))
		if err != nil {
			return err
		}
	}

	return nil
}

func EnsureSqliteExists() error {
	if err := os.MkdirAll(sqlitePath, 0755); err != nil {
		return err
	}

	var err error
	DB, err = sql.Open("sqlite3", path.Join(sqlitePath, sqliteFile))
	if err != nil {
		return err
	}

	_, err = DB.Exec(schema)
	if err != nil {
		return err
	}

	return nil
}

func CloseSqlite() error {
	return DB.Close()
}

func StoreReminder(reminder model.Reminder) error {

	datetimeString := reminder.Timestamp.DateTime.UnixMilli()
	showFromString := reminder.Timestamp.ShowFrom.UnixMilli()

	_, err := DB.Exec(`
		INSERT INTO reminders
		(timestamp, message, has_time, show_from)
		VALUES (?, ?, ?, ?)`,
		datetimeString,
		reminder.Message,
		reminder.Timestamp.HasTime,
		showFromString,
	)
	return err
}

func deserializeReminders(rows *sql.Rows) ([]model.Reminder, error) {
	dbReminders := make([]dbReminder, 0)
	for rows.Next() {
		var reminder dbReminder
		err := rows.Scan(
			&reminder.id,
			&reminder.timestamp,
			&reminder.message,
			&reminder.hasTime,
			&reminder.showFrom,
		)
		if err != nil {
			return nil, err
		}

		dbReminders = append(dbReminders, reminder)
	}

	setIdFile(dbReminders)

	reminders := make([]model.Reminder, 0)
	for _, reminder := range dbReminders {
		// datetime, err := time.Parse(time.RFC3339, reminder.timestamp)
		// if err != nil {
		// 	return nil, err
		// }
		datetime := time.UnixMilli(reminder.timestamp)
		// showFrom, err := time.Parse(time.RFC3339, reminder.showFrom)
		// if err != nil {
		// 	return nil, err
		// }
		showFrom := time.UnixMilli(reminder.showFrom)
		reminders = append(reminders, model.Reminder{
			Timestamp: model.Timestamp{
				DateTime: datetime,
				HasTime:  reminder.hasTime,
				ShowFrom: showFrom,
			},
			Message: reminder.message,
		})
	}
	return reminders, nil
}

func GetAllReminders() ([]model.Reminder, error) {
	rows, err := DB.Query("SELECT * FROM reminders ORDER BY timestamp ASC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	reminders, err := deserializeReminders(rows)
	if err != nil {
		return nil, err
	}

	return reminders, nil
}

func GetRemindersForNow(t time.Time) ([]model.Reminder, error) {
	rows, err := DB.Query("SELECT * FROM reminders WHERE show_from <= ? ORDER BY timestamp ASC", t.UnixMilli())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	reminders, err := deserializeReminders(rows)
	if err != nil {
		return nil, err
	}

	return reminders, nil
}

func DeleteReminder(id int) error {

	idMap := make(map[int]int)
	file, err := os.Open(path.Join(sqlitePath, idFile))
	if err != nil {
		return err
	}
	defer file.Close()
	for {
		var i, id int
		_, err := fmt.Fscanf(file, "%d:%d\n", &i, &id)
		if err != nil {
			break
		}
		idMap[i] = id
	}

	idToDelete, ok := idMap[id]
	if !ok {
		return remi_errors.InvalidReminderIdError
	}

	_, err = DB.Exec("DELETE FROM reminders WHERE id = ?", idToDelete)
	return err
}
