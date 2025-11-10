package main

import (
	"fmt"
	"strings"
)

type CmdRouter struct {
	service *TaskTracker
}

func (r *CmdRouter) Parse(msg RequestMessage) {
	input := msg.Text

	switch {
	case input == "/tasks":
		fmt.Println("Показать список всех задач")

	case strings.HasPrefix(input, "/new "):
		// Пример: /new XXX YYY ZZZ
		args := strings.TrimPrefix(input, "/new ")
		fmt.Printf("Создать новую задачу с параметрами: %s\n", args)

	case strings.HasPrefix(input, "/assign_"):
		id := strings.TrimPrefix(input, "/assign_")
		fmt.Printf("Назначить задачу с ID %s\n", id)

	case strings.HasPrefix(input, "/unassign_"):
		id := strings.TrimPrefix(input, "/unassign_")
		fmt.Printf("Снять назначение с задачи %s\n", id)

	case strings.HasPrefix(input, "/resolve_"):
		id := strings.TrimPrefix(input, "/resolve_")
		fmt.Printf("Закрыть задачу с ID %s\n", id)

	case input == "/my":
		fmt.Println("Показать мои задачи")

	case input == "/owner":
		fmt.Println("Показать задачи, где я владелец")

	default:
		fmt.Println("Неизвестная команда")
	}
}
