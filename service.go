package main

import (
	"fmt"
)

type Delivery interface {
	Send(chatID ID, text string)
}

type TaskTracker struct {
	delivery Delivery
	repo     Repo
}

func (t *TaskTracker) CreateTask(userID ID, text string) {
	task, err := t.repo.CreateTask(userID, text)
	if err != nil {
		t.sendError(userID, err)
		return
	}
}

func (t *TaskTracker) sendError(chatID ID, err error) {
	t.delivery.Send(chatID, fmt.Sprintf("ошибка: %s", err.Error()))
}
