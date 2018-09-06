package main

import (
	"bufio"
	"net"

	"github.com/appoptics/appoptics-apm-go/v1/ao"
	"github.com/gin-gonic/gin"
	"golang.org/x/net/context"
)

const (
	ginContextKey = "AppOptics"
	ginSpanName   = "gin"
)

func tracer() gin.HandlerFunc {
	return func(c *gin.Context) {
		t, w, _ := ao.TraceFromHTTPRequestResponse(ginSpanName, c.Writer, c.Request)
		c.Writer = &ginResponseWriter{w.(*ao.HTTPResponseWriter), c.Writer}
		t.SetTransactionName(c.HandlerName())
		defer t.End()
		// create a context.Context and bind it to the gin.Context
		c.Set(ginContextKey, ao.NewContext(context.Background(), t))
		// Pass to the next handler
		c.Next()
	}
}

// ginResponseWriter satisfies the gin.ResponseWriter interface
type ginResponseWriter struct {
	// handles Write, WriteHeader, Header (by calling wrapped gin writer)
	*ao.HTTPResponseWriter
	// handles all other gin.ResponseWriter methods
	ginWriter gin.ResponseWriter
}

func (w *ginResponseWriter) CloseNotify() <-chan bool                     { return w.ginWriter.CloseNotify() }
func (w *ginResponseWriter) Flush()                                       { w.ginWriter.Flush() }
func (w *ginResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) { return w.ginWriter.Hijack() }
func (w *ginResponseWriter) Size() int                                    { return w.ginWriter.Size() }
func (w *ginResponseWriter) Written() bool                                { return w.ginWriter.Written() }
func (w *ginResponseWriter) WriteString(s string) (int, error)            { return w.ginWriter.WriteString(s) }
func (w *ginResponseWriter) Status() int                                  { return w.StatusCode }
func (w *ginResponseWriter) WriteHeaderNow() {
	if !w.WroteHeader {
		w.WriteHeader(w.StatusCode)
	}
}
