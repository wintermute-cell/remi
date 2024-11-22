package main

import (
	"fmt"
	"os"
	"remi/pkg/actions"
	"remi/pkg/model"
	"remi/pkg/parsing"
	"remi/pkg/remi_errors"
	"remi/pkg/storage"
	"strconv"
	"time"
)

func run() error {
	if len(os.Args) == 1 {
		reminders, err := actions.ListRemindersForNow(time.Now())
		if err != nil {
			return err
		}
		for i, reminder := range reminders {
			fmt.Printf("%d: %s %s\n", i, reminder.Timestamp, reminder.Message)
		}
		return nil
	}

	switch os.Args[1] {
	case "add", "a", "+":
		timestamp := os.Args[2]
		t, err := parsing.ParseTimestamp(timestamp)
		if err != nil {
			return err
		}
		messageOrHowSoon := os.Args[3]
		howSoon, err := parsing.ParseDuration(messageOrHowSoon)
		if err == remi_errors.NotADurationError {
			// Does not have a howSoon duration, so it must be a message.
			message, err := parsing.CollectMessageFromArgs(os.Args, 3)
			if err != nil {
				return err
			}
			t.ShowFrom = time.Now()
			err = actions.AddReminder(model.Reminder{Timestamp: t, Message: message})
			if err != nil {
				return err
			}
		} else if err == nil {
			// Has a howSoon duration, so it must be a message.
			if len(os.Args) < 4 {
				return remi_errors.MissingMessageError
			}
			message, err := parsing.CollectMessageFromArgs(os.Args, 4)
			if err != nil {
				return err
			}
			t.ShowFrom = t.DateTime.Add(-howSoon)
			err = actions.AddReminder(model.Reminder{Timestamp: t, Message: message})
			if err != nil {
				return err
			}
		} else {
			return err
		}
	case "list", "l", "ls":
		reminders, err := actions.ListReminders()
		if err != nil {
			return err
		}
		for i, reminder := range reminders {
			fmt.Printf("%d: %s %s\n", i, reminder.Timestamp, reminder.Message)
		}
	case "remove", "r", "rm", "delete", "d", "-":
		if len(os.Args) < 3 {
			return remi_errors.MissingReminderIdError
		}
		removeIdx, err := strconv.Atoi(os.Args[2])
		if err != nil {
			return err
		}
		err = actions.RemoveReminder(removeIdx)
		if err != nil {
			return err
		}
	default:
		fmt.Fprintln(os.Stderr, "Invalid command")
		printUsage()
	}

	return nil
}

func printUsage() {
	fmt.Println("Usage: go run main.go <input>")
	// TODO: print proper usage
}

func main() {
	if err := storage.EnsureSqliteExists(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		printUsage()
	}

	storage.CloseSqlite()
}
