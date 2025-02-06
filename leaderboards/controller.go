package leaderboards

import (
	"net/http"
	"strconv"

	"base/core/packages/gamification/models"
	"base/core/storage"

	"github.com/gin-gonic/gin"
)

type LeaderboardController struct {
	Service *LeaderboardService
	Storage *storage.ActiveStorage
}

func NewLeaderboardController(service *LeaderboardService, storage *storage.ActiveStorage) *LeaderboardController {
	return &LeaderboardController{
		Service: service,
		Storage: storage,
	}
}

func (c *LeaderboardController) Routes(router *gin.RouterGroup) {
	// Main CRUD endpoints
	router.GET("/leaderboards", c.List)        // Paginated list
	router.GET("/leaderboards/all", c.ListAll) // Unpaginated list
	router.GET("/leaderboards/:id", c.Get)
	router.POST("/leaderboards", c.Create)
	router.PUT("/leaderboards/:id", c.Update)
	router.DELETE("/leaderboards/:id", c.Delete)

	// File/Image attachment endpoints

	// HasMany relation endpoints
}

// CreateLeaderboard godoc
// @Summary Create a new Leaderboard
// @Description Create a new Leaderboard with the input payload
// @Tags Leaderboard
// @Security ApiKeyAuth
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param leaderboards body models.CreateLeaderboardRequest true "Create Leaderboard request"
// @Success 201 {object} models.LeaderboardResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /leaderboards [post]
func (c *LeaderboardController) Create(ctx *gin.Context) {
	var req models.CreateLeaderboardRequest
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

// GetLeaderboard godoc
// @Summary Get a Leaderboard
// @Description Get a Leaderboard by its id
// @Tags Leaderboard
// @Security ApiKeyAuth
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Leaderboard id"
// @Success 200 {object} models.LeaderboardResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /leaderboards/{id} [get]
func (c *LeaderboardController) Get(ctx *gin.Context) {
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

// ListLeaderboards godoc
// @Summary List leaderboards
// @Description Get a list of leaderboards
// @Tags Leaderboard
// @Security ApiKeyAuth
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param limit query int false "Number of items per page"
// @Success 200 {object} types.PaginatedResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /leaderboards [get]
func (c *LeaderboardController) List(ctx *gin.Context) {
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

// ListAllLeaderboards godoc
// @Summary List all leaderboards without pagination
// @Description Get a list of all leaderboards without pagination
// @Tags Leaderboard
// @Security ApiKeyAuth
// @Security BearerAuth
// @Accept json
// @Produce json
// @Success 200 {object} types.PaginatedResponse
// @Failure 500 {object} ErrorResponse
// @Router /leaderboards/all [get]
func (c *LeaderboardController) ListAll(ctx *gin.Context) {
	paginatedResponse, err := c.Service.GetAll(nil, nil)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to fetch all items: " + err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, paginatedResponse)
}

// UpdateLeaderboard godoc
// @Summary Update a Leaderboard
// @Description Update a Leaderboard by its id
// @Tags Leaderboard
// @Security ApiKeyAuth
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Leaderboard id"
// @Param leaderboards body models.UpdateLeaderboardRequest true "Update Leaderboard request"
// @Success 200 {object} models.LeaderboardResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /leaderboards/{id} [put]
func (c *LeaderboardController) Update(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid id format"})
		return
	}

	var req models.UpdateLeaderboardRequest
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

// DeleteLeaderboard godoc
// @Summary Delete a Leaderboard
// @Description Delete a Leaderboard by its id
// @Tags Leaderboard
// @Security ApiKeyAuth
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Leaderboard id"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /leaderboards/{id} [delete]
func (c *LeaderboardController) Delete(ctx *gin.Context) {
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
