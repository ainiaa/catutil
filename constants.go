package catutil

const (
	CatCtx               = "cat_ctx"
	CatCtxRootTran       = "cat_root_tran"
	CatCtxRedisTran      = "cat_redis_tran"
	CatCtxMysqlTran      = "cat_mysql_tran"
	CatCtxMongoDbTran    = "cat_mongo_db_tran"
	CatCtxHttpRemoteTran = "cat_http_remote_tran"

	TypeHttpRemote    = "HttpRemoteCall"
	TypeMongoDbPrefix = "MongoDb."
	TypeUrlStatus     = "URL.status"

	RemoteCallMethod = "RemoteCall.Method"
	RemoteCallErr    = "RemoteCall.Err"
	RemoteCallStatus = "RemoteCall.Status"
	RemoteCallScheme = "RemoteCall.Scheme"
	TraceIdHeaderName = "x-trace-id"
	ContextKeyRequestId = "requestId"
	ContextKeyTraceId   = "traceId"
	ContextKeyStartTime = "startTime"
)