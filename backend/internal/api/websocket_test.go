package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

// O requestId é o que permite à UI saber que uma resposta é do pedido que ela
// fez. Se o backend parar de ecoar, o botão "Rodar Fluxo" volta a marcar
// concluído sozinho (o bug que este contrato existe para impedir).

func dialHandler(t *testing.T) (*websocket.Conn, func()) {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(HandleWebSocket))
	conn, _, err := websocket.DefaultDialer.Dial("ws"+strings.TrimPrefix(srv.URL, "http"), nil)
	if err != nil {
		srv.Close()
		t.Fatalf("falha ao conectar no WebSocket de teste: %v", err)
	}
	return conn, func() { conn.Close(); srv.Close() }
}

func readResponse(t *testing.T, conn *websocket.Conn) map[string]any {
	t.Helper()
	conn.SetReadDeadline(time.Now().Add(5 * time.Second))
	_, data, err := conn.ReadMessage()
	if err != nil {
		t.Fatalf("erro lendo resposta: %v", err)
	}
	var out map[string]any
	if err := json.Unmarshal(data, &out); err != nil {
		t.Fatalf("resposta não é JSON: %v (%s)", err, data)
	}
	return out
}

func brightnessNode(id string) map[string]any {
	return map[string]any{"id": id, "data": map[string]any{"originalType": "Brightness & Contrast"}}
}

// Ciclo no grafo: erro, e com o requestId ecoado.
func TestRunFlowEchoesRequestIdOnGraphCycle(t *testing.T) {
	conn, cleanup := dialHandler(t)
	defer cleanup()

	payload := map[string]any{
		"action":    "RUN_FLOW",
		"requestId": "req-do-teste",
		"flow": map[string]any{
			"nodes": []map[string]any{brightnessNode("a"), brightnessNode("b")},
			"edges": []map[string]any{
				{"source": "a", "target": "b"},
				{"source": "b", "target": "a"},
			},
		},
	}
	if err := conn.WriteJSON(payload); err != nil {
		t.Fatalf("erro enviando RUN_FLOW: %v", err)
	}

	resp := readResponse(t, conn)
	if resp["status"] != "error" {
		t.Fatalf("grafo com ciclo deveria responder error, veio: %+v", resp)
	}
	if resp["requestId"] != "req-do-teste" {
		t.Fatalf("requestId não foi ecoado na resposta de erro: %+v", resp)
	}
}

// Tipo de nó desconhecido: antes era ignorado em silêncio e a execução
// terminava como sucesso sem ter processado o nó.
func TestRunFlowWithUnknownNodeTypeFailsInsteadOfSilentSuccess(t *testing.T) {
	conn, cleanup := dialHandler(t)
	defer cleanup()

	payload := map[string]any{
		"action":    "RUN_FLOW",
		"requestId": "req-no-desconhecido",
		"flow": map[string]any{
			"nodes": []map[string]any{{"id": "x", "data": map[string]any{"originalType": "Nó Que Não Existe"}}},
			"edges": []map[string]any{},
		},
	}
	if err := conn.WriteJSON(payload); err != nil {
		t.Fatalf("erro enviando RUN_FLOW: %v", err)
	}

	resp := readResponse(t, conn)
	if resp["status"] != "error" {
		t.Fatalf("nó de tipo desconhecido deveria falhar, veio: %+v", resp)
	}
	if resp["requestId"] != "req-no-desconhecido" {
		t.Fatalf("requestId não foi ecoado: %+v", resp)
	}
	if msg, _ := resp["error"].(string); !strings.Contains(msg, "não suportado") {
		t.Fatalf("erro deveria dizer que o nó não é suportado, veio: %q", msg)
	}
}

// Preview (autosave) manda RUN_FLOW sem requestId: a resposta não pode carregar
// id nenhum, senão ela passaria por resposta do pedido do botão.
func TestRunFlowWithoutRequestIdDoesNotInventOne(t *testing.T) {
	conn, cleanup := dialHandler(t)
	defer cleanup()

	payload := map[string]any{
		"action": "RUN_FLOW",
		"flow":   map[string]any{"nodes": []map[string]any{}, "edges": []map[string]any{}},
	}
	if err := conn.WriteJSON(payload); err != nil {
		t.Fatalf("erro enviando RUN_FLOW: %v", err)
	}

	resp := readResponse(t, conn)
	if _, existe := resp["requestId"]; existe {
		t.Fatalf("resposta de preview não deveria trazer requestId: %+v", resp)
	}
}
