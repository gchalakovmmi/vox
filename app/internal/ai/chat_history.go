package ai

import (
	"context"
	"database/sql"
	"fmt"
	"vox/internal/config"
)

type ChatHistory struct {
	db		*sql.DB
	userID		uint
	ConversationID	int
}

func NewChatHistory(db *sql.DB, userID uint, conversationID int) *ChatHistory {
	return &ChatHistory{ db: db, userID: userID, ConversationID: conversationID }
}

func (h *ChatHistory) CreateConversation(ctx context.Context, cfg *config.Config) (int, error) {
	tx, err := h.db.BeginTx(ctx, nil)
	if err != nil {
		return -1, fmt.Errorf("chatHistory.CreateConversation begin: %w", err)
	}
	defer tx.Rollback()

	var convID int
	err = tx.QueryRowContext(ctx, `
		INSERT INTO conversations (user_id, title)
		VALUES ($1, 'New Conversation')
		RETURNING id
	`, h.userID).Scan(&convID)

	if err != nil {
		return -1, fmt.Errorf("chatHistory.CreateConversation ensure conv: %w", err)
	}

	// insert messages
	stmt, err := tx.PrepareContext(ctx, `
		INSERT INTO conversation_messages (conversation_id, role, content)
		VALUES ($1, $2, $3)
	`)
	
	if err != nil {
		return -1, err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, convID, "assistant", cfg.GetConversationDefaultGreeting())
	if err != nil {
		return -1, fmt.Errorf("chatHistory.CreateConversation insert: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return -1, err
	}
	return convID, nil
}

// Load returns the full JSON-encoded message array for the conversation with the provided id.
func (h *ChatHistory) Load(ctx context.Context) (string, error) {
	const q = `
		SELECT m.content, m.role
		FROM conversation_messages m
		JOIN conversations c ON c.id = m.conversation_id
		WHERE c.user_id = $1
			AND c.id = $2
		ORDER BY m.created_at ASC;
	`
	rows, err := h.db.QueryContext(ctx, q, h.userID, h.ConversationID)
	if err != nil {
		return "", fmt.Errorf("chatHistory.LoadLatest query: %w", err)
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
func (h *ChatHistory) Save(ctx context.Context, cfg *config.Config, raw string) error {
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
			AND id = $2
		LIMIT 1
	`, h.userID, h.ConversationID).Scan(&convID)
	if err == sql.ErrNoRows {
		_, err := h.CreateConversation(ctx, cfg)
		if err != nil {
			return fmt.Errorf("chatHistory.Save CreateConversation: %w", err)
		}
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
