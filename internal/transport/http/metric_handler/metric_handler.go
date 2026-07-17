package metric_handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	domaincontact "laboratory-internet-ai-test/internal/domain/contact"
	"laboratory-internet-ai-test/internal/pkg/apperrors"
	usecasemetrics "laboratory-internet-ai-test/internal/usecase/metrics"
	"log/slog"
	"net/http"
)

type Handler struct {
	usecase usecasemetrics.MetricUsecase
	logger  *slog.Logger
}

func New(usecase usecasemetrics.MetricUsecase, logger *slog.Logger) *Handler {
	return &Handler{
		usecase: usecase,
		logger:  logger,
	}
}

// CreateContact Запрос метрики
// @Summary      Запрос метрики
// @Description  Запрос метрики
// @Tags         GetMetric
// @Produce      json
// @Success      200 {object} metricDTO
// @Failure      400 {object} ErrorDTO
// @Failure      500 {object} ErrorDTO
// @Router       /metric [get]
func (h *Handler) GetMetric(ctx *gin.Context) {

	metric, err := h.usecase.GetMetric(ctx)
	if err != nil {
		WriteMetricError(ctx, err)
		return
	}

	ctx.JSON(200, newMetricDTO(metric))
}

func WriteMetricError(ctx *gin.Context, err error) {
	switch {
	case errors.Is(err, domaincontact.ErrInvalidInput):
		ctx.JSON(http.StatusBadRequest, newErrorDTO(err))
	default:
		ctx.JSON(http.StatusInternalServerError, newErrorDTO(apperrors.ErrServerCritical))
	}
}
