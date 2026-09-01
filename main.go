package main

import (
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"
)

var priorityScore map[string]int = map[string]int{
	"High":   3,
	"Medium": 2,
	"Low":    1,
}

func readFile() ([]Task, error) {
	file, err := os.Open("./data/2026-08-16.json")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var taskList []Task
	err = json.NewDecoder(file).Decode(&taskList)
	if err != nil {
		return nil, err
	}

	return taskList, nil
}

func displayTaskList(taskList []Task) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)

	fmt.Fprintf(w, "|\t%s\t|\t%s\t|\t%s\t|\t%s\t|\t%s\t|\t%s\t|\n", "ID", "Name", "Status", "Priority", "Start Time", "End Time")
	for _, task := range taskList {
		fmt.Fprintf(w, "|\t%s\t|\t%s\t|\t%s\t|\t%s\t|\t%s\t|\t%s\t|\n", task.ID, task.Name, task.Status, task.Priority, task.StartTime, task.EndTime)
	}

	w.Flush()
}

func sortPriority(taskList []Task) []Task {
	for i := 0; i < len(taskList)-1; i++ {
		max := i
		for j := i + 1; j < len(taskList); j++ {
			if priorityScore[taskList[max].Priority] <= priorityScore[taskList[j].Priority] {
				max = j
			}
		}
		if max != i {
			taskList[i], taskList[max] = taskList[max], taskList[i]
		}
	}
	return taskList
}

func sortStartDate(taskList []Task) []Task {
	for i := 0; i < len(taskList)-1; i++ {
		min := i
		for j := i + 1; j < len(taskList); j++ {
			if taskList[min].StartTime.After(taskList[j].StartTime) {
				min = j
			}
		}
		if min != i {
			taskList[i], taskList[min] = taskList[min], taskList[i]
		}
	}
	return taskList
}

func sortPriorityAndStartDate(taskList []Task) []Task {
	for i := 0; i < len(taskList) - 1; i++ {
		best := i;
		for j := i + 1; j < len(taskList); j++ {
			scoreBest := priorityScore[taskList[best].Priority];
			scoreJ := priorityScore[taskList[j].Priority];
			if scoreBest < scoreJ {
				best = j;
			} else if scoreBest == scoreJ {
				if taskList[best].StartTime.After(taskList[j].StartTime) {
					best = j;
				}
			}
		}
		if best != i {
			taskList[i], taskList[best] = taskList[best], taskList[i];
		}
	}
	return taskList
}

func displayTaskListByToday(taskList []Task) {
	taskList = sortPriorityAndStartDate(taskList);
	displayTaskList(taskList);
}

func displayTaskListByDate(taskList []Task) {
	taskList = sortPriorityAndStartDate(taskList);
	displayTaskList(taskList);
}

func main() {
	taskList, err := readFile();
	if err != nil {
		fmt.Println(err);
	}

	displayTaskListByToday(taskList);
}
