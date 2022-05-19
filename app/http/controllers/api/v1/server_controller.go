package v1

import (
	"archive/tar"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/net"

	"api/pkg/file"
	"api/pkg/ip"
	"api/pkg/logger"
	"api/pkg/response"
	"api/pkg/str"
)

const (
	B  = 1
	KB = 1024 * B
	MB = 1024 * KB
	GB = 1024 * MB
)

var (
	excludeNetInterfaces = []string{
		"lo", "tun", "docker", "veth", "br-", "vmbr", "vnet", "kube",
	}
)

var (
	netInSpeed, netOutSpeed, netInTransfer, netOutTransfer, lastUpdateNetStats uint64
	cachedBootTime                                                             time.Time
)

// GetHourDiffer 获取相差时间
func GetHourDiffer(startTime, endTime string) int64 {
	var hour int64
	t1, err := time.ParseInLocation("2006-01-02 15:04:05", startTime, time.Local)
	t2, err := time.ParseInLocation("2006-01-02 15:04:05", endTime, time.Local)
	if err == nil && t1.Before(t2) {
		diff := t2.Unix() - t1.Unix() //
		hour = diff / 3600
		return hour
	} else {
		return hour
	}
}

type ServerController struct {
	BaseAPIController
}

//ServerInfo 系统信息http
func (ctrl *ServerController) ServerInfo(c *gin.Context) {
	sysInfo, err := host.Info()
	osDic := make(map[string]interface{}, 0)
	osDic["goOs"] = runtime.GOOS
	osDic["arch"] = runtime.GOARCH
	osDic["mem"] = runtime.MemProfileRate
	osDic["compiler"] = runtime.Compiler
	osDic["version"] = runtime.Version()
	osDic["numGoroutine"] = runtime.NumGoroutine()
	osDic["ip"] = ip.GetLocationHost()
	osDic["projectDir"] = file.GetCurrentPath()
	osDic["hostName"] = sysInfo.Hostname
	osDic["time"] = time.Now().Format("2006-01-02 15:04:05")

	memory, _ := mem.VirtualMemory()
	memDic := make(map[string]interface{}, 0)
	memDic["used"] = memory.Used / MB
	memDic["total"] = memory.Total / MB

	memDic["percent"] = str.Round(memory.UsedPercent, 2)

	swapDic := make(map[string]interface{}, 0)
	swapDic["used"] = memory.SwapTotal - memory.SwapFree
	swapDic["total"] = memory.SwapTotal

	cpuDic := make(map[string]interface{}, 0)
	cpuDic["cpuInfo"], _ = cpu.Info()
	percent, _ := cpu.Percent(0, false)
	cpuDic["percent"] = str.Round(percent[0], 2)
	cpuDic["cpuNum"], _ = cpu.Counts(false)

	//服务器磁盘信息
	disklist := make([]disk.UsageStat, 0)
	//所有分区
	var diskTotal, diskUsed, diskUsedPercent float64
	diskInfo, err := disk.Partitions(true)
	if err == nil {
		for _, p := range diskInfo {
			diskDetail, err := disk.Usage(p.Mountpoint)
			if err == nil {
				diskDetail.UsedPercent, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", diskDetail.UsedPercent), 64)
				diskDetail.Total = diskDetail.Total / 1024 / 1024
				diskDetail.Used = diskDetail.Used / 1024 / 1024
				diskDetail.Free = diskDetail.Free / 1024 / 1024
				disklist = append(disklist, *diskDetail)

			}
		}
	}

	d, _ := disk.Usage("/")

	diskTotal = float64(d.Total / GB)
	diskUsed = float64(d.Used / GB)
	diskUsedPercent, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", d.UsedPercent), 64)

	diskDic := make(map[string]interface{}, 0)
	diskDic["total"] = diskTotal
	diskDic["used"] = diskUsed
	diskDic["percent"] = diskUsedPercent

	bootTime, _ := host.BootTime()
	cachedBootTime = time.Unix(int64(bootTime), 0)

	TrackNetworkSpeed()
	netDic := make(map[string]interface{}, 0)
	netDic["in"] = str.Round(float64(netInSpeed/KB), 2)
	netDic["out"] = str.Round(float64(netOutSpeed/KB), 2)
	response.Data(c, gin.H{
		"os":       osDic,
		"mem":      memDic,
		"cpu":      cpuDic,
		"disk":     diskDic,
		"net":      netDic,
		"swap":     swapDic,
		"location": "Aliyun",
		"bootTime": GetHourDiffer(cachedBootTime.Format("2006-01-02 15:04:05"), time.Now().Format("2006-01-02 15:04:05")),
	})
}

//GetServerInfo 系统信息ws
func (ctrl *ServerController) GetServerInfo() []byte {
	sysInfo, err := host.Info()
	osDic := make(map[string]interface{}, 0)
	osDic["goOs"] = runtime.GOOS
	osDic["arch"] = runtime.GOARCH
	osDic["mem"] = runtime.MemProfileRate
	osDic["compiler"] = runtime.Compiler
	osDic["version"] = runtime.Version()
	osDic["numGoroutine"] = runtime.NumGoroutine()
	osDic["ip"] = ip.GetLocationHost()
	osDic["projectDir"] = file.GetCurrentPath()
	osDic["hostName"] = sysInfo.Hostname
	osDic["time"] = time.Now().Format("2006-01-02 15:04:05")

	memory, _ := mem.VirtualMemory()
	memDic := make(map[string]interface{}, 0)
	memDic["used"] = memory.Used / MB
	memDic["total"] = memory.Total / MB

	memDic["percent"] = str.Round(memory.UsedPercent, 2)

	swapDic := make(map[string]interface{}, 0)
	swapDic["used"] = memory.SwapTotal - memory.SwapFree
	swapDic["total"] = memory.SwapTotal

	cpuDic := make(map[string]interface{}, 0)
	cpuDic["cpuInfo"], _ = cpu.Info()
	percent, _ := cpu.Percent(0, false)
	cpuDic["percent"] = str.Round(percent[0], 2)
	cpuDic["cpuNum"], _ = cpu.Counts(false)

	//服务器磁盘信息
	disklist := make([]disk.UsageStat, 0)
	//所有分区
	var diskTotal, diskUsed, diskUsedPercent float64
	diskInfo, err := disk.Partitions(true)
	if err == nil {
		for _, p := range diskInfo {
			diskDetail, err := disk.Usage(p.Mountpoint)
			if err == nil {
				diskDetail.UsedPercent, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", diskDetail.UsedPercent), 64)
				diskDetail.Total = diskDetail.Total / 1024 / 1024
				diskDetail.Used = diskDetail.Used / 1024 / 1024
				diskDetail.Free = diskDetail.Free / 1024 / 1024
				disklist = append(disklist, *diskDetail)

			}
		}
	}

	d, _ := disk.Usage("/")

	diskTotal = float64(d.Total / GB)
	diskUsed = float64(d.Used / GB)
	diskUsedPercent, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", d.UsedPercent), 64)

	diskDic := make(map[string]interface{}, 0)
	diskDic["total"] = diskTotal
	diskDic["used"] = diskUsed
	diskDic["percent"] = diskUsedPercent

	bootTime, _ := host.BootTime()
	cachedBootTime = time.Unix(int64(bootTime), 0)

	TrackNetworkSpeed()
	netDic := make(map[string]interface{}, 0)
	netDic["in"] = str.Round(float64(netInSpeed/KB), 2)
	netDic["out"] = str.Round(float64(netOutSpeed/KB), 2)

	serveMsgStr, _ := json.Marshal(gin.H{
		"os":       osDic,
		"mem":      memDic,
		"cpu":      cpuDic,
		"disk":     diskDic,
		"net":      netDic,
		"swap":     swapDic,
		"location": "Aliyun",
		"bootTime": GetHourDiffer(cachedBootTime.Format("2006-01-02 15:04:05"), time.Now().Format("2006-01-02 15:04:05")),
	})
	return serveMsgStr
}

func TrackNetworkSpeed() {
	var innerNetInTransfer, innerNetOutTransfer uint64
	nc, err := net.IOCounters(true)
	if err == nil {
		for _, v := range nc {
			if isListContainsStr(excludeNetInterfaces, v.Name) {
				continue
			}
			innerNetInTransfer += v.BytesRecv
			innerNetOutTransfer += v.BytesSent
		}
		now := uint64(time.Now().Unix())
		diff := now - lastUpdateNetStats
		if diff > 0 {
			netInSpeed = (innerNetInTransfer - netInTransfer) / diff
			netOutSpeed = (innerNetOutTransfer - netOutTransfer) / diff
		}
		netInTransfer = innerNetInTransfer
		netOutTransfer = innerNetOutTransfer
		lastUpdateNetStats = now
	}
}

func isListContainsStr(list []string, str string) bool {
	for i := 0; i < len(list); i++ {
		if strings.Contains(str, list[i]) {
			return true
		}
	}
	return false
}

//DownloadLog log文件压缩下载
func (ctrl *ServerController) DownloadLog(c *gin.Context) {
	date := c.DefaultQuery("date", time.Now().Format("2006-01-02"))

	path, tarballName := "storage/logs/", "storage/logs/"+date+".tar.gz"
	// 读取目录下数据和规则文件
	files, err := ioutil.ReadDir(path)
	logger.LogIf(err)

	// 打包
	fw, err := os.Create(tarballName)
	logger.LogIf(err)

	gw := gzip.NewWriter(fw)
	tw := tar.NewWriter(gw)

	defer func(fw *os.File, gw *gzip.Writer, tw *tar.Writer, name string) {
		_ = tw.Close()
		_ = gw.Close()
		_ = fw.Close()
		//需文件生成后才能下载
		c.File(tarballName)
		_ = os.Remove(name)
	}(fw, gw, tw, tarballName)

	for _, f := range files {
		//筛选时间
		if !strings.Contains(f.Name(), date) {
			continue
		}

		hdr := &tar.Header{
			Name: f.Name(),
			Mode: 0600,
			Size: f.Size(),
		}

		if err = tw.WriteHeader(hdr); err != nil {
			logger.Error(err.Error())
		}

		pa := path + f.Name()
		tf, err := ioutil.ReadFile(pa)
		if err != nil {
			logger.Error(err.Error())
		}
		if _, err = tw.Write(tf); err != nil {
			logger.Error(err.Error())
		}
		err = tw.Flush()
		logger.LogIf(err)
	}

	//设置文件类型
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Transfer-Encoding", "binary")
	//设置文件名称File Transfer
	c.Header("Content-Disposition", "attachment; filename="+date+".tar.gz")

	return
}