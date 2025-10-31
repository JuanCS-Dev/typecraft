package domain

import (
	"context"
	"errors"
)

// Common errors
var (
	ErrProjectNotFound = errors.New("project not found")
	ErrInvalidInput    = errors.New("invalid input")
)

// ProjectRepository defines the interface for project persistence
type ProjectRepository interface {
	GetByID(ctx context.Context, id uint) (*Project, error)
	Create(ctx context.Context, project *Project) error
	Update(ctx context.Context, project *Project) error
	Delete(ctx context.Context, id uint) error
	List(ctx context.Context) ([]*Project, error)
}
