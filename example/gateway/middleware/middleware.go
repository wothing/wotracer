/**
 * Copyright 2015-2016, Wothing Co., Ltd.
 * All rights reserved.
 *
 * Created by Elvizlai on 2016/06/06 09:58
 */

package middleware

import (
	"io/ioutil"
	"net/http"

	gctx "github.com/gorilla/context"
	"golang.org/x/net/context"

	"github.com/wothing/wotracer"
	"strings"
)

type LogMiddleware struct {
}

func (*LogMiddleware) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	body, err := ioutil.ReadAll(r.Body)
	r.Body.Close()
	if err != nil {
		return
	}

	span, ctx := wotracer.InjectRPC(context.Background(), r.RequestURI)
	gctx.Set(r, "ctx", ctx)
	rw.Header().Set("X-Request-Id", wotracer.GetTraceID(span))

	if realIp := r.Header.Get("X-Real-IP"); realIp != "" {
		r.RemoteAddr = realIp
	} else {
		r.RemoteAddr = strings.Split(r.RemoteAddr, ":")[0]
	}
	span.SetTag("Reqeust.Address", r.RemoteAddr)

	gctx.Set(r, "body", body)
	next(rw, r)
	span.Finish()
}
