package jobs

import (
	"fmt"
	"time"

	"api/pkg/config"
	"api/pkg/hotlist"
	"api/pkg/logger"
)

// InitJob 需要将定义的struct 添加到字典中；
// 字典 key 可以配置到 自动任务 调用目标 中；
func InitJob() {
	JobList = map[string]JobsExec{
		"ExamplesOne":  ExamplesOne{},
		"ExamplesNews": ExamplesNews{},
		// ...
	}
}

// ExamplesOne 新添加的job 必须按照以下格式定义，并实现Exec函数
type ExamplesOne struct {
}

func (t ExamplesOne) Exec(arg interface{}) error {
	str := time.Now().Format(timeFormat) + " [INFO] JobCore ExamplesOne exec success"
	//这里需要注意 Examples 传入参数是 string 所以 arg.(string)；请根据对应的类型进行转化；
	switch arg.(type) {

	case string:
		if arg.(string) != "" {
			//logger重新每天0点初始化
			logger.Logger = nil
			logger.InitLogger(
				config.GetString("log.filename"),
				config.GetInt("log.max_size"),
				config.GetInt("log.max_backup"),
				config.GetInt("log.max_age"),
				config.GetBool("log.compress"),
				config.GetString("log.type"),
				config.GetString("log.level"),
			)
			fmt.Println(str, arg.(string))
		} else {
			fmt.Println(str, "arg is nil")
		}
		break
	}

	return nil
}

type ExamplesNews struct {
}

func (t ExamplesNews) Exec(arg interface{}) error {
	str := time.Now().Format(timeFormat) + " [INFO] JobCore ExamplesNews exec success"
	fmt.Println(str, arg.(string))
	switch arg.(type) {
	case string:
		hotlist.All(true)
		break
	}

	return nil
}
