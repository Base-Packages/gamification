package user_levels

import (
	"net/http"
	"strconv"

	"base/core/packages/gamification/models"
	"base/core/storage"

	"github.com/gin-gonic/gin"
)

type UserLevelController struct {
	Service *UserLevelService
	Storage *storage.ActiveStorage
}

func NewUserLevelController(service *UserLevelService, storage *storage.ActiveStorage) *UserLevelController {
	return &UserLevelController{
		Service: service,
		Storage: storage,
	}
}

func (c *UserLevelController) Routes(router *gin.RouterGroup) {
	// Main CRUD endpoints
	router.GET("/user-levels", c.List)        // Paginated list
	router.GET("/user-levels/all", c.ListAll) // Unpaginated list
	router.GET("/user-levels/:id", c.Get)
	router.POST("/user-levels", c.Create)
	router.PUT("/user-levels/:id", c.Update)
	router.DELETE("/user-levels/:id", c.Delete)

	// File/Image attachment endpoints

	// HasMany relation endpoints
}

// CreateUserLevel godoc
// @Summary Create a new UserLevel
// @Description Create a new UserLevel with the input payload
// @Tags UserLevel
// @Security ApiKeyAuth
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param user-levels body models.CreateUserLevelRequest true "Create UserLevel request"
// @Success 201 {object} models.UserLevelResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /user-levels [post]
func (c *UserLevelController) Create(ctx *gin.Context) {
	var req models.CreateUserLevelRequest
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

// GetUserLevel godoc
// @Summary Get a UserLevel
// @Description Get a UserLevel by its id
// @Tags UserLevel
// @Security ApiKeyAuth
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "UserLevel id"
// @Success 200 {object} models.UserLevelResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /user-levels/{id} [get]
func (c *UserLevelController) Get(ctx *gin.Context) {
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

// ListUserLevels godoc
// @Summary List user-levels
// @Description Get a list of user-levels
// @Tags UserLevel
// @Security ApiKeyAuth
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param limit query int false "Number of items per page"
// @Success 200 {object} types.PaginatedResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /user-levels [get]
func (c *UserLevelController) List(ctx *gin.Context) {
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

// ListAllUserLevels godoc
// @Summary List all user-levels without pagination
// @Description Get a list of all user-levels without pagination
// @Tags UserLevel
// @Security ApiKeyAuth
// @Security BearerAuth
// @Accept json
// @Produce json
// @Success 200 {object} types.PaginatedResponse
// @Failure 500 {object} ErrorResponse
// @Router /user-levels/all [get]
func (c *UserLevelController) ListAll(ctx *gin.Context) {
	paginatedResponse, err := c.Service.GetAll(nil, nil)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to fetch all items: " + err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, paginatedResponse)
}

// UpdateUserLevel godoc
// @Summary Update a UserLevel
// @Description Update a UserLevel by its id
// @Tags UserLevel
// @Security ApiKeyAuth
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "UserLevel id"
// @Param user-levels body models.UpdateUserLevelRequest true "Update UserLevel request"
// @Success 200 {object} models.UserLevelResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /user-levels/{id} [put]
func (c *UserLevelController) Update(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid id format"})
		return
	}

	var req models.UpdateUserLevelRequest
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

// DeleteUserLevel godoc
// @Summary Delete a UserLevel
// @Description Delete a UserLevel by its id
// @Tags UserLevel
// @Security ApiKeyAuth
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "UserLevel id"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /user-levels/{id} [delete]
func (c *UserLevelController) Delete(ctx *gin.Context) {
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
