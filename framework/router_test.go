// Copyright © 2016 Tom Maiaroto <tom@shift8creative.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package framework

import (
	"context"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestRouter(t *testing.T) {
	testFallThroughHandler := func(ctx context.Context, d *HandlerDependencies, req *APIGatewayProxyRequest, res *APIGatewayProxyResponse, params url.Values) error {
		return nil
	}
	testHandler := func(ctx context.Context, d *HandlerDependencies, req *APIGatewayProxyRequest, res *APIGatewayProxyResponse, params url.Values) error {
		return nil
	}
	testMiddleware := func(ctx context.Context, d *HandlerDependencies, req *APIGatewayProxyRequest, res *APIGatewayProxyResponse, params url.Values) bool {
		return true
	}
	testMiddlewareStop := func(ctx context.Context, d *HandlerDependencies, req *APIGatewayProxyRequest, res *APIGatewayProxyResponse, params url.Values) bool {
		return false
	}
	testParams := url.Values{}
	testRouter := NewRouter(testFallThroughHandler)
	Convey("NewRouter", t, func() {
		Convey("Should create a new Router", func() {
			So(testRouter, ShouldNotBeNil)
		})
	})

	Convey("Should be able to add Router level middleware", t, func() {
		testRouter.Use(testMiddleware)
		So(testRouter.middleware, ShouldHaveLength, 1)
	})

	testRouter.GET("/path", testHandler)
	testRouter.POST("/path", testHandler)
	testRouter.PUT("/path", testHandler)
	testRouter.PATCH("/path", testHandler)
	testRouter.DELETE("/path", testHandler)
	testRouter.HEAD("/path", testHandler)
	testRouter.OPTIONS("/path", testHandler)

	node, _ := testRouter.tree.traverse(strings.Split("/path", "/")[1:], testParams)

	Convey("Should handle GET", t, func() {
		So(node.methods, ShouldContainKey, "GET")
	})

	Convey("Should handle POST", t, func() {
		So(node.methods, ShouldContainKey, "POST")
	})

	Convey("Should handle PUT", t, func() {
		So(node.methods, ShouldContainKey, "PUT")
	})

	Convey("Should handle PATCH", t, func() {
		So(node.methods, ShouldContainKey, "PATCH")
	})

	Convey("Should handle DELETE", t, func() {
		So(node.methods, ShouldContainKey, "DELETE")
	})

	Convey("Should handle HEAD", t, func() {
		So(node.methods, ShouldContainKey, "HEAD")
	})

	Convey("Should handle OPTIONS", t, func() {
		So(node.methods, ShouldContainKey, "OPTIONS")
	})

	Convey("runMiddleware", t, func() {
		Convey("Should handle middleware", func() {
			ctx := context.Background()
			req := APIGatewayProxyRequest{}
			res := APIGatewayProxyResponse{}
			params := url.Values{}
			d := HandlerDependencies{}
			next := runMiddleware(ctx, &d, &req, &res, params, testMiddleware)
			So(next, ShouldBeTrue)

			noNext := runMiddleware(ctx, &d, &req, &res, params, testMiddlewareStop, testMiddleware)
			So(noNext, ShouldBeFalse)
		})
	})

	Convey("requestToProxyRequest", t, func() {
		Convey("Should take an HTTP request and format a Lambda Event", func() {
			gwHandler := gatewayHandler{}
			r := httptest.NewRequest("GET", "/?foo=bar", strings.NewReader("some body to be read"))
			r.Header.Set("User-Agent", "aegis-test")

			_, req := gwHandler.requestToProxyRequest(r)

			So(req.Body, ShouldEqual, "some body to be read")
			So(req.Headers, ShouldContainKey, "User-Agent")
			So(req.QueryStringParameters, ShouldContainKey, "foo")

		})
	})

	Convey("proxyResponseToHTTPResponse", t, func() {
		Convey("Should take a Lambda Proxy response and format an HTTP response", func() {
			gwHandler := gatewayHandler{}
			res := APIGatewayProxyResponse{
				StatusCode: 200,
				Headers:    map[string]string{"Content-Type": "application/json"},
			}
			rw := httptest.NewRecorder()
			gwHandler.proxyResponseToHTTPResponse(&res, rw)

			result := rw.Result()
			rw.Flush()
			So(result.StatusCode, ShouldEqual, 200)
			So(result.Header.Get("Content-Type"), ShouldEqual, "application/json")
		})
	})
}
