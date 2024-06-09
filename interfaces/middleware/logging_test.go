package middleware

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/tusmasoma/connectHub-backend/internal/log"
)

func Test_Logging(t *testing.T) {
	// ログの出力をキャプチャするためのバッファ
	var logBuf bytes.Buffer
	log.SetOutput(&logBuf)

	tests := []struct {
		name            string
		handler         http.HandlerFunc
		expectedStatus  int
		expectAccessLog bool
		expectErrorLog  bool
	}{
		{
			name: "successful request",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("OK")) //nolint:errcheck // ignore error
			},
			expectedStatus:  http.StatusOK,
			expectAccessLog: true,
			expectErrorLog:  false,
		},
		{
			name: "client error request",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("Bad Request")) //nolint:errcheck // ignore error
			},
			expectedStatus:  http.StatusBadRequest,
			expectAccessLog: true,
			expectErrorLog:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// ロギングミドルウェアを適用
			handlerToTest := Logging(tt.handler)

			// テストリクエストを作成
			req := httptest.NewRequest(http.MethodGet, "http://example.com/foo", nil)
			w := httptest.NewRecorder()

			// ハンドラを呼び出す
			handlerToTest.ServeHTTP(w, req)

			// レスポンスステータスコードのチェック
			if status := w.Result().StatusCode; status != tt.expectedStatus {
				t.Errorf("expected status code %d, got %d", tt.expectedStatus, status)
			}

			// ログの出力を確認
			logOutput := logBuf.String()
			if tt.expectAccessLog && !strings.Contains(logOutput, "Access log") {
				t.Errorf("expected access log, got %s", logOutput)
			}

			if tt.expectErrorLog && !strings.Contains(logOutput, "Error log") {
				t.Errorf("expected error log, got %s", logOutput)
			}

			// ログバッファをリセット
			logBuf.Reset()
		})
	}
}
