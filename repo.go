package main

import "fmt"

type ID int

type User struct {
	ID       ID
	UserName string
}

type Task struct {
	ID         ID
	Text       string
	OwnerID    ID
	AssignedID ID
	Completed  bool
}

type Repo struct {
	tasks []Task
	users []User
}

func NewRepo() *Repo {
	return &Repo{
		tasks: make([]Task, 0),
	}
}

func (s *Repo) CreateTask(ownerID ID, text string) (Task, error) {
	nextId := s.getLastID() + 1

	task := Task{
		ID:      nextId,
		Text:    text,
		OwnerID: ownerID,
	}

	s.tasks = append(s.tasks, task)

	return task, nil
}

func (s *Repo) GetList() ([]Task, error) {
	return s.tasks, nil
}

func (s *Repo) GetTask(id ID) *Task {
	for _, task := range s.tasks {
		if task.ID == id {
			return &task
		}
	}

	return nil
}

func (s *Repo) UpdateTaskAssignedID(id, assignedID ID) error {
	for i, task := range s.tasks {
		if task.ID == id {
			s.tasks[i].AssignedID = assignedID
			return nil
		}
	}

	return fmt.Errorf("задача с id=%d не найдена", id)
}

func (s *Repo) GetUser(id ID) *User {
	for _, u := range s.users {
		if u.ID == id {
			return &u
		}
	}

	return nil
}

func (s *Repo) UpdateOrCreateUser(user User) {
	for i, u := range s.users {
		if u.ID == user.ID {
			s.users[i] = user
			return
		}
	}

	s.users = append(s.users, user)
}

func (s *Repo) getLastID() ID {
	if len(s.tasks) == 0 {
		return 0
	}

	lastTask := s.tasks[len(s.tasks)-1]

	return lastTask.ID
}
