package rate_limiter

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/gin-gonic/gin"
	"laboratory-internet-ai-test/internal/pkg/apperrors"
	usecasemetrics "laboratory-internet-ai-test/internal/usecase/metrics"
	"time"
)

func (rl *RateLimiter) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientIP := c.ClientIP()

		hashIp := rl.hashIP(clientIP)

		timeRequest := time.Now()
		isAllow, err := rl.usecase.RateLimiterAllow(c, hashIp, timeRequest)
		if err != nil {
			c.AbortWithStatusJSON(500, apperrors.ErrServerCritical)
			return
		}

		if !isAllow {
			c.AbortWithStatusJSON(429, gin.H{
				"error": "Too many requests",
			})
			return
		}

		c.Next()
	}
}

type RateLimiter struct {
	usecase usecasemetrics.RateLimiterUsecase
}

func NewRateLimiter(usecase usecasemetrics.RateLimiterUsecase) *RateLimiter {
	rl := &RateLimiter{
		usecase: usecase,
	}

	return rl
}

func (rl *RateLimiter) hashIP(ip string) string {
	hash := md5.Sum([]byte(ip))
	return hex.EncodeToString(hash[:])
}
