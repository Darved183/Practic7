package main

import (
	"fmt"
	"sync"
)

type Task struct {
	ID        int
	Title     string
	Completed bool
}

type TaskManager struct {
	tasks  []Task
	nextID int
	mux    sync.RWMutex
}

func (tm *TaskManager) AddTask(title string) {
	tm.mux.Lock()
	defer tm.mux.Unlock()

	task := Task{
		ID:    tm.nextID,
		Title: title,
	}
	tm.tasks = append(tm.tasks, task)
	tm.nextID++
}

func (tm *TaskManager) DeleteTask(id int) {
	tm.mux.Lock()
	defer tm.mux.Unlock()

	for i, task := range tm.tasks {
		if task.ID == id {
			tm.tasks = append(tm.tasks[:i], tm.tasks[i+1:]...)
			return
		}
	}
}

func (tm *TaskManager) MarkCompleted(id int) {
	tm.mux.Lock()
	defer tm.mux.Unlock()

	for i, task := range tm.tasks {
		if task.ID == id {
			tm.tasks[i].Completed = true
			return
		}
	}
}

func (tm *TaskManager) GetActiveTasks() []Task {
	tm.mux.RLock()
	defer tm.mux.RUnlock()

	var result []Task
	for _, task := range tm.tasks {
		if !task.Completed {
			result = append(result, task)
		}
	}
	return result
}

func (tm *TaskManager) GetCompletedTasks() []Task {
	tm.mux.RLock()
	defer tm.mux.RUnlock()

	var result []Task
	for _, task := range tm.tasks {
		if task.Completed {
			result = append(result, task)
		}
	}
	return result
}

func (tm *TaskManager) GetAllTasks() []Task {
	tm.mux.RLock()
	defer tm.mux.RUnlock()

	return tm.tasks
}

func main() {
	tm := TaskManager{}

	//ИИ

	tm.AddTask("Купить продукты")
	tm.AddTask("Сделать уроки")
	tm.AddTask("Позвонить маме")

	fmt.Println("Все задачи:")
	for _, task := range tm.GetAllTasks() {
		status := "активна"
		if task.Completed {
			status = "выполнена"
		}
		fmt.Printf("%d. %s (%s)\n", task.ID, task.Title, status)
	}

	tm.MarkCompleted(1)
	fmt.Println("\nПосле выполнения задачи 1:")

	fmt.Println("Активные задачи:")
	for _, task := range tm.GetActiveTasks() {
		fmt.Printf("- %s\n", task.Title)
	}

	fmt.Println("Выполненные задачи:")
	for _, task := range tm.GetCompletedTasks() {
		fmt.Printf("- %s\n", task.Title)
	}

	tm.DeleteTask(2)
	fmt.Println("\nПосле удаления задачи 2:")
	for _, task := range tm.GetAllTasks() {
		fmt.Printf("- %s\n", task.Title)
	}
}
