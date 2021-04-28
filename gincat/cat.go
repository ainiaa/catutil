package gincat

import (
	"strconv"

	"github.com/cat-go/cat"
	"github.com/gin-gonic/gin"

	"github.com/ainiaa/catutil"
)

//监控与链路追踪
func Cat() gin.HandlerFunc {
	return func(c *gin.Context) {
		if cat.IsEnabled() {
			tran := cat.NewTransaction(cat.TypeUrl, c.Request.URL.Path)
			defer func() {
				tran.AddData(catutil.TypeUrlStatus, strconv.Itoa(c.Writer.Status()))
				tran.Complete()
			}()
			catutil.SetTraceId(c, tran)
			tran.AddData(cat.TypeUrlMethod, c.Request.Method, c.FullPath())
			tran.AddData(cat.TypeUrlClient, c.ClientIP())
			c.Set(catutil.CatCtxRootTran, tran)
		}
		c.Next()
	}
}

