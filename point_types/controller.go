package point_types

import (
	"net/http"
	"strconv"

	"base/core/storage"
	"base/packages/gamification/models"

	"github.com/gin-gonic/gin"
)

type PointTypeController struct {
	Service *PointTypeService
	Storage *storage.ActiveStorage
}

func NewPointTypeController(service *PointTypeService, storage *storage.ActiveStorage) *PointTypeController {
	return &PointTypeController{
		Service: service,
		Storage: storage,
	}
}

func (c *PointTypeController) Routes(router *gin.RouterGroup) {
	// Main CRUD endpoints
	router.GET("/point-types", c.List)        // Paginated list
	router.GET("/point-types/all", c.ListAll) // Unpaginated list
	router.GET("/point-types/:id", c.Get)
	router.POST("/point-types", c.Create)
	router.PUT("/point-types/:id", c.Update)
	router.DELETE("/point-types/:id", c.Delete)

	// File/Image attachment endpoints

	// HasMany relation endpoints
}

// CreatePointType godoc
// @Summary Create a new PointType
// @Description Create a new PointType with the input payload
// @Tags PointType
// @Security ApiKeyAuth
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param point-types body models.CreatePointTypeRequest true "Create PointType request"
// @Success 201 {object} models.PointTypeResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /point-types [post]
func (c *PointTypeController) Create(ctx *gin.Context) {
	var req models.CreatePointTypeRequest
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

// GetPointType godoc
// @Summary Get a PointType
// @Description Get a PointType by its id
// @Tags PointType
// @Security ApiKeyAuth
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "PointType id"
// @Success 200 {object} models.PointTypeResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /point-types/{id} [get]
func (c *PointTypeController) Get(ctx *gin.Context) {
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

// ListPointTypes godoc
// @Summary List point-types
// @Description Get a list of point-types
// @Tags PointType
// @Security ApiKeyAuth
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param limit query int false "Number of items per page"
// @Success 200 {object} types.PaginatedResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /point-types [get]
func (c *PointTypeController) List(ctx *gin.Context) {
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

// ListAllPointTypes godoc
// @Summary List all point-types without pagination
// @Description Get a list of all point-types without pagination
// @Tags PointType
// @Security ApiKeyAuth
// @Security BearerAuth
// @Accept json
// @Produce json
// @Success 200 {object} types.PaginatedResponse
// @Failure 500 {object} ErrorResponse
// @Router /point-types/all [get]
func (c *PointTypeController) ListAll(ctx *gin.Context) {
	paginatedResponse, err := c.Service.GetAll(nil, nil)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to fetch all items: " + err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, paginatedResponse)
}

// UpdatePointType godoc
// @Summary Update a PointType
// @Description Update a PointType by its id
// @Tags PointType
// @Security ApiKeyAuth
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "PointType id"
// @Param point-types body models.UpdatePointTypeRequest true "Update PointType request"
// @Success 200 {object} models.PointTypeResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /point-types/{id} [put]
func (c *PointTypeController) Update(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid id format"})
		return
	}

	var req models.UpdatePointTypeRequest
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

// DeletePointType godoc
// @Summary Delete a PointType
// @Description Delete a PointType by its id
// @Tags PointType
// @Security ApiKeyAuth
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "PointType id"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /point-types/{id} [delete]
func (c *PointTypeController) Delete(ctx *gin.Context) {
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
