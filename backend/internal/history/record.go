package history

import "time"

type DeployRecord struct {
	ID              int64
	ProjectID       string
	Timestamp       time.Time
	Status          string
	ArtifactHash    string
	MigrationResult string
	CommitHash      string
	RolledBackFromID *int64
}

const (
	StatusSuccess    = "SUCCESS"
	StatusFailed     = "FAILED"
	StatusRolledBack = "ROLLED_BACK"
)

func NewSuccessRecord(projectID, artifactHash, migrationResult, commitHash string) *DeployRecord {
	return &DeployRecord{
		ProjectID:       projectID,
		Status:          StatusSuccess,
		ArtifactHash:    artifactHash,
		MigrationResult: migrationResult,
		CommitHash:      commitHash,
	}
}

func NewFailedRecord(projectID, artifactHash, migrationResult string) *DeployRecord {
	return &DeployRecord{
		ProjectID:       projectID,
		Status:          StatusFailed,
		ArtifactHash:    artifactHash,
		MigrationResult: migrationResult,
	}
}