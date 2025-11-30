package main

import (
	"fmt"
	"strconv"
	"strings"
)

type CmdRouter struct {
	service *TaskTracker
}

func (r *CmdRouter) Parse(msg RequestMessage) {
	input := msg.Text

	chatID := ID(msg.ChatID)

	r.service.UpdateOrCreateUser(msg.User)

	switch {
	case input == "/tasks":
		r.service.ListTasks(chatID)
	case strings.HasPrefix(input, "/new "):
		r.service.CreateTask(chatID, strings.TrimPrefix(input, "/new "))
	case strings.HasPrefix(input, "/assign_"):
		id := strings.TrimPrefix(input, "/assign_")
		intID, err := strconv.Atoi(id)
		if err != nil {
			return
		}
		r.service.AssignTask(ID(intID), chatID)

	case strings.HasPrefix(input, "/unassign_"):
		id := strings.TrimPrefix(input, "/unassign_")
		intID, err := strconv.Atoi(id)
		if err != nil {
			return
		}
		r.service.UnassignTask(ID(intID), chatID)

	case strings.HasPrefix(input, "/resolve_"):
		id := strings.TrimPrefix(input, "/resolve_")
		intID, err := strconv.Atoi(id)
		if err != nil {
			return
		}
		r.service.ResolveTask(ID(intID), chatID)
	case input == "/my":
		r.service.ListTasksByAssignee(chatID)

	case input == "/owner":
		r.service.ListTasksByOwner(chatID)

	default:
		fmt.Println("Неизвестная команда")
	}
}
