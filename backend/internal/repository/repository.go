package repository

import (
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/Roshan-Baghwar/taskflow-roshan/backend/internal/model"
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

// User
func (r *Repository) CreateUser(user *model.User) error {
	query := `INSERT INTO users (name, email, password) VALUES ($1, $2, $3) RETURNING id, created_at`
	return r.db.QueryRowx(query, user.Name, user.Email, user.Password).StructScan(user)
}

func (r *Repository) GetUserByEmail(email string) (*model.User, error) {
	var user model.User
	err := r.db.Get(&user, "SELECT * FROM users WHERE email = $1", email)
	return &user, err
}

// Project
func (r *Repository) CreateProject(p *model.Project) error {
	query := `INSERT INTO projects (name, description, owner_id) VALUES ($1, $2, $3) RETURNING id, created_at`
	return r.db.QueryRowx(query, p.Name, p.Description, p.OwnerID).StructScan(p)
}

func (r *Repository) GetProjectsForUser(userID uuid.UUID) ([]model.Project, error) {
	query := `
		SELECT DISTINCT p.* FROM projects p 
		LEFT JOIN tasks t ON p.id = t.project_id 
		WHERE p.owner_id = $1 OR t.assignee_id = $1
		ORDER BY p.created_at DESC`
	var projects []model.Project
	err := r.db.Select(&projects, query, userID)
	return projects, err
}

func (r *Repository) GetProjectByID(id uuid.UUID) (*model.Project, error) {
	var p model.Project
	err := r.db.Get(&p, "SELECT * FROM projects WHERE id = $1", id)
	return &p, err
}

func (r *Repository) UpdateProject(id, ownerID uuid.UUID, name, desc string) error {
	_, err := r.db.Exec("UPDATE projects SET name=$1, description=$2 WHERE id=$3 AND owner_id=$4", name, desc, id, ownerID)
	return err
}

func (r *Repository) DeleteProject(id, ownerID uuid.UUID) error {
	_, err := r.db.Exec("DELETE FROM projects WHERE id=$1 AND owner_id=$2", id, ownerID)
	return err
}

// Task
func (r *Repository) CreateTask(task *model.Task) error {
	query := `INSERT INTO tasks (title, description, status, priority, project_id, assignee_id, due_date) 
	          VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id, created_at, updated_at`
	return r.db.QueryRowx(query, task.Title, task.Description, task.Status, task.Priority, task.ProjectID, task.AssigneeID, task.DueDate).StructScan(task)
}

func (r *Repository) GetTasksByProject(projectID uuid.UUID, status, assignee string) ([]model.Task, error) {
	query := "SELECT * FROM tasks WHERE project_id = $1"
	args := []interface{}{projectID}
	if status != "" {
		query += " AND status = $2"
		args = append(args, status)
	}
	if assignee != "" {
		assigneeID, _ := uuid.Parse(assignee)
		query += " AND assignee_id = $3"
		args = append(args, assigneeID)
	}
	var tasks []model.Task
	err := r.db.Select(&tasks, query, args...)
	return tasks, err
}

func (r *Repository) GetTaskByID(id uuid.UUID) (*model.Task, error) {
	var t model.Task
	err := r.db.Get(&t, "SELECT * FROM tasks WHERE id = $1", id)
	return &t, err
}

func (r *Repository) UpdateTask(task *model.Task) error {
	query := `UPDATE tasks SET title=$1, description=$2, status=$3, priority=$4, assignee_id=$5, due_date=$6, updated_at=NOW() 
	          WHERE id=$7 RETURNING updated_at`
	return r.db.QueryRowx(query, task.Title, task.Description, task.Status, task.Priority, task.AssigneeID, task.DueDate, task.ID).Scan(&task.UpdatedAt)
}

func (r *Repository) DeleteTask(id uuid.UUID) error {
	_, err := r.db.Exec("DELETE FROM tasks WHERE id=$1", id)
	return err
}

func (r *Repository) GetProjectStats(projectID uuid.UUID) (model.Stats, error) {
	var stats model.Stats
	stats.ByStatus = make(map[string]int)
	stats.ByAssignee = make(map[string]int)

	// Total tasks
	r.db.Get(&stats.TotalTasks, "SELECT COUNT(*) FROM tasks WHERE project_id = $1", projectID)

	// By status
	rows, _ := r.db.Query("SELECT status, COUNT(*) FROM tasks WHERE project_id = $1 GROUP BY status", projectID)
	for rows.Next() {
		var s string
		var c int
		rows.Scan(&s, &c)
		stats.ByStatus[s] = c
	}
	rows.Close()

	// By assignee
	rows, _ = r.db.Query(`
		SELECT COALESCE(u.name, 'Unassigned'), COUNT(*) 
		FROM tasks t LEFT JOIN users u ON t.assignee_id = u.id 
		WHERE t.project_id = $1 GROUP BY u.name`, projectID)
	for rows.Next() {
		var name string
		var c int
		rows.Scan(&name, &c)
		stats.ByAssignee[name] = c
	}
	rows.Close()

	return stats, nil
}