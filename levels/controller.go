package levels

import (
	"net/http"
	"strconv"

	"base/core/packages/gamification/models"
	"base/core/storage"

	"github.com/gin-gonic/gin"
)

type LevelController struct {
	Service *LevelService
	Storage *storage.ActiveStorage
}

func NewLevelController(service *LevelService, storage *storage.ActiveStorage) *LevelController {
	return &LevelController{
		Service: service,
		Storage: storage,
	}
}

func (c *LevelController) Routes(router *gin.RouterGroup) {
	// Main CRUD endpoints
	router.GET("/levels", c.List)        // Paginated list
	router.GET("/levels/all", c.ListAll) // Unpaginated list
	router.GET("/levels/:id", c.Get)
	router.POST("/levels", c.Create)
	router.PUT("/levels/:id", c.Update)
	router.DELETE("/levels/:id", c.Delete)

	// File/Image attachment endpoints
	router.PUT("/levels/:id/icon", c.UploadIcon)
	router.DELETE("/levels/:id/icon", c.DeleteIcon)

	// HasMany relation endpoints
}

// CreateLevel godoc
// @Summary Create a new Level
// @Description Create a new Level with the input payload
// @Tags Level
// @Security ApiKeyAuth
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param levels body models.CreateLevelRequest true "Create Level request"
// @Success 201 {object} models.LevelResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /levels [post]
func (c *LevelController) Create(ctx *gin.Context) {
	var req models.CreateLevelRequest
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

// GetLevel godoc
// @Summary Get a Level
// @Description Get a Level by its id
// @Tags Level
// @Security ApiKeyAuth
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Level id"
// @Success 200 {object} models.LevelResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /levels/{id} [get]
func (c *LevelController) Get(ctx *gin.Context) {
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

// ListLevels godoc
// @Summary List levels
// @Description Get a list of levels
// @Tags Level
// @Security ApiKeyAuth
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param limit query int false "Number of items per page"
// @Success 200 {object} types.PaginatedResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /levels [get]
func (c *LevelController) List(ctx *gin.Context) {
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

// ListAllLevels godoc
// @Summary List all levels without pagination
// @Description Get a list of all levels without pagination
// @Tags Level
// @Security ApiKeyAuth
// @Security BearerAuth
// @Accept json
// @Produce json
// @Success 200 {object} types.PaginatedResponse
// @Failure 500 {object} ErrorResponse
// @Router /levels/all [get]
func (c *LevelController) ListAll(ctx *gin.Context) {
	paginatedResponse, err := c.Service.GetAll(nil, nil)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to fetch all items: " + err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, paginatedResponse)
}

// UpdateLevel godoc
// @Summary Update a Level
// @Description Update a Level by its id
// @Tags Level
// @Security ApiKeyAuth
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Level id"
// @Param levels body models.UpdateLevelRequest true "Update Level request"
// @Success 200 {object} models.LevelResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /levels/{id} [put]
func (c *LevelController) Update(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid id format"})
		return
	}

	var req models.UpdateLevelRequest
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

// DeleteLevel godoc
// @Summary Delete a Level
// @Description Delete a Level by its id
// @Tags Level
// @Security ApiKeyAuth
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Level id"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /levels/{id} [delete]
func (c *LevelController) Delete(ctx *gin.Context) {
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

// UploadIcon godoc
// @Summary Upload Icon for a Level
// @Description Upload or update the Icon of a Level
// @Tags Level
// @Security ApiKeyAuth
// @Security BearerAuth
// @Accept multipart/form-data
// @Produce json
// @Param id path int true "Level id"
// @Param file formData file true "File to upload"
// @Success 200 {object} models.LevelResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /levels/{id}/icon [put]
func (c *LevelController) UploadIcon(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid id format"})
		return
	}

	file, err := ctx.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{Error: "No file uploaded"})
		return
	}

	// Get the item first
	item, err := c.Service.GetById(uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, ErrorResponse{Error: "Item not found"})
		return
	}

	// Upload the file using storage service
	_, err = c.Storage.Attach(item, "icon", file)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to upload file: " + err.Error()})
		return
	}

	// Update the item with the new attachment
	updatedItem, err := c.Service.UploadIcon(uint(id), file)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to update item: " + err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, updatedItem.ToResponse())
}

// DeleteIcon godoc
// @Summary Delete Icon from a Level
// @Description Delete the Icon of a Level
// @Tags Level
// @Security ApiKeyAuth
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Level id"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /levels/{id}/icon [delete]
func (c *LevelController) DeleteIcon(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid id format"})
		return
	}

	// Get the item first
	item, err := c.Service.GetById(uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, ErrorResponse{Error: "Item not found"})
		return
	}

	// Delete the file using storage service
	if err := c.Storage.Delete(item.Icon); err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to delete file: " + err.Error()})
		return
	}

	// Update the item to remove the attachment reference
	updateReq := &models.UpdateLevelRequest{
		Icon: nil,
	}

	_, err = c.Service.Update(uint(id), updateReq)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to update item: " + err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, SuccessResponse{Message: "File deleted successfully"})
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type SuccessResponse struct {
	Message string `json:"message"`
}
