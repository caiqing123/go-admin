package jobs

import (
	"time"

	"github.com/robfig/cron/v3"

	"api/app/models/job"
	"api/pkg/database"
	"api/pkg/file"
	"api/pkg/http"
	"api/pkg/logger"
)

var timeFormat = "2006-01-02 15:04:05"
var retryCount = 3

var JobList map[string]JobsExec

var secondParser = cron.NewParser(cron.Second | cron.Minute |
	cron.Hour | cron.Dom | cron.Month | cron.DowOptional | cron.Descriptor)
var crontab = cron.New(cron.WithParser(secondParser), cron.WithChain())

type JobCore struct {
	InvokeTarget   string
	Name           string
	JobId          int
	EntryId        int
	CronExpression string
	Args           string
}

// HttpJob 任务类型 http
type HttpJob struct {
	JobCore
}

type ExecJob struct {
	JobCore
}

func (e *ExecJob) Run() {
	startTime := time.Now()
	var obj = JobList[e.InvokeTarget]
	if obj == nil {
		file.CronLog(time.Now().Format(timeFormat) + "[Job] ExecJob Run job nil")
		return
	}
	err := CallExec(obj.(JobsExec), e.Args)
	if err != nil {
		// 如果失败暂停一段时间重试
		file.CronLog(time.Now().Format(timeFormat)+" [ERROR] mission failed! %v", err)
	}
	// 结束时间
	endTime := time.Now()

	// 执行时间
	latencyTime := endTime.Sub(startTime)

	file.CronLog(time.Now().Format(timeFormat)+"[Job] JobCore %s exec success , spend :%v , invokeTarget :%s", e.Name, latencyTime, e.InvokeTarget)
	return
}

// Run http 任务接口
func (h *HttpJob) Run() {

	startTime := time.Now()
	var count = 0
	var err error
	var str string
	/* 循环 */
LOOP:
	if count < retryCount {
		/* 跳过迭代 */
		str, err = http.Get(h.InvokeTarget)
		if err != nil {
			// 如果失败暂停一段时间重试
			file.CronLog(time.Now().Format(timeFormat)+" [ERROR] mission failed! %v", err)
			file.CronLog(time.Now().Format(timeFormat)+" [INFO] Retry after the task fails %d seconds! %s \n", (count+1)*5, str)
			time.Sleep(time.Duration(count+1) * 5 * time.Second)
			count = count + 1
			goto LOOP
		}
	}
	// 结束时间
	endTime := time.Now()

	// 执行时间
	latencyTime := endTime.Sub(startTime)

	file.CronLog(time.Now().Format(timeFormat)+"[Job] JobCore %s http success , spend :%v , invokeTarget :%s", h.Name, latencyTime, h.InvokeTarget)
	return
}

func Setup() {
	if !database.DB.Migrator().HasTable("jobs") {
		return
	}
	jobList := job.GetList()
	if len(jobList) == 0 {
		logger.Printf("JobCore total:0")
	}
	job.RemoveAllEntryID()

	for i := 0; i < len(jobList); i++ {
		if jobList[i].JobType == 1 {
			j := &HttpJob{}
			j.InvokeTarget = jobList[i].InvokeTarget
			j.CronExpression = jobList[i].CronExpression
			j.JobId = jobList[i].JobId
			j.Name = jobList[i].JobName

			jobList[i].EntryId, _ = AddJob(crontab, j)
		} else if jobList[i].JobType == 2 {
			j := &ExecJob{}
			j.InvokeTarget = jobList[i].InvokeTarget
			j.CronExpression = jobList[i].CronExpression
			j.JobId = jobList[i].JobId
			j.Name = jobList[i].JobName
			j.Args = jobList[i].Args
			jobList[i].EntryId, _ = AddJob(crontab, j)
		}
		job.Save(jobList[i])
	}

	// 其中任务
	crontab.Start()
	logger.Printf("JobCore start success.")
	// 关闭任务
	defer crontab.Stop()
	select {}
}

// AddJob 添加任务 AddJob(invokeTarget string, jobId int, jobName string, cronExpression string)
func AddJob(c *cron.Cron, job Job) (int, error) {
	if job == nil {
		logger.Error("AddJob unknown")
		return 0, nil
	}
	return job.addJob(c)
}

func (h *HttpJob) addJob(c *cron.Cron) (int, error) {
	id, err := c.AddJob(h.CronExpression, h)
	if err != nil {
		logger.Printf("JobCore AddJob error %v", err)
		return 0, err
	}
	EntryId := int(id)
	return EntryId, nil
}

func (e *ExecJob) addJob(c *cron.Cron) (int, error) {
	id, err := c.AddJob(e.CronExpression, e)
	if err != nil {
		logger.Printf("JobCore AddJob error %v", err)
		return 0, err
	}
	EntryId := int(id)
	return EntryId, nil
}

// Remove 移除任务
func Remove(c *cron.Cron, entryID int) chan bool {
	ch := make(chan bool)
	go func() {
		c.Remove(cron.EntryID(entryID))
		logger.Printf("JobCore Remove success ,info entryID :%v", entryID)
		ch <- true
	}()
	return ch
}

//StartJob 启动任务
func StartJob(data job.Job) {
	if data.MisfirePolicy == 2 || data.Status == 1 {
		return
	}
	if data.JobType == 1 {
		var j = &HttpJob{}
		j.InvokeTarget = data.InvokeTarget
		j.CronExpression = data.CronExpression
		j.JobId = data.JobId
		j.Name = data.JobName
		data.EntryId, _ = AddJob(crontab, j)
	} else {
		var j = &ExecJob{}
		j.InvokeTarget = data.InvokeTarget
		j.CronExpression = data.CronExpression
		j.JobId = data.JobId
		j.Name = data.JobName
		j.Args = data.Args
		data.EntryId, _ = AddJob(crontab, j)
	}
	data.Save()
}

// RemoveJob 删除job
func RemoveJob(data job.Job) {
	if data.EntryId == 0 {
		return
	}
	cn := Remove(crontab, data.EntryId)
	select {
	case res := <-cn:
		if res {
			err := database.DB.Table("jobs").Where("entry_id = ?", data.EntryId).Update("entry_id", 0).Error
			if err != nil {
				logger.Printf("RemoveJob db error: %s", err)
			}
		}
	case <-time.After(time.Second * 1):
		logger.Printf("RemoveJob db error: %s", "操作超时！")
	}
}
