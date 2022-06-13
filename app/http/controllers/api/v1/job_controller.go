package v1

import (
	"strconv"

	"api/app/models/job"
	"api/app/requests"
	cron "api/pkg/job"
	"api/pkg/response"

	"github.com/gin-gonic/gin"
)

type JobController struct {
	BaseAPIController
}

// Index 数据列表
func (ctrl *JobController) Index(c *gin.Context) {
	id := c.DefaultQuery("id", "")
	if id != "" {
		data := job.GetId(id)
		response.JSON(c, gin.H{
			"data": data,
		})
		return
	}

	request := requests.JobPaginationRequest{}
	if ok := requests.Validate(c, &request, requests.JobPagination); !ok {
		return
	}
	data, pager := job.Paginate(c, 10, request)

	response.JSON(c, gin.H{
		"data":  data,
		"pager": pager,
	})
}

// Update 修改
func (ctrl *JobController) Update(c *gin.Context) {
	request := requests.JobRequest{}
	if ok := requests.Validate(c, &request, requests.JobSave); !ok {
		return
	}

	if request.JobId == 0 {
		response.NormalVerificationError(c, "id为空")
		return
	}
	jobs := job.GetId(strconv.Itoa(request.JobId))
	jobs.JobName = request.JobName
	jobs.JobGroup = request.JobGroup
	jobs.JobType = request.JobType
	jobs.CronExpression = request.CronExpression
	jobs.InvokeTarget = request.InvokeTarget
	jobs.Args = request.Args
	jobs.MisfirePolicy = request.MisfirePolicy
	jobs.Concurrent = request.Concurrent
	jobs.Status = request.Status

	rowsAffected := jobs.Save()
	if rowsAffected > 0 {
		cron.RemoveJob(jobs)
		cron.StartJob(jobs)
		response.Success(c)
	} else {
		response.Abort500(c, "更新失败，请稍后尝试~")
	}
}

// Add 添加
func (ctrl *JobController) Add(c *gin.Context) {
	request := requests.JobRequest{}
	if ok := requests.Validate(c, &request, requests.JobSave); !ok {
		return
	}
	CronModel := job.Job{
		JobName:        request.JobName,
		JobGroup:       request.JobGroup,
		JobType:        request.JobType,
		CronExpression: request.CronExpression,
		InvokeTarget:   request.InvokeTarget,
		Args:           request.Args,
		MisfirePolicy:  request.MisfirePolicy,
		Concurrent:     request.Concurrent,
		Status:         request.Status,
	}
	CronModel.Create()

	if CronModel.JobId > 0 {
		cron.StartJob(CronModel)
		response.Success(c)
	} else {
		response.Abort500(c, "创建失败，请稍后尝试~")
	}
}

func (ctrl *JobController) Delete(c *gin.Context) {
	request := requests.JobDeleteRequest{}
	if ok := requests.Validate(c, &request, requests.JobDelete); !ok {
		return
	}
	for _, v := range request.Ids {
		jobs := job.GetId(strconv.Itoa(v))
		cron.RemoveJob(jobs)
	}
	rowsAffected := job.DeleteIds(request.Ids, job.Job{})
	if rowsAffected > 0 {
		response.Success(c)
		return
	}
	response.Abort500(c, "删除失败，请稍后尝试~")
}

func (ctrl *JobController) StartJob(c *gin.Context) {
	jobId := c.DefaultQuery("job_id", "")
	if jobId == "" {
		response.NormalVerificationError(c, "任务id为空")
		return
	}

	jobs := job.GetId(jobId)
	if jobs.EntryId > 0 {
		response.Abort500(c, "启动失败,任务已启动")
		return
	}
	jobs.Status = 2
	jobs.MisfirePolicy = 1
	cron.StartJob(jobs)
	response.Success(c)
}

func (ctrl *JobController) RemoveJob(c *gin.Context) {
	jobId := c.DefaultQuery("job_id", "")
	if jobId == "" {
		response.NormalVerificationError(c, "任务id为空")
		return
	}
	jobs := job.GetId(jobId)
	if jobs.EntryId == 0 {
		response.Abort500(c, "关闭失败,任务未启动")
		return
	}
	jobs.Status = 1
	jobs.Save()
	cron.RemoveJob(jobs)
	response.Success(c)
}

func (ctrl *JobController) ParseCron(c *gin.Context) {
	expr := c.DefaultQuery("expr", "")
	if expr == "" {
		response.NormalVerificationError(c, "表达式为空")
		return
	}
	respTimes, err := cron.ParseCronList(expr)
	if err != nil {
		response.Abort500(c, "表达式解析错误")
		return
	}
	response.Data(c, respTimes)
}
