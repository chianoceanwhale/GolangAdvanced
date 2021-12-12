// #1、使用 redis benchmark 工具, 测试 10 20 50 100 200 1k 5k 字节 value 大小，redis get set 性能。
// #命令:
// redis-benchmark -d 10 -t get,set
// #SET

// -	执行次数和耗时	每秒请求次数
// 10	100000 requests completed in 3.78 seconds	26420.06 requests per second
// 20	100000 requests completed in 3.77 seconds	26454.03 requests per second
// 50	100000 requests completed in 3.92 seconds	25545.31 requests per second
// 100	100000 requests completed in 3.95 seconds	25238.26 requests per second
// 200	100000 requests completed in 3.79 seconds	26388.12 requests per second
// 1k	100000 requests completed in 3.95 seconds	25054.38 requests per second
// 5k	100000 requests completed in 4.94 seconds	20250.63 requests per second
// #GET		
// -	执行次数和耗时	每秒请求次数
// ----	----	----
// 10	100000 requests completed in 3.72 seconds	26733.97 requests per second
// 20	100000 requests completed in 3.94 seconds	25125.63 requests per second
// 50	100000 requests completed in 3.82 seconds	25863.53 requests per second
// 100	100000 requests completed in 3.87 seconds	25704.94 requests per second
// 200	100000 requests completed in 3.89 seconds	25795.79 requests per second
// 1k	100000 requests completed in 4.34 seconds	25026.03 requests per second
// 5k	100000 requests completed in 4.76 seconds	20857.41 requests per second
// 100000 requests completed in 3.79 seconds
// 26419.08 requests per second


#2、写入一定量的 kv 数据, 根据数据大小 1w-50w 自己评估, 结合写入前后的 info memory 信息 ,
分析上述不同 value 大小下，平均每个 key 的占用内存空间。 ###代码链接github ###代码工作原理，
写入不同数量不同长度的value, 分析内存占用, 导出结果到csv文件 ###结论 相同长度的value在写入数量越多情况下，
平均每个value占用内存更多
package main

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/hhxsv5/go-redis-memory-analysis"
)

var client redis.UniversalClient
var ctx context.Context

const (
	ip   string = "127.0.0.1"
	port uint16 = 6380
)

func init() {
	client = redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%v:%v", ip, port),
		Password:     "",
		DB:           0,
		PoolSize:     128,
		MinIdleConns: 100,
		MaxRetries:   5,
	})

	ctx = context.Background()
}

func main() {
  var dataMap = {
    10000："len10k",
    50000："len50k",
    500000："len500k",
  }
	var dataSlice = [10,1000,5000]
  for  _,v := range dataSlice{
    for key,value := range dataMap{
      write(key,value,generateValue(v))
    }
  }

	analysis()
}

func write(num int, key, value string) {
	for i := 0; i < num; i++ {
		k := fmt.Sprintf("%s:%v", key, i)
		cmd := client.Set(ctx, k, value, -1)
		err := cmd.Err()
		if err != nil {
			fmt.Println(cmd.String())
		}
	}
}

func analysis() {
	analysis, err := gorma.NewAnalysisConnection(ip, port, "")
	if err != nil {
		fmt.Println("something wrong:", err)
		return
	}
	defer analysis.Close()

	analysis.Start([]string{":"})

	err = analysis.SaveReports("./reports")
	if err == nil {
		fmt.Println("done")
	} else {
		fmt.Println("error:", err)
	}
}

func generateValue(size int) string {
	arr := make([]byte, size)
	for i := 0; i < size; i++ {
		arr[i] = 'a'
	}
	return string(arr)
}
