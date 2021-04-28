package main

type Config struct {
	ConnType uint           `json:"conn_type"` // 0：cluster 1: alone 2: sentinel
	Alone    AloneConfig    `json:"alone"`
	Cluster  ClusterConfig  `json:"cluster"`
}

type MongoConf struct {
	Uri      string `json:"uri"`
	Database string `json:"database"`
}