package storage

import (
	"api/pkg/queue/que/job"
	"api/pkg/queue/que/log"
)

type Bucket interface {
	Push(value []byte) bool
	PushJob(jobid, value []byte) bool
	PushJobs(jobs []*job.Job) bool
	Shift() []byte
	Length() uint64
	GetLogger() log.Logger
	Close()
}
