package harness

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"create-harness-app/internal/blueprint"
)

type NodeState struct {
	ID          string   `json:"id"`
	File        string   `json:"file"`
	Phase       int      `json:"phase"`
	Status      string   `json:"status"` // "pending", "in_progress", "completed", "blocked"
	DependsOn   []string `json:"depends_on,omitempty"`
	Inputs      []string `json:"inputs,omitempty"`
	Outputs     []string `json:"outputs,omitempty"`
	Description string   `json:"description"`
}

type HarnessState struct {
	AppName      string               `json:"app_name"`
	TargetAgent  string               `json:"target_agent"`
	CurrentPhase int                  `json:"current_phase"`
	Nodes        map[string]NodeState `json:"nodes"`
}

func GetStatePath(targetDir string) string {
	return filepath.Join(targetDir, "harness", "state.json")
}

func EnsureStateDir(targetDir string) error {
	stateDir := filepath.Join(targetDir, "harness")
	if err := os.MkdirAll(stateDir, 0755); err != nil {
		return err
	}

	return nil
}

func InitState(bp *blueprint.Blueprint, targetDir string) error {
	if err := EnsureStateDir(targetDir); err != nil {
		return err
	}

	state := &HarnessState{
		AppName:      filepath.Base(targetDir),
		TargetAgent:  bp.TargetAgent,
		CurrentPhase: 1,
		Nodes:        make(map[string]NodeState),
	}

	for _, wf := range bp.Workflows {
		for _, node := range wf.Nodes {
			id := node.ID
			if id == "" {
				id = node.File
			}
			state.Nodes[id] = NodeState{
				ID:          id,
				File:        filepath.Join(wf.Dir, node.File),
				Phase:       wf.Phase,
				Status:      "pending",
				DependsOn:   node.DependsOn,
				Inputs:      node.Inputs,
				Outputs:     node.Outputs,
				Description: node.Description,
			}
		}
	}

	return SaveState(targetDir, state)
}

func LoadState(targetDir string) (*HarnessState, error) {
	path := GetStatePath(targetDir)
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read harness state file: %w", err)
	}

	var state HarnessState
	if err := json.Unmarshal(data, &state); err != nil {
		return nil, fmt.Errorf("failed to parse harness state file: %w", err)
	}

	return &state, nil
}

func SaveState(targetDir string, state *HarnessState) error {
	path := GetStatePath(targetDir)
	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

func UpdateNodeStatus(targetDir string, nodeID string, newStatus string) error {
	state, err := LoadState(targetDir)
	if err != nil {
		return err
	}

	node, exists := state.Nodes[nodeID]
	if !exists {
		return fmt.Errorf("node ID '%s' not found in harness state", nodeID)
	}

	if newStatus == "completed" {
		for _, out := range node.Outputs {
			outAbsPath := filepath.Join(targetDir, out)
			if _, err := os.Stat(outAbsPath); os.IsNotExist(err) {
				return fmt.Errorf("validation failed: output file '%s' does not exist on disk", out)
			}
		}
	}

	node.Status = newStatus
	state.Nodes[nodeID] = node

	return SaveState(targetDir, state)
}
