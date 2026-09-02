package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"text/tabwriter"
	"time"
)

var priorityScore map[string]int = map[string]int{
	"High":   3,
	"Medium": 2,
	"Low":    1,
}

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

func formatDate(dateString string) (string, error) {
	// Format date from dd/mm/yyyy to yyyy/mm/dd
	dateTime, err := time.Parse("02/01/2006", dateString)
	if err != nil {
		return "", err
	}
	return dateTime.Format(time.DateOnly), nil
}

func sortPriorityAndStartDate(taskList []Task) []Task {
	// Sort task list priority and start time
	for i := 0; i < len(taskList)-1; i++ {
		best := i
		for j := i + 1; j < len(taskList); j++ {
			scoreBest := priorityScore[taskList[best].Priority]
			scoreJ := priorityScore[taskList[j].Priority]
			if scoreBest < scoreJ {
				best = j
			} else if scoreBest == scoreJ {
				if taskList[best].StartTime.After(taskList[j].StartTime) {
					best = j
				}
			}
		}
		if best != i {
			taskList[i], taskList[best] = taskList[best], taskList[i]
		}
	}
	return taskList
}

func displayTaskList(taskList []Task) {
	// Display task list
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)

	fmt.Fprintf(w, "|\t%s\t|\t%s\t|\t%s\t|\t%s\t|\t%s\t|\t%s\t|\n", "ID", "Name", "Status", "Priority", "Start Time", "End Time")
	for _, task := range taskList {
		fmt.Fprintf(w, "|\t%s\t|\t%s\t|\t%s\t|\t%s\t|\t%s\t|\t%s\t|\n", task.ID, task.Name, task.Status, task.Priority, task.StartTime, task.EndTime)
	}

	w.Flush()
}

func displayTaskListByToday() {
	fileName := time.Now().Format(time.DateOnly)
	dateString := time.Now().Format("02/01/2006")
	taskList, err := readFile(fileName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Can't read file!\nError: %s\n", err)
		return
	}
	if len(taskList) == 0 {
		fmt.Printf("Task list was empty in %s!\n", dateString)
		return
	}
	taskList = sortPriorityAndStartDate(taskList)
	fmt.Printf("Task List (%s)\n", dateString)
	displayTaskList(taskList)
}

func displayTaskListByDate(dateString string) {
	fileName, err := formatDate(dateString)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Invalid format date!\nError: %s\n", err)
		return
	}
	taskList, err := readFile(fileName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Can't read file! !\nError: %s\n", err)
		return
	}
	if len(taskList) == 0 {
		fmt.Printf("Task list was empty in %s!\n", dateString)
		return
	}
	taskList = sortPriorityAndStartDate(taskList)
	fmt.Printf("Task List (%s)\n", dateString)
	displayTaskList(taskList)
}

func controller() {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Println("-- TODO LIST CLI (1.0.0) ----------")
		fmt.Println("1. Task by today")
		fmt.Println("2. Task by date")
		fmt.Println("0. Exit")
		fmt.Print(">. Input your option: ")
		option, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, "Failed to read input!")
			continue
		}
		option = strings.TrimSpace(option)

		switch option {
		case "1":
			displayTaskListByToday()
		case "2":
			fmt.Print(">. Input date: ")
			dateString, err := reader.ReadString('\n')
			if err != nil {
				fmt.Fprintln(os.Stderr, "Failed to read input!")
				continue
			}
			dateString = strings.TrimSpace(dateString)
			if dateString == "" {
				fmt.Fprintln(os.Stderr, "Date was empty!")
				continue
			}
			displayTaskListByDate(dateString)
		case "0":
			fmt.Println("See you soon!")
			return
		default:
			fmt.Fprintln(os.Stderr, "Invalid option!")
		}

	}
}

func main() {
	controller()
}
