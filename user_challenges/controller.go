package user_challenges

import (
	"net/http"
	"strconv"

	"base/core/storage"
	"base/packages/gamification/models"

	"github.com/gin-gonic/gin"
)

type UserChallengeController struct {
	Service *UserChallengeService
	Storage *storage.ActiveStorage
}

func NewUserChallengeController(service *UserChallengeService, storage *storage.ActiveStorage) *UserChallengeController {
	return &UserChallengeController{
		Service: service,
		Storage: storage,
	}
}

func (c *UserChallengeController) Routes(router *gin.RouterGroup) {
	// Main CRUD endpoints
	router.GET("/user-challenges", c.List)        // Paginated list
	router.GET("/user-challenges/all", c.ListAll) // Unpaginated list
	router.GET("/user-challenges/:id", c.Get)
	router.POST("/user-challenges", c.Create)
	router.PUT("/user-challenges/:id", c.Update)
	router.DELETE("/user-challenges/:id", c.Delete)

	// File/Image attachment endpoints

	// HasMany relation endpoints
}

// CreateUserChallenge godoc
// @Summary Create a new UserChallenge
// @Description Create a new UserChallenge with the input payload
// @Tags UserChallenge
// @Security ApiKeyAuth
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param user-challenges body models.CreateUserChallengeRequest true "Create UserChallenge request"
// @Success 201 {object} models.UserChallengeResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /user-challenges [post]
func (c *UserChallengeController) Create(ctx *gin.Context) {
	var req models.CreateUserChallengeRequest
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

// GetUserChallenge godoc
// @Summary Get a UserChallenge
// @Description Get a UserChallenge by its id
// @Tags UserChallenge
// @Security ApiKeyAuth
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "UserChallenge id"
// @Success 200 {object} models.UserChallengeResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /user-challenges/{id} [get]
func (c *UserChallengeController) Get(ctx *gin.Context) {
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

// ListUserChallenges godoc
// @Summary List user-challenges
// @Description Get a list of user-challenges
// @Tags UserChallenge
// @Security ApiKeyAuth
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param limit query int false "Number of items per page"
// @Success 200 {object} types.PaginatedResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /user-challenges [get]
func (c *UserChallengeController) List(ctx *gin.Context) {
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

// ListAllUserChallenges godoc
// @Summary List all user-challenges without pagination
// @Description Get a list of all user-challenges without pagination
// @Tags UserChallenge
// @Security ApiKeyAuth
// @Security BearerAuth
// @Accept json
// @Produce json
// @Success 200 {object} types.PaginatedResponse
// @Failure 500 {object} ErrorResponse
// @Router /user-challenges/all [get]
func (c *UserChallengeController) ListAll(ctx *gin.Context) {
	paginatedResponse, err := c.Service.GetAll(nil, nil)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to fetch all items: " + err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, paginatedResponse)
}

// UpdateUserChallenge godoc
// @Summary Update a UserChallenge
// @Description Update a UserChallenge by its id
// @Tags UserChallenge
// @Security ApiKeyAuth
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "UserChallenge id"
// @Param user-challenges body models.UpdateUserChallengeRequest true "Update UserChallenge request"
// @Success 200 {object} models.UserChallengeResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /user-challenges/{id} [put]
func (c *UserChallengeController) Update(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid id format"})
		return
	}

	var req models.UpdateUserChallengeRequest
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

// DeleteUserChallenge godoc
// @Summary Delete a UserChallenge
// @Description Delete a UserChallenge by its id
// @Tags UserChallenge
// @Security ApiKeyAuth
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "UserChallenge id"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /user-challenges/{id} [delete]
func (c *UserChallengeController) Delete(ctx *gin.Context) {
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
