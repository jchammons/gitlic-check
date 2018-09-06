// Copyright (C) 2016 Librato, Inc. All rights reserved.

package ao_test

import (
	"testing"

	"github.com/appoptics/appoptics-apm-go/v1/ao"
	g "github.com/appoptics/appoptics-apm-go/v1/ao/internal/graphtest"
	"github.com/appoptics/appoptics-apm-go/v1/ao/internal/reporter"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
)

func testProf(ctx context.Context) {
	ao.BeginProfile(ctx, "testProf")
}

func TestBeginProfile(t *testing.T) {
	r := reporter.SetTestReporter()
	ctx := ao.NewContext(context.Background(), ao.NewTrace("testSpan"))
	testProf(ctx)

	r.Close(2)
	g.AssertGraph(t, r.EventBufs, 2, g.AssertNodeMap{
		{"testSpan", "entry"}: {},
		{"", "profile_entry"}: {Edges: g.Edges{{"testSpan", "entry"}}, Callback: func(n g.Node) {
			assert.Equal(t, n.Map["Language"], "go")
			assert.Equal(t, n.Map["ProfileName"], "testProf")
			assert.Equal(t, n.Map["FunctionName"], "github.com/appoptics/appoptics-apm-go/v1/ao_test.testProf")
			assert.Contains(t, n.Map["File"], "/appoptics-apm-go/v1/ao/profile_test.go")
		}},
	})
}

func testSpanProf(ctx context.Context) {
	l1, _ := ao.BeginSpan(ctx, "L1")
	p := l1.BeginProfile("testSpanProf")
	p.End()
	l1.End()
	ao.EndTrace(ctx)
}

func TestBeginSpanProfile(t *testing.T) {
	r := reporter.SetTestReporter()
	ctx := ao.NewContext(context.Background(), ao.NewTrace("testSpan"))
	testSpanProf(ctx)

	r.Close(6)
	g.AssertGraph(t, r.EventBufs, 6, g.AssertNodeMap{
		{"testSpan", "entry"}: {},
		{"L1", "entry"}:       {Edges: g.Edges{{"testSpan", "entry"}}},
		{"", "profile_entry"}: {Edges: g.Edges{{"L1", "entry"}}, Callback: func(n g.Node) {
			assert.Equal(t, n.Map["Language"], "go")
			assert.Equal(t, n.Map["ProfileName"], "testSpanProf")
			assert.Equal(t, n.Map["FunctionName"], "github.com/appoptics/appoptics-apm-go/v1/ao_test.testSpanProf")
			assert.Contains(t, n.Map["File"], "/appoptics-apm-go/v1/ao/profile_test.go")
		}},
		{"", "profile_exit"}: {Edges: g.Edges{{"", "profile_entry"}}, Callback: func(n g.Node) {
			assert.Equal(t, n.Map["Language"], "go")
			assert.Equal(t, n.Map["ProfileName"], "testSpanProf")
		}},
		{"L1", "exit"}:       {Edges: g.Edges{{"", "profile_exit"}, {"L1", "entry"}}},
		{"testSpan", "exit"}: {Edges: g.Edges{{"L1", "exit"}, {"testSpan", "entry"}}},
	})

}

// ensure above tests run smoothly with no events reported when a context has no trace
func TestNoTraceBeginProfile(t *testing.T) {
	r := reporter.SetTestReporter()
	ctx := context.Background()
	testProf(ctx)
	r.Close(0)
	assert.Len(t, r.EventBufs, 0)
}
func TestTraceErrorBeginProfile(t *testing.T) {
	// simulate reporter error on second event: prevents Span from being reported
	r := reporter.SetTestReporter()
	r.ErrorEvents = map[int]bool{1: true}
	testProf(ao.NewContext(context.Background(), ao.NewTrace("testSpan")))
	r.Close(1)
	g.AssertGraph(t, r.EventBufs, 1, g.AssertNodeMap{
		{"testSpan", "entry"}: {},
	})
}

func TestNoTraceBeginSpanProfile(t *testing.T) {
	r := reporter.SetTestReporter()
	ctx := context.Background()
	testSpanProf(ctx)
	r.Close(0)
	assert.Len(t, r.EventBufs, 0)
}
func TestTraceErrorBeginSpanProfile(t *testing.T) {
	// simulate reporter error on second event: prevents nested Span & Profile spans
	r := reporter.SetTestReporter()
	r.ErrorEvents = map[int]bool{1: true}
	testSpanProf(ao.NewContext(context.Background(), ao.NewTrace("testSpan")))
	r.Close(2)
	g.AssertGraph(t, r.EventBufs, 2, g.AssertNodeMap{
		{"testSpan", "entry"}: {},
		{"testSpan", "exit"}:  {Edges: g.Edges{{"testSpan", "entry"}}},
	})
}
