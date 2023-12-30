package app

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/dashotv/minion"
)

type CleanupJobs struct {
	minion.WorkerDefaults[*CleanupJobs]
}

func (j *CleanupJobs) Kind() string { return "cleanup_jobs" }
func (j *CleanupJobs) Work(ctx context.Context, job *minion.Job[*CleanupJobs]) error {
	if _, err := app.DB.Minion.Collection.DeleteMany(context.Background(), bson.M{"created_at": bson.M{"$lt": time.Now().UTC().AddDate(0, 0, -3)}}); err != nil {
		return errors.Wrap(err, "deleting messages")
	}
	return nil
}
