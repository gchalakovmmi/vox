package ai

import (
	"context"
	"database/sql"
	"fmt"
)

type ChatHistory struct {
	db	 *sql.DB
	userID uint
}

func NewChatHistory(db *sql.DB, userID uint) *ChatHistory {
	return &ChatHistory{db: db, userID: userID}
}

// Load returns the full JSON-encoded message array for the *latest* conversation.
func (h *ChatHistory) Load(ctx context.Context) (string, error) {
	const q = `
		SELECT m.content, m.role
		FROM conversation_messages m
		JOIN conversations c ON c.id = m.conversation_id
		WHERE c.user_id = $1
		ORDER BY c.created_at DESC, m.created_at ASC
		LIMIT 1
	`
	rows, err := h.db.QueryContext(ctx, q, h.userID)
	if err != nil {
		return "", fmt.Errorf("chatHistory.Load query: %w", err)
	}
	defer rows.Close()

	var msgs []ChatMessage
	for rows.Next() {
		var content string
		var role string
		if err := rows.Scan(&content, &role); err != nil {
			return "", err
		}
		msgs = append(msgs, ChatMessage{Role: role, Content: content})
	}
	raw, _ := MessagesToJSON(msgs)
	return raw, rows.Err()
}

// Save appends messages to the *latest* conversation; creates one if none exists.
func (h *ChatHistory) Save(ctx context.Context, raw string) error {
	tx, err := h.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("chatHistory.Save begin: %w", err)
	}
	defer tx.Rollback()

	// ensure a conversation exists
	var convID int64
	err = tx.QueryRowContext(ctx, `
		SELECT id FROM conversations
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT 1
	`, h.userID).Scan(&convID)
	if err == sql.ErrNoRows {
		// create new conversation
		err = tx.QueryRowContext(ctx, `
			INSERT INTO conversations (user_id, title)
			VALUES ($1, '')
			RETURNING id
		`, h.userID).Scan(&convID)
	}
	if err != nil {
		return fmt.Errorf("chatHistory.Save ensure conv: %w", err)
	}

	// insert messages
	msgs, err := MessagesFromJSON(raw)
	if err != nil {
		return fmt.Errorf("chatHistory.Save parse: %w", err)
	}
	stmt, err := tx.PrepareContext(ctx, `
		INSERT INTO conversation_messages (conversation_id, role, content)
		VALUES ($1, $2, $3)
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, m := range msgs {
		if _, err := stmt.ExecContext(ctx, convID, m.Role, m.Content); err != nil {
			return fmt.Errorf("chatHistory.Save insert: %w", err)
		}
	}
	return tx.Commit()
}
