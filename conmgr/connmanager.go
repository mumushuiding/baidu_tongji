package conmgr

import (
	"bufio"
	"errors"
	"log"
	"math"
	"os"
	"strings"
	"sync/atomic"
	"time"

	"github.com/mumushuiding/baidu_tongji/model"
	"github.com/mumushuiding/baidu_tongji/service"
	"github.com/mumushuiding/util"
)

// Conmgr 程序唯一的一个连接管理器
// 处理定时任务
var Conmgr *ConnManager

var numbersPertime = 200

// 消息类型

const (
	// MSGURL url
	MSGURL = "URL数组"
)

// ConnManager 连接管理器
type ConnManager struct {
	start int32
	stop  int32
	quit  chan struct{}
	// 消息通道
	msgchan chan interface{}
	// 纪录通道
	recordchan chan model.BaiduRecord
}
type chanStruct struct {
	api model.BaiduAPI
}

// Start 启动连接管理器
func (cm *ConnManager) Start() {
	// 是否已经启动
	if atomic.AddInt32(&cm.start, 1) != 1 {
		return
	}
	log.Println("启动连接管理器")
	go cm.chanHandler()
	// 定时任务
	go cronTaskStart(cm)
}

// Stop 停止连接管理器
func (cm *ConnManager) Stop() {
	if atomic.AddInt32(&cm.stop, 1) != 1 {
		log.Println("连接管理器已经关闭")
		return
	}
	close(cm.quit)
	defer close(cm.msgchan)
	defer close(cm.recordchan)

	log.Println("关闭连接管理器")
}

// New 新建一个连接管理器
func New() {
	cm := ConnManager{
		quit:       make(chan struct{}),
		msgchan:    make(chan interface{}, 50),
		recordchan: make(chan model.BaiduRecord, 20),
	}
	Conmgr = &cm
	Conmgr.Start()
}

// // MsgReq MsgReq
// type MsgReq struct {
// 	Type string
// 	Data interface{}
// }

// // ToString ToString
// func (m *MsgReq) ToString() string {
// 	str, _ := util.ToJSONStr(m)
// 	return str
// }

// 处理通道 msgchan 和 recordchan 中的数据
func (cm *ConnManager) chanHandler() {
out:
	for {
		select {
		case req := <-cm.msgchan:
			switch msg := req.(type) {
			case *model.FzManuscript:
				// 存储PageId同编辑的对应关系
				go func() {
					editor := msg.Transform2URLEditor()
					err := editor.SaveOrUpdate()
					if err != nil {
						sendRecord(cm, "保存BaiduURLEditor失败", editor.ToString(), 0, err)
					}
				}()
			case *model.BaiduData:
				go func() {
					ufs := msg.Transform2URLFlow()
					// URL数组
					for _, uf := range ufs {
						select {
						case cm.msgchan <- uf:
							// log.Println("发送数据:", uf.Name)
						case <-cm.quit:
							return
						}
					}
					// URL对应PageId对象数组
					urls := msg.GetBItems()
					for _, url := range urls {
						select {
						case cm.msgchan <- &url:
						case <-cm.quit:
							return
						}
					}

				}()
			case *model.BItems:
				go func() {
					// 查询PageId对应的编辑
					urls := strings.Split(msg.Name, "?")
					paper, err := service.FindFzManuscriptByURL(urls[0])
					if err != nil {
						sendRecord(cm, "查询FzManuscript稿件失败", msg.Name, 0, err)
						return
					}
					editor := paper.Transform2URLEditor()
					editor.PageID = msg.PageID
					select {
					case cm.msgchan <- &editor:
					case <-cm.quit:
						return
					}
				}()
			case *model.BaiduURLEditor:
				go func() {
					err := msg.SaveOrUpdate()
					if err != nil {
						sendRecord(cm, "存储BaiduURLEditor失败", msg.ToString(), 0, err)
					}
				}()
			case *model.BaiduURLFlow:
				go func() {
					err := msg.SaveOrUpdate()
					if err != nil {
						sendRecord(cm, "保存受访页面流量失败", msg.ToString(), 0, err)
					}
				}()
			}

		case record := <-cm.recordchan:
			go record.Save()
		case <-cm.quit:
			break out
		}
	}
}

// cronTaskStart 启动定时任务
func cronTaskStart(cm *ConnManager) {
	log.Println("启动定时任务")
out:
	for {
		now := time.Now()
		next := now.Add(time.Hour * 24)
		next = time.Date(next.Year(), next.Month(), next.Day(), 0, 0, 0, 0, next.Location())
		// next := now.Add(time.Second * 10)
		t := time.NewTimer(next.Sub(now))
		select {
		// 连接管理器终止时退出
		case <-cm.quit:
			break out
		case <-t.C:
			// 执行定时任务
			go GetBaiduData(cm)
		}
	}
}
func generateDatesArrByTimeSpan(startDate, endDate string) ([]string, error) {
	// 先转化成time
	startTime, err := util.ParseDate1(startDate)
	if err != nil {
		return nil, err
	}
	endTime, err := util.ParseDate1(endDate)
	if err != nil {
		return nil, err
	}
	// 确定大小
	today := time.Now()
	today = time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, today.Location())
	if startTime.Sub(today) > 0 || endTime.Sub(today) > 0 {
		return nil, errors.New("不能大于当前日期")
	}
	if endTime.Sub(startTime) < 0 {
		return nil, errors.New("开始日期startDate不能大于结束日期endDate")
	}
	sub := int(math.Ceil(endTime.Sub(startTime).Hours() / 24))
	// 生成数组
	var res []string
	res = append(res, startDate)
	for i := 1; i <= sub; i++ {
		res = append(res, util.FormatDate1(time.Date(startTime.Year(), startTime.Month(), startTime.Day()+i, 0, 0, 0, 0, startTime.Location())))
	}
	return res, nil
}

// GetBaiduDataByTimeSpan 根据时间区间获取数据
func GetBaiduDataByTimeSpan(startDate, endDate string) error {
	// 生成时间区间
	dates, err := generateDatesArrByTimeSpan(startDate, endDate)
	if err != nil {
		return err
	}
	cm := Conmgr
	apis, err := service.GetTaskFromAPI()
	if err != nil {
		return err
	}
	for _, api := range apis {
		err = util.Str2Struct(api.Par, &api.Params)
		// fmt.Println("par:", api.Par)
		if err != nil {
			sendRecord(cm, "GetURLFlow", "无法将参数：["+api.Par+"]转换成Params对象", 0, err)
		}
		switch api.Type {
		case "受访页面":
			for _, date := range dates {
				go GetURLFlowByDate(date, *api, cm)
			}
			break
		}
	}
	return nil
}

// GetBaiduData GetBaiduData
func GetBaiduData(cm *ConnManager) {
	// 从api表获取定时任务
	apis, err := service.GetTaskFromAPI()
	if err != nil {
		panic(err)
	}
	for _, api := range apis {
		err = util.Str2Struct(api.Par, &api.Params)
		// fmt.Println("par:", api.Par)
		if err != nil {
			sendRecord(cm, "GetURLFlow", "无法将参数：["+api.Par+"]转换成Params对象", 0, err)
		}
		switch api.Type {
		case "受访页面":
			go GetURLFlow(api, cm)
			break
		}
	}
}

// GetURLFlow 获取受访页面流量数据
func GetURLFlow(api *model.BaiduAPI, conmgr *ConnManager) {

	// 查询上次查询的日期
	dates, err := service.GetDatesToFind(api.Type)
	log.Println("dates:", dates)
	if err != nil {
		sendRecord(conmgr, "GetURLFlow", "查询最后查询的日期", 0, err)
	}
	for _, date := range dates {
		go GetURLFlowByDate(date, *api, conmgr)
	}

}

// GetURLFlowByDate 根据日期获取受访页面流量
func GetURLFlowByDate(date string, api model.BaiduAPI, conmgr *ConnManager) {

	// 设置查询参数
	api.Params.Body.StartDate = date
	api.Params.Body.EndDate = date
	api.Params.Body.StartIndex = 0
	api.Params.Body.MaxResults = numbersPertime
	// 先查询本次查询的总条数
	total, err := service.GetTotalNumOfTargetAPI(api.URL, api.Params2Str())
	if err != nil {
		sendRecord(conmgr, "service.GetTotalNumOfTargetAPI", "查询总条数", 0, err)
	}
	// 确定查询线程数
	// log.Println("总条数：", total)
	processnum := int(math.Ceil(float64(total) / float64(numbersPertime)))
	// log.Println("线程数:", processnum)

	// 多线程查询
	for processIndex := 0; processIndex < processnum; processIndex++ {
		go func(i int, api model.BaiduAPI) {
			// 查询i~i+99的数据

			api.Params.Body.StartIndex = i * numbersPertime

			// log.Println("查询参数:", api.Params2Str())
			// 获取百度统计数据
			baidudata, err := service.GetDatasOfTargetAPI(api.URL, api.Params2Str())
			// log.Println("百度数据:", baidudata.ToString())
			// go writeToFile(fmt.Sprintf("test/output%d-%d-%d-%s", i, api.Params.Body.StartIndex, api.Params.Body.MaxResults, api.Params.Body.StartDate), "")
			if err != nil {
				sendRecord(conmgr, api.Type, api.ToString(), 0, err)
			}

			select {
			case conmgr.msgchan <- baidudata:
			case <-conmgr.quit:
				return
			}
		}(processIndex, api)
	}

}

// sendRecord 发送纪录
func sendRecord(conmgr *ConnManager, typename, data string, flag uint8, err error) {
	record := model.BaiduRecord{
		Data: data,
		Type: typename,
		Flag: flag,
		Err:  err.Error(),
	}
	select {
	case conmgr.recordchan <- record:
	case <-conmgr.quit:
		return
	}
}

func writeToFile(filepath string, content string) {
	f, _ := os.Create(filepath)
	defer f.Close()
	w := bufio.NewWriter(f)
	w.WriteString(content)
	w.Flush()
	f.Close()

}
