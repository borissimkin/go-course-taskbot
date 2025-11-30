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

func (r *Repo) CreateTask(ownerID ID, text string) (Task, error) {
	nextId := r.getLastID() + 1

	task := Task{
		ID:      nextId,
		Text:    text,
		OwnerID: ownerID,
	}

	r.tasks = append(r.tasks, task)

	return task, nil
}

func (r *Repo) GetList() ([]Task, error) {
	tasks, err := r.GetListWithFilter(r.tasks, func(task Task) bool {
		return !task.Completed
	})

	return tasks, err
}

func (r *Repo) GetListByAssignee(assignedID ID) ([]Task, error) {
	tasks, err := r.GetList()
	if err != nil {
		return tasks, err
	}

	return r.GetListWithFilter(tasks, func(task Task) bool {
		return task.AssignedID == assignedID
	})
}

func (r *Repo) GetListByOwner(ownerID ID) ([]Task, error) {
	tasks, err := r.GetList()
	if err != nil {
		return tasks, err
	}

	return r.GetListWithFilter(tasks, func(task Task) bool {
		return task.OwnerID == ownerID
	})
}

func (r *Repo) GetListWithFilter(tasks []Task, filter func(Task) bool) ([]Task, error) {
	filtered := make([]Task, 0, len(r.tasks))
	for _, task := range tasks {
		if filter(task) {
			filtered = append(filtered, task)
		}
	}

	return filtered, nil
}

func (r *Repo) GetTask(id ID) *Task {
	for _, task := range r.tasks {
		if task.ID == id {
			return &task
		}
	}

	return nil
}

func (r *Repo) UpdateTaskAssignedID(id, assignedID ID) error {
	for i, task := range r.tasks {
		if task.ID == id {
			r.tasks[i].AssignedID = assignedID
			return nil
		}
	}

	return fmt.Errorf("задача с id=%d не найдена", id)
}

func (r *Repo) GetUser(id ID) *User {
	for _, u := range r.users {
		if u.ID == id {
			return &u
		}
	}

	return nil
}

func (r *Repo) UpdateOrCreateUser(user User) {
	for i, u := range r.users {
		if u.ID == user.ID {
			r.users[i] = user
			return
		}
	}

	r.users = append(r.users, user)
}

func (r *Repo) CompteteTask(id ID) error {
	for i, task := range r.tasks {
		if task.ID == id {
			r.tasks[i].Completed = true
			return nil
		}
	}

	return fmt.Errorf("Задача не найдена")
}

func (r *Repo) getLastID() ID {
	if len(r.tasks) == 0 {
		return 0
	}

	lastTask := r.tasks[len(r.tasks)-1]

	return lastTask.ID
}
