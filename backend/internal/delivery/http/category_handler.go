package http

import (
	"github.com/aswinsreeraj/evntx/internal/domain"
	"github.com/aswinsreeraj/evntx/internal/usecase"
	pkgErrors "github.com/aswinsreeraj/evntx/pkg/errors"
	"github.com/aswinsreeraj/evntx/pkg/response"
	"github.com/gin-gonic/gin"
)

type CategoryHandler struct {
	categoryUsecase *usecase.CategoryUsecase
	auditUsecase    *usecase.AuditUsecase
}

func NewCategoryHandler(categoryUsecase *usecase.CategoryUsecase, auditUsecase *usecase.AuditUsecase) *CategoryHandler {
	return &CategoryHandler{
		categoryUsecase: categoryUsecase,
		auditUsecase:    auditUsecase,
	}
}

func (h *CategoryHandler) ListCategories(c *gin.Context) {
	categories, err := h.categoryUsecase.GetAllCategories()
	if err != nil {
		response.AppError(c, pkgErrors.ErrInternalServerError)
		return
	}

	response.Success(c, "Categories retrieved successfully", categories)
}

func (h *CategoryHandler) CreateCategory(c *gin.Context) {
	var req struct {
		Name string `json:"name" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.AppError(c, pkgErrors.ErrInvalidRequestBody)
		return
	}

	category, err := h.categoryUsecase.CreateCategory(req.Name)
	if err != nil {
		response.AppError(c, pkgErrors.New(400, "BAD_REQUEST", err.Error()))
		return
	}

	if h.auditUsecase != nil {
		adminID := c.GetString("user_id")
		go h.auditUsecase.LogAction(adminID, "Event category '"+req.Name+"' created", domain.ActionTagSettings, map[string]interface{}{"category_id": category.ID}, c.ClientIP())
	}

	response.Success(c, "Category created successfully", category)
}

func (h *CategoryHandler) UpdateCategory(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		Name string `json:"name" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.AppError(c, pkgErrors.ErrInvalidRequestBody)
		return
	}

	category, err := h.categoryUsecase.UpdateCategory(id, req.Name)
	if err != nil {
		response.AppError(c, pkgErrors.New(400, "BAD_REQUEST", err.Error()))
		return
	}

	if h.auditUsecase != nil {
		adminID := c.GetString("user_id")
		go h.auditUsecase.LogAction(adminID, "Event category updated", domain.ActionTagSettings, map[string]interface{}{"category_id": category.ID, "new_name": req.Name}, c.ClientIP())
	}

	response.Success(c, "Category updated successfully", category)
}

func (h *CategoryHandler) DeleteCategory(c *gin.Context) {
	id := c.Param("id")

	err := h.categoryUsecase.DeleteCategory(id)
	if err != nil {
		response.AppError(c, pkgErrors.New(400, "BAD_REQUEST", err.Error()))
		return
	}

	if h.auditUsecase != nil {
		adminID := c.GetString("user_id")
		go h.auditUsecase.LogAction(adminID, "Event category deleted", domain.ActionTagSettings, map[string]interface{}{"category_id": id}, c.ClientIP())
	}

	response.Success(c, "Category deleted successfully", nil)
}
