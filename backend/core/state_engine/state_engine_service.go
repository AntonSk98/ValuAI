package core

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"valuai/core/state_engine/persistence"
)

// StateEngineService provides high-level operations for managing conversation state
type StateEngineService struct {
	queries *persistence.Queries
}

// NewStateEngineService creates a new state engine service with the provided database connection
func NewStateEngineService(db *sql.DB) *StateEngineService {
	// Create queries instance with the provided database connection
	// *sql.DB implements the DBTX interface
	queries := persistence.New(db)

	return &StateEngineService{
		queries: queries,
	}
}

// GetConversationState retrieves the current state for a given email
func (s *StateEngineService) GetConversationState(ctx context.Context, email string) (*string, error) {
	state, err := s.queries.FindStateByEmail(ctx, email)
	if err != nil {
		// If no state found, return empty string (this is not an error for new conversations)
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get conversation state for email %s: %w", email, err)
	}

	return &state.State, nil
}

// UpdateConversationState updates or creates the conversation state for a given email
func (s *StateEngineService) UpdateConversationState(ctx context.Context, email, state string) error {
	params := persistence.UpdateStateParams{
		Email: email,
		State: state,
	}

	err := s.queries.UpsertState(ctx, params)
	if err != nil {
		return fmt.Errorf("failed to update conversation state for email %s: %w", email, err)
	}

	return nil
}
