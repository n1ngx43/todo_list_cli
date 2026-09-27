package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

func readFile(fileName string) ([]Task, error) {
	// Read file
	file, err := os.Open(filepath.Join("data", fileName+".json"))
	if err != nil {
		return nil, err
	}
	defer file.Close()
	var taskList []Task
	err = json.NewDecoder(file).Decode(&taskList)
	if err != nil {
		if err == io.EOF {
			return []Task{}, nil
		}
		return nil, err
	}
	return taskList, nil
}

func writeFile(taskList []Task, fileName string) error {
	file, err := os.OpenFile(filepath.Join("data", fileName+".json"), os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	defer file.Close()
	err = json.NewEncoder(file).Encode(taskList)
	if err != nil {
		return err
	}
	return nil
}

func timeToString(timeTime time.Time) string {
	return timeTime.Format(time.RFC3339)
}

func stringToTimeWithDate(timeString string, targetDate time.Time) (time.Time, error) {
	parsedTime, err := time.Parse(time.TimeOnly, timeString)
	if err != nil {
		return parsedTime, err
	}
	date := time.Date(targetDate.Year(), targetDate.Month(), targetDate.Day(), parsedTime.Hour(), parsedTime.Minute(), parsedTime.Second(), 0, targetDate.Location())
	return date, nil
}

func stringToTimeToday(timeString string) (time.Time, error) {
	return stringToTimeWithDate(timeString, time.Now())
}

func formatDate(dateString string) (string, error) {
	// Format date from dd/mm/yyyy to yyyy-mm-dd
	dateTime, err := time.Parse("02/01/2006", dateString)
	if err != nil {
		return "", err
	}
	return dateTime.Format(time.DateOnly), nil
}

func generateID(taskList []Task) string {
	return fmt.Sprintf("t%d", len(taskList)+1)
}

func carryOverYesterdayTasks() error {
	now := time.Now()
	yesterday := now.AddDate(0, 0, -1)
	yesterdayFileName := yesterday.Format(time.DateOnly)
	todayFileName := now.Format(time.DateOnly)

	yesterdayTasks, err := readFile(yesterdayFileName)
	if err != nil || len(yesterdayTasks) == 0 {
		return nil
	}

	var tasksToCarryOver []Task
	hasModifiedYesterday := false

	for _, t := range yesterdayTasks {
		if t.Status == "To-do" || t.Status == "Inprogress" {
			tasksToCarryOver = append(tasksToCarryOver, t)
			t.Status = "Done"
			hasModifiedYesterday = true
		}
	}

	if !hasModifiedYesterday {
		return nil
	}

	err = writeFile(yesterdayTasks, yesterdayFileName)
	if err != nil {
		return errors.New("failed to update yesterday tasks!")
	}

	todayTasks, err := readFile(todayFileName)
	if err != nil {
		todayTasks = []Task{}
	}

	for _, task := range tasksToCarryOver {
		newStartTime := time.Date(now.Year(), now.Month(), now.Day(), task.StartTime.Hour(), task.StartTime.Minute(), task.StartTime.Second(), 0, now.Location())
		newEndTime := time.Date(now.Year(), now.Month(), now.Day(), task.EndTime.Hour(), task.EndTime.Minute(), task.EndTime.Second(), 0, now.Location())
		newTask := Task{
			ID:        generateID(todayTasks),
			Name:      task.Name,
			Status:    task.Status,
			Priority:  task.Priority,
			StartTime: newStartTime,
			EndTime:   newEndTime,
		}
		todayTasks = append(todayTasks, newTask)
	}

	err = writeFile(todayTasks, todayFileName)
	if err != nil {
		return errors.New("failed to save today tasks!")
	}
	return nil
}

func startMidnightWatcher() {
	go func() {
		for {
			now := time.Now()
			nextMidnight := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 1, 0, now.Location())
			time.Sleep(time.Until(nextMidnight))
			_ = carryOverYesterdayTasks()
		}
	}()
}