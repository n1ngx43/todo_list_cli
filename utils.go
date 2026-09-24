package main

import (
	"encoding/json"
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
