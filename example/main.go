package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/ainiaa/catutil/gincat"
	"github.com/ainiaa/catutil/gormcat"
	"github.com/ainiaa/catutil/mongocat"
	"github.com/ainiaa/catutil/rediscat"

	gredis "github.com/go-redis/redis/v7"
)

func initGinCat() {
	r := gin.New()
	//健康检测
	r.HEAD("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"code":    0,
			"message": "",
			"data":    nil,
		})
	})
	address := ":8080"
	r.Use(gincat.Cat())

	srv := http.Server{
		Addr:    address,
		Handler: r,
	}
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("listen:%s error:%s", address, err.Error())
	}
}

func initMongoCat() {

}

func initRedisSingleCat(c *Config) {
	client := gredis.NewClient(&gredis.Options{
		Network:            c.Alone.Network,
		Addr:               c.Alone.Addr,
		Username:           c.Alone.Username,
		Password:           c.Alone.Password,
		DB:                 c.Alone.DB,
		MaxRetries:         c.Alone.MaxRetries,
		MinRetryBackoff:    time.Duration(c.Alone.MinRetryBackoff) * time.Millisecond,
		MaxRetryBackoff:    time.Duration(c.Alone.MaxRetryBackoff) * time.Millisecond,
		DialTimeout:        time.Duration(c.Alone.DialTimeout) * time.Millisecond,
		ReadTimeout:        time.Duration(c.Alone.ReadTimeout) * time.Millisecond,
		WriteTimeout:       time.Duration(c.Alone.WriteTimeout) * time.Millisecond,
		PoolSize:           c.Alone.PoolSize,
		MinIdleConns:       c.Alone.MinIdleConns,
		MaxConnAge:         time.Duration(c.Alone.MaxConnAge) * time.Millisecond,
		PoolTimeout:        time.Duration(c.Alone.PoolTimeout) * time.Millisecond,
		IdleTimeout:        time.Duration(c.Alone.IdleTimeout) * time.Millisecond,
		IdleCheckFrequency: time.Duration(c.Alone.IdleCheckFrequency) * time.Millisecond,
	})
	pong, err := client.Ping().Result()
	if pong != "PONG" || err != nil {
		panic(fmt.Sprintf("alone redis conn error: %s", err))
	}
	client.AddHook(rediscat.RedisTraceHook{})
}

func initRedisClusterCat(c *Config) {
	client := gredis.NewClusterClient(&gredis.ClusterOptions{
		Addrs:              c.Cluster.Addrs,
		MaxRedirects:       c.Cluster.MaxRedirects,
		ReadOnly:           c.Cluster.ReadOnly,
		RouteByLatency:     c.Cluster.RouteByLatency,
		RouteRandomly:      c.Cluster.RouteRandomly,
		Username:           c.Cluster.Username,
		Password:           c.Cluster.Password,
		MaxRetries:         c.Cluster.MaxRetries,
		MinRetryBackoff:    time.Duration(c.Cluster.MinRetryBackoff) * time.Millisecond,
		MaxRetryBackoff:    time.Duration(c.Cluster.MaxRetryBackoff) * time.Millisecond,
		DialTimeout:        time.Duration(c.Cluster.DialTimeout) * time.Millisecond,
		ReadTimeout:        time.Duration(c.Cluster.ReadTimeout) * time.Millisecond,
		WriteTimeout:       time.Duration(c.Cluster.WriteTimeout) * time.Millisecond,
		PoolSize:           c.Cluster.PoolSize,
		MinIdleConns:       c.Cluster.MinIdleConns,
		MaxConnAge:         time.Duration(c.Cluster.MaxConnAge) * time.Millisecond,
		PoolTimeout:        time.Duration(c.Cluster.PoolTimeout) * time.Millisecond,
		IdleTimeout:        time.Duration(c.Cluster.IdleTimeout) * time.Millisecond,
		IdleCheckFrequency: time.Duration(c.Cluster.IdleCheckFrequency) * time.Millisecond,
	})
	pong, err := client.Ping().Result()
	if pong != "PONG" || err != nil {
		panic(fmt.Sprintf("cluster redis conn error: %s", err))
	}
	client.AddHook(rediscat.RedisTraceHook{})
}

var mgdb *mongo.Database
func initMongo(c *MongoConf) *mongo.Database {
	if mgdb == nil {
		client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(c.Uri))

		if err != nil {
			panic(fmt.Sprintf("mongo Connect error %s", err.Error()))
		}
		if e := client.Ping(context.Background(), nil); e != nil {
			panic(fmt.Sprintf("mongo ping error %s", e.Error()))
		}

		mgdb = client.Database(c.Database)
	}
	return mgdb
}

func GetUserCollection(collectionName string) mongocat.ConnHandler {
	return mongocat.WithCat(mgdb.Collection(collectionName))
}

func initGormCat() {
	dns := "" // mysql dns
	masterDialect := mysql.Open(dns)
	orm, err := gorm.Open(masterDialect, &gorm.Config{})
	if err != nil {
		panic("orm conn err")
	}
	gormcat.AddGormCallbacks(orm)
}

func main() {
	initGinCat()

	c := &Config{}

	initRedisClusterCat(c)

	initRedisSingleCat(c)

	initGormCat()

	// mongo
	m := &MongoConf{}
	initMongo(m)
	id :=""
	objId, _ := primitive.ObjectIDFromHex(id)
	GetUserCollection("user").FindOne(context.Background(),bson.D{{"_id", objId}},)

}