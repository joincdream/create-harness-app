package cli

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"create-harness-app/internal/blueprint"
	"create-harness-app/internal/hub"
	"create-harness-app/internal/template"
)

func HandleHubSubCommand(args []string) {
	if len(args) < 2 {
		PrintSubCommandHelp("hub")
		os.Exit(1)
	}
	subCmd := args[1]
	switch subCmd {
	case "list", "ls":
		var agent, query string
		for i := 2; i < len(args); i++ {
			if args[i] == "--agent" && i+1 < len(args) {
				agent = args[i+1]
			}
			if (args[i] == "--query" || args[i] == "-q") && i+1 < len(args) {
				query = args[i+1]
			}
		}
		HandleHubList(query, agent)

	case "info":
		if len(args) < 3 {
			fmt.Fprintln(os.Stderr, "Usage: create-harness-app hub info <name> [version]")
			os.Exit(1)
		}
		version := "v1.0.0"
		if len(args) >= 4 {
			version = args[3]
		}
		HandleHubInfo(args[2], version)

	case "pull", "install":
		if len(args) < 3 {
			fmt.Fprintln(os.Stderr, "Usage: create-harness-app hub pull <name> [version]")
			os.Exit(1)
		}
		version := "v1.0.0"
		if len(args) >= 4 {
			version = args[3]
		}
		HandleHubPull(args[2], version)

	case "push", "publish":
		if len(args) < 3 {
			fmt.Fprintln(os.Stderr, "Usage: create-harness-app hub push <name> [version] [--description <desc>] [--agent <agent>]")
			os.Exit(1)
		}
		name := args[2]
		version := "v1.0.0"
		description := "Custom Harness Template published from CLI"
		targetAgent := "antigravity"

		for i := 3; i < len(args); i++ {
			if !strings.HasPrefix(args[i], "--") && i == 3 {
				version = args[i]
			}
			if args[i] == "--description" && i+1 < len(args) {
				description = args[i+1]
			}
			if args[i] == "--agent" && i+1 < len(args) {
				targetAgent = args[i+1]
			}
		}
		HandleHubPush(name, version, description, targetAgent)

	case "delete", "rm":
		if len(args) < 3 {
			fmt.Fprintln(os.Stderr, "Usage: create-harness-app hub delete <name> [version]")
			os.Exit(1)
		}
		version := "v1.0.0"
		if len(args) >= 4 {
			version = args[3]
		}
		HandleHubDelete(args[2], version)

	default:
		PrintSubCommandHelp("hub")
	}
}

func HandleHubList(query, agent string) {
	client := hub.NewClient("")
	list, err := client.ListTemplates(context.Background(), query, agent)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching templates from HarnessHub: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("🌐 HarnessHub Registry Templates (%s):\n\n", client.BaseURL)
	if len(list) == 0 {
		fmt.Println("  (No templates found)")
		return
	}

	for _, item := range list {
		ver := item.Version
		if !strings.HasPrefix(ver, "v") {
			ver = "v" + ver
		}
		fmt.Printf("  • %s (%s) [agent: %s]\n", item.Name, ver, item.TargetAgent)
		if item.Description != "" {
			lines := strings.Split(item.Description, "\n")
			fmt.Printf("    Description: %s\n", lines[0])
		}
		fmt.Println("")
	}
}

func HandleHubInfo(name, version string) {
	client := hub.NewClient("")
	detail, err := client.GetTemplateDetail(context.Background(), name, version)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching template info: %v\n", err)
		os.Exit(1)
	}

	ver := detail.Version
	if !strings.HasPrefix(ver, "v") {
		ver = "v" + ver
	}
	fmt.Printf("📦 HarnessHub Template Info: %s (%s)\n", detail.Name, ver)
	fmt.Printf("  • Target Agent: %s\n", detail.TargetAgent)
	fmt.Printf("  • Description: %s\n", detail.Description)
	fmt.Printf("  • Archive File: %s (%d bytes)\n", detail.FilePath, detail.FileSize)
	if detail.BlueprintJSON != "" {
		fmt.Println("\n📄 Blueprint Specification:")
		fmt.Println(detail.BlueprintJSON)
	}
}

func HandleHubPull(name, version string) {
	client := hub.NewClient("")
	fmt.Printf("📥 Pulling template '%s' (%s) from HarnessHub (%s)...\n", name, version, client.BaseURL)

	stream, err := client.DownloadTemplate(context.Background(), name, version)
	if err != nil {
		detail, detailErr := client.GetTemplateDetail(context.Background(), name, version)
		if detailErr == nil && detail.BlueprintJSON != "" {
			destDir, saveErr := template.SaveRawBlueprint(name, detail.BlueprintJSON)
			if saveErr == nil {
				fmt.Printf("✅ Successfully pulled and cached template '%s' to %s\n", name, destDir)
				return
			}
		}
		fmt.Fprintf(os.Stderr, "Pull failed: %v\n", err)
		os.Exit(1)
	}
	defer stream.Close()

	destDir, err := template.InstallFromTarGz(stream, name)
	if err != nil {
		detail, detailErr := client.GetTemplateDetail(context.Background(), name, version)
		if detailErr == nil && detail.BlueprintJSON != "" {
			destDir, saveErr := template.SaveRawBlueprint(name, detail.BlueprintJSON)
			if saveErr == nil {
				fmt.Printf("✅ Successfully pulled and cached template '%s' to %s\n", name, destDir)
				return
			}
		}
		fmt.Fprintf(os.Stderr, "Extraction failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✅ Successfully pulled and extracted template '%s' to %s\n", name, destDir)
}

func HandleHubPush(name, version, description, targetAgent string) {
	client := hub.NewClient("")
	fmt.Printf("🚀 Packaging workspace and pushing template '%s' (%s) to HarnessHub (%s)...\n", name, version, client.BaseURL)

	cwd, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to get working directory: %v\n", err)
		os.Exit(1)
	}

	archiveBytes, err := packWorkspaceTarGz(cwd)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to pack workspace files: %v\n", err)
		os.Exit(1)
	}

	blueprintJSON := "{}"
	statePath := filepath.Join(cwd, ".harness", "state.json")
	if bytes, readErr := os.ReadFile(statePath); readErr == nil {
		var stateMap map[string]interface{}
		if jsonErr := json.Unmarshal(bytes, &stateMap); jsonErr == nil {
			if verVal, ok := stateMap["version"].(string); ok && verVal != "" {
				version = verVal
			}
			substituteTemplateVariables(stateMap, cwd)

			if descVal, ok := stateMap["description"].(string); ok {
				description = descVal
			}

			// Transpile stateMap into standard Registry BlueprintJSON schema
			blueprintJSON = buildStandardBlueprintJSON(name, version, targetAgent, description, stateMap)
		} else {
			blueprintJSON = string(bytes)
		}
	}

	// Warning if description remains empty
	if strings.TrimSpace(description) == "" {
		fmt.Fprintln(os.Stderr, "⚠️ Warning: Template description is missing or empty.")
		fmt.Fprintln(os.Stderr, "   Please set \"description\": \".harness/description.md\" in .harness/state.json.")
	}

	detail, err := client.PublishTemplate(context.Background(), name, version, targetAgent, description, archiveBytes, blueprintJSON)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Push failed: %v\n", err)
		os.Exit(1)
	}

	ver := detail.Version
	if !strings.HasPrefix(ver, "v") {
		ver = "v" + ver
	}
	fmt.Printf("✅ Successfully published template '%s' (%s) to HarnessHub!\n", detail.Name, ver)
}

func HandleHubDelete(name, version string) {
	client := hub.NewClient("")
	fmt.Printf("🗑️ Deleting template '%s' (%s) from HarnessHub (%s)...\n", name, version, client.BaseURL)

	err := client.DeleteTemplate(context.Background(), name, version)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Delete failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✅ Successfully deleted template '%s' (%s) from HarnessHub!\n", name, version)
}

func packWorkspaceTarGz(rootDir string) ([]byte, error) {
	var buf bytes.Buffer
	gw := gzip.NewWriter(&buf)
	tw := tar.NewWriter(gw)

	targets := []string{"AGENTS.md", "docs", ".harness", ".agents"}

	for _, target := range targets {
		targetPath := filepath.Join(rootDir, target)
		info, err := os.Stat(targetPath)
		if err != nil {
			continue // skip if optional file/directory doesn't exist
		}

		if !info.IsDir() {
			header, err := tar.FileInfoHeader(info, info.Name())
			if err != nil {
				continue
			}
			header.Name = target
			if err := tw.WriteHeader(header); err != nil {
				continue
			}
			content, err := os.ReadFile(targetPath)
			if err == nil {
				_, _ = tw.Write(content)
			}
		} else {
			_ = filepath.Walk(targetPath, func(path string, walkInfo os.FileInfo, walkErr error) error {
				if walkErr != nil || walkInfo.IsDir() {
					return nil
				}
				relPath, relErr := filepath.Rel(rootDir, path)
				if relErr != nil {
					return nil
				}
				header, headerErr := tar.FileInfoHeader(walkInfo, walkInfo.Name())
				if headerErr != nil {
					return nil
				}
				header.Name = relPath
				if err := tw.WriteHeader(header); err != nil {
					return nil
				}
				content, err := os.ReadFile(path)
				if err == nil {
					_, _ = tw.Write(content)
				}
				return nil
			})
		}
	}

	_ = tw.Close()
	_ = gw.Close()

	return buf.Bytes(), nil
}

func substituteTemplateVariables(node interface{}, rootDir string) {
	switch v := node.(type) {
	case map[string]interface{}:
		for k, val := range v {
			if strVal, ok := val.(string); ok {
				if strings.Contains(strVal, "${") {
					v[k] = expandVariable(strVal, rootDir, k)
				}
			} else {
				substituteTemplateVariables(val, rootDir)
			}
		}
	case []interface{}:
		for i, item := range v {
			if strVal, ok := item.(string); ok {
				if strings.Contains(strVal, "${") {
					v[i] = expandVariable(strVal, rootDir, "")
				}
			} else {
				substituteTemplateVariables(item, rootDir)
			}
		}
	}
}

func buildStandardBlueprintJSON(name, version, targetAgent, description string, stateMap map[string]interface{}) string {
	bp := blueprint.Blueprint{
		Name:        name,
		Version:     version,
		TargetAgent: targetAgent,
		Description: description,
		Workflows:   make([]blueprint.Workflow, 0),
	}

	nodesRaw, ok := stateMap["nodes"].(map[string]interface{})
	if !ok {
		bytes, _ := json.MarshalIndent(bp, "", "  ")
		return string(bytes)
	}

	// Group nodes by phase
	phaseMap := make(map[int][]blueprint.Node)
	for nodeID, nodeData := range nodesRaw {
		nMap, nOk := nodeData.(map[string]interface{})
		if !nOk {
			continue
		}

		phase := 1
		if pVal, pOk := nMap["phase"].(float64); pOk {
			phase = int(pVal)
		} else if pVal, pOk := nMap["phase"].(int); pOk {
			phase = pVal
		}

		var inputs []string
		if inRaw, inOk := nMap["inputs"].([]interface{}); inOk {
			for _, inItem := range inRaw {
				if s, sOk := inItem.(string); sOk {
					inputs = append(inputs, s)
				}
			}
		}

		var outputs []string
		if outRaw, outOk := nMap["outputs"].([]interface{}); outOk {
			for _, outItem := range outRaw {
				if s, sOk := outItem.(string); sOk {
					outputs = append(outputs, s)
				}
			}
		}

		nodeDesc, _ := nMap["description"].(string)

		bNode := blueprint.Node{
			ID:          nodeID,
			FileName:    fmt.Sprintf("%s.md", nodeID),
			Description: nodeDesc,
			Inputs:      inputs,
			Outputs:     outputs,
		}

		phaseMap[phase] = append(phaseMap[phase], bNode)
	}

	for p := 1; p <= 10; p++ {
		jobs, exists := phaseMap[p]
		if !exists || len(jobs) == 0 {
			continue
		}

		dirName := fmt.Sprintf("stage_%d", p)
		if p == 1 {
			dirName = "docs/01_planning"
		} else if p == 2 {
			dirName = "docs/02_design"
		} else if p == 3 {
			dirName = "docs/03_development"
		}

		wf := blueprint.Workflow{
			Phase:     p,
			Dir:       dirName,
			Directory: dirName,
			Jobs:      jobs,
		}
		bp.Workflows = append(bp.Workflows, wf)
	}

	bytes, err := json.MarshalIndent(bp, "", "  ")
	if err != nil {
		return "{}"
	}
	return string(bytes)
}

func expandVariable(input, rootDir, keyName string) string {
	return os.Expand(input, func(key string) string {
		if strings.HasPrefix(key, "file:") {
			filePath := strings.TrimPrefix(key, "file:")
			absPath := filepath.Join(rootDir, filePath)
			content, err := os.ReadFile(absPath)
			if err == nil {
				return string(content)
			}
			fmt.Fprintf(os.Stderr, "⚠️ Warning: Variable file '%s' not found for key '%s'\n", filePath, keyName)
		}
		return "${" + key + "}"
	})
}
