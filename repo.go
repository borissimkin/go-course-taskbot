package main

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

func (s *Repo) GetUser(id ID) *User {
	for _, u := range s.users {
		if u.ID == id {
			return &u
		}
	}

	return nil
}

func (s *Repo) UpdateOrCreateUser(user User) {
	for _, u := range s.users {
		if u.ID == user.ID {
			u = user
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