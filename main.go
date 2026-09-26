package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"
	"time"
)

var priorityScore map[string]int = map[string]int{
	"High":   3,
	"Medium": 2,
	"Low":    1,
}

func searchTask(taskList []Task, id string) *Task {
	for i, t := range taskList {
		if strings.EqualFold(id, t.ID) {
			return &taskList[i]
		}
	}
	return nil
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
	fmt.Printf("-- Task List (%s) ----------\n", dateString)
	displayTaskList(taskList)
}

func displayTaskListByDate(dateString string) {
	fileName, _ := formatDate(dateString)
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
	fmt.Printf("-- Task List (%s) ----------\n", dateString)
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
		if name == "" {
			name = task.Name
		}

		statusList := []string{"To-do", "Inprogress", "Done"}
		for {
			fmt.Printf(">. Enter status [%s] (hit enter to skip): ", task.Status)
			status, _ = reader.ReadString('\n')
			status = strings.TrimSpace(status)
			if status == "" {
				status = task.Status
			}
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
			if priority == "" {
				priority = task.Priority
			}
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
			if startTime == "" {
				startTime = task.StartTime.Format(time.TimeOnly)
			}
			_, err := time.Parse("15:04:05", startTime)
			if err == nil {
				break
			}
		}

		for {
			fmt.Printf(">. Enter end time [%s] (hit enter to skip): ", task.EndTime)
			endTime, _ = reader.ReadString('\n')
			endTime = strings.TrimSpace(endTime)
			if endTime == "" {
				endTime = task.EndTime.Format(time.TimeOnly)
			}
			_, err := time.Parse("15:04:05", endTime)
			if err == nil {
				break
			}
		}
	}
	if endTime == "" {
		endTime = task.EndTime.String()
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
	startTime, err := stringToTimeToday(startTimeStr)
	if err != nil {
		return err
	}
	endTime, err := stringToTimeToday(endTimeStr)
	if err != nil {
		return err
	}
	if startTime.After(endTime) {
		return errors.New("start time cannot be after end time")
	}
	newTassk := Task{
		ID:        id,
		Name:      name,
		Status:    status,
		Priority:  priority,
		StartTime: startTime,
		EndTime:   endTime,
	}
	taskList = append(taskList, newTassk)
	err = writeFile(taskList, fileName)
	if err != nil {
		return err
	}
	return nil
}

func editTask(dateString string, id string) error {
	fileName, err := formatDate(dateString)
	if err != nil {
		return err
	}
	taskList, err := readFile(fileName)
	if err != nil {
		return err
	}
	foundedTask := searchTask(taskList, id)
	if foundedTask == nil {
		return errors.New("task not found")
	}

	name, status, priority, startTimeStr, endTimeStr := inputFromKeyboard(foundedTask)
	startTime, err := stringToTimeWithDate(startTimeStr, foundedTask.StartTime)
	if err != nil {
		return err
	}
	endTime, err := stringToTimeWithDate(endTimeStr, foundedTask.EndTime)
	if err != nil {
		return err
	}
	if startTime.After(endTime) {
		return errors.New("start time cannot be after end time")
	}

	foundedTask.Name = name
	foundedTask.Status = status
	foundedTask.Priority = priority
	foundedTask.StartTime = startTime
	foundedTask.EndTime = endTime

	err = writeFile(taskList, fileName)
	if err != nil {
		return err
	}
	return nil
}

func updateStatusTask(status string, id string) error {
	fileName := time.Now().Format(time.DateOnly)
	taskList, err := readFile(fileName)
	if err != nil {
		return err
	}
	foundedTask := searchTask(taskList, id)
	foundedTask.Status = status
	err = writeFile(taskList, fileName)
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
		fmt.Println("4. Edit task")
		fmt.Println("5. Update status for today task")
		fmt.Println("0. Exit")
		fmt.Print(">. Input your option: ")
		option, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to read input!\nError: %s\n", err)
			continue
		}
		option = strings.TrimSpace(option)

		switch option {
		case "1":
			displayTaskListByToday()
		case "2":
			var dateString string
			for {
				isValid := true
				fmt.Print(">. Input date: ")
				dateString, err = reader.ReadString('\n')
				if err != nil {
					isValid = false
					fmt.Fprintf(os.Stderr, "Failed to read input!\nError: %s\n", err)
				}
				dateString = strings.TrimSpace(dateString)
				if dateString == "" {
					isValid = false
					fmt.Println("Date was empty!")
				}
				if isValid {
					break
				}
			}
			displayTaskListByDate(dateString)
		case "3":
			err := addTask()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Failed to add task!\nError: %s\n", err)
			}
		case "4":
			var dateString string
			for {
				isValid := true
				fmt.Print(">. Input date: ")
				dateString, err = reader.ReadString('\n')
				if err != nil {
					isValid = false
					fmt.Fprintf(os.Stderr, "Failed to read input!\nError: %s\n", err)
				}
				dateString = strings.TrimSpace(dateString)
				_, err := formatDate(dateString)
				if err != nil {
					isValid = false
					fmt.Fprintf(os.Stderr, "Invalid format date!\nError: %s\n", err)
				}
				if isValid {
					break
				}
			}
			displayTaskListByDate(dateString)
			fileName, _ := formatDate(dateString)
			taskList, err := readFile(fileName)
			if err != nil || len(taskList) == 0 {
				continue
			}
			var id string
			for {
				isValid := true
				fmt.Print(">. Input task id to edit: ")
				id, err = reader.ReadString('\n')
				if err != nil {
					isValid = false
					fmt.Fprintf(os.Stderr, "Invalid ID\nError: %s\n", err)
				}
				id = strings.TrimSpace(id)
				if searchTask(taskList, id) == nil {
					isValid = false
					fmt.Printf("Task ID '%s' not found! Please try again.\n", id)
				}
				if isValid {
					break
				}
			}
			err = editTask(dateString, id)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Failed to edit task!\nError: %s\n", err)
			}
		case "5":
			displayTaskListByToday()
			fileName := time.Now().Format(time.DateOnly)
			taskList, err := readFile(fileName)
			if err != nil || len(taskList) == 0 {
				continue
			}
			var id string
			for {
				isValid := true
				fmt.Print(">. Input task id to edit: ")
				id, err = reader.ReadString('\n')
				if err != nil {
					isValid = false
					fmt.Fprintf(os.Stderr, "Invalid ID\nError: %s\n", err)
				}
				id = strings.TrimSpace(id)
				if searchTask(taskList, id) == nil {
					isValid = false
					fmt.Printf("Task ID '%s' not found! Please try again.\n", id)
				}
				if isValid {
					break
				}
			}
			statusList := []string{"To-do", "Inprogress", "Done"}
			var status string
			for {
				isValid := false
				fmt.Print(">. Enter status: ")
				status, err = reader.ReadString('\n')
				if err != nil {
					isValid = false
					fmt.Fprintf(os.Stderr, "Invalid Status!\nError: %s\n", err)
				}
				status = strings.TrimSpace(status)

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
				fmt.Println("Invalid Status!\nPlease enter To-do, Inprogress, or Done.")
			}
			err = updateStatusTask(status, id)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Failed to update status!\nError: %s\n", err)
			}
		case "0":
			fmt.Println("See you soon!")
			return
		default:
			fmt.Println("Invalid option!")
		}
	}
}

func main() {
	controller()
}
