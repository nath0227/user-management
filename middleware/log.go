package middleware

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"time"

	echo "github.com/labstack/echo/v4"
)

type CustomResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

func (w CustomResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}
func LoggingMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		start := time.Now()
		reqBody := new(bytes.Buffer)
		reqBody.ReadFrom(c.Request().Body)
		log.Printf("[REQUEST]: %s\n", reqBody.String())
		c.Request().Body = io.NopCloser(reqBody)
		resBody := new(bytes.Buffer)
		multiWriter := io.MultiWriter(c.Response().Writer, resBody)
		writer := &CustomResponseWriter{Writer: multiWriter, ResponseWriter: c.Response().Writer}
		c.Response().Writer = writer
		err := next(c)
		end := time.Now()

		log.Printf("[RESPONSE]: %s\n", resBody.String())
		log.Printf("HTTP Method: %s, API path: %s, Latency: %s\n", c.Request().Method, c.Path(), end.Sub(start))
		return err
	}
}
