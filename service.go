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

	t.delivery.Send(userID, t.getTaskListMsg(tasks, userID, true))
}

func (t *TaskTracker) ListTasksByAssignee(userID ID) {
	tasks, err := t.repo.GetListByAssignee(userID)
	if err != nil {
		t.sendError(userID, err)
		return
	}

	t.delivery.Send(userID, t.getTaskListMsg(tasks, userID, false))
}

func (t *TaskTracker) ListTasksByOwner(userID ID) {
	tasks, err := t.repo.GetListByOwner(userID)
	if err != nil {
		t.sendError(userID, err)
		return
	}

	t.delivery.Send(userID, t.getTaskListMsg(tasks, userID, false))
}

func (t *TaskTracker) getTaskListMsg(tasks []Task, userID ID, withAssignee bool) string {
	if len(tasks) == 0 {
		return "Нет задач"
	}

	listItems := make([]string, 0, len(tasks))
	for _, task := range tasks {
		listItems = append(listItems, t.getTaskListItem(userID, task, withAssignee))
	}

	result := strings.Join(listItems, "\n\n")

	return result
}

func (t *TaskTracker) getTaskListItem(userID ID, task Task, withAssignee bool) string {
	owner := t.repo.GetUser(task.OwnerID)

	taskTmpl := fmt.Sprintf("%d. %s by %s", task.ID, task.Text, t.getUserName(owner.UserName))
	if task.AssignedID == 0 {
		taskTmpl += fmt.Sprintf("\n/assign_%d", task.ID)
		return taskTmpl
	}

	if withAssignee {
		taskTmpl = t.addAssignee(taskTmpl, task, userID)
	}

	if task.AssignedID == userID {
		taskTmpl += fmt.Sprintf("\n/unassign_%d /resolve_%d", task.ID, task.ID)
	}

	return taskTmpl
}

func (t *TaskTracker) addAssignee(taskListItem string, task Task, userID ID) string {
	assigneeName := "я"
	if task.AssignedID != userID {
		assigned := t.repo.GetUser(task.AssignedID)
		if assigned != nil {
			assigneeName = t.getUserName(assigned.UserName)
		}
	}

	return taskListItem + fmt.Sprintf("\nassignee: %s", assigneeName)
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

func (t *TaskTracker) UnassignTask(taskID, userID ID) {
	task := t.repo.GetTask(taskID)
	if task == nil {
		t.delivery.Send(userID, "Задачи нет")
		return
	}

	if task.AssignedID != userID {
		t.delivery.Send(userID, "Задача не на вас")
		return
	}

	err := t.repo.UpdateTaskAssignedID(taskID, 0)
	if err != nil {
		t.sendError(userID, err)
		return
	}

	msgToUnassigned := "Принято"
	msgToOwner := fmt.Sprintf("Задача \"%s\" осталась без исполнителя", task.Text)

	t.delivery.Send(userID, msgToUnassigned)
	t.delivery.Send(task.OwnerID, msgToOwner)
}

func (t *TaskTracker) ResolveTask(taskID, userID ID) {
	task := t.repo.GetTask(taskID)
	if task == nil {
		t.delivery.Send(userID, "Задачи нет")
		return
	}

	err := t.repo.UpdateTaskAssignedID(taskID, 0)
	if err != nil {
		t.sendError(userID, err)
		return
	}

	err = t.repo.CompteteTask(task.ID)
	if err != nil {
		t.sendError(userID, err)
		return
	}

	baseMsg := fmt.Sprintf("Задача \"%s\" выполнена", task.Text)

	t.delivery.Send(userID, baseMsg)

	assignedUser := t.repo.GetUser(userID)
	if assignedUser != nil {
		msgToOwner := fmt.Sprintf("%s %s", baseMsg, t.getUserName(assignedUser.UserName))
		t.delivery.Send(task.OwnerID, msgToOwner)
	}
}

func (t *TaskTracker) sendError(chatID ID, err error) {
	t.delivery.Send(chatID, fmt.Sprintf("ошибка: %s", err))
}

func (t *TaskTracker) getUserName(userName string) string {
	return fmt.Sprintf("@%s", userName)
}
