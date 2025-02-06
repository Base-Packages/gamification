package activity_types

import (
	"net/http"
	"strconv"

	"base/core/packages/gamification/models"
	"base/core/storage"

	"github.com/gin-gonic/gin"
)

type ActivityTypeController struct {
	Service *ActivityTypeService
	Storage *storage.ActiveStorage
}

func NewActivityTypeController(service *ActivityTypeService, storage *storage.ActiveStorage) *ActivityTypeController {
	return &ActivityTypeController{
		Service: service,
		Storage: storage,
	}
}

func (c *ActivityTypeController) Routes(router *gin.RouterGroup) {
	// Main CRUD endpoints
	router.GET("/activity-types", c.List)        // Paginated list
	router.GET("/activity-types/all", c.ListAll) // Unpaginated list
	router.GET("/activity-types/:id", c.Get)
	router.POST("/activity-types", c.Create)
	router.PUT("/activity-types/:id", c.Update)
	router.DELETE("/activity-types/:id", c.Delete)

	// File/Image attachment endpoints

	// HasMany relation endpoints
}

// CreateActivityType godoc
// @Summary Create a new ActivityType
// @Description Create a new ActivityType with the input payload
// @Tags ActivityType
// @Security ApiKeyAuth
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param activity-types body models.CreateActivityTypeRequest true "Create ActivityType request"
// @Success 201 {object} models.ActivityTypeResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /activity-types [post]
func (c *ActivityTypeController) Create(ctx *gin.Context) {
	var req models.CreateActivityTypeRequest
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

// GetActivityType godoc
// @Summary Get a ActivityType
// @Description Get a ActivityType by its id
// @Tags ActivityType
// @Security ApiKeyAuth
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "ActivityType id"
// @Success 200 {object} models.ActivityTypeResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /activity-types/{id} [get]
func (c *ActivityTypeController) Get(ctx *gin.Context) {
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

// ListActivityTypes godoc
// @Summary List activity-types
// @Description Get a list of activity-types
// @Tags ActivityType
// @Security ApiKeyAuth
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param limit query int false "Number of items per page"
// @Success 200 {object} types.PaginatedResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /activity-types [get]
func (c *ActivityTypeController) List(ctx *gin.Context) {
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

// ListAllActivityTypes godoc
// @Summary List all activity-types without pagination
// @Description Get a list of all activity-types without pagination
// @Tags ActivityType
// @Security ApiKeyAuth
// @Security BearerAuth
// @Accept json
// @Produce json
// @Success 200 {object} types.PaginatedResponse
// @Failure 500 {object} ErrorResponse
// @Router /activity-types/all [get]
func (c *ActivityTypeController) ListAll(ctx *gin.Context) {
	paginatedResponse, err := c.Service.GetAll(nil, nil)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to fetch all items: " + err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, paginatedResponse)
}

// UpdateActivityType godoc
// @Summary Update a ActivityType
// @Description Update a ActivityType by its id
// @Tags ActivityType
// @Security ApiKeyAuth
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "ActivityType id"
// @Param activity-types body models.UpdateActivityTypeRequest true "Update ActivityType request"
// @Success 200 {object} models.ActivityTypeResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /activity-types/{id} [put]
func (c *ActivityTypeController) Update(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid id format"})
		return
	}

	var req models.UpdateActivityTypeRequest
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

// DeleteActivityType godoc
// @Summary Delete a ActivityType
// @Description Delete a ActivityType by its id
// @Tags ActivityType
// @Security ApiKeyAuth
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "ActivityType id"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /activity-types/{id} [delete]
func (c *ActivityTypeController) Delete(ctx *gin.Context) {
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
