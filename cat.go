package catutil

import (
	"context"

	"github.com/cat-go/cat"
	"github.com/cat-go/cat/message"
	"github.com/gin-gonic/gin"
)

func GetRootTran(ctx context.Context) message.Transactor {
	if ctx == nil || !cat.IsEnabled() {
		return nil
	}
	if rootTranRaw := ctx.Value(CatCtxRootTran); rootTranRaw != nil {
		if rootTran, ok := rootTranRaw.(message.Transactor); ok && rootTran != nil {
			return rootTran
		}
	}
	return nil
}

func GetRequestId(c context.Context) string {
	ctx, ok := c.(*gin.Context)
	if !ok {
		return ""
	}
	return ctx.GetString(ContextKeyRequestId)
}

func SetRequestId(c context.Context, requestId string) {
	ctx, ok := c.(*gin.Context)
	if !ok {
		return
	}
	ctx.Set(ContextKeyRequestId, requestId)
}



func SetTraceId(c *gin.Context, tran message.Transactor) {
	var root, parent, child string
	root = c.Request.Header.Get(cat.RootId)
	parent = c.Request.Header.Get(cat.ParentId)
	child = c.Request.Header.Get(cat.ChildId)
	if root == "" {
		root = cat.MessageId()
	}
	if parent == "" {
		parent = root
	}
	if child == "" {
		child = cat.MessageId()
	}
	tran.SetRootMessageId(root)
	tran.SetParentMessageId(parent)
	tran.SetMessageId(child)
}

func GetTraceId(ctx context.Context) string {
	if ctx, ok := ctx.(*gin.Context); ok {
		if tran, ok := ctx.Value(CatCtxRootTran).(message.Transactor); ok {
			return tran.GetRootMessageId()
		}
	}
	return ""
}

