package api

import (
	"encoding/json"
	"log"
	"net/http"
	"image"
	"strings"

	"foto-facil-backend/internal/dag"
	"foto-facil-backend/internal/nodes"
	"foto-facil-backend/internal/storage"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  4096,
	WriteBufferSize: 4096,
	CheckOrigin: func(r *http.Request) bool { return true },
}
var Store *storage.SQLiteStore

type ReactFlowPayload struct {
	Action string `json:"action"`
	ID     string `json:"id,omitempty"`
	Name   string `json:"name,omitempty"`
	Flow   struct {
		Nodes []struct {
			ID   string `json:"id"`
			Data struct {
				OriginalType string `json:"originalType"`
				FilePaths    string `json:"filePaths"`
				OutputDir    string `json:"outputDir"`
				Brightness   int    `json:"brightness"`
				ColorMode    string `json:"colorMode"`
				ResizeW      int    `json:"resizeW"`
				ResizeH      int    `json:"resizeH"`
				CropX        int    `json:"cropX"`
				CropY        int    `json:"cropY"`
				CropW        int    `json:"cropW"`
				CropH        int    `json:"cropH"`
				RotateAngle  int    `json:"rotateAngle"`
				FlipAxis     int    `json:"flipAxis"`
				DoFlip       bool   `json:"doFlip"`
				AITool       string `json:"aiTool"`
				AITolerance  int    `json:"aiTolerance"`
			} `json:"data"`
		} `json:"nodes"`
		Edges []struct {
			Source string `json:"source"`
			Target string `json:"target"`
		} `json:"edges"`
	} `json:"flow"`
}

type RunFlowResponse struct {
	Status     string            `json:"status"`
	Message    string            `json:"message,omitempty"`
	Error      string            `json:"error,omitempty"`
	Thumbnails map[string]string `json:"thumbnails,omitempty"`
}

func HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade error:", err)
		return
	}
	defer conn.Close()
	log.Println("Client connected via WebSocket")

	for {
		_, p, err := conn.ReadMessage()
		if err != nil {
			break
		}

		var payload ReactFlowPayload
		if err := json.Unmarshal(p, &payload); err != nil {
			log.Printf("JSON unmarshal error: %v", err)
			continue
		}

		switch payload.Action {
		case "SAVE_FLOW":
			flowJSON, err := json.Marshal(payload.Flow)
			if err != nil {
				sendJSON(conn, map[string]string{"status": "error", "error": err.Error()})
				continue
			}
			if Store != nil {
				if err := Store.SaveFlow(payload.ID, payload.Name, string(flowJSON)); err != nil {
					sendJSON(conn, map[string]string{"status": "error", "error": err.Error()})
				} else {
					sendJSON(conn, map[string]string{"status": "success", "action": "SAVE_FLOW", "id": payload.ID})
				}
			} else {
				sendJSON(conn, map[string]string{"status": "error", "error": "Database not initialized"})
			}
			continue

		case "LOAD_FLOWS":
			if Store != nil {
				flows, err := Store.GetAllFlows()
				if err != nil {
					sendJSON(conn, map[string]string{"status": "error", "error": err.Error()})
				} else {
					sendJSON(conn, map[string]interface{}{"status": "success", "action": "LOAD_FLOWS", "flows": flows})
				}
			} else {
				sendJSON(conn, map[string]string{"status": "error", "error": "Database not initialized"})
			}
			continue

		case "DELETE_FLOW":
			if Store != nil {
				if err := Store.DeleteFlow(payload.ID); err != nil {
					sendJSON(conn, map[string]string{"status": "error", "error": err.Error()})
				} else {
					sendJSON(conn, map[string]string{"status": "success", "action": "DELETE_FLOW", "id": payload.ID})
				}
			} else {
				sendJSON(conn, map[string]string{"status": "error", "error": "Database not initialized"})
			}
			continue

		case "RUN_FLOW":
			// Prosseguir com a execução do grafo
		default:
			sendJSON(conn, map[string]string{"status": "error", "error": "Unknown action: " + payload.Action})
			continue
		}

		log.Printf("RUN_FLOW received: %d nodes, %d edges", len(payload.Flow.Nodes), len(payload.Flow.Edges))
		for _, n := range payload.Flow.Nodes {
			log.Printf("Node: %s | Type: %s | Data: %+v", n.ID, n.Data.OriginalType, n.Data)
		}

		// Montar DAG
		dagNodes := make(map[string]*dag.Node)
		for _, n := range payload.Flow.Nodes {
			dagNodes[n.ID] = &dag.Node{ID: n.ID, Dependencies: []string{}}
		}
		for _, e := range payload.Flow.Edges {
			if dest, ok := dagNodes[e.Target]; ok {
				dest.Dependencies = append(dest.Dependencies, e.Source)
			}
		}

		scheduler := dag.NewScheduler()
		order, err := scheduler.Sort(dagNodes)
		if err != nil {
			sendJSON(conn, RunFlowResponse{Status: "error", Error: "Ciclo detectado no grafo"})
			continue
		}

		// Instanciar nós processadores
		instances := make(map[string]nodes.Node)
		thumbNodes := make(map[string]*nodes.ThumbnailNode)

		for _, n := range payload.Flow.Nodes {
			switch n.Data.OriginalType {
			case "Image Input":
				paths := []string{}
				for _, p := range strings.Split(n.Data.FilePaths, ",") {
					p = strings.TrimSpace(p)
					if p != "" {
						paths = append(paths, p)
					}
				}
				instances[n.ID] = &nodes.InputNode{ID: n.ID, FilePaths: paths}
			case "Image Output":
				instances[n.ID] = &nodes.OutputNode{ID: n.ID, OutputDir: n.Data.OutputDir}
			case "Brightness & Contrast":
				instances[n.ID] = &nodes.BrightnessNode{ID: n.ID, Brightness: n.Data.Brightness}
			case "Color Space":
				instances[n.ID] = &nodes.GrayscaleNode{ID: n.ID}
			case "Crop/Resize":
				instances[n.ID] = &nodes.CropResizeNode{
					ID:     n.ID,
					Width:  n.Data.ResizeW,
					Height: n.Data.ResizeH,
					CropX:  n.Data.CropX,
					CropY:  n.Data.CropY,
					CropW:  n.Data.CropW,
					CropH:  n.Data.CropH,
				}
			case "Rotate/Flip":
				angle := nodes.RotateAngle(n.Data.RotateAngle)
				instances[n.ID] = &nodes.RotateFlipNode{
					ID:       n.ID,
					Angle:    angle,
					Flip:     nodes.FlipAxis(n.Data.FlipAxis),
					DoRotate: n.Data.RotateAngle != 0,
					DoFlip:   n.Data.DoFlip,
				}
			case "Comparação Visual":
				instances[n.ID] = &nodes.CompareNode{ID: n.ID}
			case "IA/Machine Learning":
				tool := n.Data.AITool
				if tool == "" {
					tool = "background_removal" // default fallback
				}
				instances[n.ID] = &nodes.AINode{
					ID:        n.ID,
					Tool:      tool,
					Tolerance: n.Data.AITolerance,
				}
			}

			// Gerador de thumbnail para qualquer nó exceto Output
			if n.Data.OriginalType != "Image Output" {
				tn := &nodes.ThumbnailNode{ID: "thumb_" + n.ID, MaxWidth: 200, MaxHeight: 200}
				thumbNodes[n.ID] = tn
			}
		}

		// Executar em ordem topológica
		ctx := &nodes.ProcessContext{
			NodeOutputs: make(map[string][]image.Image),
		}
		thumbnails := make(map[string]string)

		success := true
		for _, id := range order {
			node := instances[id]
			if node == nil {
				continue
			}

			// Coletar imagens de entrada a partir das saídas dos nós pais diretos
			var inputImages []image.Image
			if dagNode, ok := dagNodes[id]; ok {
				for _, depID := range dagNode.Dependencies {
					if parentOut, exists := ctx.NodeOutputs[depID]; exists {
						inputImages = append(inputImages, parentOut...)
					}
				}
			}
			ctx.Images = inputImages

			if err := node.Process(ctx); err != nil {
				log.Printf("Error in node %s: %v", id, err)
				sendJSON(conn, RunFlowResponse{Status: "error", Error: "Erro no nó " + id + ": " + err.Error()})
				success = false
				break
			}

			// Armazenar a saída deste nó para consumo dos nós filhos subsequentes
			ctx.NodeOutputs[id] = ctx.Images

			// Gerar thumbnail do estado atual das imagens
			if tn, ok := thumbNodes[id]; ok && len(ctx.Images) > 0 {
				if err := tn.Process(ctx); err == nil && len(tn.Thumbnails) > 0 {
					thumbnails[id] = tn.Thumbnails[0]
				}
			}
		}

		if success {
			sendJSON(conn, RunFlowResponse{
				Status:     "success",
				Message:    "Imagens processadas com sucesso!",
				Thumbnails: thumbnails,
			})
		}
	}
}

func sendJSON(conn *websocket.Conn, v any) {
	data, _ := json.Marshal(v)
	conn.WriteMessage(websocket.TextMessage, data)
}
