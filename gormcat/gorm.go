package gormcat

import (
	"fmt"

	"github.com/cat-go/cat"
	"github.com/cat-go/cat/message"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/ainiaa/catutil"
)


// 添加gormCallback
func AddGormCallbacks(db *gorm.DB) {
	if cat.IsEnabled() {
		callbacks := newCallbacks()
		registerCallbacks(db, "create", callbacks)
		registerCallbacks(db, "query", callbacks)
		registerCallbacks(db, "update", callbacks)
		registerCallbacks(db, "delete", callbacks)
		registerCallbacks(db, "row_query", callbacks)
	}
}

type callbacks struct{}

func newCallbacks() *callbacks {
	return &callbacks{}
}

func (c *callbacks) beforeCreate(scope *gorm.DB)   { c.before(scope) }
func (c *callbacks) afterCreate(scope *gorm.DB)    { c.after(scope, "INSERT") }
func (c *callbacks) beforeQuery(scope *gorm.DB)    { c.before(scope) }
func (c *callbacks) afterQuery(scope *gorm.DB)     { c.after(scope, "SELECT") }
func (c *callbacks) beforeUpdate(scope *gorm.DB)   { c.before(scope) }
func (c *callbacks) afterUpdate(scope *gorm.DB)    { c.after(scope, "UPDATE") }
func (c *callbacks) beforeDelete(scope *gorm.DB)   { c.before(scope) }
func (c *callbacks) afterDelete(scope *gorm.DB)    { c.after(scope, "DELETE") }
func (c *callbacks) beforeRowQuery(scope *gorm.DB) { c.before(scope) }
func (c *callbacks) afterRowQuery(scope *gorm.DB)  { c.after(scope, "QUERY") }

func (c *callbacks) before(scope *gorm.DB) {
	ct := scope.Statement.Context
	ctx, ok := ct.(*gin.Context)
	if !ok {
		return
	}
	tran := cat.NewTransaction(cat.TypeSql, scope.Statement.Dialector.Name())
	ctx.Set(catutil.CatCtxMysqlTran, tran)
	if t, ok := ctx.Get(catutil.CatCtxRootTran); ok {
		if t1, ok := t.(message.Transactor); ok {
			cat.SetChildTraceId(t1, tran)
			t1.AddChild(tran)
		}
	}
}

func (c *callbacks) after(scope *gorm.DB, operation string) {
	ct := scope.Statement.Context
	ctx, ok := ct.(*gin.Context)
	if !ok {
		return
	}
	if tran, ok := ctx.Value(catutil.CatCtxMysqlTran).(message.Transactor); ok {
		tran.LogEvent(cat.TypeSqlOp, operation)
		tran.Complete()
	}
}

func registerCallbacks(db *gorm.DB, name string, c *callbacks) {
	beforeName := fmt.Sprintf("tracing:%v_before", name)
	afterName := fmt.Sprintf("tracing:%v_after", name)
	gormCallbackName := fmt.Sprintf("gorm:%v", name)
	switch name {
	case "create":
		_ = db.Callback().Create().Before(gormCallbackName).Register(beforeName, c.beforeCreate)
		_ = db.Callback().Create().After(gormCallbackName).Register(afterName, c.afterCreate)
	case "query":
		_ = db.Callback().Query().Before(gormCallbackName).Register(beforeName, c.beforeQuery)
		_ = db.Callback().Query().After(gormCallbackName).Register(afterName, c.afterQuery)
	case "update":
		_ = db.Callback().Update().Before(gormCallbackName).Register(beforeName, c.beforeUpdate)
		_ = db.Callback().Update().After(gormCallbackName).Register(afterName, c.afterUpdate)
	case "delete":
		_ = db.Callback().Delete().Before(gormCallbackName).Register(beforeName, c.beforeDelete)
		_ = db.Callback().Delete().After(gormCallbackName).Register(afterName, c.afterDelete)
	case "row_query":
		_ = db.Callback().Row().Before(gormCallbackName).Register(beforeName, c.beforeRowQuery)
		_ = db.Callback().Row().After(gormCallbackName).Register(afterName, c.afterRowQuery)
	}
}
