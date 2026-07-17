package http

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	handlercontact "laboratory-internet-ai-test/internal/transport/http/contact_handler"
	_ "laboratory-internet-ai-test/internal/transport/http/docs"
	handlermetric "laboratory-internet-ai-test/internal/transport/http/metric_handler"
)

func NewContactsController(router *gin.Engine, handler *handlercontact.Handler) {
	subGroup := router.Group("/api/v1")
	{
		subGroup.POST("/contact", handler.CreateContact)
	}

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}

func NewMetricController(router *gin.Engine, handler *handlermetric.Handler) {

	subGroup := router.Group("/api/v1")
	{
		subGroup.GET("/metric", handler.GetMetric)
	}

}
