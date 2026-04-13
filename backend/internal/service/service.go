package service

import (
	"github.com/google/uuid"
	"github.com/Roshan-Baghwar/taskflow-roshan/backend/internal/model"
	"github.com/Roshan-Baghwar/taskflow-roshan/backend/internal/repository"
)

type Service struct {
	repo *repository.Repository
}

func NewService(repo *repository.Repository) *Service {
	return &Service{repo: repo}
}

// User
func (s *Service) CreateUser(user *model.User) error {
	return s.repo.CreateUser(user)
}

func (s *Service) GetUserByEmail(email string) (*model.User, error) {
	return s.repo.GetUserByEmail(email)
}

// Project
func (s *Service) CreateProject(p *model.Project) error {
	return s.repo.CreateProject(p)
}

func (s *Service) GetProjectsForUser(userID uuid.UUID) ([]model.Project, error) {
	return s.repo.GetProjectsForUser(userID)
}

func (s *Service) GetProjectByID(id uuid.UUID) (*model.Project, error) {
	return s.repo.GetProjectByID(id)
}

func (s *Service) UpdateProject(id, ownerID uuid.UUID, name, desc string) error {
	return s.repo.UpdateProject(id, ownerID, name, desc)
}

func (s *Service) DeleteProject(id, ownerID uuid.UUID) error {
	return s.repo.DeleteProject(id, ownerID)
}

// Task
func (s *Service) CreateTask(task *model.Task) error {
	return s.repo.CreateTask(task)
}

func (s *Service) GetTasksByProject(projectID uuid.UUID, status, assignee string) ([]model.Task, error) {
	return s.repo.GetTasksByProject(projectID, status, assignee)
}

func (s *Service) UpdateTask(task *model.Task) error {
	return s.repo.UpdateTask(task)
}

func (s *Service) DeleteTask(id uuid.UUID) error {
	return s.repo.DeleteTask(id)
}

func (s *Service) GetProjectStats(projectID uuid.UUID) (model.Stats, error) {
	return s.repo.GetProjectStats(projectID)
}