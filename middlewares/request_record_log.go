package middlewares

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/sanyuanya/dongle/tools"
)

func RecordLog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Read the request body
		body, err := io.ReadAll(r.Body)
		if err != nil {
			response := tools.Response{
				Code:    49999,
				Message: fmt.Sprintf("Error reading body: %v", err),
				Result:  struct{}{},
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
			return
		}

		tools.Logger.Info("Request:",
			"method", r.Method,
			"path", r.URL.Path,
			"headers", r.Header,
			"body", string(body),
		)

		// Set the body back to the request
		r.Body = io.NopCloser(bytes.NewReader(body))

		next.ServeHTTP(w, r)
	})
}
