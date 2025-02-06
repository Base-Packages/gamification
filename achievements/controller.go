package achievements

import (
	"net/http"
	"strconv"

	"base/core/packages/gamification/models"
	"base/core/storage"

	"github.com/gin-gonic/gin"
)

type AchievementController struct {
	Service *AchievementService
	Storage *storage.ActiveStorage
}

func NewAchievementController(service *AchievementService, storage *storage.ActiveStorage) *AchievementController {
	return &AchievementController{
		Service: service,
		Storage: storage,
	}
}

func (c *AchievementController) Routes(router *gin.RouterGroup) {
	// Main CRUD endpoints
	router.GET("/achievements", c.List)        // Paginated list
	router.GET("/achievements/all", c.ListAll) // Unpaginated list
	router.GET("/achievements/:id", c.Get)
	router.POST("/achievements", c.Create)
	router.PUT("/achievements/:id", c.Update)
	router.DELETE("/achievements/:id", c.Delete)

	// File/Image attachment endpoints
	router.PUT("/achievements/:id/icon", c.UploadIcon)
	router.DELETE("/achievements/:id/icon", c.DeleteIcon)

	// HasMany relation endpoints
}

// CreateAchievement godoc
// @Summary Create a new Achievement
// @Description Create a new Achievement with the input payload
// @Tags Achievement
// @Security ApiKeyAuth
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param achievements body models.CreateAchievementRequest true "Create Achievement request"
// @Success 201 {object} models.AchievementResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /achievements [post]
func (c *AchievementController) Create(ctx *gin.Context) {
	var req models.CreateAchievementRequest
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

// GetAchievement godoc
// @Summary Get a Achievement
// @Description Get a Achievement by its id
// @Tags Achievement
// @Security ApiKeyAuth
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Achievement id"
// @Success 200 {object} models.AchievementResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /achievements/{id} [get]
func (c *AchievementController) Get(ctx *gin.Context) {
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

// ListAchievements godoc
// @Summary List achievements
// @Description Get a list of achievements
// @Tags Achievement
// @Security ApiKeyAuth
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param limit query int false "Number of items per page"
// @Success 200 {object} types.PaginatedResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /achievements [get]
func (c *AchievementController) List(ctx *gin.Context) {
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

// ListAllAchievements godoc
// @Summary List all achievements without pagination
// @Description Get a list of all achievements without pagination
// @Tags Achievement
// @Security ApiKeyAuth
// @Security BearerAuth
// @Accept json
// @Produce json
// @Success 200 {object} types.PaginatedResponse
// @Failure 500 {object} ErrorResponse
// @Router /achievements/all [get]
func (c *AchievementController) ListAll(ctx *gin.Context) {
	paginatedResponse, err := c.Service.GetAll(nil, nil)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to fetch all items: " + err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, paginatedResponse)
}

// UpdateAchievement godoc
// @Summary Update a Achievement
// @Description Update a Achievement by its id
// @Tags Achievement
// @Security ApiKeyAuth
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Achievement id"
// @Param achievements body models.UpdateAchievementRequest true "Update Achievement request"
// @Success 200 {object} models.AchievementResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /achievements/{id} [put]
func (c *AchievementController) Update(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid id format"})
		return
	}

	var req models.UpdateAchievementRequest
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

// DeleteAchievement godoc
// @Summary Delete a Achievement
// @Description Delete a Achievement by its id
// @Tags Achievement
// @Security ApiKeyAuth
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Achievement id"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /achievements/{id} [delete]
func (c *AchievementController) Delete(ctx *gin.Context) {
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
// @Summary Upload Icon for a Achievement
// @Description Upload or update the Icon of a Achievement
// @Tags Achievement
// @Security ApiKeyAuth
// @Security BearerAuth
// @Accept multipart/form-data
// @Produce json
// @Param id path int true "Achievement id"
// @Param file formData file true "File to upload"
// @Success 200 {object} models.AchievementResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /achievements/{id}/icon [put]
func (c *AchievementController) UploadIcon(ctx *gin.Context) {
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
// @Summary Delete Icon from a Achievement
// @Description Delete the Icon of a Achievement
// @Tags Achievement
// @Security ApiKeyAuth
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Achievement id"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /achievements/{id}/icon [delete]
func (c *AchievementController) DeleteIcon(ctx *gin.Context) {
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
	updateReq := &models.UpdateAchievementRequest{
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
