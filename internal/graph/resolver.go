package graph

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

import (
	"github.com/eldhosereji541/task-tracker-backend/internal/auth"
	"github.com/eldhosereji541/task-tracker-backend/internal/model"
)

// TaskRepo defines the data access contract required by graph resolvers.
// Defined in the consumer package per Go interface conventions.
type TaskRepo interface {
	CreateTask(task *model.Task) (*model.Task, error)
	GetTaskByID(id string) (*model.Task, bool)
	GetTasks(userID string, filter *model.TaskFilter) ([]*model.Task, error)
	UpdateTask(task *model.Task) (*model.Task, error)
	DeleteTask(id string) error
	CreateUser(name, email, hashedPassword string) (*model.User, error)
	GetUserByID(id string) (*model.User, bool)
	GetUserByEmail(email string) (*model.User, string, bool)
}

type Resolver struct {
	TaskRepo TaskRepo
	TokenSvc *auth.TokenService
}
