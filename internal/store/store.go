package store

import (
	"strings"
	"sync"
	"time"

	"github.com/eldhosereji541/task-tracker-backend/internal/model"
	"github.com/google/uuid"
)

type Store struct {
	mu    sync.RWMutex
	tasks map[string]*model.Task
	users map[string]*model.User
	// email => hashed password
	userPasswords map[string]string
}

func NewStore() *Store {
	return &Store{
		tasks:         make(map[string]*model.Task),
		users:         make(map[string]*model.User),
		userPasswords: make(map[string]string),
	}
}

// user related methods

func (s *Store) CreateUser(name, email, hashedPassword string) (*model.User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()
	u := &model.User{
		ID:        uuid.NewString(),
		Name:      name,
		Email:     email,
		CreatedAt: now,
		UpdatedAt: now,
		Deleted:   false,
	}
	s.users[u.ID] = u
	s.userPasswords[u.Email] = hashedPassword
	return u, nil
}

func (s *Store) GetUserByID(id string) (*model.User, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	u, exists := s.users[id]
	if !exists || u.Deleted {
		return nil, false
	}
	return u, true
}

func (s *Store) GetUserByEmail(email string) (*model.User, string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, u := range s.users {
		if u.Email == email && !u.Deleted {
			hashedPassword := s.userPasswords[email]
			return u, hashedPassword, true
		}
	}
	return nil, "", false
}

// task related methods

func (s *Store) CreateTask(task *model.Task) (*model.Task, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	task.ID = uuid.NewString()
	task.CreatedAt = time.Now()
	task.UpdatedAt = time.Now()
	s.tasks[task.ID] = task
	return task, nil
}

func (s *Store) GetTaskByID(id string) (*model.Task, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	t, exists := s.tasks[id]
	if !exists || t.Deleted {
		return nil, false
	}
	return t, true
}

func (s *Store) GetTasks(userID string, filter *model.TaskFilter) ([]*model.Task, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var tasks []*model.Task
	for _, t := range s.tasks {
		if t.Deleted || t.User.ID != userID {
			continue
		}
		if filter != nil {
			if filter.Status != nil && t.Status != *filter.Status {
				continue
			}
			if filter.Priority != nil && t.Priority != *filter.Priority {
				continue
			}
			if filter.Search != nil && !strings.Contains(t.Title, *filter.Search) &&
				(t.Description == nil || !strings.Contains(*t.Description, *filter.Search)) {
				continue
			}
		}
		tasks = append(tasks, t)
	}
	return tasks, nil
}

func (s *Store) UpdateTask(task *model.Task) (*model.Task, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	existingTask, ok := s.tasks[task.ID]
	if !ok || existingTask.Deleted {
		return nil, nil
	}
	task.UpdatedAt = time.Now()
	s.tasks[task.ID] = task
	return task, nil
}

func (s *Store) DeleteTask(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	task, ok := s.tasks[id]
	if !ok || task.Deleted {
		return nil
	}
	task.Deleted = true
	task.UpdatedAt = time.Now()
	s.tasks[id] = task
	return nil
}
