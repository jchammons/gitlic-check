// Copyright (C) 2016 Librato, Inc. All rights reserved.
// test usage example from doc.go

package ao_test

import (
	"errors"
	"testing"

	"github.com/appoptics/appoptics-apm-go/v1/ao"
	g "github.com/appoptics/appoptics-apm-go/v1/ao/internal/graphtest"
	"github.com/appoptics/appoptics-apm-go/v1/ao/internal/reporter"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
)

func testDocSpanExample() {
	// create trace and bind to new context
	ctx := ao.NewContext(context.Background(), ao.NewTrace("myApp"))
	// create new span for this trace
	l, ctxL := ao.BeginSpan(ctx, "mySpan")

	// profile a slow part of a span
	p := ao.BeginProfile(ctxL, "slowFunc")
	// slowFunc(x, y)
	p.End()

	// Start a new span, given a parent span
	db1L := l.BeginSpan("myDB1", "Query", "SELECT * FROM tbl1")
	// perform a query
	db1L.End()

	// Start a new span, given a context.Context
	db2L, _ := ao.BeginSpan(ctxL, "myDB2", "Query", "SELECT * FROM tbl2")
	// perform a query
	db2L.End()

	l.End()
	ao.Err(ctx, errors.New("Got bad error!"))
	ao.EndTrace(ctx)
}

func testDocSpanExampleCtx() {
	// create trace and bind to new context
	ctx := ao.NewContext(context.Background(), ao.NewTrace("myApp"))
	// create new span for this trace
	_, ctxL := ao.BeginSpan(ctx, "mySpan")

	// profile a nested block or function call
	slowFunc := func() {
		defer ao.BeginProfile(ctxL, "slowFunc").End()
		// ... do something slow
	}
	slowFunc()

	// Start a new span, given a parent span
	_, ctxQ1 := ao.BeginSpan(ctxL, "myDB1", "Query", "SELECT * FROM tbl1")
	// perform a query
	ao.End(ctxQ1)

	// Start a new span, given a context.Context
	_, ctxQ2 := ao.BeginSpan(ctxL, "myDB2", "Query", "SELECT * FROM tbl2")
	// perform a query
	ao.End(ctxQ2)

	ao.End(ctxL)
	ao.Err(ctx, errors.New("Got bad error!"))
	ao.EndTrace(ctx)
}

func TestDocSpanExample(t *testing.T) {
	r := reporter.SetTestReporter()
	testDocSpanExample()
	r.Close(11)
	assertDocSpanExample(t, r.EventBufs)
}
func TestDocSpanExampleCtx(t *testing.T) {
	r := reporter.SetTestReporter()
	testDocSpanExampleCtx()
	r.Close(11)
	assertDocSpanExample(t, r.EventBufs)
}
func assertDocSpanExample(t *testing.T, bufs [][]byte) {
	g.AssertGraph(t, bufs, 11, g.AssertNodeMap{
		{"myApp", "entry"}:  {},
		{"mySpan", "entry"}: {Edges: g.Edges{{"myApp", "entry"}}},
		{"", "profile_entry"}: {Edges: g.Edges{{"mySpan", "entry"}}, Callback: func(n g.Node) {
			assert.Equal(t, n.Map["ProfileName"], "slowFunc")
		}},
		{"", "profile_exit"}: {Edges: g.Edges{{"", "profile_entry"}}},
		{"myDB1", "entry"}: {Edges: g.Edges{{"mySpan", "entry"}}, Callback: func(n g.Node) {
			assert.Equal(t, "SELECT * FROM tbl1", n.Map["Query"])
		}},
		{"myDB1", "exit"}: {Edges: g.Edges{{"myDB1", "entry"}}},
		{"myDB2", "entry"}: {Edges: g.Edges{{"mySpan", "entry"}}, Callback: func(n g.Node) {
			assert.Equal(t, "SELECT * FROM tbl2", n.Map["Query"])
		}},
		{"myDB2", "exit"}:  {Edges: g.Edges{{"myDB2", "entry"}}},
		{"mySpan", "exit"}: {Edges: g.Edges{{"", "profile_exit"}, {"myDB1", "exit"}, {"myDB2", "exit"}, {"mySpan", "entry"}}},
		{"myApp", "error"}: {Edges: g.Edges{{"myApp", "entry"}}, Callback: func(n g.Node) {
			assert.Equal(t, "error", n.Map["ErrorClass"])
			assert.Equal(t, "Got bad error!", n.Map["ErrorMsg"])
		}},
		{"myApp", "exit"}: {Edges: g.Edges{{"mySpan", "exit"}, {"myApp", "error"}}},
	})
}
