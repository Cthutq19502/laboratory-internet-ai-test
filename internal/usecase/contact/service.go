package contact

import (
	"context"
	clientgigachat "laboratory-internet-ai-test/internal/clients/gigachat"
	domaincontact "laboratory-internet-ai-test/internal/domain/contact"
	"log/slog"
)

type Service struct {
	repo           Repository
	clientGigachat ClientGigachat

	logger *slog.Logger
}

func NewService(repo Repository, clientGigachat *clientgigachat.Client, logger *slog.Logger) *Service {
	return &Service{repo: repo, clientGigachat: clientGigachat, logger: logger}
}

func (s *Service) CreateContact(ctx context.Context, create Create) (domaincontact.Contact, error) {

	s.logger.InfoContext(ctx, "CreateContact", "data:", create)

	if err := create.Validate(); err != nil {
		return domaincontact.Contact{}, err
	}

	messages := make([]string, 0, 1)
	messages = append(messages, create.Comment)

	req := clientgigachat.Request{
		Messages: messages,
	}

	tonal, err := s.clientGigachat.GetTonal(ctx, &req)
	if err != nil {
		s.logger.Error("CreateContact contact [Get tonal gigachat error]", "error", err)
		tonal = domaincontact.TonalUnexpected
	}

	s.logger.InfoContext(ctx, "result", "response", tonal)

	contact, err := s.repo.CreateContact(ctx, domaincontact.Contact{
		Name:    create.Name,
		Phone:   create.Phone,
		Email:   create.Email,
		Comment: create.Comment,
		Tonal:   tonal,
	})

	if err != nil {
		s.logger.Error("CreateContact contact [CreateContact Contact Repository Error]", "error", err)
		return domaincontact.Contact{}, err
	}

	return contact, nil
}
