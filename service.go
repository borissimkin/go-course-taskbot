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
		listItems = append(listItems, t.getTaskListItem(userID, index, task))
	}

	result := strings.Join(listItems, "\n")

	t.delivery.Send(userID, result)
}

func (t *TaskTracker) getTaskListItem(userID ID, index int, task Task) string {
	owner := t.repo.GetUser(task.OwnerID)

	taskTmpl := fmt.Sprintf("%d. %s by %s", index+1, task.Text, t.getUserName(owner.UserName))
	if task.AssignedID == 0 {
		taskTmpl += fmt.Sprintf("\n/assign_%d", task.ID)
		return taskTmpl
	}

	assigneeName := "я"
	if task.AssignedID != userID {
		assigned := t.repo.GetUser(task.AssignedID)
		if assigned != nil {
			assigneeName = t.getUserName(assigned.UserName)
		}
	}

	taskTmpl += fmt.Sprintf("\nassignee: %s", assigneeName)

	if task.AssignedID == userID {
		taskTmpl += fmt.Sprintf("\n/unassign_%d /resolve_%d", task.ID, task.ID)
	}

	return taskTmpl
}

func (t *TaskTracker) AssignTask(taskID, userID ID) {
	var prevAssignedID ID = 0
	task := t.repo.GetTask(taskID)
	if task != nil {
		prevAssignedID = task.AssignedID
	}

	err := t.repo.UpdateTaskAssignedID(taskID, userID)
	if err != nil {
		t.sendError(userID, err)
		return
	}

	baseMessage := fmt.Sprintf("Задача \"%s\" назначена на ", task.Text)

	t.delivery.Send(userID, baseMessage+"вас")

	var prevUser *User
	if prevAssignedID != 0 && prevAssignedID != userID {
		prevUser = t.repo.GetUser(prevAssignedID)
	} else if userID != task.OwnerID {
		prevUser = t.repo.GetUser(task.OwnerID)
	}

	if prevUser != nil {
		user := t.repo.GetUser(userID)
		fmt.Println()
		if user != nil {
			t.delivery.Send(prevUser.ID, baseMessage+t.getUserName(user.UserName))
		}
	}
}

func (t *TaskTracker) sendError(chatID ID, err error) {
	t.delivery.Send(chatID, fmt.Sprintf("ошибка: %s", err))
}

func (t *TaskTracker) getUserName(userName string) string {
	return fmt.Sprintf("@%s", userName)
}
