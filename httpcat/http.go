package httpcat

import (
	"context"
	"encoding/json"
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

/**
 * 这个方法会返回*http.Response，理论上如果不用，处理方应该手动close掉
 */
func RequestWithCat(ctx context.Context, urlRaw string, method string,
	data string, header map[string]string, timeout int64, mapData map[string]string) (re []byte) {
	if !cat.IsEnabled() {
		if len(mapData) > 0 {
			return RequestUrlencoded(ctx, urlRaw, method, mapData, header, timeout)
		}
		return Request(ctx, urlRaw, method, data, header, timeout)
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
	re = request(ctx, urlRaw, method, data, header, timeout, mapData)
	return
}

func request(ctx context.Context, uri string, method string,
	data string, header map[string]string, timeout int64, mapData map[string]string) (re []byte) {
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
	if header["Content-Type"] == "application/x-www-form-urlencoded" && mapData != nil {
		dataUrlVal := url.Values{}
		for key, val := range mapData {
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
	if timeout == 0 {
		//默认超时2000ms
		timeout = 2000
	}
	clt := http.Client{
		Timeout: time.Duration(timeout) * time.Millisecond,
	}
	resp, err := clt.Do(req)
	if err != nil {
		if tran != nil && cat.IsEnabled() {
			tran.AddData(catutil.RemoteCallErr, err.Error())
			tran.SetStatus(message.CatError)
		}
		return re
	}
	body, err := ioutil.ReadAll(resp.Body)
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
	return body
}

type Resp struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func GetBody(req string) *Resp {
	r := &Resp{}
	_ = json.Unmarshal([]byte(req), &r)
	return r
}

func PostWithTime(ctx context.Context, url string, data string, timeout int64) (re []byte) {
	re = RequestWithCat(ctx, url, http.MethodPost, data, map[string]string{
		"Content-Type":            "application/json",
		catutil.TraceIdHeaderName: catutil.GetTraceId(ctx),
	}, timeout, nil)
	return re
}
