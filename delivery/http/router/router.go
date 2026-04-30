package router

import (
	"Inventory-Management-Erajaya/delivery/http/handler"
	"Inventory-Management-Erajaya/delivery/http/middleware"

	"github.com/gin-gonic/gin"
)

func NewRouter(productHandler *handler.ProductHandler) *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.LoggerMiddleware())
	productHandler.RegisterRoutes(r)
	return r
}
