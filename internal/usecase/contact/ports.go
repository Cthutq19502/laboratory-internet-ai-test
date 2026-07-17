package contact

import (
	"context"
	"fmt"
	"laboratory-internet-ai-test/internal/clients/gigachat"
	domaincontact "laboratory-internet-ai-test/internal/domain/contact"
	"laboratory-internet-ai-test/internal/pkg/utils"
	"strings"
	"unicode/utf8"
)

type Repository interface {
	CreateContact(ctx context.Context, contact domaincontact.Contact) (domaincontact.Contact, error)
}

type ClientGigachat interface {
	GetTonal(ctx context.Context, request *gigachat.Request) (domaincontact.Tonal, error)
}

type Usecase interface {
	CreateContact(ctx context.Context, create Create) (domaincontact.Contact, error)
}

type Create struct {
	Name    string `json:"name"`
	Phone   string `json:"phone"`
	Email   string `json:"email"`
	Comment string `json:"comment"`
}

func (c *Create) Validate() error {
	mass := make([]string, 0)
	if utf8.RuneCountInString(c.Name) < 3 {
		mass = append(mass, "name must collect 3 symbols min")
	}

	if utf8.RuneCountInString(c.Phone) != 11 {
		mass = append(mass, "phone must collect 11 symbols")
	}

	if !utils.IsValidEmail(c.Email) {
		mass = append(mass, "email is not valid")
	}

	if utf8.RuneCountInString(c.Comment) < 10 {
		mass = append(mass, "comment must collect 10 symbols min")
	}

	//---------------------

	if len(mass) != 0 {
		return fmt.Errorf("%w: %s", domaincontact.ErrInvalidInput, strings.Join(mass, ",\n"))
	}

	return nil
}
