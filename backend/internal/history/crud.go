package history

import (
	"database/sql"
	"fmt"
)

type CRUD struct {
	db *DB
}

func NewCRUD(db *DB) *CRUD {
	return &CRUD{db: db}
}

func (c *CRUD) Insert(rec *DeployRecord) (int64, error) {
	query := `
		INSERT INTO deployments (project_id, status, artifact_hash, migration_result, commit_hash, rolled_back_from_id)
		VALUES (?, ?, ?, ?, ?, ?)
	`

	var rolledBackID interface{}
	if rec.RolledBackFromID != nil {
		rolledBackID = *rec.RolledBackFromID
	} else {
		rolledBackID = nil
	}

	result, err := c.db.conn.Exec(query,
		rec.ProjectID,
		rec.Status,
		rec.ArtifactHash,
		rec.MigrationResult,
		rec.CommitHash,
		rolledBackID,
	)
	if err != nil {
		return 0, fmt.Errorf("failed to insert record: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get last insert ID: %w", err)
	}

	return id, nil
}

func (c *CRUD) GetByProject(projectID string, limit int) ([]*DeployRecord, error) {
	if limit <= 0 {
		limit = 20
	}

	query := `
		SELECT id, project_id, timestamp, status, artifact_hash, migration_result, commit_hash, rolled_back_from_id
		FROM deployments
		WHERE project_id = ?
		ORDER BY id DESC
		LIMIT ?
	`

	rows, err := c.db.conn.Query(query, projectID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query records: %w", err)
	}
	defer rows.Close()

	var records []*DeployRecord
	for rows.Next() {
		rec := &DeployRecord{}
		var artifactHash, migrationResult, commitHash sql.NullString
		var rolledBackFromID sql.NullInt64

		err := rows.Scan(
			&rec.ID,
			&rec.ProjectID,
			&rec.Timestamp,
			&rec.Status,
			&artifactHash,
			&migrationResult,
			&commitHash,
			&rolledBackFromID,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		rec.ArtifactHash = artifactHash.String
		rec.MigrationResult = migrationResult.String
		rec.CommitHash = commitHash.String
		if rolledBackFromID.Valid {
			rec.RolledBackFromID = &rolledBackFromID.Int64
		}

		records = append(records, rec)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return records, nil
}

func (c *CRUD) GetLastSuccessful(projectID string) (*DeployRecord, error) {
	query := `
		SELECT id, project_id, timestamp, status, artifact_hash, migration_result, commit_hash, rolled_back_from_id
		FROM deployments
		WHERE project_id = ? AND status = ?
		ORDER BY id DESC
		LIMIT 1
	`

	rec := &DeployRecord{}
	var artifactHash, migrationResult, commitHash sql.NullString
	var rolledBackFromID sql.NullInt64

	err := c.db.conn.QueryRow(query, projectID, StatusSuccess).Scan(
		&rec.ID,
		&rec.ProjectID,
		&rec.Timestamp,
		&rec.Status,
		&artifactHash,
		&migrationResult,
		&commitHash,
		&rolledBackFromID,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to query last successful: %w", err)
	}

	rec.ArtifactHash = artifactHash.String
	rec.MigrationResult = migrationResult.String
	rec.CommitHash = commitHash.String
	if rolledBackFromID.Valid {
		rec.RolledBackFromID = &rolledBackFromID.Int64
	}

	return rec, nil
}

func (c *CRUD) GetByID(id int64) (*DeployRecord, error) {
	query := `
		SELECT id, project_id, timestamp, status, artifact_hash, migration_result, commit_hash, rolled_back_from_id
		FROM deployments
		WHERE id = ?
	`

	rec := &DeployRecord{}
	var artifactHash, migrationResult, commitHash sql.NullString
	var rolledBackFromID sql.NullInt64

	err := c.db.conn.QueryRow(query, id).Scan(
		&rec.ID,
		&rec.ProjectID,
		&rec.Timestamp,
		&rec.Status,
		&artifactHash,
		&migrationResult,
		&commitHash,
		&rolledBackFromID,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to query by ID: %w", err)
	}

	rec.ArtifactHash = artifactHash.String
	rec.MigrationResult = migrationResult.String
	rec.CommitHash = commitHash.String
	if rolledBackFromID.Valid {
		rec.RolledBackFromID = &rolledBackFromID.Int64
	}

	return rec, nil
}

func (c *CRUD) CountByProject(projectID string) (int, error) {
	var count int
	err := c.db.conn.QueryRow(
		"SELECT COUNT(*) FROM deployments WHERE project_id = ?",
		projectID,
	).Scan(&count)

	if err != nil {
		return 0, fmt.Errorf("failed to count records: %w", err)
	}

	return count, nil
}
