// Copyright (C) 2016 Librato, Inc. All rights reserved.

package ao_test

import (
	"github.com/appoptics/appoptics-apm-go/v1/ao"
	"golang.org/x/net/context"
)

func ExampleNewTrace() {
	f0 := func(ctx context.Context) { // example span
		l, _ := ao.BeginSpan(ctx, "myDB",
			"Query", "SELECT * FROM tbl1",
			"RemoteHost", "db1.com")
		// ... run a query ...
		l.End()
	}

	// create a new trace, and a context to carry it around
	ctx := ao.NewContext(context.Background(), ao.NewTrace("myExample"))
	// do some work
	f0(ctx)
	// end the trace
	ao.EndTrace(ctx)
}

func ExampleBeginSpan() {
	// create trace and bind to context, reporting first event
	ctx := ao.NewContext(context.Background(), ao.NewTrace("baseSpan"))
	// ... do something ...

	// instrument a DB query
	l, _ := ao.BeginSpan(ctx, "DBx", "Query", "SELECT * FROM tbl")
	// .. execute query ..
	l.End()

	// end trace
	ao.EndTrace(ctx)
}
