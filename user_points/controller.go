package user_points

import (
	"net/http"
	"strconv"

	"base/core/storage"
	"base/packages/gamification/models"

	"github.com/gin-gonic/gin"
)

type UserPointController struct {
	Service *UserPointService
	Storage *storage.ActiveStorage
}

func NewUserPointController(service *UserPointService, storage *storage.ActiveStorage) *UserPointController {
	return &UserPointController{
		Service: service,
		Storage: storage,
	}
}

func (c *UserPointController) Routes(router *gin.RouterGroup) {
	// Main CRUD endpoints
	router.GET("/user-points", c.List)        // Paginated list
	router.GET("/user-points/all", c.ListAll) // Unpaginated list
	router.GET("/user-points/:id", c.Get)
	router.POST("/user-points", c.Create)
	router.PUT("/user-points/:id", c.Update)
	router.DELETE("/user-points/:id", c.Delete)

	// File/Image attachment endpoints

	// HasMany relation endpoints
}

// CreateUserPoint godoc
// @Summary Create a new UserPoint
// @Description Create a new UserPoint with the input payload
// @Tags UserPoint
// @Security ApiKeyAuth
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param user-points body models.CreateUserPointRequest true "Create UserPoint request"
// @Success 201 {object} models.UserPointResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /user-points [post]
func (c *UserPointController) Create(ctx *gin.Context) {
	var req models.CreateUserPointRequest
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

// GetUserPoint godoc
// @Summary Get a UserPoint
// @Description Get a UserPoint by its id
// @Tags UserPoint
// @Security ApiKeyAuth
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "UserPoint id"
// @Success 200 {object} models.UserPointResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /user-points/{id} [get]
func (c *UserPointController) Get(ctx *gin.Context) {
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

// ListUserPoints godoc
// @Summary List user-points
// @Description Get a list of user-points
// @Tags UserPoint
// @Security ApiKeyAuth
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param limit query int false "Number of items per page"
// @Success 200 {object} types.PaginatedResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /user-points [get]
func (c *UserPointController) List(ctx *gin.Context) {
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

// ListAllUserPoints godoc
// @Summary List all user-points without pagination
// @Description Get a list of all user-points without pagination
// @Tags UserPoint
// @Security ApiKeyAuth
// @Security BearerAuth
// @Accept json
// @Produce json
// @Success 200 {object} types.PaginatedResponse
// @Failure 500 {object} ErrorResponse
// @Router /user-points/all [get]
func (c *UserPointController) ListAll(ctx *gin.Context) {
	paginatedResponse, err := c.Service.GetAll(nil, nil)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to fetch all items: " + err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, paginatedResponse)
}

// UpdateUserPoint godoc
// @Summary Update a UserPoint
// @Description Update a UserPoint by its id
// @Tags UserPoint
// @Security ApiKeyAuth
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "UserPoint id"
// @Param user-points body models.UpdateUserPointRequest true "Update UserPoint request"
// @Success 200 {object} models.UserPointResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /user-points/{id} [put]
func (c *UserPointController) Update(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid id format"})
		return
	}

	var req models.UpdateUserPointRequest
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

// DeleteUserPoint godoc
// @Summary Delete a UserPoint
// @Description Delete a UserPoint by its id
// @Tags UserPoint
// @Security ApiKeyAuth
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "UserPoint id"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /user-points/{id} [delete]
func (c *UserPointController) Delete(ctx *gin.Context) {
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
