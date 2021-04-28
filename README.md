# catutil
catutil

支持 gorm、mongodb、http 的 cat 监控

golang cat client: github.com/cat-go/cat

PS 使用之前需要先初始化 cat
```golang
type CatConf struct {
	Flag       bool   `json:"flag"`
	AppId      string `json:"app_id"`
	Port       int    `json:"port"`
	HttpPort   int    `json:"http_port"`
	ServerAddr string `json:"server_addr"`
	IsDebug    bool   `json:"is_debug"`
}

func initCat(c *CatConf) {
	if c.Flag {
		if c.Cat.IsDebug {
			cat.DebugOn()
		}
		cat.Init(&cat.Options{
			AppId:      c.AppId,
			Port:       c.Port,
			HttpPort:   c.HttpPort,
			ServerAddr: c.ServerAddr,
		})
	}
}
```