package queue

import (
	"encoding/json"
	"sync"

	"api/app/models/opera_log"
	"api/pkg/logger"
	"api/pkg/queue/que"
	"api/pkg/queue/que/job"
)

var QueueFile *que.Queue

var onceFile sync.Once

func ConnectNsq(name string) {
	onceFile.Do(func() {
		QueueFile = que.NewQueue(8, 32, "storage/queue", name, nil)
	})
}

//Producers 生产
func Producers(task interface{}, name string) {
	tasks, _ := json.Marshal(task)
	jo := &job.Job{
		Name:    name,
		Payload: tasks,
	}
	QueueFile.Dispatch(jo)
}

func Oplog(job *job.Job) {
	var log opera_log.OperaLog
	if err := job.Unmarshal(&log); err != nil {
		logger.Error(err.Error())
		return
	}
	log.Create()
	if log.ID <= 0 {
		Producers(log, "oplog")
	}
}
