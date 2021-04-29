package httpcat

import (
	"net"
	"net/http"
	"time"
)

type HttpClientConf struct {
	ConnTimeoutMS       int  `json:"conn_timeout_ms"`         // 写超时
	ReadWriteTimeoutMS  int  `json:"read_write_timeout_ms"`   // 读超时
	KeepAliveTime       int  `json:"keep_alive_time"`         // keepalive 时间
	DisableKeepAlives   bool `json:"disable_keep_alives"`     // 是否启用 keepalive
	MaxIdleConnsPerHost int  `json:"max_idle_conns_per_host"` //连接池大小
	RetryTimes          int  `json:"retry_times"`             // 重试次数
	BodyMaxLength       int  `json:"body_max_length"`         // 接收 body 的 最大长度
	BodyReadTimeoutMS   int  `json:"body_read_timeout_ms"`    //  body 最大读取超时时间
}

var c HttpClientConf

func init() {
	c = HttpClientConf{
		ConnTimeoutMS:       100,
		ReadWriteTimeoutMS:  200,
		KeepAliveTime:       120,
		DisableKeepAlives:   false,
		MaxIdleConnsPerHost: 10,
		RetryTimes:          0,
		BodyMaxLength:       0,
		BodyReadTimeoutMS:   0,
	}
}
func NewHttpClientWrapper() *http.Client {
	return NewHttpClientWithTimeout(0, 0)
}

//NewHttpClientWithTimeout 生成带有timeout的 http client
//connTimeout time.Duration 连接超时时间
//readWriteTimeout time.Duration 读写超时时间
func NewHttpClientWithTimeout(connTimeoutMS, readWriteTimeoutMS int) *http.Client {

	cc := c
	cc.ConnTimeoutMS = connTimeoutMS
	cc.ReadWriteTimeoutMS = readWriteTimeoutMS

	return NewHttpClient(cc)
}

//NewHttpClientWithTimeout 生成带有timeout的 http client
//connTimeout time.Duration 连接超时时间
//readWriteTimeout time.Duration 读写超时时间
func NewHttpClient(cc HttpClientConf) *http.Client {
	var connTimeoutMS int
	var readWriteTimeoutMS int
	var keepAliveTime int
	var maxIdleConnsPerHost int
	if cc.ConnTimeoutMS == 0 {
		connTimeoutMS = c.ConnTimeoutMS
	}
	if cc.ReadWriteTimeoutMS == 0 {
		readWriteTimeoutMS = c.ReadWriteTimeoutMS
	}
	if cc.KeepAliveTime == 0 {
		keepAliveTime = c.KeepAliveTime
	}
	if cc.MaxIdleConnsPerHost == 0 {
		maxIdleConnsPerHost = c.MaxIdleConnsPerHost
	}
	return &http.Client{
		Timeout: time.Duration(readWriteTimeoutMS) * time.Millisecond,
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   time.Duration(connTimeoutMS) * time.Millisecond,
				KeepAlive: time.Duration(keepAliveTime) * time.Second,
			}).DialContext,
			DisableKeepAlives:   cc.DisableKeepAlives,
			MaxIdleConnsPerHost: maxIdleConnsPerHost,
		},
	}
}
