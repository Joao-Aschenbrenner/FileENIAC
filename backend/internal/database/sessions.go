// SPDX-License-Identifier: MIT

package database

import (
	"database/sql"
	"fmt"
	"time"
)

type Session struct {
	ID            int64  `json:"id"`
	Name          string `json:"name"`
	WorkspacePath string `json:"workspace_path"`
	Description   string `json:"description,omitempty"`
	GitHubUser    string `json:"github_user,omitempty"`
	GitHubToken   string `json:"github_token,omitempty"`
	IsActive      bool   `json:"is_active"`
	CreatedAt     string `json:"created_at"`
	UpdatedAt     string `json:"updated_at"`
}

type SessionStore struct {
	db *DB
}

func NewSessionStore(db *DB) *SessionStore {
	return &SessionStore{db: db}
}

func (s *SessionStore) List() ([]Session, error) {
	rows, err := s.db.conn.Query(`
		SELECT id, name, workspace_path, description, github_user, github_token, is_active, created_at, updated_at
		FROM sessions
		ORDER BY updated_at DESC
	`)
	if err != nil {
		return nil, fmt.Errorf("list sessions: %w", err)
	}
	defer rows.Close()

	var sessions []Session
	for rows.Next() {
		var sess Session
		var desc, ghUser, ghToken sql.NullString
		err := rows.Scan(&sess.ID, &sess.Name, &sess.WorkspacePath, &desc, &ghUser, &ghToken, &sess.IsActive, &sess.CreatedAt, &sess.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("scan session: %w", err)
		}
		sess.Description = desc.String
		sess.GitHubUser = ghUser.String
		sess.GitHubToken = ghToken.String
		sessions = append(sessions, sess)
	}
	return sessions, nil
}

func (s *SessionStore) Get(id int64) (*Session, error) {
	var sess Session
	var desc, ghUser, ghToken sql.NullString
	err := s.db.conn.QueryRow(`
		SELECT id, name, workspace_path, description, github_user, github_token, is_active, created_at, updated_at
		FROM sessions WHERE id = ?
	`, id).Scan(&sess.ID, &sess.Name, &sess.WorkspacePath, &desc, &ghUser, &ghToken, &sess.IsActive, &sess.CreatedAt, &sess.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("session not found")
	}
	if err != nil {
		return nil, fmt.Errorf("get session: %w", err)
	}
	sess.Description = desc.String
	sess.GitHubUser = ghUser.String
	sess.GitHubToken = ghToken.String
	return &sess, nil
}

func (s *SessionStore) Create(name, workspacePath, description string) (*Session, error) {
	now := time.Now().Format(time.RFC3339)
	result, err := s.db.conn.Exec(`
		INSERT INTO sessions (name, workspace_path, description, is_active, created_at, updated_at)
		VALUES (?, ?, ?, 0, ?, ?)
	`, name, workspacePath, description, now, now)
	if err != nil {
		return nil, fmt.Errorf("create session: %w", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("last insert id: %w", err)
	}
	return &Session{
		ID:            id,
		Name:          name,
		WorkspacePath: workspacePath,
		Description:   description,
		IsActive:      false,
		CreatedAt:     now,
		UpdatedAt:     now,
	}, nil
}

func (s *SessionStore) Update(id int64, name, workspacePath, description, githubUser, githubToken string, isActive *bool) error {
	sess, err := s.Get(id)
	if err != nil {
		return err
	}
	if name != "" {
		sess.Name = name
	}
	if workspacePath != "" {
		sess.WorkspacePath = workspacePath
	}
	if description != "" {
		sess.Description = description
	}
	if githubUser != "" {
		sess.GitHubUser = githubUser
	}
	if githubToken != "" {
		sess.GitHubToken = githubToken
	}
	if isActive != nil {
		sess.IsActive = *isActive
	}
	sess.UpdatedAt = time.Now().Format(time.RFC3339)
	_, err = s.db.conn.Exec(`
		UPDATE sessions
		SET name = ?, workspace_path = ?, description = ?, github_user = ?, is_active = ?, updated_at = ?
		WHERE id = ?
	`, sess.Name, sess.WorkspacePath, sess.Description, sess.GitHubUser, sess.IsActive, sess.UpdatedAt, id)
	if err != nil {
		return fmt.Errorf("update session: %w", err)
	}
	return nil
}

func (s *SessionStore) Delete(id int64) error {
	result, err := s.db.conn.Exec("DELETE FROM sessions WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("delete session: %w", err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("session not found")
	}
	return nil
}

func (s *SessionStore) Activate(id int64) error {
	tx, err := s.db.conn.Begin()
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback()

	_, err = tx.Exec("UPDATE sessions SET is_active = 0")
	if err != nil {
		return fmt.Errorf("deactivate all: %w", err)
	}
	_, err = tx.Exec("UPDATE sessions SET is_active = 1, updated_at = ? WHERE id = ?", time.Now().Format(time.RFC3339), id)
	if err != nil {
		return fmt.Errorf("activate session: %w", err)
	}
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("commit: %w", err)
	}
	return nil
}

func (s *SessionStore) ClearWorkspace(id int64) error {
	_, err := s.db.conn.Exec("UPDATE sessions SET workspace_path = '', updated_at = ? WHERE id = ?", time.Now().Format(time.RFC3339), id)
	if err != nil {
		return fmt.Errorf("clear workspace: %w", err)
	}
	return nil
}

func (s *SessionStore) GetActive() (*Session, error) {
	var sess Session
	var desc, ghUser, ghToken sql.NullString
	err := s.db.conn.QueryRow(`
		SELECT id, name, workspace_path, description, github_user, github_token, is_active, created_at, updated_at
		FROM sessions WHERE is_active = 1
		LIMIT 1
	`).Scan(&sess.ID, &sess.Name, &sess.WorkspacePath, &desc, &ghUser, &ghToken, &sess.IsActive, &sess.CreatedAt, &sess.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get active session: %w", err)
	}
	sess.Description = desc.String
	sess.GitHubUser = ghUser.String
	sess.GitHubToken = ghToken.String
	return &sess, nil
}
