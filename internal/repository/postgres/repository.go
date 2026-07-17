package postgres

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	domaincontact "laboratory-internet-ai-test/internal/domain/contact"
	"log/slog"
)

type Repository struct {
	pool   *pgxpool.Pool
	logger *slog.Logger
}

func New(pool *pgxpool.Pool, logger *slog.Logger) *Repository {
	return &Repository{pool: pool, logger: logger}
}

func (r *Repository) CreateContact(ctx context.Context, contact domaincontact.Contact) (domaincontact.Contact, error) {
	const query = `
		INSERT INTO contacts (name, phone, email, comment, tonal) 
		values ($1,$2,$3,$4,$5) returning id, name, phone, email, comment, date_create, tonal
	`

	row := r.pool.QueryRow(ctx, query, contact.Name, contact.Phone, contact.Email, contact.Comment, contact.Tonal)

	created, err := scanContact(row)
	if err != nil {
		r.logger.Error("Repository CreateContact", "error", err)
		return domaincontact.Contact{}, err
	}

	return created, nil
}

type taskScanner interface {
	Scan(dest ...any) error
}

func scanContact(scanner taskScanner) (domaincontact.Contact, error) {
	var (
		sub domaincontact.Contact
	)

	if err := scanner.Scan(
		&sub.ID,
		&sub.Name,
		&sub.Phone,
		&sub.Email,
		&sub.Comment,
		&sub.DateCreate,
		&sub.Tonal,
	); err != nil {
		return domaincontact.Contact{}, err
	}

	return sub, nil
}
