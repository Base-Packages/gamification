package user_achievements

import (
	"net/http"
	"strconv"

	"base/core/packages/gamification/models"
	"base/core/storage"

	"github.com/gin-gonic/gin"
)

type UserAchievementController struct {
	Service *UserAchievementService
	Storage *storage.ActiveStorage
}

func NewUserAchievementController(service *UserAchievementService, storage *storage.ActiveStorage) *UserAchievementController {
	return &UserAchievementController{
		Service: service,
		Storage: storage,
	}
}

func (c *UserAchievementController) Routes(router *gin.RouterGroup) {
	// Main CRUD endpoints
	router.GET("/user-achievements", c.List)        // Paginated list
	router.GET("/user-achievements/all", c.ListAll) // Unpaginated list
	router.GET("/user-achievements/:id", c.Get)
	router.POST("/user-achievements", c.Create)
	router.PUT("/user-achievements/:id", c.Update)
	router.DELETE("/user-achievements/:id", c.Delete)

	// File/Image attachment endpoints

	// HasMany relation endpoints
}

// CreateUserAchievement godoc
// @Summary Create a new UserAchievement
// @Description Create a new UserAchievement with the input payload
// @Tags UserAchievement
// @Security ApiKeyAuth
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param user-achievements body models.CreateUserAchievementRequest true "Create UserAchievement request"
// @Success 201 {object} models.UserAchievementResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /user-achievements [post]
func (c *UserAchievementController) Create(ctx *gin.Context) {
	var req models.CreateUserAchievementRequest
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

// GetUserAchievement godoc
// @Summary Get a UserAchievement
// @Description Get a UserAchievement by its id
// @Tags UserAchievement
// @Security ApiKeyAuth
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "UserAchievement id"
// @Success 200 {object} models.UserAchievementResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /user-achievements/{id} [get]
func (c *UserAchievementController) Get(ctx *gin.Context) {
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

// ListUserAchievements godoc
// @Summary List user-achievements
// @Description Get a list of user-achievements
// @Tags UserAchievement
// @Security ApiKeyAuth
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param limit query int false "Number of items per page"
// @Success 200 {object} types.PaginatedResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /user-achievements [get]
func (c *UserAchievementController) List(ctx *gin.Context) {
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

// ListAllUserAchievements godoc
// @Summary List all user-achievements without pagination
// @Description Get a list of all user-achievements without pagination
// @Tags UserAchievement
// @Security ApiKeyAuth
// @Security BearerAuth
// @Accept json
// @Produce json
// @Success 200 {object} types.PaginatedResponse
// @Failure 500 {object} ErrorResponse
// @Router /user-achievements/all [get]
func (c *UserAchievementController) ListAll(ctx *gin.Context) {
	paginatedResponse, err := c.Service.GetAll(nil, nil)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to fetch all items: " + err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, paginatedResponse)
}

// UpdateUserAchievement godoc
// @Summary Update a UserAchievement
// @Description Update a UserAchievement by its id
// @Tags UserAchievement
// @Security ApiKeyAuth
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "UserAchievement id"
// @Param user-achievements body models.UpdateUserAchievementRequest true "Update UserAchievement request"
// @Success 200 {object} models.UserAchievementResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /user-achievements/{id} [put]
func (c *UserAchievementController) Update(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid id format"})
		return
	}

	var req models.UpdateUserAchievementRequest
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

// DeleteUserAchievement godoc
// @Summary Delete a UserAchievement
// @Description Delete a UserAchievement by its id
// @Tags UserAchievement
// @Security ApiKeyAuth
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "UserAchievement id"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /user-achievements/{id} [delete]
func (c *UserAchievementController) Delete(ctx *gin.Context) {
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
