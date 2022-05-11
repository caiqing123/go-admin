package que

import (
	"fmt"
	"log"
	"strconv"
	"testing"
	"time"

	"api/pkg/queue/que/job"
)

func TestNewQueue(t *testing.T) {

	var (
		count    uint
		maxItems uint
	)
	count = 4
	maxItems = 16
	path := "storage/queue"
	q := NewQueue(count, maxItems, path, "mfworker", nil)
	q.Handler("Test", func(job *job.Job) {
		time.Sleep(time.Second)
		log.Printf("the job name %s, job body %s ", job.Name, job.Payload)
	})
	q.Start()
	defer q.Stop()

	go func() {
		for a := 0; a < 10; a++ {
			for i := 0; i < 100; i++ {
				jo := &job.Job{
					Name:    "Test",
					Payload: []byte("body " + strconv.Itoa(i)),
				}
				if (i % 2) != 0 {
					jo.Id = strconv.Itoa(a) + "/" + strconv.Itoa(i)
				}
				q.Dispatch(jo)
			}
		}
		var jobs []*job.Job
		for i := 0; i < 64; i++ {
			jo := &job.Job{
				Name:    "Test",
				Payload: []byte("body " + strconv.Itoa(i)),
			}
			if (i % 2) != 0 {
				jo.Id = strconv.Itoa(i)
			}
			jobs = append(jobs, jo)
		}
		q.DispatchJobs(jobs)
		num := q.CountPendingJobs()
		fmt.Println(num)
		if num <= 0 {
			t.Errorf("Queue jobs should not empty.")
		}
	}()
	<-time.After(10 * time.Second)
	return
}
