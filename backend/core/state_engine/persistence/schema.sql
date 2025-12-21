-- Create conversation_state table
CREATE TABLE IF NOT EXISTS conversation_state (
    email TEXT PRIMARY KEY,
    state TEXT NOT NULL
);