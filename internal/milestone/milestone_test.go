package milestone

import (
	"testing"
)

func TestNew(t *testing.T) {
	m := New()
	if m == nil {
		t.Error("New() returned nil Maker instance")
	}
	if m.milestones == nil {
		t.Error("New() returned Maker with nil milestones map")
	}
}

func TestCreate(t *testing.T) {
	tests := []struct {
		name        string
		id          string
		mileName    string
		description string
		wantErr     bool
		errMsg      string
	}{
		{
			name:        "valid milestone",
			id:          "m1",
			mileName:    "First Milestone",
			description: "Description of the first milestone",
			wantErr:     false,
		},
		{
			name:        "empty ID",
			id:          "",
			mileName:    "First Milestone",
			description: "Description",
			wantErr:     true,
			errMsg:      "milestone ID cannot be empty",
		},
		{
			name:        "empty name",
			id:          "m1",
			mileName:    "",
			description: "Description",
			wantErr:     true,
			errMsg:      "milestone name cannot be empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := New()
			milestone, err := m.Create(tt.id, tt.mileName, tt.description)

			if (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && err != nil && err.Error() != tt.errMsg {
				t.Errorf("Create() error = %v, want %v", err.Error(), tt.errMsg)
				return
			}

			if !tt.wantErr {
				if milestone.ID != tt.id {
					t.Errorf("Create() ID = %v, want %v", milestone.ID, tt.id)
				}
				if milestone.Name != tt.mileName {
					t.Errorf("Create() Name = %v, want %v", milestone.Name, tt.mileName)
				}
				if milestone.Status != StatusPending {
					t.Errorf("Create() Status = %v, want %v", milestone.Status, StatusPending)
				}
			}
		})
	}
}

func TestCreateDuplicate(t *testing.T) {
	m := New()
	_, err := m.Create("m1", "First Milestone", "Description")
	if err != nil {
		t.Fatalf("First Create() error = %v", err)
	}

	_, err = m.Create("m1", "Second Milestone", "Another Description")
	if err == nil {
		t.Error("Expected error for duplicate ID, got nil")
	}
	if err.Error() != "milestone with this ID already exists" {
		t.Errorf("Create() error = %v, want 'milestone with this ID already exists'", err.Error())
	}
}

func TestGet(t *testing.T) {
	m := New()
	_, err := m.Create("m1", "First Milestone", "Description")
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	tests := []struct {
		name    string
		id      string
		wantErr bool
		errMsg  string
	}{
		{
			name:    "existing milestone",
			id:      "m1",
			wantErr: false,
		},
		{
			name:    "non-existing milestone",
			id:      "m2",
			wantErr: true,
			errMsg:  "milestone not found",
		},
		{
			name:    "empty ID",
			id:      "",
			wantErr: true,
			errMsg:  "milestone ID cannot be empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			milestone, err := m.Get(tt.id)

			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && err != nil && err.Error() != tt.errMsg {
				t.Errorf("Get() error = %v, want %v", err.Error(), tt.errMsg)
				return
			}

			if !tt.wantErr && milestone.ID == "" {
				t.Error("Get() returned empty milestone for existing ID")
			}
		})
	}
}

func TestStart(t *testing.T) {
	m := New()
	_, err := m.Create("m1", "First Milestone", "Description")
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	// Start the milestone
	err = m.Start("m1")
	if err != nil {
		t.Errorf("Start() error = %v", err)
	}

	// Verify status
	milestone, _ := m.Get("m1")
	if milestone.Status != StatusInProgress {
		t.Errorf("Start() Status = %v, want %v", milestone.Status, StatusInProgress)
	}

	// Test starting non-existing milestone
	err = m.Start("m2")
	if err == nil {
		t.Error("Expected error for non-existing milestone, got nil")
	}

	// Test starting with empty ID
	err = m.Start("")
	if err == nil {
		t.Error("Expected error for empty ID, got nil")
	}
}

func TestComplete(t *testing.T) {
	m := New()
	_, err := m.Create("m1", "First Milestone", "Description")
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	// Complete the milestone
	err = m.Complete("m1")
	if err != nil {
		t.Errorf("Complete() error = %v", err)
	}

	// Verify status
	milestone, _ := m.Get("m1")
	if milestone.Status != StatusCompleted {
		t.Errorf("Complete() Status = %v, want %v", milestone.Status, StatusCompleted)
	}
	if milestone.CompletedAt == nil {
		t.Error("Complete() CompletedAt should not be nil")
	}

	// Test completing already completed milestone
	err = m.Complete("m1")
	if err == nil {
		t.Error("Expected error for already completed milestone, got nil")
	}

	// Test completing non-existing milestone
	err = m.Complete("m2")
	if err == nil {
		t.Error("Expected error for non-existing milestone, got nil")
	}

	// Test completing with empty ID
	err = m.Complete("")
	if err == nil {
		t.Error("Expected error for empty ID, got nil")
	}
}

func TestStartCompletedMilestone(t *testing.T) {
	m := New()
	_, err := m.Create("m1", "First Milestone", "Description")
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	// Complete the milestone
	err = m.Complete("m1")
	if err != nil {
		t.Fatalf("Complete() error = %v", err)
	}

	// Try to start the completed milestone
	err = m.Start("m1")
	if err == nil {
		t.Error("Expected error for starting completed milestone, got nil")
	}
	if err.Error() != "cannot start a completed milestone" {
		t.Errorf("Start() error = %v, want 'cannot start a completed milestone'", err.Error())
	}
}

func TestList(t *testing.T) {
	m := New()

	// Empty list
	milestones := m.List()
	if len(milestones) != 0 {
		t.Errorf("List() returned %d milestones, want 0", len(milestones))
	}

	// Create milestones
	_, _ = m.Create("m1", "First Milestone", "Description 1")
	_, _ = m.Create("m2", "Second Milestone", "Description 2")

	milestones = m.List()
	if len(milestones) != 2 {
		t.Errorf("List() returned %d milestones, want 2", len(milestones))
	}
}

func TestCount(t *testing.T) {
	m := New()

	// Empty count
	if m.Count() != 0 {
		t.Errorf("Count() = %d, want 0", m.Count())
	}

	// Create milestones
	_, _ = m.Create("m1", "First Milestone", "Description 1")
	_, _ = m.Create("m2", "Second Milestone", "Description 2")

	if m.Count() != 2 {
		t.Errorf("Count() = %d, want 2", m.Count())
	}
}

func TestCountByStatus(t *testing.T) {
	m := New()

	// Create milestones with different statuses
	_, _ = m.Create("m1", "First Milestone", "Description 1")
	_, _ = m.Create("m2", "Second Milestone", "Description 2")
	_, _ = m.Create("m3", "Third Milestone", "Description 3")

	_ = m.Start("m2")
	_ = m.Complete("m3")

	if m.CountByStatus(StatusPending) != 1 {
		t.Errorf("CountByStatus(Pending) = %d, want 1", m.CountByStatus(StatusPending))
	}

	if m.CountByStatus(StatusInProgress) != 1 {
		t.Errorf("CountByStatus(InProgress) = %d, want 1", m.CountByStatus(StatusInProgress))
	}

	if m.CountByStatus(StatusCompleted) != 1 {
		t.Errorf("CountByStatus(Completed) = %d, want 1", m.CountByStatus(StatusCompleted))
	}
}
