CREATE TABLE messages (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    message TEXT NOT NULL,
    scheduled_sending_at BIGINT,
    sent_at BIGINT,             
    retried_count INT DEFAULT 0,
    status VARCHAR(50),
    reason TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
