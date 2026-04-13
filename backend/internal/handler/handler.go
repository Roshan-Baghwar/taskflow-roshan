package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/Roshan-Baghwar/taskflow-roshan/backend/internal/middleware"
	"github.com/Roshan-Baghwar/taskflow-roshan/backend/internal/model"
	"github.com/Roshan-Baghwar/taskflow-roshan/backend/internal/service"
	"golang.org/x/crypto/bcrypt"
)

type Handler struct {
	service *service.Service
}

func NewHandler(s *service.Service) *Handler {
	return &Handler{service: s}
}

func (h *Handler) Register(c *gin.Context) {
	var req struct {
		Name     string `json:"name" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=6"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "validation failed", "fields": err.Error()})
		return
	}

	hashed, _ := bcrypt.GenerateFromPassword([]byte(req.Password), 12)
	user := &model.User{Name: req.Name, Email: req.Email, Password: string(hashed)}
	if err := h.service.CreateUser(user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email already exists"})
		return
	}

	token, _ := middleware.GenerateToken(user.ID, user.Email)
	c.JSON(http.StatusCreated, gin.H{"token": token, "user": gin.H{"id": user.ID, "name": user.Name, "email": user.Email}})
}

func (h *Handler) Login(c *gin.Context) {
	var req struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "validation failed"})
		return
	}

	user, err := h.service.GetUserByEmail(req.Email)
	if err != nil || bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)) != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	token, _ := middleware.GenerateToken(user.ID, user.Email)
	c.JSON(http.StatusOK, gin.H{"token": token, "user": gin.H{"id": user.ID, "name": user.Name, "email": user.Email}})
}

// Projects
func (h *Handler) GetProjects(c *gin.Context) {
	userID := middleware.GetUserID(c)
	projects, _ := h.service.GetProjectsForUser(userID)
	c.JSON(http.StatusOK, gin.H{"projects": projects})
}

func (h *Handler) CreateProject(c *gin.Context) {
	userID := middleware.GetUserID(c)
	var req struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "validation failed"})
		return
	}
	p := &model.Project{Name: req.Name, Description: req.Description, OwnerID: userID}
	h.service.CreateProject(p)
	c.JSON(http.StatusCreated, p)
}

func (h *Handler) GetProject(c *gin.Context) {
	id := uuid.MustParse(c.Param("id"))
	p, err := h.service.GetProjectByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.JSON(http.StatusOK, p)
}

func (h *Handler) UpdateProject(c *gin.Context) {
	id := uuid.MustParse(c.Param("id"))
	userID := middleware.GetUserID(c)
	var req struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}
	c.ShouldBindJSON(&req)
	if err := h.service.UpdateProject(id, userID, req.Name, req.Description); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "updated"})
}

func (h *Handler) DeleteProject(c *gin.Context) {
	id := uuid.MustParse(c.Param("id"))
	userID := middleware.GetUserID(c)
	if err := h.service.DeleteProject(id, userID); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}
	c.Status(http.StatusNoContent)
}

// Tasks
func (h *Handler) GetProjectTasks(c *gin.Context) {
	projectID := uuid.MustParse(c.Param("id"))
	status := c.Query("status")
	assignee := c.Query("assignee")

	tasks, _ := h.service.GetTasksByProject(projectID, status, assignee)

	// Bonus pagination (simple)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	start := (page - 1) * limit
	end := start + limit
	if end > len(tasks) {
		end = len(tasks)
	}
	if start > len(tasks) {
		tasks = []model.Task{}
	} else {
		tasks = tasks[start:end]
	}

	c.JSON(http.StatusOK, gin.H{"tasks": tasks})
}

func (h *Handler) CreateTask(c *gin.Context) {
	projectID := uuid.MustParse(c.Param("id"))
	var req model.Task
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "validation failed"})
		return
	}
	req.ProjectID = projectID
	if req.Status == "" {
		req.Status = "todo"
	}
	if req.Priority == "" {
		req.Priority = "medium"
	}
	h.service.CreateTask(&req)
	c.JSON(http.StatusCreated, req)
}

func (h *Handler) UpdateTask(c *gin.Context) {
	id := uuid.MustParse(c.Param("id"))
	var req model.Task
	c.ShouldBindJSON(&req)
	req.ID = id
	if err := h.service.UpdateTask(&req); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.JSON(http.StatusOK, req)
}

func (h *Handler) DeleteTask(c *gin.Context) {
	id := uuid.MustParse(c.Param("id"))
	h.service.DeleteTask(id)
	c.Status(http.StatusNoContent)
}

func (h *Handler) GetProjectStats(c *gin.Context) {
	id := uuid.MustParse(c.Param("id"))
	stats, _ := h.service.GetProjectStats(id)
	c.JSON(http.StatusOK, stats)
}