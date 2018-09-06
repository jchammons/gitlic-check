// +build go1.7
// Copyright (C) 2016 Librato, Inc. All rights reserved.

package ao_test

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"net"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"os"

	"github.com/appoptics/appoptics-apm-go/v1/ao"
	"github.com/appoptics/appoptics-apm-go/v1/ao/internal/config"
	g "github.com/appoptics/appoptics-apm-go/v1/ao/internal/graphtest"
	"github.com/appoptics/appoptics-apm-go/v1/ao/internal/reporter"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/context"
)

func handler404(w http.ResponseWriter, r *http.Request)      { w.WriteHeader(404) }
func handler403(w http.ResponseWriter, r *http.Request)      { w.WriteHeader(403) }
func handler200(w http.ResponseWriter, r *http.Request)      { checkAOContextAndSetCustomTxnName(w, r) }
func handlerPanic(w http.ResponseWriter, r *http.Request)    { panic("panicking!") }
func handlerDelay200(w http.ResponseWriter, r *http.Request) { time.Sleep(httpSpanSleep) }
func handlerDelay503(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(503)
	time.Sleep(httpSpanSleep)
}

// checkAOContext checks if the AO context is attached
func checkAOContextAndSetCustomTxnName(w http.ResponseWriter, r *http.Request) {
	xtrace := ""
	var t ao.Trace
	if t = ao.TraceFromContext(r.Context()); t == nil {
		return
	}

	defer func() {
		fmt.Fprint(w, xtrace)
	}()

	// Concurrently set custom transaction names
	for i := 0; i < 10; i++ {
		go func(i int) {
			time.Sleep(time.Duration(rand.Intn(5)) * time.Millisecond)
			ao.SetTransactionName(r.Context(), "my-custom-transaction-name-"+strconv.Itoa(i))
		}(i)
	}
	time.Sleep(10 * time.Millisecond)
	ao.SetTransactionName(r.Context(), "final-"+ao.GetTransactionName(r.Context()))
	xtrace = t.MetadataString()
}

func handlerDoubleWrapped(w http.ResponseWriter, r *http.Request) {
	t, _, _ := ao.TraceFromHTTPRequestResponse("myHandler", w, r)
	ao.NewContext(context.Background(), t)
	defer t.End()
}

func httpTestWithEndpoint(f http.HandlerFunc, ep string) *httptest.ResponseRecorder {
	return httpTestWithEndpointWithHeaders(f, ep, nil)
}

func httpTestWithEndpointWithHeaders(f http.HandlerFunc, ep string, hd map[string]string) *httptest.ResponseRecorder {
	h := http.HandlerFunc(ao.HTTPHandler(f))
	// test a single GET request
	req, _ := http.NewRequest("GET", ep, nil)
	for k, v := range hd {
		req.Header.Set(k, v)
	}
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)
	return w
}

func httpTest(f http.HandlerFunc) *httptest.ResponseRecorder {
	return httpTestWithEndpoint(f, "http://test.com/hello?testq")
}

func TestHTTPHandler404(t *testing.T) {
	r := reporter.SetTestReporter() // set up test reporter
	response := httpTest(handler404)
	assert.Len(t, response.HeaderMap[ao.HTTPHeaderName], 1)

	r.Close(2)
	g.AssertGraph(t, r.EventBufs, 2, g.AssertNodeMap{
		// entry event should have no edges
		{"http.HandlerFunc", "entry"}: {Edges: g.Edges{}, Callback: func(n g.Node) {
			assert.Equal(t, "/hello", n.Map["URL"])
			assert.Equal(t, "test.com", n.Map["HTTP-Host"])
			assert.Equal(t, "GET", n.Map["Method"])
			assert.Equal(t, "testq", n.Map["Query-String"])
		}},
		{"http.HandlerFunc", "exit"}: {Edges: g.Edges{{"http.HandlerFunc", "entry"}}, Callback: func(n g.Node) {
			// assert that response X-Trace header matches trace exit event
			assert.Equal(t, response.HeaderMap.Get(ao.HTTPHeaderName), n.Map[ao.HTTPHeaderName])
			assert.EqualValues(t, response.Code, n.Map["Status"])
			assert.EqualValues(t, 404, n.Map["Status"])
			assert.Equal(t, "ao_test", n.Map["Controller"])
			assert.Equal(t, "handler404", n.Map["Action"])
		}},
	})
}

func TestHTTPHandler200(t *testing.T) {
	os.Setenv("APPOPTICS_PREPEND_DOMAIN", "false")
	config.Refresh()
	r := reporter.SetTestReporter() // set up test reporter
	response := httpTestWithEndpoint(handler200, "http://test.com/hello world/one/two/three?testq")

	r.Close(2)
	g.AssertGraph(t, r.EventBufs, 2, g.AssertNodeMap{
		// entry event should have no edges
		{"http.HandlerFunc", "entry"}: {Edges: g.Edges{}, Callback: func(n g.Node) {
			assert.Equal(t, "/hello%20world/one/two/three", n.Map["URL"])
			assert.Equal(t, "test.com", n.Map["HTTP-Host"])
			assert.Equal(t, "GET", n.Map["Method"])
			assert.Equal(t, "testq", n.Map["Query-String"])
		}},
		{"http.HandlerFunc", "exit"}: {Edges: g.Edges{{"http.HandlerFunc", "entry"}}, Callback: func(n g.Node) {
			// assert that response X-Trace header matches trace exit event
			assert.Len(t, response.HeaderMap[ao.HTTPHeaderName], 1)
			assert.Equal(t, response.HeaderMap[ao.HTTPHeaderName][0], n.Map[ao.HTTPHeaderName])
			assert.EqualValues(t, response.Code, n.Map["Status"])
			assert.EqualValues(t, 200, n.Map["Status"])
			assert.Equal(t, "ao_test", n.Map["Controller"])
			assert.Equal(t, "handler200", n.Map["Action"])
			assert.True(t, strings.HasPrefix(n.Map["TransactionName"].(string), "final-my-custom-transaction-name"))
			assert.True(t, reporter.ValidMetadata(response.Body.String()))
		}},
	})
}

func TestHTTPHandlerNoTrace(t *testing.T) {
	r := reporter.SetTestReporter(reporter.TestReporterDisableTracing())
	httpTest(handler404)

	// tracing disabled, shouldn't report anything
	assert.Len(t, r.EventBufs, 0)
}

var httpSpanSleep time.Duration

func TestHTTPSpan(t *testing.T) {
	r := reporter.SetTestReporter(reporter.TestReporterDisableDefaultSetting(true)) // set up test reporter

	httpSpanSleep = time.Duration(0) // fire off first request just as preparation for the following requests
	httpTest(handlerDelay200)
	httpSpanSleep = time.Duration(0)
	httpTest(handlerDelay200)
	httpSpanSleep = time.Duration(25 * time.Millisecond)
	httpTest(handlerDelay200)
	httpSpanSleep = time.Duration(456 * time.Millisecond)
	httpTest(handlerDelay200)
	httpSpanSleep = time.Duration(54 * time.Millisecond)
	httpTest(handlerDelay503)

	r.Close(5)

	require.Len(t, r.SpanMessages, 5)

	m, ok := r.SpanMessages[1].(*reporter.HTTPSpanMessage)
	assert.True(t, ok)
	nullDuration := m.Duration

	m, ok = r.SpanMessages[2].(*reporter.HTTPSpanMessage)
	assert.True(t, ok)
	assert.Equal(t, "ao_test.handlerDelay200", m.Transaction)
	assert.Equal(t, "/hello", m.Path)
	assert.Equal(t, 200, m.Status)
	assert.Equal(t, "GET", m.Method)
	assert.False(t, m.HasError)
	assert.InDelta(t, (25*time.Millisecond + nullDuration).Seconds(), m.Duration.Seconds(), (10 * time.Millisecond).Seconds())

	m, ok = r.SpanMessages[3].(*reporter.HTTPSpanMessage)
	assert.True(t, ok)
	assert.InDelta(t, (456*time.Millisecond + nullDuration).Seconds(), m.Duration.Seconds(), (10 * time.Millisecond).Seconds())

	m, ok = r.SpanMessages[4].(*reporter.HTTPSpanMessage)
	assert.True(t, ok)
	assert.Equal(t, "ao_test.handlerDelay503", m.Transaction)
	assert.Equal(t, 503, m.Status)
	assert.True(t, m.HasError)
	assert.InDelta(t, (54*time.Millisecond + nullDuration).Seconds(), m.Duration.Seconds(), (10 * time.Millisecond).Seconds())
}

func TestSingleHTTPSpan(t *testing.T) {
	r := reporter.SetTestReporter(reporter.TestReporterDisableDefaultSetting(true)) // set up test reporter
	httpTest(handlerDoubleWrapped)
	r.Close(1)

	assert.Equal(t, 1, len(r.SpanMessages))
}

// testServer tests creating a span/trace from inside an HTTP handler (using ao.TraceFromHTTPRequest)
func testServer(t *testing.T, list net.Listener) {
	s := &http.Server{Handler: http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		// create span from incoming HTTP Request headers, if trace exists
		tr, w, req := ao.TraceFromHTTPRequestResponse("myHandler", w, req)
		defer tr.End()

		tr.AddEndArgs("NotReported") // odd-length args, should have no effect

		t.Logf("server: got request %v", req)
		l2 := tr.BeginSpan("DBx", "Query", "SELECT *", "RemoteHost", "db.net")
		// Run a query ...
		l2.End()

		w.WriteHeader(403) // return Forbidden
	})}
	assert.NoError(t, s.Serve(list))
}

// same as testServer, but with external ao.HTTPHandler() handler wrapping
func testDoubleWrappedServer(t *testing.T, list net.Listener) {
	s := &http.Server{Handler: http.HandlerFunc(ao.HTTPHandler(func(writer http.ResponseWriter, req *http.Request) {
		// create span from incoming HTTP Request headers, if trace exists
		tr, w, req := ao.TraceFromHTTPRequestResponse("myHandler", writer, req)
		defer tr.End()

		t.Logf("server: got request %v", req)
		l2 := tr.BeginSpan("DBx", "Query", "SELECT *", "RemoteHost", "db.net")
		// Run a query ...
		l2.End()

		w.WriteHeader(403) // return Forbidden
	}))}
	assert.NoError(t, s.Serve(list))
}

// testServer200 does not trace and returns a 200.
func testServer200(t *testing.T, list net.Listener) {
	s := &http.Server{Handler: http.HandlerFunc(handler200)}
	assert.NoError(t, s.Serve(list))
}

// testServer403 does not trace and returns a 403.
func testServer403(t *testing.T, list net.Listener) {
	s := &http.Server{Handler: http.HandlerFunc(handler403)}
	assert.NoError(t, s.Serve(list))
}

// simulate panic-catching middleware wrapping ao.HTTPHandler(handlerPanic)
func panicCatchingMiddleware(t *testing.T, f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				t.Logf("panicCatcher caught panic %v", err)
				w.WriteHeader(500)
				w.Write([]byte(fmt.Sprintf("500 Error: %v", err)))
			}
		}()
		f(w, r)
	}
}

// testServerPanic traces a wrapped http.HandlerFunc that panics
func testServerPanic(t *testing.T, list net.Listener) {
	s := &http.Server{Handler: http.HandlerFunc(
		panicCatchingMiddleware(t, ao.HTTPHandler(handlerPanic)))}
	assert.NoError(t, s.Serve(list))
}

// begin an HTTP client span, make an HTTP request, and propagate the trace context manually
func testHTTPClient(t *testing.T, ctx context.Context, method, url string) (*http.Response, error) {
	httpClient := &http.Client{}
	httpReq, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}
	l, _ := ao.BeginSpan(ctx, "http.Client", "IsService", true, "RemoteURL", url)
	defer l.End()
	httpReq.Header.Set(ao.HTTPHeaderName, l.MetadataString())

	resp, err := httpClient.Do(httpReq)
	if err != nil {
		l.Err(err)
		return resp, err
	}
	defer resp.Body.Close()

	l.AddEndArgs("Edge", resp.Header.Get(ao.HTTPHeaderName))
	return resp, err
}

// create an HTTP client span, make an HTTP request, and propagate the trace using HTTPClientSpan
func testHTTPClientA(t *testing.T, ctx context.Context, method, url string) (*http.Response, error) {
	httpClient := &http.Client{}
	httpReq, err := http.NewRequest(method, url, nil)
	l := ao.BeginHTTPClientSpan(ctx, httpReq)
	defer l.End()
	if err != nil {
		l.Err(err)
		return nil, err
	}

	resp, err := httpClient.Do(httpReq)
	l.AddHTTPResponse(resp, err)
	if err != nil {
		t.Logf("JoinResponse err: %v", err)
		return resp, err
	}
	defer resp.Body.Close()

	return resp, err
}

// create an HTTP client span, make an HTTP request, and propagate the trace using HTTPClientSpan
// and a different exception-handling flow
func testHTTPClientB(t *testing.T, ctx context.Context, method, url string) (*http.Response, error) {
	httpClient := &http.Client{}
	httpReq, err := http.NewRequest(method, url, nil)
	l := ao.BeginHTTPClientSpan(ctx, httpReq)
	if err != nil {
		l.Err(err)
		l.End()
		return nil, err
	}

	resp, err := httpClient.Do(httpReq)
	l.AddHTTPResponse(resp, err)
	l.End()
	if err != nil {
		t.Logf("JoinResponse err: %v", err)
		return resp, err
	}
	defer resp.Body.Close()

	return resp, err
}

type testClientFn func(t *testing.T, ctx context.Context, method, url string) (*http.Response, error)
type testServerFn struct {
	serverFn func(t *testing.T, list net.Listener)
	assertFn func(t *testing.T, bufs [][]byte, resp *http.Response, url, method string, port, status int)
	numBufs  int
	status   int
}

var testHTTPSvr = testServerFn{testServer, assertHTTPRequestGraph, 8, 403}
var testHTTPSvr200 = testServerFn{testServer200, assertHTTPRequestUntracedGraph, 4, 200}
var testHTTPSvr403 = testServerFn{testServer403, assertHTTPRequestUntracedGraph, 4, 403}
var testHTTPSvrPanic = testServerFn{testServerPanic, assertHTTPRequestPanic, 7, 200}

var badURL = "%gh&%ij" // url.Parse() will return error
var invalidPortURL = "http://0.0.0.0:888888"

func TestTraceHTTP(t *testing.T)              { testHTTP(t, "GET", false, testHTTPClient, testHTTPSvr) }
func TestTraceHTTPHelperA(t *testing.T)       { testHTTP(t, "GET", false, testHTTPClientA, testHTTPSvr) }
func TestTraceHTTPHelperB(t *testing.T)       { testHTTP(t, "GET", false, testHTTPClientB, testHTTPSvr) }
func TestTraceHTTP200(t *testing.T)           { testHTTP(t, "GET", false, testHTTPClient, testHTTPSvr200) }
func TestTraceHTTPHelperA200(t *testing.T)    { testHTTP(t, "GET", false, testHTTPClientA, testHTTPSvr200) }
func TestTraceHTTPHelperB200(t *testing.T)    { testHTTP(t, "GET", false, testHTTPClientB, testHTTPSvr200) }
func TestTraceHTTP403(t *testing.T)           { testHTTP(t, "GET", false, testHTTPClient, testHTTPSvr403) }
func TestTraceHTTPHelperA403(t *testing.T)    { testHTTP(t, "GET", false, testHTTPClientA, testHTTPSvr403) }
func TestTraceHTTPHelperB403(t *testing.T)    { testHTTP(t, "GET", false, testHTTPClientB, testHTTPSvr403) }
func TestTraceHTTPPost(t *testing.T)          { testHTTP(t, "POST", false, testHTTPClient, testHTTPSvr) }
func TestTraceHTTPHelperPostA(t *testing.T)   { testHTTP(t, "POST", false, testHTTPClientA, testHTTPSvr) }
func TestTraceHTTPHelperPostB(t *testing.T)   { testHTTP(t, "POST", false, testHTTPClientB, testHTTPSvr) }
func TestTraceHTTPBadRequest(t *testing.T)    { testHTTP(t, "GET", true, testHTTPClient, testHTTPSvr) }
func TestTraceHTTPHelperBadReqA(t *testing.T) { testHTTP(t, "GET", true, testHTTPClientA, testHTTPSvr) }
func TestTraceHTTPHelperBadReqB(t *testing.T) { testHTTP(t, "GET", true, testHTTPClientB, testHTTPSvr) }
func TestTraceHTTPPanic(t *testing.T)         { testHTTP(t, "GET", false, testHTTPClient, testHTTPSvrPanic) }
func TestTraceHTTPPanicA(t *testing.T)        { testHTTP(t, "GET", false, testHTTPClientA, testHTTPSvrPanic) }
func TestTraceHTTPPanicB(t *testing.T)        { testHTTP(t, "GET", false, testHTTPClientB, testHTTPSvrPanic) }

// launch a test HTTP server and trace an HTTP request to it
func testHTTP(t *testing.T, method string, badReq bool, clientFn testClientFn, server testServerFn) {
	ln, err := net.Listen("tcp", ":0") // pick an unallocated port
	assert.NoError(t, err)
	port := ln.Addr().(*net.TCPAddr).Port
	go server.serverFn(t, ln) // start test server

	r := reporter.SetTestReporter() // set up test reporter
	ctx := ao.NewContext(context.Background(), ao.NewTrace("httpTest"))
	// make request to URL of test server
	url := fmt.Sprintf("http://127.0.0.1:%d/test?qs=1", port)
	if badReq {
		url = badURL // causes url.Parse() in http.NewRequest() to fail
	}
	resp, err := clientFn(t, ctx, method, url)
	ao.EndTrace(ctx)

	if badReq { // handle case where http.NewRequest() returned nil
		assert.Error(t, err)
		assert.Nil(t, resp)
		r.Close(2)
		g.AssertGraph(t, r.EventBufs, 2, g.AssertNodeMap{
			{"httpTest", "entry"}: {},
			{"httpTest", "exit"}:  {Edges: g.Edges{{"httpTest", "entry"}}},
		})
		return
	}
	// handle case where http.Client.Do() did not return an error
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	r.Close(server.numBufs)
	server.assertFn(t, r.EventBufs, resp, url, method, port, server.status)
}

// assert traces that hit testServer, which uses the HTTP server instrumentation.
func assertHTTPRequestGraph(t *testing.T, bufs [][]byte, resp *http.Response, url, method string, port, status int) {
	assert.Len(t, resp.Header[ao.HTTPHeaderName], 1)
	assert.Equal(t, status, resp.StatusCode)

	g.AssertGraph(t, bufs, 8, g.AssertNodeMap{
		{"httpTest", "entry"}: {},
		{"http.Client", "entry"}: {Edges: g.Edges{{"httpTest", "entry"}}, Callback: func(n g.Node) {
			assert.Equal(t, true, n.Map["IsService"])
			assert.Equal(t, url, n.Map["RemoteURL"])
		}},
		{"http.Client", "exit"}: {Edges: g.Edges{{"myHandler", "exit"}, {"http.Client", "entry"}}},
		{"myHandler", "entry"}: {Edges: g.Edges{{"http.Client", "entry"}}, Callback: func(n g.Node) {
			assert.Equal(t, "/test", n.Map["URL"])
			assert.Equal(t, fmt.Sprintf("127.0.0.1:%d", port), n.Map["HTTP-Host"])
			assert.Equal(t, "qs=1", n.Map["Query-String"])
			assert.Equal(t, method, n.Map["Method"])
		}},
		{"myHandler", "exit"}: {Edges: g.Edges{{"DBx", "exit"}, {"myHandler", "entry"}}, Callback: func(n g.Node) {
			assert.Equal(t, status, n.Map["Status"])
		}},
		{"DBx", "entry"}: {Edges: g.Edges{{"myHandler", "entry"}}, Callback: func(n g.Node) {
			assert.Equal(t, "SELECT *", n.Map["Query"])
			assert.Equal(t, "db.net", n.Map["RemoteHost"])
		}},
		{"DBx", "exit"}:      {Edges: g.Edges{{"DBx", "entry"}}},
		{"httpTest", "exit"}: {Edges: g.Edges{{"http.Client", "exit"}, {"httpTest", "entry"}}},
	})
}

// assert traces of an HTTP client to untraced servers testServer200 and testServer403.
func assertHTTPRequestUntracedGraph(t *testing.T, bufs [][]byte, resp *http.Response, url, method string, port, status int) {
	assert.NotContains(t, resp.Header[ao.HTTPHeaderName], "Header")
	assert.Equal(t, status, resp.StatusCode)

	g.AssertGraph(t, bufs, 4, g.AssertNodeMap{
		{"httpTest", "entry"}: {},
		{"http.Client", "entry"}: {Edges: g.Edges{{"httpTest", "entry"}}, Callback: func(n g.Node) {
			assert.Equal(t, true, n.Map["IsService"])
			assert.Equal(t, url, n.Map["RemoteURL"])
		}},
		{"http.Client", "exit"}: {Edges: g.Edges{{"http.Client", "entry"}}},
		{"httpTest", "exit"}:    {Edges: g.Edges{{"http.Client", "exit"}, {"httpTest", "entry"}}},
	})
}

// assert traces that hit an AO-wrapped, panicking http Handler.
func assertHTTPRequestPanic(t *testing.T, bufs [][]byte, resp *http.Response, url, method string, port, status int) {

	g.AssertGraph(t, bufs, 7, g.AssertNodeMap{
		{"httpTest", "entry"}: {},
		{"http.Client", "entry"}: {Edges: g.Edges{{"httpTest", "entry"}}, Callback: func(n g.Node) {
			assert.Equal(t, true, n.Map["IsService"])
			assert.Equal(t, url, n.Map["RemoteURL"])
		}},
		{"http.Client", "exit"}: {Edges: g.Edges{{"http.HandlerFunc", "exit"}, {"http.Client", "entry"}}},
		{"http.HandlerFunc", "entry"}: {Edges: g.Edges{{"http.Client", "entry"}}, Callback: func(n g.Node) {
			assert.Equal(t, "/test", n.Map["URL"])
			assert.Equal(t, fmt.Sprintf("127.0.0.1:%d", port), n.Map["HTTP-Host"])
			assert.Equal(t, "qs=1", n.Map["Query-String"])
			assert.Equal(t, method, n.Map["Method"])
		}},
		{"http.HandlerFunc", "error"}: {Edges: g.Edges{{"http.HandlerFunc", "entry"}}, Callback: func(n g.Node) {
			assert.Equal(t, "panic", n.Map["ErrorClass"])
			assert.Equal(t, "panicking!", n.Map["ErrorMsg"])
		}},
		{"http.HandlerFunc", "exit"}: {Edges: g.Edges{{"http.HandlerFunc", "error"}}, Callback: func(n g.Node) {
			assert.Equal(t, "ao_test", n.Map["Controller"])
			assert.Equal(t, "handlerPanic", n.Map["Action"])
			assert.Equal(t, status, n.Map["Status"])
		}},
		{"httpTest", "exit"}: {Edges: g.Edges{{"http.Client", "exit"}, {"httpTest", "entry"}}},
	})
}

func TestTraceHTTPError(t *testing.T)            { testTraceHTTPError(t, "GET", false, testHTTPClient) }
func TestTraceHTTPErrorA(t *testing.T)           { testTraceHTTPError(t, "GET", false, testHTTPClientA) }
func TestTraceHTTPErrorB(t *testing.T)           { testTraceHTTPError(t, "GET", false, testHTTPClientB) }
func TestTraceHTTPErrorBadRequest(t *testing.T)  { testTraceHTTPError(t, "GET", true, testHTTPClient) }
func TestTraceHTTPErrorABadRequest(t *testing.T) { testTraceHTTPError(t, "GET", true, testHTTPClientA) }
func TestTraceHTTPErrorBBadRequest(t *testing.T) { testTraceHTTPError(t, "GET", true, testHTTPClientB) }

// test making an HTTP request that causes http.Client.Do() to fail
func testTraceHTTPError(t *testing.T, method string, badReq bool, clientFn testClientFn) {
	r := reporter.SetTestReporter() // set up test reporter
	ctx := ao.NewContext(context.Background(), ao.NewTrace("httpTest"))
	url := invalidPortURL // make HTTP req to invalid port
	if badReq {
		url = badURL // causes url.Parse() in http.NewRequest() to fail
	}
	resp, err := clientFn(t, ctx, method, url)
	ao.EndTrace(ctx)

	assert.Error(t, err)
	assert.Nil(t, resp)

	if badReq { // handle case where http.NewRequest() returned nil
		r.Close(2)
		g.AssertGraph(t, r.EventBufs, 2, g.AssertNodeMap{
			{"httpTest", "entry"}: {},
			{"httpTest", "exit"}:  {Edges: g.Edges{{"httpTest", "entry"}}},
		})
		return
	}
	// handle case where http.Client.Do() returned an error
	r.Close(5)
	g.AssertGraph(t, r.EventBufs, 5, g.AssertNodeMap{
		{"httpTest", "entry"}: {},
		{"http.Client", "entry"}: {Edges: g.Edges{{"httpTest", "entry"}}, Callback: func(n g.Node) {
			assert.Equal(t, true, n.Map["IsService"])
			assert.Equal(t, url, n.Map["RemoteURL"])
		}},
		{"http.Client", "error"}: {Edges: g.Edges{{"http.Client", "entry"}}, Callback: func(n g.Node) {
			assert.Equal(t, "error", n.Map["ErrorClass"])
			assert.Contains(t, n.Map["ErrorMsg"], "dial tcp:")
			assert.Contains(t, n.Map["ErrorMsg"], "invalid port")
		}},
		{"http.Client", "exit"}: {Edges: g.Edges{{"http.Client", "error"}}},
		{"httpTest", "exit"}:    {Edges: g.Edges{{"http.Client", "exit"}, {"httpTest", "entry"}}},
	})
}

func TestDoubleWrappedHTTPRequest(t *testing.T) {
	list, err := net.Listen("tcp", ":0") // pick an unallocated port
	assert.NoError(t, err)
	port := list.Addr().(*net.TCPAddr).Port
	go testDoubleWrappedServer(t, list) // start test server

	r := reporter.SetTestReporter() // set up test reporter
	ctx := ao.NewContext(context.Background(), ao.NewTrace("httpTest"))
	url := fmt.Sprintf("http://127.0.0.1:%d/test?qs=1", port)
	resp, err := testHTTPClient(t, ctx, "GET", url)
	t.Logf("response: %v", resp)
	ao.EndTrace(ctx)

	assert.NoError(t, err)
	assert.Len(t, resp.Header[ao.HTTPHeaderName], 1)
	assert.Equal(t, 403, resp.StatusCode)

	r.Close(10)
	g.AssertGraph(t, r.EventBufs, 10, g.AssertNodeMap{
		{"httpTest", "entry"}: {},
		{"http.Client", "entry"}: {Edges: g.Edges{{"httpTest", "entry"}}, Callback: func(n g.Node) {
			assert.Equal(t, true, n.Map["IsService"])
			assert.Equal(t, url, n.Map["RemoteURL"])
		}},
		{"http.Client", "exit"}: {Edges: g.Edges{{"http.HandlerFunc", "exit"}, {"http.Client", "entry"}}},
		{"http.HandlerFunc", "entry"}: {Edges: g.Edges{{"http.Client", "entry"}}, Callback: func(n g.Node) {
			assert.Equal(t, "/test", n.Map["URL"])
			assert.Equal(t, fmt.Sprintf("127.0.0.1:%d", port), n.Map["HTTP-Host"])
			assert.Equal(t, "qs=1", n.Map["Query-String"])
			assert.Equal(t, "GET", n.Map["Method"])
		}},
		{"http.HandlerFunc", "exit"}: {Edges: g.Edges{{"myHandler", "exit"}, {"http.HandlerFunc", "entry"}}, Callback: func(n g.Node) {
			assert.Equal(t, 403, n.Map["Status"])
			assert.Equal(t, "ao_test", n.Map["Controller"])
			assert.Equal(t, "testDoubleWrappedServer.func1", n.Map["Action"])
		}},
		{"myHandler", "entry"}: {Edges: g.Edges{{"http.HandlerFunc", "entry"}}, Callback: func(n g.Node) {
			assert.Equal(t, "/test", n.Map["URL"])
			assert.Equal(t, fmt.Sprintf("127.0.0.1:%d", port), n.Map["HTTP-Host"])
			assert.Equal(t, "qs=1", n.Map["Query-String"])
			assert.Equal(t, "GET", n.Map["Method"])
		}},
		{"myHandler", "exit"}: {Edges: g.Edges{{"DBx", "exit"}, {"myHandler", "entry"}}, Callback: func(n g.Node) {
			assert.Equal(t, 403, n.Map["Status"])
		}},
		{"DBx", "entry"}: {Edges: g.Edges{{"myHandler", "entry"}}, Callback: func(n g.Node) {
			assert.Equal(t, "SELECT *", n.Map["Query"])
			assert.Equal(t, "db.net", n.Map["RemoteHost"])
		}},
		{"DBx", "exit"}:      {Edges: g.Edges{{"DBx", "entry"}}},
		{"httpTest", "exit"}: {Edges: g.Edges{{"http.Client", "exit"}, {"httpTest", "entry"}}},
	})
}

// based on examples/distributed_app
func AliceHandler(w http.ResponseWriter, r *http.Request) {
	// trace this request, overwriting w with wrapped ResponseWriter
	t, w, _ := ao.TraceFromHTTPRequestResponse("aliceHandler", w, r)
	ctx := ao.NewContext(context.Background(), t)
	defer t.End()

	// call an HTTP endpoint and propagate the distributed trace context
	url := "http://localhost:8081/bob"

	// create HTTP client and set trace metadata header
	httpClient := &http.Client{}
	httpReq, _ := http.NewRequest("GET", url, nil)
	// begin span for the client side of the HTTP service request
	l := ao.BeginHTTPClientSpan(ctx, httpReq)

	// make HTTP request to external API
	resp, err := httpClient.Do(httpReq)
	l.AddHTTPResponse(resp, err)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(fmt.Sprintf("err: %v", err)))
		l.End() // end HTTP client timing
		return
	}

	// read response body
	defer resp.Body.Close()
	buf, err := ioutil.ReadAll(resp.Body)
	l.End() // end HTTP client timing
	//w.WriteHeader(200)
	if err != nil {
		w.Write([]byte(`{"error":true}`))
	} else {
		w.Write(buf) // return API response to caller
	}
}

func BobHandler(w http.ResponseWriter, r *http.Request) {
	t, w, _ := ao.TraceFromHTTPRequestResponse("bobHandler", w, r)
	defer t.End()
	w.Write([]byte(`{"result":"hello from bob"}`))
}

func TestDistributedApp(t *testing.T) {
	r := reporter.SetTestReporter() // set up test reporter

	aliceLn, err := net.Listen("tcp", ":8080")
	assert.NoError(t, err)
	require.NotNil(t, aliceLn, "can't open port 8080")
	bobLn, err := net.Listen("tcp", ":8081")
	assert.NoError(t, err)
	require.NotNil(t, bobLn, "can't open port 8081")
	go func() {
		s := &http.Server{Handler: http.HandlerFunc(ao.HTTPHandler(AliceHandler))}
		assert.NoError(t, s.Serve(aliceLn))
	}()
	go func() {
		s := &http.Server{Handler: http.HandlerFunc(ao.HTTPHandler(BobHandler))}
		assert.NoError(t, s.Serve(bobLn))
	}()

	resp, err := http.Get("http://localhost:8080/alice")
	assert.NoError(t, err)
	defer resp.Body.Close()
	buf, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)
	t.Logf("Response: %v BUF %s", resp, buf)

	r.Close(10)
	g.AssertGraph(t, r.EventBufs, 10, g.AssertNodeKVMap{
		{"http.HandlerFunc", "entry", "URL", "/alice"}:         {},
		{"aliceHandler", "entry", "URL", "/alice"}:             {Edges: g.Edges{{"http.HandlerFunc", "entry"}}},
		{"http.Client", "entry", "", ""}:                       {Edges: g.Edges{{"aliceHandler", "entry"}}, Callback: func(n g.Node) {}},
		{"http.HandlerFunc", "entry", "URL", "/bob"}:           {Edges: g.Edges{{"http.Client", "entry"}}},
		{"bobHandler", "entry", "URL", "/bob"}:                 {Edges: g.Edges{{"http.HandlerFunc", "entry"}}},
		{"bobHandler", "exit", "", ""}:                         {Edges: g.Edges{{"bobHandler", "entry"}}},
		{"http.HandlerFunc", "exit", "Action", "BobHandler"}:   {Edges: g.Edges{{"bobHandler", "exit"}, {"http.HandlerFunc", "entry"}}},
		{"http.Client", "exit", "", ""}:                        {Edges: g.Edges{{"http.HandlerFunc", "exit"}, {"http.Client", "entry"}}},
		{"aliceHandler", "exit", "", ""}:                       {Edges: g.Edges{{"http.Client", "exit"}, {"aliceHandler", "entry"}}},
		{"http.HandlerFunc", "exit", "Action", "AliceHandler"}: {Edges: g.Edges{{"aliceHandler", "exit"}, {"http.HandlerFunc", "entry"}}},
	})
}

func concurrentAliceHandler(w http.ResponseWriter, r *http.Request) {
	// trace this request, overwriting w with wrapped ResponseWriter
	t, w, _ := ao.TraceFromHTTPRequestResponse("aliceHandler", w, r)
	ctx := ao.NewContext(context.Background(), t)
	t.SetAsync(true)
	defer t.End()

	// call an HTTP endpoint and propagate the distributed trace context
	urls := []string{
		"http://localhost:8083/A",
		"http://localhost:8083/B",
		"http://localhost:8083/C",
	}

	var wg sync.WaitGroup
	wg.Add(len(urls))
	var out []byte
	outCh := make(chan []byte)
	doneCh := make(chan struct{})
	go func() {
		for buf := range outCh {
			out = append(out, buf...)
		}
		close(doneCh)
	}()
	for _, u := range urls {
		go func(url string) {
			// create HTTP client and set trace metadata header
			client := &http.Client{}
			req, _ := http.NewRequest("GET", url, nil)
			// begin span for the client side of the HTTP service request
			l := ao.BeginHTTPClientSpan(ctx, req)

			// make HTTP request to external API
			resp, err := client.Do(req)
			l.AddHTTPResponse(resp, err)
			if err != nil {
				l.End() // end HTTP client timing
				w.WriteHeader(500)
				return
			}
			// read response body
			defer resp.Body.Close()
			buf, err := ioutil.ReadAll(resp.Body)
			l.End() // end HTTP client timing
			if err != nil {
				outCh <- []byte(fmt.Sprintf(`{"error":"%v"}`, err))
			} else {
				outCh <- buf
			}
			wg.Done()
		}(u)
	}
	wg.Wait()
	close(outCh)
	<-doneCh

	w.Write(out)
}

func TestConcurrentApp(t *testing.T) {
	r := reporter.SetTestReporter() // set up test reporter

	aliceLn, err := net.Listen("tcp", ":8082")
	assert.NoError(t, err)
	bobLn, err := net.Listen("tcp", ":8083")
	assert.NoError(t, err)
	go func() {
		s := &http.Server{Handler: http.HandlerFunc(concurrentAliceHandler)}
		assert.NoError(t, s.Serve(aliceLn))
	}()
	go func() {
		s := &http.Server{Handler: http.HandlerFunc(BobHandler)}
		assert.NoError(t, s.Serve(bobLn))
	}()

	resp, err := http.Get("http://localhost:8082/alice")
	assert.NoError(t, err)
	defer resp.Body.Close()
	buf, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)
	t.Logf("Response: %v BUF %s", resp, buf)

	r.Close(14)
	g.AssertGraph(t, r.EventBufs, 14, g.AssertNodeKVMap{
		{"aliceHandler", "entry", "URL", "/alice"}:                       {},
		{"http.Client", "entry", "RemoteURL", "http://localhost:8083/A"}: {Edges: g.Edges{{"aliceHandler", "entry"}}},
		{"http.Client", "entry", "RemoteURL", "http://localhost:8083/B"}: {Edges: g.Edges{{"aliceHandler", "entry"}}},
		{"http.Client", "entry", "RemoteURL", "http://localhost:8083/C"}: {Edges: g.Edges{{"aliceHandler", "entry"}}},
		{"bobHandler", "entry", "URL", "/A"}:                             {Edges: g.Edges{{"http.Client", "entry"}}},
		{"bobHandler", "entry", "URL", "/B"}:                             {Edges: g.Edges{{"http.Client", "entry"}}},
		{"bobHandler", "entry", "URL", "/C"}:                             {Edges: g.Edges{{"http.Client", "entry"}}},
		{"bobHandler", "exit", "", ""}:                                   {Edges: g.Edges{{"bobHandler", "entry"}}, Count: 3},
		{"http.Client", "exit", "", ""}: {
			Edges: g.Edges{{"bobHandler", "exit"}, {"http.Client", "entry"}}, Count: 3, Callback: func(n g.Node) {
				assert.EqualValues(t, 200, n.Map["RemoteStatus"])
			}},
		{"aliceHandler", "exit", "", ""}: {
			Edges: g.Edges{{"http.Client", "exit"}, {"http.Client", "exit"}, {"http.Client", "exit"}, {"aliceHandler", "entry"}}, Callback: func(n g.Node) {
				assert.Equal(t, true, n.Map["Async"])
				assert.EqualValues(t, 200, n.Map["Status"])
			}},
	})
}

func TestConcurrentAppNoTrace(t *testing.T) {
	r := reporter.SetTestReporter(reporter.TestReporterDisableTracing())

	aliceLn, err := net.Listen("tcp", ":8084")
	assert.NoError(t, err)
	bobLn, err := net.Listen("tcp", ":8085")
	assert.NoError(t, err)
	go func() {
		s := &http.Server{Handler: http.HandlerFunc(concurrentAliceHandler)}
		assert.NoError(t, s.Serve(aliceLn))
	}()
	go func() {
		s := &http.Server{Handler: http.HandlerFunc(BobHandler)}
		assert.NoError(t, s.Serve(bobLn))
	}()

	resp, err := http.Get("http://localhost:8084/alice")
	assert.NoError(t, err)
	defer resp.Body.Close()
	buf, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)
	assert.NotNil(t, buf)

	// shouldn't report anything
	assert.Len(t, r.EventBufs, 0)
}
