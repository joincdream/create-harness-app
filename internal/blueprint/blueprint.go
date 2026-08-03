package blueprint

import (
	"encoding/json"
	"fmt"
)

// Node represents a single file or task within a workflow
type Node struct {
	ID          string   `json:"id"`
	File        string   `json:"file,omitempty"`
	FilePath    string   `json:"file_path,omitempty"`
	FileName    string   `json:"file_name,omitempty"`
	Type        string   `json:"type,omitempty"`
	Description string   `json:"description,omitempty"`
	DependsOn   []string `json:"depends_on,omitempty"`
	Inputs      []string `json:"inputs,omitempty"`
	Outputs     []string `json:"outputs,omitempty"`
}

// Workflow represents an SDLC Phase directory and its nodes
type Workflow struct {
	Phase       int    `json:"phase"`
	Dir         string `json:"dir,omitempty"`
	Directory   string `json:"directory,omitempty"`
	Description string `json:"description,omitempty"`
	SkillPath   string `json:"skill_path,omitempty"`
	Nodes       []Node `json:"nodes,omitempty"`
	Jobs        []Node `json:"jobs,omitempty"`
}

// Author represents creator information for a template
type Author struct {
	Name  string `json:"name,omitempty"`
	Email string `json:"email,omitempty"`
}

// Blueprint represents a full declarative harness template
type Blueprint struct {
	Name        string     `json:"name"`
	Version     string     `json:"version"`
	Description string     `json:"description"`
	Author      *Author    `json:"author,omitempty"`
	License     string     `json:"license,omitempty"`
	Keywords    []string   `json:"keywords,omitempty"`
	TargetAgent string     `json:"target_agent,omitempty"`
	Repository  string     `json:"repository,omitempty"`
	Readme      string     `json:"readme,omitempty"`
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
