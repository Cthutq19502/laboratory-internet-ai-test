package contact_handler

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	domaincontact "laboratory-internet-ai-test/internal/domain/contact"
	"laboratory-internet-ai-test/internal/pkg/apperrors"
	usecasecontact "laboratory-internet-ai-test/internal/usecase/contact"
	"log/slog"
	"net/http"
)

type Handler struct {
	usecase usecasecontact.Usecase
	logger  *slog.Logger
}

func New(usecase usecasecontact.Usecase, logger *slog.Logger) *Handler {
	return &Handler{
		usecase: usecase,
		logger:  logger,
	}
}

// CreateContact Создать контакт
// @Summary      Создать контакт
// @Description  Создать контакт
// @Tags         createContact
// @Accept       json
// @Produce      json
// @Param        request body createContactDTO true "Создание контакта"
// @Success      201 {object} contactDTO
// @Failure      400 {object} ErrorDTO
// @Failure      500 {object} ErrorDTO
// @Router       /contact [post]
func (h *Handler) CreateContact(ctx *gin.Context) {

	var create createContactDTO

	if err := ctx.ShouldBindJSON(&create); err != nil {
		h.logger.Error("CreateContact", "error", err)
		WriteContactError(ctx, fmt.Errorf("%w: %s", domaincontact.ErrInvalidInput, "wrong json"))
		return
	}

	contact, err := h.usecase.CreateContact(ctx, usecasecontact.Create{
		Name:    create.Name,
		Phone:   create.Phone,
		Email:   create.Email,
		Comment: create.Comment,
	})
	if err != nil {
		WriteContactError(ctx, err)
		return
	}

	ctx.JSON(http.StatusCreated, newContactDTO(contact))

}

func WriteContactError(ctx *gin.Context, err error) {
	switch {
	case errors.Is(err, domaincontact.ErrInvalidInput):
		ctx.JSON(http.StatusBadRequest, newErrorDTO(err))
	default:
		ctx.JSON(http.StatusInternalServerError, newErrorDTO(apperrors.ErrServerCritical))
	}
}
