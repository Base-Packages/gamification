package leaderboard_entries

import (
	"net/http"
	"strconv"

	"base/core/storage"
	"base/packages/gamification/models"

	"github.com/gin-gonic/gin"
)

type LeaderboardEntryController struct {
	Service *LeaderboardEntryService
	Storage *storage.ActiveStorage
}

func NewLeaderboardEntryController(service *LeaderboardEntryService, storage *storage.ActiveStorage) *LeaderboardEntryController {
	return &LeaderboardEntryController{
		Service: service,
		Storage: storage,
	}
}

func (c *LeaderboardEntryController) Routes(router *gin.RouterGroup) {
	// Main CRUD endpoints
	router.GET("/leaderboard-entries", c.List)        // Paginated list
	router.GET("/leaderboard-entries/all", c.ListAll) // Unpaginated list
	router.GET("/leaderboard-entries/:id", c.Get)
	router.POST("/leaderboard-entries", c.Create)
	router.PUT("/leaderboard-entries/:id", c.Update)
	router.DELETE("/leaderboard-entries/:id", c.Delete)

	// File/Image attachment endpoints

	// HasMany relation endpoints
}

// CreateLeaderboardEntry godoc
// @Summary Create a new LeaderboardEntry
// @Description Create a new LeaderboardEntry with the input payload
// @Tags LeaderboardEntry
// @Security ApiKeyAuth
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param leaderboard-entries body models.CreateLeaderboardEntryRequest true "Create LeaderboardEntry request"
// @Success 201 {object} models.LeaderboardEntryResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /leaderboard-entries [post]
func (c *LeaderboardEntryController) Create(ctx *gin.Context) {
	var req models.CreateLeaderboardEntryRequest
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

// GetLeaderboardEntry godoc
// @Summary Get a LeaderboardEntry
// @Description Get a LeaderboardEntry by its id
// @Tags LeaderboardEntry
// @Security ApiKeyAuth
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "LeaderboardEntry id"
// @Success 200 {object} models.LeaderboardEntryResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /leaderboard-entries/{id} [get]
func (c *LeaderboardEntryController) Get(ctx *gin.Context) {
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

// ListLeaderboardEntries godoc
// @Summary List leaderboard-entries
// @Description Get a list of leaderboard-entries
// @Tags LeaderboardEntry
// @Security ApiKeyAuth
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param limit query int false "Number of items per page"
// @Success 200 {object} types.PaginatedResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /leaderboard-entries [get]
func (c *LeaderboardEntryController) List(ctx *gin.Context) {
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

// ListAllLeaderboardEntries godoc
// @Summary List all leaderboard-entries without pagination
// @Description Get a list of all leaderboard-entries without pagination
// @Tags LeaderboardEntry
// @Security ApiKeyAuth
// @Security BearerAuth
// @Accept json
// @Produce json
// @Success 200 {object} types.PaginatedResponse
// @Failure 500 {object} ErrorResponse
// @Router /leaderboard-entries/all [get]
func (c *LeaderboardEntryController) ListAll(ctx *gin.Context) {
	paginatedResponse, err := c.Service.GetAll(nil, nil)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to fetch all items: " + err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, paginatedResponse)
}

// UpdateLeaderboardEntry godoc
// @Summary Update a LeaderboardEntry
// @Description Update a LeaderboardEntry by its id
// @Tags LeaderboardEntry
// @Security ApiKeyAuth
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "LeaderboardEntry id"
// @Param leaderboard-entries body models.UpdateLeaderboardEntryRequest true "Update LeaderboardEntry request"
// @Success 200 {object} models.LeaderboardEntryResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /leaderboard-entries/{id} [put]
func (c *LeaderboardEntryController) Update(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid id format"})
		return
	}

	var req models.UpdateLeaderboardEntryRequest
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

// DeleteLeaderboardEntry godoc
// @Summary Delete a LeaderboardEntry
// @Description Delete a LeaderboardEntry by its id
// @Tags LeaderboardEntry
// @Security ApiKeyAuth
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "LeaderboardEntry id"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /leaderboard-entries/{id} [delete]
func (c *LeaderboardEntryController) Delete(ctx *gin.Context) {
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
