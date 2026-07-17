package metric

import (
	"github.com/gin-gonic/gin"
	usecasemetrics "laboratory-internet-ai-test/internal/usecase/metrics"
	"log"
	"time"
)

type Stats struct {
	usecase usecasemetrics.MetricUsecase
}

func NewMetric(usecase usecasemetrics.MetricUsecase) *Stats {
	return &Stats{usecase: usecase}
}

func (s *Stats) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		statusCode := c.Writer.Status()

		err := s.usecase.SaveMetrics(c, usecasemetrics.Metric{
			Now:        time.Now(),
			StatusCode: statusCode,
		})
		if err != nil {
			log.Println(err)
		}
		return
	}
}
