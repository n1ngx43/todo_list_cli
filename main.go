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

// -- UTIL FUNCTIONS ----------
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

func writeFile(taskList []Task) error {
	fileName := time.Now().Format(time.DateOnly)
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

func stringToTime(timeString string) (time.Time, error) {
	now := time.Now()
	parsedTime, err := time.Parse(time.TimeOnly, timeString)
	if err != nil {
		return parsedTime, err
	}
	date := time.Date(now.Year(), now.Month(), now.Day(), parsedTime.Hour(), parsedTime.Minute(), parsedTime.Second(), 0, now.Location())
	return date, nil
}

func formatDate(dateString string) (string, error) {
	// Format date from dd/mm/yyyy to yyyy/mm/dd
	dateTime, err := time.Parse("02/01/2006", dateString)
	if err != nil {
		return "", err
	}
	return dateTime.Format(time.DateOnly), nil
}

func generateID(taskList []Task) string {
	return fmt.Sprintf("t%d", len(taskList)+1)
}

// -- END ----------

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

func inputFromKeyboard(task *Task) (string, string, string, string, string) {
	reader := bufio.NewReader(os.Stdin)
	var name, status, priority, startTime, endTime string

	if task == nil || *task == (Task{}) {
		fmt.Print(">. Enter name: ")
		name, _ = reader.ReadString('\n')
		name = strings.TrimSpace(name)

		statusList := []string{"To-do", "Inprogress", "Done"}
		for {
			fmt.Print(">. Enter status: ")
			status, _ = reader.ReadString('\n')
			status = strings.TrimSpace(status)
			isValid := false
			for _, s := range statusList {
				if strings.EqualFold(s, status) {
					status = s
					isValid = true
					break
				}
			}
			if isValid {
				break
			}
		}

		priorityList := []string{"High", "Medium", "Low"}
		for {
			fmt.Print(">. Enter priority: ")
			priority, _ = reader.ReadString('\n')
			priority = strings.TrimSpace(priority)
			isValid := false
			for _, p := range priorityList {
				if strings.EqualFold(p, priority) {
					priority = p
					isValid = true
					break
				}
			}
			if isValid {
				break
			}
		}

		for {
			fmt.Print(">. Enter start time: ")
			startTime, _ = reader.ReadString('\n')
			startTime = strings.TrimSpace(startTime)
			_, err := time.Parse("15:04:05", startTime)
			if err == nil {
				break
			}
		}

		for {
			fmt.Print(">. Enter end time: ")
			endTime, _ = reader.ReadString('\n')
			endTime = strings.TrimSpace(endTime)
			_, err := time.Parse("15:04:05", endTime)
			if err == nil {
				break
			}
		}
	} else {
		fmt.Printf(">. Enter name [%s] (hit enter to skip): ", task.Name)
		name, _ = reader.ReadString('\n')
		name = strings.TrimSpace(name)

		statusList := []string{"To-do", "Inprogress", "Done"}
		for {
			fmt.Printf(">. Enter status [%s] (hit enter to skip): ", task.Status)
			status, _ = reader.ReadString('\n')
			status = strings.TrimSpace(status)
			isValid := false
			for _, s := range statusList {
				if strings.EqualFold(s, status) {
					status = s
					isValid = true
					break

				}
			}
			if isValid {
				break
			}
		}

		priorityList := []string{"High", "Medium", "Low"}
		for {
			fmt.Printf(">. Enter priority [%s] (hit enter to skip): ", task.Priority)
			priority, _ = reader.ReadString('\n')
			priority = strings.TrimSpace(priority)
			isValid := false
			for _, p := range priorityList {
				if strings.EqualFold(p, priority) {
					priority = p
					isValid = true
					break
				}
			}
			if isValid {
				break
			}
		}

		for {
			fmt.Printf(">. Enter start time [%s] (hit enter to skip): ", task.StartTime)
			startTime, _ = reader.ReadString('\n')
			startTime = strings.TrimSpace(startTime)
			_, err := time.Parse("15:04:05", startTime)
			if err == nil {
				break
			}
		}

		for {
			fmt.Printf(">. Enter end time [%s] (hit enter to skip): ", task.EndTime)
			endTime, _ = reader.ReadString('\n')
			endTime = strings.TrimSpace(endTime)
			_, err := time.Parse("15:04:05", endTime)
			if err == nil {
				break
			}
		}
	}

	return name, status, priority, startTime, endTime
}

func addTask() error {
	fileName := time.Now().Format(time.DateOnly)
	taskList, err := readFile(fileName)
	if err != nil {
		taskList = []Task{}
	}
	id := generateID(taskList)
	name, status, priority, startTimeStr, endTimeStr := inputFromKeyboard(nil)
	startTime, err := stringToTime(startTimeStr)
	if err != nil {
		return err
	}
	endTime, err := stringToTime(endTimeStr)
	if err != nil {
		return err
	}
	newTassk := Task{
		ID: id,
		Name: name,
		Status: status,
		Priority: priority,
		StartTime: startTime,
		EndTime: endTime,
	}
	taskList = append(taskList, newTassk)
	err = writeFile(taskList)
	if err != nil {
		return err
	}
	return nil
}

func controller() {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Println("-- TODO LIST CLI (1.0.0) ----------")
		fmt.Println("1. Task by today")
		fmt.Println("2. Task by date")
		fmt.Println("3. Add task")
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
		case "3":
			err := addTask()
			if err != nil {
				fmt.Fprintln(os.Stderr, "Failed to add task!")
			}
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
