package httpcat

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/cat-go/cat"
	"github.com/cat-go/cat/message"
	"github.com/gin-gonic/gin"

	"github.com/ainiaa/catutil"
)

type RequestData struct {
	Method      string            // http.MethodPost http.MethodGet ...
	Header      map[string]string // header 信息
	KVData      map[string]string // data 数据
	RequestBody string
	ClientConf  HttpClientConf
}

func GetWithCat(ctx context.Context, uri string, data *RequestData) (res []byte) {
	if data == nil {
		data = &RequestData{
			Method: http.MethodGet,
			ClientConf: c,
		}
	} else {
		data.Method = http.MethodGet
	}
	return RequestWithCat(ctx, uri, data)
}

func PostWithCat(ctx context.Context, uri string, data *RequestData) (res []byte) {
	if data == nil {
		data = &RequestData{
			Method: http.MethodPost,
		}
	} else {
		data.Method = http.MethodPost
	}
	return RequestWithCat(ctx, uri, data)
}

/**
 * 这个方法会返回*http.Response，理论上如果不用，处理方应该手动close掉
 */
func RequestWithCat(ctx context.Context, urlRaw string, data *RequestData) (re []byte) {
	var method = data.Method
	var header = data.Header
	var timeout = data.ClientConf.ReadWriteTimeoutMS
	var kvData = data.KVData
	if !cat.IsEnabled() {
		if len(kvData) > 0 {
			return RequestUrlencoded(ctx, urlRaw, method, kvData, header, int64(timeout))
		}
		return Request(ctx, urlRaw, method, data.RequestBody, header, int64(timeout))

	}
	var tran message.Transactor

	u, err := url.Parse(urlRaw)
	if ctx == nil {
		ctx = context.Background()
	}
	if err == nil {
		if method == "" {
			method = http.MethodPost
		}
		tran = cat.NewTransaction(catutil.TypeHttpRemote, u.Host+u.Path)
		tran.SetMessageId(cat.MessageId())
		defer tran.Complete()
		tran.AddData(catutil.RemoteCallMethod, method)
		tran.AddData(catutil.RemoteCallScheme, u.Scheme)
		if rootTran := catutil.GetRootTran(ctx); rootTran != nil {
			cat.SetChildTraceId(rootTran, tran)
			rootTran.AddChild(tran)
		}
	}
	switch ctx := ctx.(type) {
	case *gin.Context:
		ctx.Set(catutil.CatCtxHttpRemoteTran, tran)
	default:
		ctx = context.WithValue(ctx, catutil.CatCtxHttpRemoteTran, tran)
	}
	re = request(ctx, urlRaw, data)
	return
}

func request(ctx context.Context, uri string, reqData *RequestData) (re []byte) {
	var method = reqData.Method
	var header = reqData.Header
	var kvData = reqData.KVData
	var data = reqData.RequestBody
	if method == "" {
		method = http.MethodPost
	}
	var tran message.Transactor
	if cat.IsEnabled() {
		if ctx != nil {
			if tranRaw := ctx.Value(catutil.CatCtxHttpRemoteTran); tranRaw != nil {
				if tranInner, ok := tranRaw.(message.Transactor); ok && tranInner != nil {
					tran = tranInner
				}
			}
		} else {
			ctx = context.Background()
		}
	}
	var reqBody io.Reader
	if header["Content-Type"] == "application/x-www-form-urlencoded" && kvData != nil {
		dataUrlVal := url.Values{}
		for key, val := range kvData {
			dataUrlVal.Add(key, val)
		}
		reqBody = strings.NewReader(dataUrlVal.Encode())
	} else {
		reqBody = strings.NewReader(data)
	}
	req, err := http.NewRequest(method, uri, reqBody)
	if err != nil {
		if tran != nil && cat.IsEnabled() {
			tran.AddData(catutil.RemoteCallErr, err.Error())
			tran.SetStatus(message.CatError)
		}
		return re
	}
	defer req.Body.Close()
	for k, v := range header {
		req.Header.Add(k, v)
	}
	if tran != nil && cat.IsEnabled() {
		childId := cat.MessageId()
		tran.LogEvent(cat.TypeRemoteCall, "rpc remote", cat.SUCCESS, childId)
		req.Header.Set(cat.RootId, tran.GetRootMessageId())
		req.Header.Set(cat.ParentId, tran.GetParentMessageId())
		req.Header.Set(cat.ChildId, childId)
	}
	if req.Header.Get("Content-Type") == "" {
		req.Header.Add("Content-Type", "application/json")
	}
	clt := NewHttpClient(reqData.ClientConf)
	resp, err := clt.Do(req)
	if err != nil {
		if tran != nil && cat.IsEnabled() {
			tran.AddData(catutil.RemoteCallErr, err.Error())
			tran.SetStatus(message.CatError)
		}
		return re
	}
	if reqData.ClientConf.BodyReadTimeoutMS == 0 {
		re, err = ioutil.ReadAll(resp.Body)
	} else {
		re, err = readAllWithTimout(resp, reqData.ClientConf)
	}

	if err != nil {
		if tran != nil && cat.IsEnabled() {
			tran.AddData(catutil.RemoteCallErr, err.Error())
			tran.SetStatus(message.CatError)
		}
		return
	}
	if tran != nil && cat.IsEnabled() {
		tran.AddData(catutil.RemoteCallStatus, resp.Status)
	}
	if resp.StatusCode != http.StatusOK {
		if tran != nil && cat.IsEnabled() {
			tran.AddData(catutil.RemoteCallErr, fmt.Sprintf("StatusCode:%d", resp.StatusCode))
			tran.SetStatus(message.CatError)
		}
	} else {
		if tran != nil && cat.IsEnabled() {
			tran.SetStatus(message.CatSuccess)
		}
	}
	if resp != nil && resp.Body != nil {
		_ = resp.Body.Close()
	}
	return
}

func readAllWithTimout(resp *http.Response, cc HttpClientConf) (body []byte, err error) {
	if resp.StatusCode == http.StatusOK { // 请求成功
		bodyReadTimeoutMS := 3000
		if cc.BodyReadTimeoutMS > 0 {
			bodyReadTimeoutMS = cc.BodyReadTimeoutMS
		}

		ctx, _ := context.WithTimeout(context.Background(), time.Duration(bodyReadTimeoutMS)*time.Millisecond)
		finish := make(chan struct{}, 1)
		go func() {
			body, err = ioutil.ReadAll(resp.Body)
			finish <- struct{}{}
		}()
		select {
		case <-ctx.Done(): // read body 超时
			return body, fmt.Errorf("readBody timeout bodyReadTimeoutMS:%d error", bodyReadTimeoutMS)
		case <-finish: // 正常获取返回值
			return
		}
	} else {
		return nil, nil
	}
}
