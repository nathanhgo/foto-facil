package api

import (
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gorilla/websocket"
)

// dialTestWebSocket spins up the real HandleWebSocket handler behind an
// httptest server and returns a connected client plus a cleanup function.
func dialTestWebSocket(t *testing.T) (*websocket.Conn, func()) {
	t.Helper()

	server := httptest.NewServer(http.HandlerFunc(HandleWebSocket))
	dialURL := "ws" + strings.TrimPrefix(server.URL, "http")

	conn, _, err := websocket.DefaultDialer.Dial(dialURL, nil)
	if err != nil {
		server.Close()
		t.Fatalf("failed to dial test websocket: %v", err)
	}

	return conn, func() {
		conn.Close()
		server.Close()
	}
}

func writeRunFlow(t *testing.T, conn *websocket.Conn, requestID, inputPath, outputDir string) {
	t.Helper()

	request := map[string]any{
		"action":    "RUN_FLOW",
		"requestId": requestID,
		"flow": map[string]any{
			"nodes": []map[string]any{
				{"id": "input", "data": map[string]any{"originalType": "Image Input", "filePaths": inputPath}},
				{"id": "output", "data": map[string]any{"originalType": "Image Output", "outputDir": outputDir}},
			},
			"edges": []map[string]any{
				{"source": "input", "target": "output"},
			},
		},
	}

	raw, err := json.Marshal(request)
	if err != nil {
		t.Fatalf("failed to marshal RUN_FLOW request: %v", err)
	}
	if err := conn.WriteMessage(websocket.TextMessage, raw); err != nil {
		t.Fatalf("failed to write RUN_FLOW request: %v", err)
	}
}

func readRunFlowResponse(t *testing.T, conn *websocket.Conn) RunFlowResponse {
	t.Helper()

	_, raw, err := conn.ReadMessage()
	if err != nil {
		t.Fatalf("failed to read RUN_FLOW response: %v", err)
	}

	var response RunFlowResponse
	if err := json.Unmarshal(raw, &response); err != nil {
		t.Fatalf("failed to unmarshal RUN_FLOW response %q: %v", string(raw), err)
	}
	return response
}

func writeTestPNG(t *testing.T, dir, name string) string {
	t.Helper()

	img := image.NewRGBA(image.Rect(0, 0, 4, 4))
	for y := 0; y < 4; y++ {
		for x := 0; x < 4; x++ {
			img.Set(x, y, color.RGBA{10, 20, 30, 255})
		}
	}

	path := filepath.Join(dir, name)
	f, err := os.Create(path)
	if err != nil {
		t.Fatalf("failed to create test image: %v", err)
	}
	defer f.Close()
	if err := png.Encode(f, img); err != nil {
		t.Fatalf("failed to encode test image: %v", err)
	}
	return path
}

func TestRunFlowEchoesRequestIDOnSuccess(t *testing.T) {
	conn, cleanup := dialTestWebSocket(t)
	defer cleanup()

	tmp := t.TempDir()
	input := writeTestPNG(t, tmp, "input.png")

	writeRunFlow(t, conn, "req-success", input, filepath.Join(tmp, "out"))
	response := readRunFlowResponse(t, conn)

	if response.Action != "RUN_FLOW" {
		t.Errorf("expected action RUN_FLOW, got %q", response.Action)
	}
	if response.RequestID != "req-success" {
		t.Errorf("expected requestId echo %q, got %q", "req-success", response.RequestID)
	}
	if response.Status != "success" {
		t.Errorf("expected success status, got %q (error: %q)", response.Status, response.Error)
	}
}

func TestRunFlowEchoesRequestIDOnExecutionError(t *testing.T) {
	conn, cleanup := dialTestWebSocket(t)
	defer cleanup()

	tmp := t.TempDir()
	missing := filepath.Join(tmp, "does-not-exist.png")

	writeRunFlow(t, conn, "req-error", missing, filepath.Join(tmp, "out"))
	response := readRunFlowResponse(t, conn)

	if response.Action != "RUN_FLOW" {
		t.Errorf("expected action RUN_FLOW, got %q", response.Action)
	}
	if response.RequestID != "req-error" {
		t.Errorf("expected requestId echo %q, got %q", "req-error", response.RequestID)
	}
	if response.Status != "error" {
		t.Errorf("expected error status, got %q", response.Status)
	}
	if response.Error == "" {
		t.Error("expected a non-empty error message for a failed run")
	}
}

func TestRunFlowWithoutRequestIDStillResponds(t *testing.T) {
	conn, cleanup := dialTestWebSocket(t)
	defer cleanup()

	tmp := t.TempDir()
	input := writeTestPNG(t, tmp, "input.png")

	writeRunFlow(t, conn, "", input, filepath.Join(tmp, "out"))
	response := readRunFlowResponse(t, conn)

	if response.Action != "RUN_FLOW" {
		t.Errorf("expected action RUN_FLOW, got %q", response.Action)
	}
	if response.Status != "success" {
		t.Errorf("expected success status, got %q", response.Status)
	}
	if response.RequestID != "" {
		t.Errorf("expected empty requestId for preview run, got %q", response.RequestID)
	}
}

// Ensure the test helper actually produces a decodable image and that the
// response type stays JSON-compatible with the frontend contract.
func TestRunFlowResponseJSONFields(t *testing.T) {
	raw, err := json.Marshal(RunFlowResponse{Action: "RUN_FLOW", RequestID: "abc", Status: "success"})
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}
	var decoded map[string]any
	if err := json.Unmarshal(raw, &decoded); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}
	for _, key := range []string{"action", "requestId", "status"} {
		if _, ok := decoded[key]; !ok {
			t.Errorf("expected response JSON to contain %q, got %s", key, fmt.Sprint(decoded))
		}
	}
}
