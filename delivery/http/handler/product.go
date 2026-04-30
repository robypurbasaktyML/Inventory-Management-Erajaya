package handler

import (
	"errors"
	"net/http"

	"Inventory-Management-Erajaya/domain"

	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	usecase domain.ProductUsecase
}

func NewProductHandler(uc domain.ProductUsecase) *ProductHandler {
	return &ProductHandler{usecase: uc}
}

func (h *ProductHandler) RegisterRoutes(r *gin.Engine) {
	products := r.Group("/products")
	{
		products.POST("", h.Create)
		products.GET("", h.GetAll)
		products.GET("/:sku", h.GetBySKU)
		products.PUT("/:sku", h.Update)
		products.DELETE("/:sku", h.Delete)
		products.POST("/:sku/purchase", h.Purchase)
	}
}

type createProductRequest struct {
	SKU   string  `json:"sku" binding:"required"`
	Name  string  `json:"name" binding:"required"`
	Qty   int     `json:"qty"`
	Price float64 `json:"price" binding:"required"`
}

func (h *ProductHandler) Create(c *gin.Context) {
	var req createProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse(err.Error()))
		return
	}

	product := &domain.Product{
		SKU:   req.SKU,
		Name:  req.Name,
		Qty:   req.Qty,
		Price: req.Price,
	}

	if err := h.usecase.Create(c.Request.Context(), product); err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, domain.SuccessResponse("data created", product))
}

func (h *ProductHandler) GetAll(c *gin.Context) {
	products, err := h.usecase.GetAll(c.Request.Context())
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, domain.SuccessResponse("data retrieved", products))
}

func (h *ProductHandler) GetBySKU(c *gin.Context) {
	sku := c.Param("sku")
	product, err := h.usecase.GetBySKU(c.Request.Context(), sku)
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, domain.SuccessResponse("data retrieved", product))
}

type updateProductRequest struct {
	Name  string  `json:"name" binding:"required"`
	Qty   int     `json:"qty"`
	Price float64 `json:"price" binding:"required"`
}

func (h *ProductHandler) Update(c *gin.Context) {
	sku := c.Param("sku")
	var req updateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse(err.Error()))
		return
	}

	product := &domain.Product{
		Name:  req.Name,
		Qty:   req.Qty,
		Price: req.Price,
	}

	updated, err := h.usecase.Update(c.Request.Context(), sku, product)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, domain.SuccessResponse("data updated", updated))
}

func (h *ProductHandler) Delete(c *gin.Context) {
	sku := c.Param("sku")
	if err := h.usecase.Delete(c.Request.Context(), sku); err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, domain.SuccessResponse("data deleted", nil))
}

type purchaseRequest struct {
	Qty int `json:"qty" binding:"required"`
}

func (h *ProductHandler) Purchase(c *gin.Context) {
	sku := c.Param("sku")
	var req purchaseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse(err.Error()))
		return
	}

	if err := h.usecase.Purchase(c.Request.Context(), sku, req.Qty); err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, domain.SuccessResponse("purchase successful", nil))
}

func handleError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, domain.ErrProductNotFound):
		c.JSON(http.StatusNotFound, domain.ErrorResponse(err.Error()))
	case errors.Is(err, domain.ErrDuplicateSKU):
		c.JSON(http.StatusConflict, domain.ErrorResponse(err.Error()))
	case errors.Is(err, domain.ErrInsufficientStock):
		c.JSON(http.StatusBadRequest, domain.ErrorResponse(err.Error()))
	case errors.Is(err, domain.ErrInvalidQty):
		c.JSON(http.StatusBadRequest, domain.ErrorResponse(err.Error()))
	case errors.Is(err, domain.ErrInvalidPrice):
		c.JSON(http.StatusBadRequest, domain.ErrorResponse(err.Error()))
	default:
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse("internal server error"))
	}
}
