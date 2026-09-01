package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"create-harness-app/internal/harness"
)

type MCPConfig struct {
	ServerName string    `json:"server_name"`
	Version    string    `json:"version"`
	Tools      []MCPTool `json:"tools"`
}

type MCPTool struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	InputSchema interface{} `json:"inputSchema"`
}

type JSONRPCRequest struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      interface{}     `json:"id"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
}

type JSONRPCResponse struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      interface{} `json:"id"`
	Result  interface{} `json:"result,omitempty"`
}

type CallToolParams struct {
	Name      string                 `json:"name"`
	Arguments map[string]interface{} `json:"arguments"`
}

func logDebug(format string, v ...interface{}) {
	logFile := "/tmp/harness_mcp_debug.log"
	f, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return
	}
	defer f.Close()

	timestamp := time.Now().Format(time.RFC3339)
	msg := fmt.Sprintf("[%s] ", timestamp) + fmt.Sprintf(format, v...) + "\n"
	_, _ = f.WriteString(msg)
}

func HandleMCPSubCommand(args []string) {
	cwd, _ := os.Getwd()
	if len(args) >= 2 && (args[1] == "inspect" || args[1] == "test" || args[1] == "info") {
		InspectMCPConfig(cwd)
		return
	}

	RunMCPServer(cwd)
}

func InspectMCPConfig(cwd string) {
	mcpPath := filepath.Join(cwd, ".harness", "mcp.json")
	data, err := os.ReadFile(mcpPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "⚠️ Error reading %s: %v\n", mcpPath, err)
		os.Exit(1)
	}

	var cfg MCPConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		fmt.Fprintf(os.Stderr, "⚠️ Invalid JSON in mcp.json: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("🧠 ========================================================")
	fmt.Println("🧠 MCP Semantic Layer Inspection (Gemini 3.5 Flash Injection)")
	fmt.Println("🧠 ========================================================")
	fmt.Printf(" Server Name : %s (v%s)\n\n", cfg.ServerName, cfg.Version)

	fmt.Printf("📋 Injected Semantic Tools (%d count):\n\n", len(cfg.Tools))
	for i, t := range cfg.Tools {
		fmt.Printf(" [%d] Tool Name : %s\n", i+1, t.Name)
		fmt.Printf("     Description: \"%s\"\n", t.Description)
		schemaJSON, _ := json.MarshalIndent(t.InputSchema, "     ", "  ")
		fmt.Printf("     Input Schema:\n     %s\n\n", string(schemaJSON))
	}
	fmt.Println("==========================================================")
}

func RunMCPServer(cwd string) {
	logDebug("🚀 create-harness-app mcp started. CWD: %s, Args: %v", cwd, os.Args)

	mcpPath := filepath.Join(cwd, "harness", "mcp.json")
	data, err := os.ReadFile(mcpPath)
	if err != nil {
		logDebug("⚠️ Failed to read harness/mcp.json: %v", err)
	} else {
		logDebug("📄 Successfully loaded harness/mcp.json (%d bytes)", len(data))
	}

	var cfg MCPConfig
	_ = json.Unmarshal(data, &cfg)

	decoder := json.NewDecoder(os.Stdin)
	encoder := json.NewEncoder(os.Stdout)

	for {
		var req JSONRPCRequest
		if err := decoder.Decode(&req); err != nil {
			logDebug("🔌 STDIN decoder finished/closed: %v", err)
			break
		}

		logDebug("📩 [REQ Received] Method: '%s', ID: %v, Params: %s", req.Method, req.ID, string(req.Params))

		if req.ID == nil {
			logDebug("ℹ️ Notification message received (no ID). Skipping response.")
			continue
		}

		switch req.Method {
		case "server/discover", "initialize":
			resp := JSONRPCResponse{
				JSONRPC: "2.0",
				ID:      req.ID,
				Result: map[string]interface{}{
					"protocolVersion": "2026-07-28",
					"capabilities": map[string]interface{}{
						"tools":     map[string]interface{}{},
						"resources": map[string]interface{}{},
						"prompts":   map[string]interface{}{},
					},
					"serverInfo": map[string]interface{}{
						"name":    "harness-mcp-server",
						"version": "1.0.0",
					},
				},
			}
			_ = encoder.Encode(resp)
			os.Stdout.Sync()
			logDebug("📤 [RESP Sent] %s response sent successfully", req.Method)

		case "ping":
			resp := JSONRPCResponse{
				JSONRPC: "2.0",
				ID:      req.ID,
				Result:  map[string]interface{}{},
			}
			_ = encoder.Encode(resp)
			os.Stdout.Sync()
			logDebug("📤 [RESP Sent] ping response sent")

		case "resources/list":
			resp := JSONRPCResponse{
				JSONRPC: "2.0",
				ID:      req.ID,
				Result: map[string]interface{}{
					"resources": []interface{}{},
				},
			}
			_ = encoder.Encode(resp)
			os.Stdout.Sync()
			logDebug("📤 [RESP Sent] resources/list response sent")

		case "prompts/list":
			resp := JSONRPCResponse{
				JSONRPC: "2.0",
				ID:      req.ID,
				Result: map[string]interface{}{
					"prompts": []interface{}{},
				},
			}
			_ = encoder.Encode(resp)
			os.Stdout.Sync()
			logDebug("📤 [RESP Sent] prompts/list response sent")

		case "tools/list":
			resp := JSONRPCResponse{
				JSONRPC: "2.0",
				ID:      req.ID,
				Result: map[string]interface{}{
					"tools": cfg.Tools,
				},
			}
			_ = encoder.Encode(resp)
			os.Stdout.Sync()
			logDebug("📤 [RESP Sent] tools/list response sent (%d tools)", len(cfg.Tools))

		case "tools/call":
			var callParams CallToolParams
			if err := json.Unmarshal(req.Params, &callParams); err == nil {
				logDebug("🔧 [Tool Execution Call] Tool: '%s', Args: %v", callParams.Name, callParams.Arguments)
				content := handleToolCallExecution(cwd, callParams)
				resp := JSONRPCResponse{
					JSONRPC: "2.0",
					ID:      req.ID,
					Result: map[string]interface{}{
						"content": content,
					},
				}
				_ = encoder.Encode(resp)
				os.Stdout.Sync()
				logDebug("📤 [RESP Sent] tools/call execution response sent")
			}
		}
	}
}

func handleToolCallExecution(cwd string, params CallToolParams) []map[string]interface{} {
	st, err := harness.LoadState(cwd)
	if err != nil {
		return []map[string]interface{}{
			{"type": "text", "text": fmt.Sprintf("Error loading state: %v", err)},
		}
	}

	switch params.Name {
	case "harness_next":
		var pendingNodeID string
		var pendingNode harness.NodeState
		for id, node := range st.Nodes {
			if node.Status != "completed" {
				pendingNodeID = id
				pendingNode = node
				break
			}
		}
		if pendingNodeID != "" {
			msg := fmt.Sprintf("🎯 [Harness Next Target] App: %s (Phase %d)\n👉 Pending Node: [%s]\n   Skill File: %s\n   Inputs: %v\n   Outputs: %v",
				st.AppName, st.CurrentPhase, pendingNodeID, pendingNode.File, pendingNode.Inputs, pendingNode.Outputs)
			return []map[string]interface{}{{"type": "text", "text": msg}}
		}
		return []map[string]interface{}{{"type": "text", "text": "🎉 All SDLC harness nodes completed!"}}

	case "harness_status":
		data, _ := json.MarshalIndent(st, "", "  ")
		return []map[string]interface{}{{"type": "text", "text": string(data)}}

	case "harness_update_state":
		nodeID, _ := params.Arguments["node_id"].(string)
		status, _ := params.Arguments["status"].(string)
		if nodeID != "" && status != "" {
			if err := harness.UpdateNodeStatus(cwd, nodeID, status); err == nil {
				return []map[string]interface{}{{"type": "text", "text": fmt.Sprintf("✅ State updated: node '%s' -> '%s'", nodeID, status)}}
			}
		}
		return []map[string]interface{}{{"type": "text", "text": "❌ Failed to update node status"}}
	}

	return []map[string]interface{}{{"type": "text", "text": "Unknown tool"}}
}
