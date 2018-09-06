// Copyright (C) 2016 Librato, Inc. All rights reserved.

package ao

import (
	"reflect"
	"testing"

	g "github.com/appoptics/appoptics-apm-go/v1/ao/internal/graphtest"
	"github.com/appoptics/appoptics-apm-go/v1/ao/internal/reporter"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
)

func TestContext(t *testing.T) {
	r := reporter.SetTestReporter()

	ctx := context.Background()
	assert.Empty(t, MetadataString(ctx))
	tr := NewTrace("test").(*aoTrace)

	xt := tr.aoCtx.MetadataString()
	//	assert.True(t, IsSampled(ctx), "%T", tr.aoCtx)

	var traceKey = struct{}{}

	ctx2 := context.WithValue(ctx, traceKey, tr)
	assert.Equal(t, ctx2.Value(traceKey), tr)
	assert.Equal(t, ctx2.Value(traceKey).(*aoTrace).aoCtx.MetadataString(), xt)

	ctxx := tr.aoCtx.Copy()
	lbl := spanLabeler{"L1"}
	tr2 := &aoTrace{layerSpan: layerSpan{span: span{aoCtx: ctxx, labeler: lbl}}}
	ctx3 := context.WithValue(ctx2, traceKey, tr2)
	assert.Equal(t, ctx3.Value(traceKey), tr2)

	ctxx2 := tr2.aoCtx.Copy()
	tr3 := &aoTrace{layerSpan: layerSpan{span: span{aoCtx: ctxx2}}}
	ctx4 := context.WithValue(ctx3, traceKey, tr3)
	assert.Equal(t, ctx4.Value(traceKey), tr3)

	r.Close(1)
	g.AssertGraph(t, r.EventBufs, 1, g.AssertNodeMap{{"test", "entry"}: {}})
}

func TestTraceFromContext(t *testing.T) {
	r := reporter.SetTestReporter()
	tr := NewTrace("TestTFC")
	ctx := NewContext(context.Background(), tr)
	trFC := TraceFromContext(ctx)
	assert.Equal(t, tr.ExitMetadata(), trFC.ExitMetadata())
	assert.Len(t, tr.ExitMetadata(), 60)

	trN := TraceFromContext(context.Background()) // no trace bound to this ctx
	assert.Len(t, trN.ExitMetadata(), 0)

	r.Close(1)
	g.AssertGraph(t, r.EventBufs, 1, g.AssertNodeMap{{"TestTFC", "entry"}: {}})
}

func TestContextIsSampled(t *testing.T) {
	// no context: not sampled
	assert.False(t, IsSampled(context.Background()))
	// sampled context
	_ = reporter.SetTestReporter()
	tr := NewTrace("TestTFC")
	ctx := NewContext(context.Background(), tr)
	assert.True(t, IsSampled(ctx))
}

func TestNullSpan(t *testing.T) {
	// enable reporting to test reporter
	r := reporter.SetTestReporter()

	ctx := NewContext(context.Background(), NewTrace("TestNullSpan")) // reports event
	l1, ctxL := BeginSpan(ctx, "L1")                                  // reports event
	assert.True(t, l1.IsReporting())
	assert.Equal(t, l1.MetadataString(), MetadataString(ctxL))
	assert.Len(t, l1.MetadataString(), 60)

	l1.End() // reports event
	assert.False(t, l1.IsReporting())
	assert.Empty(t, l1.MetadataString())

	p1 := l1.BeginProfile("P2") // try to start profile after end: no effect
	p1.End()

	c1 := l1.BeginSpan("C1") // child after parent ended
	assert.IsType(t, c1, nullSpan{})
	assert.False(t, c1.IsReporting())
	assert.False(t, c1.IsSampled())
	assert.False(t, c1.ok())
	assert.Empty(t, c1.MetadataString())
	c1.addChildEdge(l1.aoContext())
	c1.addProfile(p1)

	nctx := c1.aoContext()
	assert.Equal(t, reflect.TypeOf(nctx).Elem().Name(), "nullContext")
	assert.IsType(t, reflect.TypeOf(nctx.Copy()).Elem().Name(), "nullContext")

	r.Close(3)
	g.AssertGraph(t, r.EventBufs, 3, g.AssertNodeMap{
		{"TestNullSpan", "entry"}: {},
		{"L1", "entry"}:           {Edges: g.Edges{{"TestNullSpan", "entry"}}},
		{"L1", "exit"}:            {Edges: g.Edges{{"L1", "entry"}}},
	})
}
