package rediscat

import (
	"context"
	"strconv"

	"github.com/cat-go/cat"
	"github.com/cat-go/cat/message"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"

	"github.com/ainiaa/catutil/v2"
)

type RedisTraceHook struct {
}

func (t RedisTraceHook) BeforeProcess(ctx context.Context, cmd redis.Cmder) (context.Context, error) {
	if cat.IsEnabled() {
		if ctx == nil {
			ctx = context.Background()
		}
		tran := cat.NewTransaction(cat.TypeRedis, cmd.Name())
		tran.AddData(cat.TypeRedisCmd, cmd.String())
		if c, ok := ctx.(*gin.Context); ok {
			c.Set(catutil.CatCtxRedisTran, tran)
		} else {
			ctx = context.WithValue(ctx, catutil.CatCtxRedisTran, tran)
		}
		if rootTran := catutil.GetRootTran(ctx); rootTran != nil {
			cat.SetChildTraceId(rootTran, tran)
			rootTran.AddChild(tran)
		}
	}
	return ctx, nil
}

func (t RedisTraceHook) AfterProcess(ctx context.Context, cmd redis.Cmder) error {
	if cat.IsEnabled() {
		if tranRaw := ctx.Value(catutil.CatCtxRedisTran); tranRaw != nil {
			if tran, ok := tranRaw.(message.Transactor); ok && tran != nil {
				if cmd != nil && cmd.Err() != nil && cmd.Err() != redis.Nil {
					tran.AddData("err", cmd.Err().Error())
				}
				tran.Complete()
			}
		}
	}
	return nil
}

func (t RedisTraceHook) BeforeProcessPipeline(ctx context.Context, cmds []redis.Cmder) (context.Context, error) {
	if cat.IsEnabled() {
		if ctx == nil {
			ctx = context.Background()
		}
		tran := cat.NewTransaction(cat.TypeRedis, "redis.pipeline")
		for _, cmd := range cmds {
			tran.AddData(cat.TypeRedisCmd, cmd.String())
		}
		if c, ok := ctx.(*gin.Context); ok {
			c.Set(catutil.CatCtxRedisTran, tran)
		} else {
			ctx = context.WithValue(ctx, catutil.CatCtxRedisTran, tran)
		}
		if rootTran := catutil.GetRootTran(ctx); rootTran != nil {
			cat.SetChildTraceId(rootTran, tran)
			rootTran.AddChild(tran)
		}
	}
	return ctx, nil
}

func (t RedisTraceHook) AfterProcessPipeline(ctx context.Context, cmds []redis.Cmder) error {
	if cat.IsEnabled() {
		if tranRaw := ctx.Value(catutil.CatCtxRedisTran); tranRaw != nil {
			if tran, ok := tranRaw.(message.Transactor); ok && tran != nil {
				if len(cmds) > 0 {
					for idx, cmd := range cmds {
						if cmd != nil && cmd.Err() != nil && cmd.Err() != redis.Nil {
							tran.AddData("err"+strconv.Itoa(idx), cmd.Err().Error())
						}
					}
				}
				tran.Complete()
			}
		}
	}
	return nil
}
