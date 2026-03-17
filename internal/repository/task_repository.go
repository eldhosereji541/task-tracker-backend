package repository

import (
	"context"
	"strings"
	"sync"

	"github.com/eldhosereji541/task-tracker-backend/internal/model"
)

// Defines the interface for task data operations.
type TaskRepository interface {
	// retrieves all tasks, optionally filtered by status and search string.
	GetTasks(ctx context.Context, filter *model.TaskFilter) ([]*model.Task, error)
	CreateTask(ctx context.Context, task *model.Task) (*model.Task, error)
}

type TaskRepositoryImpl struct {
	mu sync.Mutex
	// in-memory store.
	tasks map[string]*model.Task
}

// create a new instance of TaskRepositoryImpl with empty task list.
func NewTaskRepository() TaskRepository {
	return &TaskRepositoryImpl{
		tasks: make(map[string]*model.Task, 0),
	}
}

func (r *TaskRepositoryImpl) CreateTask(ctx context.Context, task *model.Task) (*model.Task, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.tasks[task.ID] = task
	return task, nil
}

func (r *TaskRepositoryImpl) GetTasks(ctx context.Context, filter *model.TaskFilter) ([]*model.Task, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	var filteredTasks []*model.Task
	for _, task := range r.tasks {
		if task.Deleted {
			continue
		}
		if filter.Status != nil && task.Status != *filter.Status {
			continue
		}
		if filter.Search != nil && !strings.Contains(strings.ToLower(task.Title), strings.ToLower(*filter.Search)) &&
			(task.Description == nil || !strings.Contains(strings.ToLower(*task.Description), strings.ToLower(*filter.Search))) {
			continue
		}
		filteredTasks = append(filteredTasks, task)
	}
	return filteredTasks, nil
}
