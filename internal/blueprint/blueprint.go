package blueprint

import (
	"encoding/json"
	"fmt"
)

// Node represents a single file or task within a workflow
type Node struct {
	ID          string   `json:"id"`
	File        string   `json:"file"`
	Type        string   `json:"type"` // "input", "spec", "guardrail"
	Description string   `json:"description"`
	DependsOn   []string `json:"depends_on,omitempty"`
	Inputs      []string `json:"inputs,omitempty"`
	Outputs     []string `json:"outputs,omitempty"`
}

// Workflow represents an SDLC Phase directory and its nodes
type Workflow struct {
	Phase       int    `json:"phase"`
	Dir         string `json:"dir"`
	Description string `json:"description"`
	SkillPath   string `json:"skill_path,omitempty"`
	Nodes       []Node `json:"nodes"`
}

// Blueprint represents a full declarative harness template
type Blueprint struct {
	Name        string     `json:"name"`
	Version     string     `json:"version"`
	Description string     `json:"description"`
	TargetAgent string     `json:"target_agent,omitempty"`
	Workflows   []Workflow `json:"workflows"`
}

// ParseBlueprint parses JSON data into a Blueprint struct
func ParseBlueprint(data []byte) (*Blueprint, error) {
	var bp Blueprint
	if err := json.Unmarshal(data, &bp); err != nil {
		return nil, fmt.Errorf("failed to parse blueprint.json: %w", err)
	}
	return &bp, nil
}
