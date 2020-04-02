package model_test

import (
	"fmt"
	"testing"

	redis "github.com/go-redis/redis"
)

var c client

type client interface {
	Ping() *redis.StatusCmd
	// HExists 判断是否存在
	HExists(key, field string) *redis.BoolCmd
	// HMset 设置值
	HMSet(key string, values ...interface{}) *redis.BoolCmd
	// HMGet 获取值
	HMGet(key string, fields ...string) *redis.SliceCmd
	// HScan 分页查询
	HScan(key string, cursor uint64, match string, count int64) *redis.ScanCmd
	// Sort 排序
	Sort(key string, sort *redis.Sort) *redis.StringSliceCmd
	// SAdd 添加值到集合
	SAdd(key string, members ...interface{}) *redis.IntCmd
	Watch(fn func(*redis.Tx) error, keys ...string) error
	Get(key string) *redis.StringCmd
	HLen(key string) *redis.IntCmd
}
type Pageinfo struct {
	Dm  string `json:"dm"`  // 域名
	URL string `json:"url"` // 网址
	// Keywords 关键词
	Keywords      string `json:"keywords"`
	Description   string `json:"description"`
	Filetype      int8   `json:"filetype"`
	Publishedtype int8   `json:"publishedtype"`
	Pagetype      int8   `json:"pagetype"`
	Catalogs      string `json:"catalogs"`
	Contentid     string `json:"contentid"`
	Publishdate   string `json:"publishdate"`
	Author        string `json:"author"`
	Source        string `json:"source"`
}
type WebFlow struct {
	DM     string `json:"dm"`     // 域名
	URL    string `json:"url"`    // 网址
	PV     int    `json:"pv"`     // 页面浏览量
	IP     int    `json:"ip"`     // 访问ip数
	UV     int    `json:"uv"`     // 独立访问者数
	Visits int    `json:"visits"` // 访问次数(半个小时内多次算一次)
	BR     int    `json:"br"`     // Bounce Rate 跳出率,只访问一次就跳出
}

func init() {
	c = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
	})
}

// func TestHScan(t *testing.T) {
// 	for i := 1; i < 10; i++ {
// 		x := map[string]interface{}{
// 			fmt.Sprintf("F-%d", i): "hello",
// 		}
// 		fmt.Println(x)
// 		c.HMSet("test", x)
// 	}
// 	result := c.HScan("test", 0, "*", 10)
// 	r, _ := result.Val()
// 	for j := 0; j < len(r); j += 2 {
// 		fmt.Printf("key:%s,value:%s\n", r[j], r[j+1])
// 	}
// }
// func TestSAdd(t *testing.T) {
// 	for i := 1; i < 10; i++ {
// 		x := &WebFlow{
// 			PV: i,
// 		}
// 		c.SAdd("test3", x)
// 	}
// 	c.SAdd("test3", 11)
// }

// func TestWatch(t *testing.T) {
// 	const routineCount = 1000
// 	increment := func(key string) error {
// 		txf := func(tx *redis.Tx) error {
// 			n, err := tx.Get(key).Int()
// 			if err != nil && err != redis.Nil {
// 				return err
// 			}
// 			n++
// 			_, err = tx.Pipelined(func(pipe redis.Pipeliner) error {
// 				pipe.Set(key, n, 0)
// 				return nil
// 			})
// 			return err
// 		}
// 		i := 0
// 		j := 0
// 		for {
// 			i++
// 			if i > 2 {
// 				time.Sleep(1 * time.Second)
// 				j++
// 				fmt.Println("沉睡1秒")
// 				if j > 100 {
// 					break
// 				}
// 				i = 0
// 			}
// 			err := c.Watch(txf, key)
// 			if err != redis.TxFailedErr {
// 				return err
// 			}
// 		}
// 		return errors.New("increment reached maximum number of retries")
// 	}
// 	var wg sync.WaitGroup
// 	wg.Add(routineCount)
// 	for i := 0; i < routineCount; i++ {
// 		go func() {
// 			defer wg.Done()
// 			if err := increment("counter3"); err != nil {
// 				fmt.Println("increment error:", err)
// 			}
// 		}()
// 	}
// 	wg.Wait()
// 	fmt.Println("ended with", c.Get("counter3").Val())
// }
func TestTest(t *testing.T) {
	fmt.Println(c.HLen("abc").Err())
}
