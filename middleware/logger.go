package middleware

import (
	"bytes"
	"fmt"
	"io"
	"time"

	"github.com/gin-gonic/gin"
)

func Logger() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(params gin.LogFormatterParams) string {
		var statusColor, methodColor, resetColor string
		if params.IsOutputColor() {
			statusColor = params.StatusCodeColor()
			methodColor = params.MethodColor()
			resetColor = params.ResetColor()
		}

		var buf bytes.Buffer

		// Request info
		buf.WriteString(fmt.Sprintf("%s - [%s] \"%s %s %s\" ",
			params.ClientIP,
			params.TimeStamp.Format(time.RFC3339),
			methodColor, params.Method, resetColor,
		))
		buf.WriteString(fmt.Sprintf("%s ", params.Path))

		// Response info
		buf.WriteString(fmt.Sprintf("%s %3d %s ",
			statusColor, params.StatusCode, resetColor,
		))

		// Latency
		buf.WriteString(fmt.Sprintf("%13v ", params.Latency))

		if params.ErrorMessage != "" {
			buf.WriteString(fmt.Sprintf("\"%s\" ", params.ErrorMessage))
		}

		buf.WriteString(fmt.Sprintf("%s %s",
			params.Request.UserAgent(),
			params.Request.Header.Get("X-Request-Id"),
		))

		if params.Request.Method == "POST" || params.Request.Method == "PUT" {
			body, _ := io.ReadAll(params.Request.Body)
			params.Request.Body = io.NopCloser(bytes.NewBuffer(body))
			if len(body) > 0 {
				buf.WriteString(fmt.Sprintf("\nRequest Body: %s", string(body)))
			}
		}

		return buf.String() + "\n"
	})
}
