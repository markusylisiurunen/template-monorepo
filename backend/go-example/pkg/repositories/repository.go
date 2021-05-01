package repositories

import (
	"database/sql"

	"github.com/markusylisiurunen/template-monorepo/backend/go-example/pkg/logger"
)

type Repository interface {
	// TODO: ...
}

type Repositories struct {
	Messages *MessagesRepository
}

func NewRepositories(db *sql.DB, log logger.Logger) (*Repositories, error) {
	messages, err := NewMessagesRepository(db, log)
	if err != nil {
		return nil, err
	}

	return &Repositories{
		Messages: messages,
	}, nil
}
