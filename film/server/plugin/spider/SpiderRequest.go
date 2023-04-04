package spider

import (
	"fmt"
	"github.com/gocolly/colly/v2"
	"log"
	"net/http"
	"net/url"
	"time"
)

var (
	Client = CreateClient()
)

// RequestInfo 请求参数结构体
type RequestInfo struct {
	Uri    string      `json:"uri"`    // 请求url地址
	Params url.Values  `json:"params"` // 请求参数
	Header http.Header `json:"header"` // 请求头数据
	Resp   []byte      `json:"resp"`   // 响应结果数据
}

// CreateClient 初始化请求客户端
func CreateClient() *colly.Collector {
	c := colly.NewCollector()
	// 设置代理信息
	//if proxy, err := proxy.RoundRobinProxySwitcher("127.0.0.1:7890"); err != nil {
	//	c.SetProxyFunc(proxy)
	//}
	// 设置并发数量控制
	//c.Async = true
	// 访问深度
	c.MaxDepth = 1
	//可重复访问
	c.AllowURLRevisit = true
	// 设置超时时间 默认10s
	c.SetRequestTimeout(20 * time.Second)
	// 发起请求之前会调用的方法
	c.OnRequest(func(request *colly.Request) {
		// 设置一些请求头信息
		request.Headers.Set("Content-Type", "application/json;charset=UTF-8")
		request.Headers.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/111.0.0.0 Safari/537.36")
		//request.Headers.Set("cookie", "ge_ua_key=sxo%2Bz4kkS7clWpEtg2m7HioRfIo%3D")
		request.Headers.Set("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	})
	// 请求期间报错的回调
	c.OnError(func(response *colly.Response, err error) {
		log.Printf("请求异常: URL: %s Error: %s\n", response.Request.URL, err)
	})
	return c
}

// ApiGet 请求数据的方法
func ApiGet(r *RequestInfo) {
	// 请求成功后的响应
	Client.OnResponse(func(response *colly.Response) {
		// 将响应结构封装到 RequestInfo.Resp中
		r.Resp = response.Body
		// 拿到response后输出请求url
		//log.Println("\n请求成功: ", response.Request.URL)
	})
	// 处理请求参数
	err := Client.Visit(fmt.Sprintf("%s?%s", r.Uri, r.Params.Encode()))
	if err != nil {
		log.Println("获取数据失败: ", err)
	}
}
