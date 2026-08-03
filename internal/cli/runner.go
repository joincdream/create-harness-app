package cli

import (
	"encoding/json"
	"fmt"
	"os"

	"create-harness-app/internal/config"
	"create-harness-app/internal/engine"
	"create-harness-app/internal/harness"
	"create-harness-app/internal/template"
)

const Version = "0.2.0-go"

func Run(args []string) {
	if len(args) > 0 {
		switch args[0] {
		case "create":
			RunCreateCommand(args[1:])
			return
		case "state", "status", "next", "extract", "templates", "template", "list", "hub", "pull", "mcp":
			HandleHarnessCommand(args)
			return
		}
	}

	RunCreateCommand(args)
}

func RunCreateCommand(args []string) {
	cfg, err := config.ParseConfig(args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing arguments: %v\n", err)
		os.Exit(1)
	}

	if cfg.ShowVersion {
		fmt.Printf("create-harness-app version %s (Go 1.22 Native)\n", Version)
		os.Exit(0)
	}

	if cfg.ListTemplates {
		PrintTemplatesList()
		os.Exit(0)
	}

	if cfg.ShowHelp {
		fmt.Println("Usage: create-harness-app [create] [target_directory] [flags]")
		fmt.Println("       create-harness-app templates [list]")
		fmt.Println("       create-harness-app templates remove <name>")
		fmt.Println("       create-harness-app hub list [--agent <agent>] [--query <q>]")
		fmt.Println("       create-harness-app hub info <name> [version]")
		fmt.Println("       create-harness-app hub pull <name> [version]")
		fmt.Println("       create-harness-app state update --node <id> --status <status>")
		fmt.Println("       create-harness-app status [--json]")
		fmt.Println("       create-harness-app next")
		fmt.Println("       create-harness-app extract <sample_app_dir> [--template <name>]")
		fmt.Println("\nFlags:")
		fmt.Println("  -template, --template <name>  Specify blueprint template (default: 'default')")
		fmt.Println("  -path, --path <dir>           Specify custom template directory path")
		fmt.Println("  -list, -l                     List available blueprint templates")
		fmt.Println("  -force, -f                    Overwrite existing files")
		fmt.Println("  -version, -v                  Show version")
		fmt.Println("  -help, -h                     Show help")
		os.Exit(0)
	}

	tmpl, err := template.ResolveTemplate(cfg.TemplateName, cfg.CustomPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error resolving template: %v\n", err)
		os.Exit(2)
	}

	if _, err := engine.Scaffold(cfg, tmpl); err != nil {
		fmt.Fprintf(os.Stderr, "Scaffolding error: %v\n", err)
		os.Exit(2)
	}
}

func HandleHarnessCommand(args []string) {
	cmd := args[0]
	cwd, _ := os.Getwd()

	if HasHelpFlag(args) {
		PrintSubCommandHelp(cmd)
		return
	}

	switch cmd {
	case "status":
		handleStatusCommand(cwd)
	case "state":
		handleStateCommand(cwd, args)
	case "next":
		handleNextCommand(cwd)
	case "extract":
		handleExtractCommand(cwd, args)
	case "pull":
		handlePullCommand(args)
	case "mcp":
		HandleMCPSubCommand(args)
	case "hub":
		HandleHubSubCommand(args)
	case "templates", "template", "list":
		HandleTemplatesSubCommand(args)
	}
}

func handleStatusCommand(cwd string) {
	st, err := harness.LoadState(cwd)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	data, _ := json.MarshalIndent(st, "", "  ")
	fmt.Println(string(data))
}

func handleStateCommand(cwd string, args []string) {
	if len(args) < 5 || args[1] != "update" {
		PrintSubCommandHelp("state")
		return
	}
	var nodeID, status string
	for i := 2; i < len(args); i++ {
		if args[i] == "--node" && i+1 < len(args) {
			nodeID = args[i+1]
		}
		if args[i] == "--status" && i+1 < len(args) {
			status = args[i+1]
		}
	}
	if nodeID == "" || status == "" {
		fmt.Fprintln(os.Stderr, "Usage: create-harness-app state update --node <id> --status <status>")
		os.Exit(1)
	}
	if err := harness.UpdateNodeStatus(cwd, nodeID, status); err != nil {
		fmt.Fprintf(os.Stderr, "State update failed: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("✅ State updated successfully: node '%s' -> '%s'\n", nodeID, status)
}

func handleNextCommand(cwd string) {
	st, err := harness.LoadState(cwd)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("🎯 [Harness Next Target] App: %s (Phase %d)\n", st.AppName, st.CurrentPhase)
	for id, node := range st.Nodes {
		if node.Status != "completed" {
			fmt.Printf("  👉 Next Pending Node: [%s] (%s)\n", id, node.File)
			fmt.Printf("     Inputs: %v\n", node.Inputs)
			fmt.Printf("     Outputs: %v\n", node.Outputs)
			break
		}
	}
}

func handleExtractCommand(cwd string, args []string) {
	if len(args) < 2 {
		PrintSubCommandHelp("extract")
		os.Exit(1)
	}
	sampleDir := args[1]
	templateName := "antigravity"
	for i := 2; i < len(args); i++ {
		if (args[i] == "--template" || args[i] == "-template") && i+1 < len(args) {
			templateName = args[i+1]
		}
	}
	if err := harness.ExtractTemplate(sampleDir, templateName, cwd); err != nil {
		fmt.Fprintf(os.Stderr, "Extract failed: %v\n", err)
		os.Exit(1)
	}
}

func handlePullCommand(args []string) {
	if len(args) < 2 {
		fmt.Fprintln(os.Stderr, "Usage: create-harness-app pull <name> [version]")
		os.Exit(1)
	}
	name := args[1]
	version := "v1.0.0"
	if len(args) >= 3 {
		version = args[2]
	}
	HandleHubPull(name, version)
}
