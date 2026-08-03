package cli

import (
	"context"
	"fmt"
	"os"
	"strings"

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
