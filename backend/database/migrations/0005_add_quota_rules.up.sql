CREATE TABLE IF NOT EXISTS quota_rules (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    min_participants INTEGER NOT NULL,
    max_participants INTEGER,
    winners_quota INTEGER NOT NULL,
    winners_and_prizers_quota INTEGER NOT NULL,
    min_winners_points_percent INTEGER NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);
