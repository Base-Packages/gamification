package user_activities

import (
	"net/http"
	"strconv"

	"base/core/storage"
	"base/packages/gamification/models"

	"github.com/gin-gonic/gin"
)

type UserActivityController struct {
	Service *UserActivityService
	Storage *storage.ActiveStorage
}

func NewUserActivityController(service *UserActivityService, storage *storage.ActiveStorage) *UserActivityController {
	return &UserActivityController{
		Service: service,
		Storage: storage,
	}
}

func (c *UserActivityController) Routes(router *gin.RouterGroup) {
	// Main CRUD endpoints
	router.GET("/user-activities", c.List)        // Paginated list
	router.GET("/user-activities/all", c.ListAll) // Unpaginated list
	router.GET("/user-activities/:id", c.Get)
	router.POST("/user-activities", c.Create)
	router.PUT("/user-activities/:id", c.Update)
	router.DELETE("/user-activities/:id", c.Delete)

	// File/Image attachment endpoints

	// HasMany relation endpoints
}

// CreateUserActivity godoc
// @Summary Create a new UserActivity
// @Description Create a new UserActivity with the input payload
// @Tags Package/Gamification/UserActivity
// @Security ApiKeyAuth
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param user-activities body models.CreateUserActivityRequest true "Create UserActivity request"
// @Success 201 {object} models.UserActivityResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /user-activities [post]
func (c *UserActivityController) Create(ctx *gin.Context) {
	var req models.CreateUserActivityRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	item, err := c.Service.Create(&req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to create item: " + err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, item.ToResponse())
}

// GetUserActivity godoc
// @Summary Get a UserActivity
// @Description Get a UserActivity by its id
// @Tags UserActivity
// @Security ApiKeyAuth
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "UserActivity id"
// @Success 200 {object} models.UserActivityResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /user-activities/{id} [get]
func (c *UserActivityController) Get(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid id format"})
		return
	}

	item, err := c.Service.GetById(uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, ErrorResponse{Error: "Item not found"})
		return
	}

	ctx.JSON(http.StatusOK, item.ToResponse())
}

// ListUserActivities godoc
// @Summary List user-activities
// @Description Get a list of user-activities
// @Tags UserActivity
// @Security ApiKeyAuth
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param limit query int false "Number of items per page"
// @Success 200 {object} types.PaginatedResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /user-activities [get]
func (c *UserActivityController) List(ctx *gin.Context) {
	var page, limit *int

	if pageStr := ctx.Query("page"); pageStr != "" {
		if pageNum, err := strconv.Atoi(pageStr); err == nil && pageNum > 0 {
			page = &pageNum
		} else {
			ctx.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid page number"})
			return
		}
	}

	if limitStr := ctx.Query("limit"); limitStr != "" {
		if limitNum, err := strconv.Atoi(limitStr); err == nil && limitNum > 0 {
			limit = &limitNum
		} else {
			ctx.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid limit number"})
			return
		}
	}

	paginatedResponse, err := c.Service.GetAll(page, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to fetch items: " + err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, paginatedResponse)
}

// ListAllUserActivities godoc
// @Summary List all user-activities without pagination
// @Description Get a list of all user-activities without pagination
// @Tags UserActivity
// @Security ApiKeyAuth
// @Security BearerAuth
// @Accept json
// @Produce json
// @Success 200 {object} types.PaginatedResponse
// @Failure 500 {object} ErrorResponse
// @Router /user-activities/all [get]
func (c *UserActivityController) ListAll(ctx *gin.Context) {
	paginatedResponse, err := c.Service.GetAll(nil, nil)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to fetch all items: " + err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, paginatedResponse)
}

// UpdateUserActivity godoc
// @Summary Update a UserActivity
// @Description Update a UserActivity by its id
// @Tags UserActivity
// @Security ApiKeyAuth
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "UserActivity id"
// @Param user-activities body models.UpdateUserActivityRequest true "Update UserActivity request"
// @Success 200 {object} models.UserActivityResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /user-activities/{id} [put]
func (c *UserActivityController) Update(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid id format"})
		return
	}

	var req models.UpdateUserActivityRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	item, err := c.Service.Update(uint(id), &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to update item: " + err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, item.ToResponse())
}

// DeleteUserActivity godoc
// @Summary Delete a UserActivity
// @Description Delete a UserActivity by its id
// @Tags UserActivity
// @Security ApiKeyAuth
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "UserActivity id"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /user-activities/{id} [delete]
func (c *UserActivityController) Delete(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid id format"})
		return
	}

	if err := c.Service.Delete(uint(id)); err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to delete item: " + err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, SuccessResponse{Message: "Item deleted successfully"})
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type SuccessResponse struct {
	Message string `json:"message"`
}
