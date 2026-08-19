CREATE TABLE IF NOT EXISTS tasks (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    value FLOAT NOT NULL DEFAULT 0,
    number INTEGER NOT NULL,
    participant_id INTEGER NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (participant_id) REFERENCES participants(id) ON DELETE CASCADE
);
CREATE UNIQUE INDEX idx_tasks_common ON tasks (participant_id, number);
