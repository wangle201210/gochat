package storage

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/wangle201210/gochat/internal/models"
)

// Database SQLite 数据库
type Database struct {
	db *sql.DB
}

// NewDatabase 创建数据库连接
func NewDatabase(dbPath string) (*Database, error) {
	// 确保数据库目录存在
	dbDir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		return nil, fmt.Errorf("创建数据库目录失败: %w", err)
	}

	// 打开数据库连接
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("打开数据库失败: %w", err)
	}

	// 测试连接
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("数据库连接测试失败: %w", err)
	}

	database := &Database{db: db}

	// 初始化表结构
	if err := database.initTables(); err != nil {
		db.Close()
		return nil, err
	}

	return database, nil
}

// initTables 初始化数据库表
func (d *Database) initTables() error {
	// 创建会话表
	createSessionsTable := `
	CREATE TABLE IF NOT EXISTS sessions (
		id TEXT PRIMARY KEY,
		title TEXT NOT NULL,
		created_at DATETIME NOT NULL,
		updated_at DATETIME NOT NULL
	);
	`

	// 创建消息表
	createMessagesTable := `
	CREATE TABLE IF NOT EXISTS messages (
		id TEXT PRIMARY KEY,
		session_id TEXT NOT NULL,
		role TEXT NOT NULL,
		content TEXT NOT NULL,
		timestamp DATETIME NOT NULL,
		FOREIGN KEY (session_id) REFERENCES sessions(id) ON DELETE CASCADE
	);
	`

	// 创建索引
	createIndexes := `
	CREATE INDEX IF NOT EXISTS idx_messages_session_id ON messages(session_id);
	CREATE INDEX IF NOT EXISTS idx_sessions_updated_at ON sessions(updated_at DESC);
	`

	// 执行建表语句
	if _, err := d.db.Exec(createSessionsTable); err != nil {
		return fmt.Errorf("创建会话表失败: %w", err)
	}

	if _, err := d.db.Exec(createMessagesTable); err != nil {
		return fmt.Errorf("创建消息表失败: %w", err)
	}

	if _, err := d.db.Exec(createIndexes); err != nil {
		return fmt.Errorf("创建索引失败: %w", err)
	}

	return nil
}

// Close 关闭数据库连接
func (d *Database) Close() error {
	return d.db.Close()
}

// SaveSession 保存会话
func (d *Database) SaveSession(session *models.Session) error {
	query := `
	INSERT OR REPLACE INTO sessions (id, title, created_at, updated_at)
	VALUES (?, ?, ?, ?)
	`

	_, err := d.db.Exec(query, session.ID, session.Title, session.CreatedAt, session.UpdatedAt)
	if err != nil {
		return fmt.Errorf("保存会话失败: %w", err)
	}

	return nil
}

// GetSession 获取会话
func (d *Database) GetSession(sessionID string) (*models.Session, error) {
	query := `
	SELECT id, title, created_at, updated_at
	FROM sessions
	WHERE id = ?
	`

	session := &models.Session{}
	err := d.db.QueryRow(query, sessionID).Scan(
		&session.ID,
		&session.Title,
		&session.CreatedAt,
		&session.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, fmt.Errorf("获取会话失败: %w", err)
	}

	return session, nil
}

// ListSessions 获取所有会话列表（按更新时间倒序）
func (d *Database) ListSessions() ([]*models.Session, error) {
	query := `
	SELECT id, title, created_at, updated_at
	FROM sessions
	ORDER BY updated_at DESC
	`

	rows, err := d.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("查询会话列表失败: %w", err)
	}
	defer rows.Close()

	sessions := make([]*models.Session, 0)
	for rows.Next() {
		session := &models.Session{}
		if err := rows.Scan(&session.ID, &session.Title, &session.CreatedAt, &session.UpdatedAt); err != nil {
			return nil, fmt.Errorf("读取会话数据失败: %w", err)
		}
		sessions = append(sessions, session)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("遍历会话列表失败: %w", err)
	}

	return sessions, nil
}

// DeleteSession 删除会话（会级联删除相关消息）
func (d *Database) DeleteSession(sessionID string) error {
	query := `DELETE FROM sessions WHERE id = ?`

	_, err := d.db.Exec(query, sessionID)
	if err != nil {
		return fmt.Errorf("删除会话失败: %w", err)
	}

	return nil
}

// UpdateSessionTitle 更新会话标题
func (d *Database) UpdateSessionTitle(sessionID, title string) error {
	query := `
	UPDATE sessions
	SET title = ?, updated_at = ?
	WHERE id = ?
	`

	_, err := d.db.Exec(query, title, time.Now(), sessionID)
	if err != nil {
		return fmt.Errorf("更新会话标题失败: %w", err)
	}

	return nil
}

// SaveMessage 保存消息
func (d *Database) SaveMessage(sessionID string, message *models.Message) error {
	query := `
	INSERT INTO messages (id, session_id, role, content, timestamp)
	VALUES (?, ?, ?, ?, ?)
	`

	_, err := d.db.Exec(query, message.ID, sessionID, message.Role, message.Content, message.Timestamp)
	if err != nil {
		return fmt.Errorf("保存消息失败: %w", err)
	}

	// 更新会话的更新时间
	updateQuery := `UPDATE sessions SET updated_at = ? WHERE id = ?`
	_, err = d.db.Exec(updateQuery, time.Now(), sessionID)
	if err != nil {
		return fmt.Errorf("更新会话时间失败: %w", err)
	}

	return nil
}

// GetMessages 获取会话的所有消息
func (d *Database) GetMessages(sessionID string) ([]*models.Message, error) {
	query := `
	SELECT id, role, content, timestamp
	FROM messages
	WHERE session_id = ?
	ORDER BY timestamp ASC
	`

	rows, err := d.db.Query(query, sessionID)
	if err != nil {
		return nil, fmt.Errorf("查询消息列表失败: %w", err)
	}
	defer rows.Close()

	messages := make([]*models.Message, 0)
	for rows.Next() {
		message := &models.Message{}
		var roleStr string
		if err := rows.Scan(&message.ID, &roleStr, &message.Content, &message.Timestamp); err != nil {
			return nil, fmt.Errorf("读取消息数据失败: %w", err)
		}
		message.Role = models.Role(roleStr)
		messages = append(messages, message)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("遍历消息列表失败: %w", err)
	}

	return messages, nil
}
