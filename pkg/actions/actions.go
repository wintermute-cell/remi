package actions

import (
	"remi/pkg/model"
	"remi/pkg/storage"
	"time"
)

// called with: "add" "a" "+"
// AddReminder adds a new reminder.
func AddReminder(reminder model.Reminder) error {
	err := storage.StoreReminder(reminder)
	return err
}

// called with: "list" "l" "ls"
// ListReminders lists all reminders.
func ListReminders() ([]model.Reminder, error) {
	reminders, err := storage.GetAllReminders()
	if err != nil {
		return []model.Reminder{}, err
	}
	return reminders, nil
}

// called with: "" (no arguments)
// ListRemindersForNow lists all reminders set to show for a specific time.
func ListRemindersForNow(day time.Time) ([]model.Reminder, error) {
	reminders, err := storage.GetRemindersForNow(day)
	if err != nil {
		return []model.Reminder{}, err
	}
	return reminders, nil
}

// called with: "remove" "r" "rm" "delete" "d" <reminderId>
// RemoveReminder removes a reminder.
func RemoveReminder(reminderId int) error {
	return storage.DeleteReminder(reminderId)
}
