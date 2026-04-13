package model

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `db:"id"`
	Name      string    `db:"name"`
	Email     string    `db:"email"`
	Password  string    `db:"password"`
	CreatedAt time.Time `db:"created_at"`
}

type Project struct {
	ID          uuid.UUID `db:"id"`
	Name        string    `db:"name"`
	Description string    `db:"description"`
	OwnerID     uuid.UUID `db:"owner_id"`
	CreatedAt   time.Time `db:"created_at"`
}

type Task struct {
	ID          uuid.UUID  `db:"id"`
	Title       string     `db:"title"`
	Description string     `db:"description"`
	Status      string     `db:"status"`
	Priority    string     `db:"priority"`
	ProjectID   uuid.UUID  `db:"project_id"`
	AssigneeID  *uuid.UUID `db:"assignee_id"`
	DueDate     *time.Time `db:"due_date"`
	CreatedAt   time.Time  `db:"created_at"`
	UpdatedAt   time.Time  `db:"updated_at"`
}

type ProjectWithTasks struct {
	Project
	Tasks []Task `db:"tasks"`
}

type Stats struct {
	TotalTasks     int            `json:"total_tasks"`
	ByStatus       map[string]int `json:"by_status"`
	ByAssignee     map[string]int `json:"by_assignee"`
}