package main

import (
	"fmt"
	"strings"
)

type Delivery interface {
	Send(chatID ID, text string)
}

type TaskTracker struct {
	delivery Delivery
	repo     *Repo
}

func (t *TaskTracker) UpdateOrCreateUser(user User) {
	t.repo.UpdateOrCreateUser(user)
}

func (t *TaskTracker) CreateTask(userID ID, text string) {
	task, err := t.repo.CreateTask(userID, text)
	if err != nil {
		t.sendError(userID, err)
		return
	}

	result := fmt.Sprintf("Задача \"%s\" создана, id=%d", task.Text, task.ID)
	t.delivery.Send(userID, result)
}

func (t *TaskTracker) ListTasks(userID ID) {
	tasks, err := t.repo.GetList()
	if err != nil {
		t.sendError(userID, err)
		return
	}

	if len(tasks) == 0 {
		t.delivery.Send(userID, "Нет задач")
		return
	}

	listItems := make([]string, 0, len(tasks))
	for index, task := range tasks {
		owner := t.repo.GetUser(task.OwnerID)
		listItems = append(listItems, fmt.Sprintf("%d. %s by @%s\n/assign_%d", index+1, task.Text, owner.UserName, task.ID))
	}

	result := strings.Join(listItems, "\n")

	t.delivery.Send(userID, result)
}

func (t *TaskTracker) AssignTask() {

}

func (t *TaskTracker) sendError(chatID ID, err error) {
	t.delivery.Send(chatID, fmt.Sprintf("ошибка: %s", err))
}
