package repositories

import (
	"database/sql"
	"errors"
	"time"

	"github.com/markusylisiurunen/template-monorepo/backend/go-example/pkg/entities"
	"github.com/markusylisiurunen/template-monorepo/backend/go-example/pkg/logger"
)

var (
	ErrMessagesCouldNotCreate = errors.New("could not create a new message")
	ErrMessagesNotFound       = errors.New("message not found")
)

type databaseMessage struct {
	ID        int64
	CreatedAt time.Time
	Text      string
}

func (m *databaseMessage) toMessage() (entities.Message, error) {
	return entities.Message{
		ID:        m.ID,
		CreatedAt: m.CreatedAt,
		Text:      m.Text,
	}, nil
}

type MessagesRepository struct {
	db  *sql.DB
	log logger.Logger
}

func (r *MessagesRepository) Create(message entities.Message) (entities.Message, error) {
	query := `
    INSERT INTO messages (
      message_created_at,
      message_text
    )
    VALUES (now(), $1)
    RETURNING message_id
  `

	result, err := r.db.Query(query, message.Text)
	if err != nil {
		r.log.Errorw("could not insert a message",
			"Error", err.Error(),
		)

		return entities.Message{}, err
	}

	var insertedID int64

	if result.Next() {
		if err := result.Scan(&insertedID); err != nil {
			r.log.Errorf("could not find inserted message")
			return entities.Message{}, ErrMessagesCouldNotCreate
		}
	} else {
		r.log.Errorf("could not find inserted message")
		return entities.Message{}, ErrMessagesCouldNotCreate
	}

	return r.FindOneById(insertedID)
}

func (r *MessagesRepository) FindOneById(id int64) (entities.Message, error) {
	query := `
    SELECT
      message_id,
      message_created_at,
      message_text
    FROM messages
    WHERE
      message_id = $1
    LIMIT 1
  `

	message := &databaseMessage{}

	err := r.db.QueryRow(query, id).Scan(
		&message.ID,
		&message.CreatedAt,
		&message.Text,
	)

	if err != nil {
		r.log.Errorw("could not find a message for given ID",
			"ID", id,
			"Error", err.Error(),
		)

		return entities.Message{}, ErrMessagesNotFound
	}

	return message.toMessage()
}

func NewMessagesRepository(db *sql.DB, log logger.Logger) (*MessagesRepository, error) {
	return &MessagesRepository{
		db:  db,
		log: log,
	}, nil
}
