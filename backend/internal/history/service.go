// SPDX-License-Identifier: MIT
package history

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/ENIACSystems/FileENIAC/backend/internal/database"
	"github.com/ENIACSystems/FileENIAC/backend/internal/log"
	"go.uber.org/zap"
)

const (
	EventDeployStarted   = "DEPLOY_STARTED"
	EventDeploySuccess   = "DEPLOY_SUCCESS"
	EventDeployFailed    = "DEPLOY_FAILED"
	EventRollbackStarted = "ROLLBACK_STARTED"
	EventRollbackSuccess = "ROLLBACK_SUCCESS"
	EventRollbackFailed  = "ROLLBACK_FAILED"
	EventVerifySuccess   = "VERIFY_SUCCESS"
	EventVerifyFailed    = "VERIFY_FAILED"
	EventSyncStarted     = "SYNC_STARTED"
	EventSyncCompleted   = "SYNC_COMPLETED"
	EventSyncFailed      = "SYNC_FAILED"
	EventProjectCreated  = "PROJECT_CREATED"
	EventProjectRemoved  = "PROJECT_REMOVED"
	EventServerAdded     = "SERVER_ADDED"
	EventServerUpdated   = "SERVER_UPDATED"
	EventServerRemoved   = "SERVER_REMOVED"
	EventAlert           = "ALERT"
	EventError           = "ERROR"
)

type Event struct {
	ID          int64  `json:"id"`
	EventType   string `json:"event_type"`
	Description string `json:"description"`
	Metadata    string `json:"metadata,omitempty"`
	CreatedAt   string `json:"created_at"`
}

type DeployLog struct {
	ID            int64  `json:"id"`
	ProjectID     int64  `json:"project_id"`
	DeployID      string `json:"deploy_id"`
	Status        string `json:"status"`
	CommitHash    string `json:"commit_hash,omitempty"`
	CommitMessage string `json:"commit_message,omitempty"`
	Branch        string `json:"branch,omitempty"`
	ArtifactHash  string `json:"artifact_hash,omitempty"`
	FilesCount    int    `json:"files_count,omitempty"`
	Manifest      string `json:"manifest,omitempty"`
	ErrorMessage  string `json:"error_message,omitempty"`
	StartedAt     string `json:"started_at"`
	CompletedAt   string `json:"completed_at"`
	CreatedAt     string `json:"created_at"`
}

type RollbackLog struct {
	ID        int64  `json:"id"`
	DeployID  string `json:"deploy_id"`
	ProjectID int64  `json:"project_id"`
	Reason    string `json:"reason"`
	Status    string `json:"status"`
	CreatedAt string `json:"created_at"`
}

type Service struct {
	db *database.DB
}

func NewService(db *database.DB) *Service {
	return &Service{db: db}
}

func (s *Service) RecordEvent(eventType, description string, metadata map[string]interface{}) error {
	metaJSON := "{}"
	if metadata != nil {
		b, err := json.Marshal(metadata)
		if err == nil {
			metaJSON = string(b)
		}
	}

	_, err := s.db.Exec(
		"INSERT INTO events (event_type, description, metadata) VALUES (?, ?, ?)",
		eventType, description, metaJSON,
	)
	if err != nil {
		return fmt.Errorf("failed to record event: %w", err)
	}

	log.L().Info("event recorded",
		zap.String("type", eventType),
		zap.String("description", description),
	)
	return nil
}

func (s *Service) RecordDeploy(log *DeployLog) (int64, error) {
	query := `INSERT INTO deploy_logs (project_id, deploy_id, status, commit_hash, commit_message, branch, artifact_hash, files_count, manifest, error_message, started_at, completed_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	result, err := s.db.Exec(query,
		log.ProjectID, log.DeployID, log.Status, log.CommitHash,
		log.CommitMessage, log.Branch, log.ArtifactHash, log.FilesCount,
		log.Manifest, log.ErrorMessage, log.StartedAt, log.CompletedAt,
	)
	if err != nil {
		return 0, fmt.Errorf("failed to record deploy: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	eventType := EventDeploySuccess
	if log.Status == "failed" {
		eventType = EventDeployFailed
	}
	s.RecordEvent(eventType, fmt.Sprintf("Deploy %s for project %d: %s", log.DeployID, log.ProjectID, log.Status),
		map[string]interface{}{"deploy_id": log.DeployID, "project_id": log.ProjectID, "status": log.Status})

	return id, nil
}

func (s *Service) RecordRollback(log *RollbackLog) (int64, error) {
	query := `INSERT INTO rollback_logs (deploy_id, project_id, reason, status)
		VALUES (?, ?, ?, ?)`

	result, err := s.db.Exec(query, log.DeployID, log.ProjectID, log.Reason, log.Status)
	if err != nil {
		return 0, fmt.Errorf("failed to record rollback: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	eventType := EventRollbackSuccess
	if log.Status == "failed" {
		eventType = EventRollbackFailed
	}
	s.RecordEvent(eventType, fmt.Sprintf("Rollback deploy %s for project %d", log.DeployID, log.ProjectID),
		map[string]interface{}{"deploy_id": log.DeployID, "project_id": log.ProjectID})

	return id, nil
}

func (s *Service) GetDeployHistory(projectID int64, limit int) ([]*DeployLog, error) {
	if limit <= 0 {
		limit = 20
	}

	rows, err := s.db.Query(
		`SELECT id, project_id, deploy_id, status, COALESCE(commit_hash,''), COALESCE(commit_message,''), COALESCE(branch,''), COALESCE(artifact_hash,''), COALESCE(files_count,0), COALESCE(manifest,''), COALESCE(error_message,''), started_at, completed_at, created_at
		FROM deploy_logs WHERE project_id = ? ORDER BY id DESC LIMIT ?`,
		projectID, limit,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query deploy history: %w", err)
	}
	defer rows.Close()

	var logs []*DeployLog
	for rows.Next() {
		l := &DeployLog{}
		err := rows.Scan(&l.ID, &l.ProjectID, &l.DeployID, &l.Status,
			&l.CommitHash, &l.CommitMessage, &l.Branch, &l.ArtifactHash,
			&l.FilesCount, &l.Manifest, &l.ErrorMessage,
			&l.StartedAt, &l.CompletedAt, &l.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan deploy log: %w", err)
		}
		logs = append(logs, l)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return logs, nil
}

func (s *Service) GetLastDeploy(projectID int64) (*DeployLog, error) {
	l := &DeployLog{}
	query := `SELECT id, project_id, deploy_id, status, COALESCE(commit_hash,''), COALESCE(commit_message,''), COALESCE(branch,''), COALESCE(artifact_hash,''), COALESCE(files_count,0), COALESCE(manifest,''), COALESCE(error_message,''), started_at, completed_at, created_at
		FROM deploy_logs WHERE project_id = ? ORDER BY id DESC LIMIT 1`

	err := s.db.QueryRow(query, projectID).Scan(
		&l.ID, &l.ProjectID, &l.DeployID, &l.Status,
		&l.CommitHash, &l.CommitMessage, &l.Branch, &l.ArtifactHash,
		&l.FilesCount, &l.Manifest, &l.ErrorMessage,
		&l.StartedAt, &l.CompletedAt, &l.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("no deploys found for project %d", projectID)
	}
	return l, nil
}

func (s *Service) GetRecentEvents(limit int) ([]*Event, error) {
	if limit <= 0 {
		limit = 50
	}

	rows, err := s.db.Query(
		"SELECT id, event_type, description, COALESCE(metadata,''), created_at FROM events ORDER BY id DESC LIMIT ?",
		limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []*Event
	for rows.Next() {
		e := &Event{}
		if err := rows.Scan(&e.ID, &e.EventType, &e.Description, &e.Metadata, &e.CreatedAt); err != nil {
			return nil, err
		}
		events = append(events, e)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return events, nil
}

func (s *Service) GetEventList(eventType string, limit, offset int) ([]*Event, error) {
	if limit <= 0 {
		limit = 50
	}

	var rows *sql.Rows
	var err error
	if eventType != "" {
		rows, err = s.db.Query(
			"SELECT id, event_type, description, COALESCE(metadata,''), created_at FROM events WHERE event_type = ? ORDER BY id DESC LIMIT ? OFFSET ?",
			eventType, limit, offset,
		)
	} else {
		rows, err = s.db.Query(
			"SELECT id, event_type, description, COALESCE(metadata,''), created_at FROM events ORDER BY id DESC LIMIT ? OFFSET ?",
			limit, offset,
		)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []*Event
	for rows.Next() {
		e := &Event{}
		if err := rows.Scan(&e.ID, &e.EventType, &e.Description, &e.Metadata, &e.CreatedAt); err != nil {
			return nil, err
		}
		events = append(events, e)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return events, nil
}
