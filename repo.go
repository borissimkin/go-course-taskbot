package main

type ID int

type Task struct {
	ID         ID
	Text       string
	OwnerID    ID
	AssignedID ID
	Completed  bool
}

type Repo struct {
	tasks []Task
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

func (s *Repo) getLastID() ID {
	lastTask := s.tasks[len(s.tasks)-1]

	return lastTask.ID
}
