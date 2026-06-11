package deployer

import (
	"context"
	"fmt"

	"github.com/eniacsystems/eniac-deploy/internal/config"
	"github.com/eniacsystems/eniac-deploy/internal/errors"
	"github.com/eniacsystems/eniac-deploy/internal/ftp"
	"github.com/eniacsystems/eniac-deploy/internal/history"
)

type Rollbacker struct {
	project *config.Project
	ftp     *ftp.Client
	history *history.CRUD
}

func NewRollbacker(project *config.Project, ftpClient *ftp.Client, historyCRUD *history.CRUD) *Rollbacker {
	return &Rollbacker{
		project: project,
		ftp:     ftpClient,
		history: historyCRUD,
	}
}

func (r *Rollbacker) Execute(ctx context.Context) error {
	if r.history == nil {
		return errors.NewDeployError("ROLLBACK_NOT_FOUND", "history not available", nil)
	}

	lastSuccess, err := r.history.GetLastSuccessful(r.project.Name)
	if err != nil {
		return errors.NewDeployError("HISTORY_READ_FAILED", "failed to get last successful deployment", err)
	}

	if lastSuccess == nil {
		return errors.ErrRollbackNotFound
	}

	fmt.Printf("Rolling back deployment #%d from %s\n", lastSuccess.ID, lastSuccess.Timestamp)

	rec := &history.DeployRecord{
		ProjectID:         r.project.Name,
		Status:            history.StatusRolledBack,
		RolledBackFromID:  &lastSuccess.ID,
	}

	if _, err := r.history.Insert(rec); err != nil {
		return errors.NewDeployError("HISTORY_WRITE_FAILED", "failed to record rollback", err)
	}

	fmt.Printf("Rollback recorded successfully\n")
	return nil
}