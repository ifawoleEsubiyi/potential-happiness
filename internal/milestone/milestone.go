// Package milestone provides milestone tracking functionality for the validator daemon.
// It allows creating, tracking, and completing milestones to monitor validator progress.
package milestone

import (
	"errors"
	"sync"
	"time"
)

// Status represents the current status of a milestone.
type Status string

const (
	// StatusPending means the milestone has not been started.
	StatusPending Status = "pending"
	// StatusInProgress means the milestone is currently being worked on.
	StatusInProgress Status = "in_progress"
	// StatusCompleted means the milestone has been completed.
	StatusCompleted Status = "completed"
)

// Milestone represents a single milestone with its metadata.
type Milestone struct {
	ID          string
	Name        string
	Description string
	Status      Status
	CreatedAt   time.Time
	CompletedAt *time.Time
}

// Maker manages milestone creation and tracking.
// This type is safe for concurrent use by multiple goroutines.
type Maker struct {
	mu         sync.RWMutex
	milestones map[string]*Milestone
}

// New creates a new milestone maker instance.
func New() *Maker {
	return &Maker{
		milestones: make(map[string]*Milestone),
	}
}

// Create creates a new milestone with the given ID, name, and description.
// Returns a copy of the created milestone to prevent external modification.
func (m *Maker) Create(id, name, description string) (Milestone, error) {
	if id == "" {
		return Milestone{}, errors.New("milestone ID cannot be empty")
	}
	if name == "" {
		return Milestone{}, errors.New("milestone name cannot be empty")
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.milestones[id]; exists {
		return Milestone{}, errors.New("milestone with this ID already exists")
	}

	milestone := &Milestone{
		ID:          id,
		Name:        name,
		Description: description,
		Status:      StatusPending,
		CreatedAt:   time.Now(),
	}

	m.milestones[id] = milestone
	return *milestone, nil
}

// Get retrieves a copy of a milestone by its ID.
// Returns a copy to prevent external modification.
func (m *Maker) Get(id string) (Milestone, error) {
	if id == "" {
		return Milestone{}, errors.New("milestone ID cannot be empty")
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	milestone, exists := m.milestones[id]
	if !exists {
		return Milestone{}, errors.New("milestone not found")
	}

	return *milestone, nil
}

// Start marks a milestone as in progress.
func (m *Maker) Start(id string) error {
	if id == "" {
		return errors.New("milestone ID cannot be empty")
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	milestone, exists := m.milestones[id]
	if !exists {
		return errors.New("milestone not found")
	}

	if milestone.Status == StatusCompleted {
		return errors.New("cannot start a completed milestone")
	}

	milestone.Status = StatusInProgress
	return nil
}

// Complete marks a milestone as completed.
func (m *Maker) Complete(id string) error {
	if id == "" {
		return errors.New("milestone ID cannot be empty")
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	milestone, exists := m.milestones[id]
	if !exists {
		return errors.New("milestone not found")
	}

	if milestone.Status == StatusCompleted {
		return errors.New("milestone is already completed")
	}

	now := time.Now()
	milestone.Status = StatusCompleted
	milestone.CompletedAt = &now
	return nil
}

// List returns copies of all milestones.
// Returns copies to prevent external modification.
func (m *Maker) List() []Milestone {
	m.mu.RLock()
	defer m.mu.RUnlock()

	milestones := make([]Milestone, 0, len(m.milestones))
	for _, milestone := range m.milestones {
		milestones = append(milestones, *milestone)
	}
	return milestones
}

// Count returns the total number of milestones.
func (m *Maker) Count() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.milestones)
}

// CountByStatus returns the number of milestones with the given status.
func (m *Maker) CountByStatus(status Status) int {
	m.mu.RLock()
	defer m.mu.RUnlock()

	count := 0
	for _, milestone := range m.milestones {
		if milestone.Status == status {
			count++
		}
	}
	return count
}
