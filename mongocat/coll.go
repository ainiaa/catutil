package mongocat

import (
	"context"
	"encoding/json"

	"github.com/cat-go/cat"
	"github.com/cat-go/cat/message"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/ainiaa/catutil"
)

type handler struct {
	coll *mongo.Collection
}
type ConnHandler interface {
	Find(ctx context.Context, filter interface{}, opts ...*options.FindOptions) (cur *mongo.Cursor, err error)
	FindOne(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) *mongo.SingleResult
	CountDocuments(ctx context.Context, filter interface{}, opts ...*options.CountOptions) (res int64, err error)
	UpdateOne(ctx context.Context, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (res *mongo.UpdateResult, err error)
	InsertOne(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (res *mongo.InsertOneResult, err error)
	FindOneAndUpdate(ctx context.Context, filter interface{}, update interface{}, opts ...*options.FindOneAndUpdateOptions) *mongo.SingleResult
	InsertMany(ctx context.Context, documents []interface{}, opts ...*options.InsertManyOptions) (res *mongo.InsertManyResult, err error)
	UpdateMany(ctx context.Context, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (res *mongo.UpdateResult, err error)
}

func WithCat(coll *mongo.Collection) ConnHandler {
	h := &handler{}
	h.coll = coll
	return h
}

func (h *handler) Find(ctx context.Context, filter interface{},
	opts ...*options.FindOptions) (cur *mongo.Cursor, err error) {
	var tran message.Transactor
	if cat.IsEnabled() {
		tran = cat.NewTransaction(catutil.TypeMongoDbPrefix+h.coll.Name(), "Find")
		defer tran.Complete()
		if filterInfo, filterInfoErr := json.Marshal(filter); filterInfoErr == nil {
			tran.AddData(string(filterInfo))
		}
		if rootTran := catutil.GetRootTran(ctx); rootTran != nil {
			cat.SetChildTraceId(rootTran, tran)
			rootTran.AddChild(tran)
		}
	}
	cur, err = h.coll.Find(ctx, filter, opts...)
	if err != nil {
		if cat.IsEnabled() && tran != nil {
			tran.AddData("err", err.Error())
		}
	}
	return
}

func (h *handler) FindOne(ctx context.Context, filter interface{},
	opts ...*options.FindOneOptions) *mongo.SingleResult {
	if cat.IsEnabled() {
		tran := cat.NewTransaction(catutil.TypeMongoDbPrefix+h.coll.Name(), "FindOne")
		defer tran.Complete()
		if filterInfo, filterInfoErr := json.Marshal(filter); filterInfoErr == nil {
			tran.AddData(string(filterInfo))
		}
		if rootTran := catutil.GetRootTran(ctx); rootTran != nil {
			cat.SetChildTraceId(rootTran, tran)
			rootTran.AddChild(tran)
		}
	}
	return h.coll.FindOne(ctx, filter, opts...)
}

func (h *handler) CountDocuments(ctx context.Context, filter interface{},
	opts ...*options.CountOptions) (res int64, err error) {
	var tran message.Transactor
	if cat.IsEnabled() {
		tran = cat.NewTransaction(catutil.TypeMongoDbPrefix+h.coll.Name(), "CountDocuments")
		defer tran.Complete()
		if filterInfo, filterInfoErr := json.Marshal(filter); filterInfoErr == nil {
			tran.AddData(string(filterInfo))
		}
		if rootTran := catutil.GetRootTran(ctx); rootTran != nil {
			cat.SetChildTraceId(rootTran, tran)
			rootTran.AddChild(tran)
		}
	}
	res, err = h.coll.CountDocuments(ctx, filter, opts...)
	if err != nil {
		if cat.IsEnabled() && tran != nil {
			tran.AddData("err", err.Error())
		}
	}
	return
}

func (h *handler) UpdateOne(ctx context.Context, filter interface{}, update interface{},
	opts ...*options.UpdateOptions) (res *mongo.UpdateResult, err error) {
	var tran message.Transactor
	if cat.IsEnabled() {
		tran = cat.NewTransaction(catutil.TypeMongoDbPrefix+h.coll.Name(), "UpdateOne")
		defer tran.Complete()
		if filterInfo, filterInfoErr := json.Marshal(filter); filterInfoErr == nil {
			tran.AddData(string(filterInfo))
		}
		if rootTran := catutil.GetRootTran(ctx); rootTran != nil {
			cat.SetChildTraceId(rootTran, tran)
			rootTran.AddChild(tran)
		}
	}
	res, err = h.coll.UpdateOne(ctx, filter, update, opts...)
	if err != nil {
		if cat.IsEnabled() && tran != nil {
			tran.AddData("err", err.Error())
		}
	}
	return
}

func (h *handler) InsertOne(ctx context.Context, document interface{},
	opts ...*options.InsertOneOptions) (res *mongo.InsertOneResult, err error) {
	var tran message.Transactor
	if cat.IsEnabled() {
		tran = cat.NewTransaction(catutil.TypeMongoDbPrefix+h.coll.Name(), "InsertOne")
		defer tran.Complete()
		if documentInfo, documentInfoErr := json.Marshal(document); documentInfoErr == nil {
			tran.AddData(string(documentInfo))
		}
		if rootTran := catutil.GetRootTran(ctx); rootTran != nil {
			cat.SetChildTraceId(rootTran, tran)
			rootTran.AddChild(tran)
		}
	}
	res, err = h.coll.InsertOne(ctx, document, opts...)
	if err != nil {
		if cat.IsEnabled() && tran != nil {
			tran.AddData("err", err.Error())
		}
	}
	return
}

func (h *handler) FindOneAndUpdate(ctx context.Context, filter interface{},
	update interface{}, opts ...*options.FindOneAndUpdateOptions) *mongo.SingleResult {
	if cat.IsEnabled() {
		tran := cat.NewTransaction(catutil.TypeMongoDbPrefix+h.coll.Name(), "FindOneAndUpdate")
		defer tran.Complete()
		if filterInfo, filterInfoErr := json.Marshal(filter); filterInfoErr == nil {
			tran.AddData("filterInfo", string(filterInfo))
		}
		if updateInfo, updateInfoErr := json.Marshal(update); updateInfoErr == nil {
			tran.AddData("updateInfo", string(updateInfo))
		}
		if rootTran := catutil.GetRootTran(ctx); rootTran != nil {
			cat.SetChildTraceId(rootTran, tran)
			rootTran.AddChild(tran)
		}
	}
	return h.coll.FindOneAndUpdate(ctx, filter, update, opts...)
}

func (h *handler) InsertMany(ctx context.Context, documents []interface{},
	opts ...*options.InsertManyOptions) (res *mongo.InsertManyResult, err error) {
	var tran message.Transactor
	if cat.IsEnabled() {
		tran = cat.NewTransaction(catutil.TypeMongoDbPrefix+h.coll.Name(), "InsertMany")
		defer tran.Complete()
		if rootTran := catutil.GetRootTran(ctx); rootTran != nil {
			cat.SetChildTraceId(rootTran, tran)
			rootTran.AddChild(tran)
		}
	}
	res, err = h.coll.InsertMany(ctx, documents, opts...)
	if err != nil {
		if cat.IsEnabled() && tran != nil {
			tran.AddData("err", err.Error())
		}
	}
	return
}

func (h *handler) UpdateMany(ctx context.Context, filter interface{}, update interface{},
	opts ...*options.UpdateOptions) (res *mongo.UpdateResult, err error) {
	var tran message.Transactor
	if cat.IsEnabled() {
		tran = cat.NewTransaction(catutil.TypeMongoDbPrefix+h.coll.Name(), "UpdateMany")
		defer tran.Complete()
		if rootTran := catutil.GetRootTran(ctx); rootTran != nil {
			cat.SetChildTraceId(rootTran, tran)
			rootTran.AddChild(tran)
		}
	}
	res, err = h.coll.UpdateMany(ctx, filter, update, opts...)
	if err != nil {
		if cat.IsEnabled() && tran != nil {
			tran.AddData("err", err.Error())
		}
	}
	return
}
