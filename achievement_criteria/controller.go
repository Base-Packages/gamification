package achievement_criteria

import (
	"net/http"
	"strconv"

	"base/core/storage"
	"base/packages/gamification/models"

	"github.com/gin-gonic/gin"
)

type AchievementCriteriaController struct {
	Service *AchievementCriteriaService
	Storage *storage.ActiveStorage
}

func NewAchievementCriteriaController(service *AchievementCriteriaService, storage *storage.ActiveStorage) *AchievementCriteriaController {
	return &AchievementCriteriaController{
		Service: service,
		Storage: storage,
	}
}

func (c *AchievementCriteriaController) Routes(router *gin.RouterGroup) {
	// Main CRUD endpoints
	router.GET("/achievement-criteria", c.List)        // Paginated list
	router.GET("/achievement-criteria/all", c.ListAll) // Unpaginated list
	router.GET("/achievement-criteria/:id", c.Get)
	router.POST("/achievement-criteria", c.Create)
	router.PUT("/achievement-criteria/:id", c.Update)
	router.DELETE("/achievement-criteria/:id", c.Delete)

	// File/Image attachment endpoints

	// HasMany relation endpoints
}

// CreateAchievementCriteria godoc
// @Summary Create a new AchievementCriteria
// @Description Create a new AchievementCriteria with the input payload
// @Tags AchievementCriteria
// @Security ApiKeyAuth
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param achievement-criteria body models.CreateAchievementCriteriaRequest true "Create AchievementCriteria request"
// @Success 201 {object} models.AchievementCriteriaResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /achievement-criteria [post]
func (c *AchievementCriteriaController) Create(ctx *gin.Context) {
	var req models.CreateAchievementCriteriaRequest
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

// GetAchievementCriteria godoc
// @Summary Get a AchievementCriteria
// @Description Get a AchievementCriteria by its id
// @Tags AchievementCriteria
// @Security ApiKeyAuth
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "AchievementCriteria id"
// @Success 200 {object} models.AchievementCriteriaResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /achievement-criteria/{id} [get]
func (c *AchievementCriteriaController) Get(ctx *gin.Context) {
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

// ListAchievementCriteria godoc
// @Summary List achievement-criteria
// @Description Get a list of achievement-criteria
// @Tags AchievementCriteria
// @Security ApiKeyAuth
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param limit query int false "Number of items per page"
// @Success 200 {object} types.PaginatedResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /achievement-criteria [get]
func (c *AchievementCriteriaController) List(ctx *gin.Context) {
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

// ListAllAchievementCriteria godoc
// @Summary List all achievement-criteria without pagination
// @Description Get a list of all achievement-criteria without pagination
// @Tags AchievementCriteria
// @Security ApiKeyAuth
// @Security BearerAuth
// @Accept json
// @Produce json
// @Success 200 {object} types.PaginatedResponse
// @Failure 500 {object} ErrorResponse
// @Router /achievement-criteria/all [get]
func (c *AchievementCriteriaController) ListAll(ctx *gin.Context) {
	paginatedResponse, err := c.Service.GetAll(nil, nil)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to fetch all items: " + err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, paginatedResponse)
}

// UpdateAchievementCriteria godoc
// @Summary Update a AchievementCriteria
// @Description Update a AchievementCriteria by its id
// @Tags AchievementCriteria
// @Security ApiKeyAuth
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "AchievementCriteria id"
// @Param achievement-criteria body models.UpdateAchievementCriteriaRequest true "Update AchievementCriteria request"
// @Success 200 {object} models.AchievementCriteriaResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /achievement-criteria/{id} [put]
func (c *AchievementCriteriaController) Update(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid id format"})
		return
	}

	var req models.UpdateAchievementCriteriaRequest
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

// DeleteAchievementCriteria godoc
// @Summary Delete a AchievementCriteria
// @Description Delete a AchievementCriteria by its id
// @Tags AchievementCriteria
// @Security ApiKeyAuth
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "AchievementCriteria id"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /achievement-criteria/{id} [delete]
func (c *AchievementCriteriaController) Delete(ctx *gin.Context) {
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
