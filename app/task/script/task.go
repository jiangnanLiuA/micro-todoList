package script

import (
	"context"
	"micro_todoList/app/task/repository/mq/task"
	log "micro_todoList/pkg/logger"
)

func TaskCreateSync(ctx context.Context) {
	tSync := new(task.SyncTask)
	err := tSync.RunTaskCreate(ctx)
	if err != nil {
		log.LogrusObj.Infof("RunTaskCreate:%s", err)
	}
}
