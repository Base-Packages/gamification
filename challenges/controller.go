package challenges

import (
	"net/http"
	"strconv"

	"base/core/packages/gamification/models"
	"base/core/storage"

	"github.com/gin-gonic/gin"
)

type ChallengeController struct {
	Service *ChallengeService
	Storage *storage.ActiveStorage
}

func NewChallengeController(service *ChallengeService, storage *storage.ActiveStorage) *ChallengeController {
	return &ChallengeController{
		Service: service,
		Storage: storage,
	}
}

func (c *ChallengeController) Routes(router *gin.RouterGroup) {
	// Main CRUD endpoints
	router.GET("/challenges", c.List)        // Paginated list
	router.GET("/challenges/all", c.ListAll) // Unpaginated list
	router.GET("/challenges/:id", c.Get)
	router.POST("/challenges", c.Create)
	router.PUT("/challenges/:id", c.Update)
	router.DELETE("/challenges/:id", c.Delete)

	// File/Image attachment endpoints

	// HasMany relation endpoints
}

// CreateChallenge godoc
// @Summary Create a new Challenge
// @Description Create a new Challenge with the input payload
// @Tags Challenge
// @Security ApiKeyAuth
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param challenges body models.CreateChallengeRequest true "Create Challenge request"
// @Success 201 {object} models.ChallengeResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /challenges [post]
func (c *ChallengeController) Create(ctx *gin.Context) {
	var req models.CreateChallengeRequest
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

// GetChallenge godoc
// @Summary Get a Challenge
// @Description Get a Challenge by its id
// @Tags Challenge
// @Security ApiKeyAuth
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Challenge id"
// @Success 200 {object} models.ChallengeResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /challenges/{id} [get]
func (c *ChallengeController) Get(ctx *gin.Context) {
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

// ListChallenges godoc
// @Summary List challenges
// @Description Get a list of challenges
// @Tags Challenge
// @Security ApiKeyAuth
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param limit query int false "Number of items per page"
// @Success 200 {object} types.PaginatedResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /challenges [get]
func (c *ChallengeController) List(ctx *gin.Context) {
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

// ListAllChallenges godoc
// @Summary List all challenges without pagination
// @Description Get a list of all challenges without pagination
// @Tags Challenge
// @Security ApiKeyAuth
// @Security BearerAuth
// @Accept json
// @Produce json
// @Success 200 {object} types.PaginatedResponse
// @Failure 500 {object} ErrorResponse
// @Router /challenges/all [get]
func (c *ChallengeController) ListAll(ctx *gin.Context) {
	paginatedResponse, err := c.Service.GetAll(nil, nil)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to fetch all items: " + err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, paginatedResponse)
}

// UpdateChallenge godoc
// @Summary Update a Challenge
// @Description Update a Challenge by its id
// @Tags Challenge
// @Security ApiKeyAuth
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Challenge id"
// @Param challenges body models.UpdateChallengeRequest true "Update Challenge request"
// @Success 200 {object} models.ChallengeResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /challenges/{id} [put]
func (c *ChallengeController) Update(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid id format"})
		return
	}

	var req models.UpdateChallengeRequest
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

// DeleteChallenge godoc
// @Summary Delete a Challenge
// @Description Delete a Challenge by its id
// @Tags Challenge
// @Security ApiKeyAuth
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Challenge id"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /challenges/{id} [delete]
func (c *ChallengeController) Delete(ctx *gin.Context) {
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
