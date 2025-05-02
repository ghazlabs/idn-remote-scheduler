CREATE TABLE messages (
    id VARCHAR(36) PRIMARY KEY,
    message TEXT NOT NULL,
    scheduled_sending_at DATETIME,
    sent_at DATETIME,
    retried_count INT DEFAULT 0,
    status VARCHAR(50),
    reason TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);