package contact_handler

import (
	domaincontact "laboratory-internet-ai-test/internal/domain/contact"
	"time"
)

type ErrorDTO struct {
	Error interface{} `json:"error"`
}

func newErrorDTO(err error) ErrorDTO {
	return ErrorDTO{Error: err.Error()}
}

//-----------------------------------

type createContactDTO struct {
	Name    string `json:"name"`
	Phone   string `json:"phone"`
	Email   string `json:"email"`
	Comment string `json:"comment"`
}

type contactDTO struct {
	Name       string    `json:"name"`
	Phone      string    `json:"phone"`
	Email      string    `json:"email"`
	Comment    string    `json:"comment"`
	Tonal      string    `json:"tonal"`
	DateCreate time.Time `json:"date_create"`
}

func newContactDTO(contact domaincontact.Contact) contactDTO {
	return contactDTO{
		Name:       contact.Name,
		Phone:      contact.Phone,
		Email:      contact.Email,
		Comment:    contact.Comment,
		Tonal:      string(contact.Tonal),
		DateCreate: contact.DateCreate,
	}
}
