package conmgr

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"math"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/mumushuiding/baidu_tongji/model"
	"github.com/mumushuiding/baidu_tongji/service"
	"github.com/mumushuiding/util"
)

// Conmgr 程序唯一的一个连接管理器
// 处理定时任务
var Conmgr *ConnManager

// 值不可大，否则会因为处理来不及，导致占满内存
var numbersPertime = 500

// 消息类型

const (
	// MSGURL url
	MSGURL = "URL数组"
	// EditorFlowAndManuscriptNumLastMonth EditorFlowAndManuscriptNumLastMonth
	EditorFlowAndManuscriptNumLastMonth = "上月编辑流量和稿件"
)

// ConnManager 连接管理器
type ConnManager struct {
	start int32
	stop  int32
	quit  chan struct{}
	// 消息通道
	msgchan chan interface{}
	// 受访页面
	editorchan     chan interface{}
	editorSaveLock sync.Mutex

	// 纪录通道
	recordchan chan model.BaiduRecord

	// 稿件通道
	fzManuscriptchan         chan interface{}
	fzManuscriptNotFound     map[string]string
	fzManuscriptNotFoundLock sync.Mutex
	cacheMap                 map[string]interface{}
	cacheMapLock             sync.RWMutex
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
	defer close(cm.editorchan)
	defer close(cm.fzManuscriptchan)
	log.Println("关闭连接管理器")
}

// New 新建一个连接管理器
func New() {
	cm := ConnManager{
		quit:                 make(chan struct{}),
		msgchan:              make(chan interface{}, 50),
		recordchan:           make(chan model.BaiduRecord, 20),
		editorchan:           make(chan interface{}, 100),
		fzManuscriptchan:     make(chan interface{}, 100),
		fzManuscriptNotFound: make(map[string]string),
		cacheMap:             make(map[string]interface{}),
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

			case *model.BaiduData:
				go func() {
					if msg.Header.Status == 3 {
						sendRecord(cm, "百度API报错", msg.ToString(), 0, nil)
						return
					}
					if msg == nil || len(msg.Body.Data) == 0 {
						return
					}
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
					select {
					case cm.editorchan <- urls:
					case <-cm.quit:
						return
					}

				}()
			case *model.BaiduURLFlow:
				go func() {
					msg.TimeSpan = strings.ReplaceAll(msg.TimeSpan, "/", "-")
					// err := msg.SaveOrUpdate()
					if len(msg.Name) > 0 {
						msg.Name = strings.Split(msg.Name, "?")[0]
					}

					err := msg.Save()
					if err != nil && !strings.Contains(err.Error(), "Error 1062: Duplicate entry") {
						sendRecord(cm, "保存受访页面流量失败", msg.ToString(), 0, err)
					}
				}()
			}
		case req1 := <-cm.editorchan:
			switch msg1 := req1.(type) {
			case []model.BItems:
				go func() {
					for _, url := range msg1 {
						select {
						case cm.editorchan <- url:
						case <-cm.quit:
							return
						}
					}
				}()
			case model.BItems:
				go func() {
					// 查询PageId对应的编辑
					urls := strings.Split(msg1.Name, "?")
					if len(urls) == 0 {
						return
					}
					paper, err := service.FindFzManuscriptByURLFromLocalDB(urls[0])
					if err != nil {
						select {
						case cm.fzManuscriptchan <- model.FzManuscriptNotFound{Filename: urls[0], PageID: msg1.PageID}:
						case <-cm.quit:
							return
						}
					} else {
						editor := paper.Transform2URLEditor()
						editor.PageID = msg1.PageID
						select {
						case cm.editorchan <- &editor:
						case <-cm.quit:
							return
						}
					}

				}()
			case *model.BaiduURLEditor:
				go func() {
					// log.Println("editor:", msg1.ToString())
					// cm.editorSaveLock.Lock()
					// defer cm.editorSaveLock.Unlock()
					// err := msg1.SaveOrUpdate()
					err := msg1.Save()
					if err != nil && !strings.Contains(err.Error(), "Error 1062: Duplicate entry") {
						sendRecord(cm, "存储BaiduURLEditor失败", msg1.ToString(), 0, err)
					}
				}()
			}
		case req2 := <-cm.fzManuscriptchan:
			switch msg2 := req2.(type) {
			case model.FzManuscriptNotFound:
				go func() {
					// cm.fzManuscriptNotFound[msg2.Filename] = msg2.PageID

					result, err := service.FindFzManuscriptByURLFromDBNews(msg2.Filename)
					if err != nil {
						sendRecord(cm, "远程查询FzManuscript失败", msg2.Filename, 0, err)
						return
					}
					result.PageID = msg2.PageID
					// log.Printf("url: %s,result: %s\n", result.Filename, result.ToString())
					select {
					case cm.fzManuscriptchan <- result:
					case <-cm.quit:
						return
					}
				}()
			case []*model.FzManuscript:
				go func() {
					for _, v := range msg2 {
						select {
						case cm.fzManuscriptchan <- v:
						case <-cm.quit:
							return
						}
					}
				}()
			case *model.FzManuscript:
				go func() {
					err := msg2.SaveOrUpdate()
					// log.Println("存储:", msg2.ToString())
					if err != nil && !strings.Contains(err.Error(), "Error 1062: Duplicate entry") {
						sendRecord(cm, "存储FzManuscript失败", msg2.ToString(), 0, err)
						return
					}
					if len(msg2.PageID) > 0 {
						editor := model.BaiduURLEditor{
							PageID:   msg2.PageID,
							Username: msg2.Editor,
							Realname: msg2.Editorname,
						}
						select {
						case cm.editorchan <- &editor:
						case <-cm.quit:
							return
						}
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
		next = time.Date(next.Year(), next.Month(), next.Day(), 0, 4, 0, 0, next.Location())
		// next := now.Add(time.Second * 10)
		t := time.NewTimer(next.Sub(now))
		select {
		// 连接管理器终止时退出
		case <-cm.quit:
			break out
		case <-t.C:
			// 刷新缓存表
			go RefreshCacheMap()
			// 先拉取最新稿件
			go cm.GetRemoteFzManuscriptNotHave()
			// 然后暂停10秒，保证存储最新稿件
			time.Sleep(time.Second * 50)
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
			continue
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
		return
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
		return
	}
	// 确定查询线程数
	// log.Println("总条数：", total)
	processnum := int(math.Ceil(float64(total) / float64(numbersPertime)))
	// log.Println("线程数:", processnum)

	// 多线程查询
	for processIndex := 0; processIndex < processnum; processIndex++ {
		// 查询i~i+99的数据

		api.Params.Body.StartIndex = processIndex * numbersPertime

		// log.Println("查询参数:", api.Params2Str())
		// 获取百度统计数据
		baidudata, err := service.GetDatasOfTargetAPI(api.URL, api.Params2Str())
		// log.Println("百度数据:", baidudata.ToString())
		// go writeToFile(fmt.Sprintf("test/output%d-%d-%d-%s", i, api.Params.Body.StartIndex, api.Params.Body.MaxResults, api.Params.Body.StartDate), "")
		if err != nil {
			sendRecord(conmgr, api.Type, api.ToString(), 0, err)
			continue
		}
		baidudata.Body.Source = api.Params.Body.Source
		baidudata.Body.Visitor = api.Params.Body.Visitor
		select {
		case conmgr.msgchan <- baidudata:
		case <-conmgr.quit:
			return
		}
		// 防止增加过快，处理来不及
		time.Sleep(time.Millisecond * time.Duration(numbersPertime*100))
		// go func(i int, api model.BaiduAPI) {
		// 	// 查询i~i+99的数据

		// 	api.Params.Body.StartIndex = i * numbersPertime

		// 	// log.Println("查询参数:", api.Params2Str())
		// 	// 获取百度统计数据
		// 	baidudata, err := service.GetDatasOfTargetAPI(api.URL, api.Params2Str())
		// 	// log.Println("百度数据:", baidudata.ToString())
		// 	// go writeToFile(fmt.Sprintf("test/output%d-%d-%d-%s", i, api.Params.Body.StartIndex, api.Params.Body.MaxResults, api.Params.Body.StartDate), "")
		// 	if err != nil {
		// 		sendRecord(conmgr, api.Type, api.ToString(), 0, err)
		// 	}

		// 	select {
		// 	case conmgr.msgchan <- baidudata:
		// 	case <-conmgr.quit:
		// 		return
		// 	}
		// }(processIndex, api)
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

func (cm *ConnManager) getFzManuscriptNotFoundLength() int {
	cm.fzManuscriptNotFoundLock.Lock()
	defer cm.fzManuscriptNotFoundLock.Unlock()
	return len(cm.fzManuscriptNotFound)
}

// GetRemoteFzManuscriptNotHave 从远程获取本地没有的稿件
func (cm *ConnManager) GetRemoteFzManuscriptNotHave() {
	// 获取最后一次 publictime 的时间，将其设置成 startDate,要覆盖一次防止遗漏
	// 然后将当前时间设置为endDate
	res, _, err := service.FindFzManuscriptFromLocal(map[string]interface{}{"order": "publictime desc", "max_results": 1})
	if err != nil && err != gorm.ErrRecordNotFound {
		sendRecord(cm, "getRemoteFzManuscriptNotHave", err.Error(), 0, err)
	}

	now := time.Now()

	end := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, now.Location())

	pub := now
	if res != nil && len(res) > 0 {
		// str, _ := util.ToJSONStr(res)
		// log.Println("res:", str)
		pub = res[0].Publictime
	}
	// log.Println("pub", pub)
	start := time.Date(pub.Year(), pub.Month(), pub.Day(), 0, 0, 0, 0, pub.Location())

	fields := make(map[string]interface{})

	where := fmt.Sprintf("publictime>='%s' and publictime<='%s'", util.FormatDate3(start), util.FormatDate3(end))

	fields["where"] = where
	fields["order"] = "publictime desc"
	fields["max_results"] = 1
	// 计算条数
	_, count, err := service.FindFzManuscriptFromLocal(fields)
	// 每200条一个线程，并发查询最后将结果发送到fzManuscriptchan通道
	processnum := int(math.Ceil(float64(count) / float64(200)))
	fields["max_results"] = 200
	for i := 0; i < processnum; i++ {
		go func(index int, fields map[string]interface{}) {
			fields["start_index"] = index * 200
			result, _, err := service.FindFzManuscriptFromDBNews(fields)
			if err != nil && err != gorm.ErrRecordNotFound {
				str, _ := util.ToJSONStr(fields)
				sendRecord(cm, "getRemoteFzManuscriptNotHave", str, 0, err)
			}
			select {
			case <-cm.quit:
			case cm.fzManuscriptchan <- result:
			}
		}(i, fields)
	}

}

// RefreshCacheMap 刷新cacheMap中的内容
func RefreshCacheMap() {
	clearCacheMap()
	// 查询上月编辑流量
	go GetFlowAndManuscriptNumLastMonth(&model.EditorTongji{})
}

// clearCacheMap 清空ClearCacheMap
func clearCacheMap() {
	Conmgr.cacheMapLock.Lock()
	defer Conmgr.cacheMapLock.Unlock()
	len := len(Conmgr.cacheMap)
	if len > 0 {
		//清空 map 的唯一办法就是重新 make 一个新的 map，不用担心垃圾回收的效率，Go语言中的并行垃圾回收效率比写一个清空函数要高效的多。
		Conmgr.cacheMap = map[string]interface{}{}
	}
}

// GetValFromCache 从缓存中获取值
func GetValFromCache(cachename string) interface{} {
	result := Conmgr.cacheMap[cachename]
	return result
}

// SetVal2Cache SetVal2Cache
func SetVal2Cache(cachename string, val interface{}) {
	Conmgr.cacheMapLock.Lock()
	Conmgr.cacheMap[cachename] = val
	Conmgr.cacheMapLock.Unlock()
}

// GetArticleFlowWithAvators 查询文章流量和编辑的头像
func GetArticleFlowWithAvators(e *model.EditorTongji) (string, error) {
	var result = GetValFromCache(e.Body.Metrics)
	if result == nil {
		var err error
		result, err = FindArticleByTimeSpan(e.Body.Metrics)
		if err != nil {
			return "", err
		}
		// 保存到缓存
		SetVal2Cache(e.Body.Metrics, result)
	}
	e.Body.Data = append(e.Body.Data, result)
	return e.ToString(), nil
}

// GetFlowAndManuscriptNumLastMonth 获取上月编辑的流量和稿件量
// 结果会缓存在conmgr.cacheMap中,key为：上月编辑流量和稿件量
func GetFlowAndManuscriptNumLastMonth(e *model.EditorTongji) (string, error) {
	var result = GetValFromCache(EditorFlowAndManuscriptNumLastMonth)
	if result == nil {
		log.Println("从mysql查询编辑上月流量")
		var err error
		start, end := util.GetLastMonthStartAndEnd()
		e.Body.StartDate = util.FormatDate4(start)
		e.Body.EndDate = util.FormatDate4(end)
		err = e.FindFlowAndManuscriptNum()
		if err != nil {
			return "", err
		}
		// 保存到缓存
		SetVal2Cache(EditorFlowAndManuscriptNumLastMonth, e.Body.Data[0])
	} else {
		log.Println("从缓存获取上月编辑的流量")
		e.Body.Data = append(e.Body.Data, result)
	}

	return e.ToString(), nil
}

// FindArticleByTimeSpan 根据时间段查询
// 如: 文章top50-最近30天,
// 如: 文章top50-最近7天,
func FindArticleByTimeSpan(timespan string) ([]*model.EURLFlow, error) {
	var start, end string

	now := time.Now()
	end = util.FormatDate3(now.Add(time.Hour * 24 * (-1)))
	switch timespan {
	case "文章top50-最近30天":
		start = util.FormatDate3(now.Add(time.Hour * 24 * (-31)))
		break
	case "文章top50-最近7天":
		start = util.FormatDate3(now.Add(time.Hour * 24 * (-8)))
		break
	case "文章top50-昨天":
		end = util.FormatDate3(now.Add(time.Hour * 24 * (-1)))
		break
	default:
		return nil, errors.New("metrics参数必须是【文章top50-最近30天,文章top50-最近7天,文章top50-昨天】中的一个")
	}
	result, err := model.FindFzManuscripFlow(map[string]interface{}{"start_date": start, "end_date": end})
	if err != nil {
		return nil, err
	}
	var usernames []string
	for _, v := range result {
		usernames = append(usernames, v.Realname)
	}
	wxs, err := model.FindWeixinStaffAvatarByUsername(usernames)
	if err != nil {
		return nil, err
	}
	avatarMap := make(map[string]string)
	for _, val := range wxs {
		avatarMap[val.Username] = val.Avatar
	}
	for _, flow := range result {
		flow.Avatar = avatarMap[flow.Realname]
	}
	return result, nil

}
