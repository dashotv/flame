package app

import (
	"context"

	"github.com/dashotv/minion"
)

type CleanupJobs struct {
	minion.WorkerDefaults[*CleanupJobs]
}

func (j *CleanupJobs) Kind() string { return "cleanup_jobs" }
func (j *CleanupJobs) Work(ctx context.Context, job *minion.Job[*CleanupJobs]) error {
	return nil
}
