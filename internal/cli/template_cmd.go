package cli

import (
	"fmt"
	"os"
	"strings"

	"create-harness-app/internal/template"
)

func HandleTemplatesSubCommand(args []string) {
	if len(args) >= 3 && (args[1] == "remove" || args[1] == "rm" || args[2] == "remove" || args[2] == "rm") {
		HandleTemplateRemove(args[len(args)-1])
	} else if len(args) >= 2 && (args[1] == "remove" || args[1] == "rm") {
		if len(args) < 3 {
			fmt.Fprintln(os.Stderr, "Usage: create-harness-app templates remove <name>")
			os.Exit(1)
		}
		HandleTemplateRemove(args[2])
	} else {
		PrintTemplatesList()
	}
}

func HandleTemplateRemove(name string) {
	if name == "default" || name == "antigravity" {
		fmt.Fprintf(os.Stderr, "⚠️ Cannot remove embedded template '%s'. Embedded templates are built into the binary.\n", name)
		os.Exit(1)
	}

	if err := template.RemoveCachedTemplate(name); err != nil {
		fmt.Fprintf(os.Stderr, "Remove failed: %v\n", err)
		os.Exit(1)
	}

	baseDir, _ := template.GetUserTemplatesDir()
	fmt.Printf("🗑️ Successfully removed cached template '%s' from %s/%s\n", name, baseDir, name)
}

func PrintTemplatesList() {
	list, err := template.ListTemplates()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error listing templates: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("📦 Available Harness Templates:")
	fmt.Println("")
	for _, t := range list {
		versionStr := ""
		if t.Version != "" {
			versionStr = fmt.Sprintf(" (v%s)", t.Version)
		}
		authorStr := ""
		if t.Author != "" {
			authorStr = fmt.Sprintf(" by %s", t.Author)
		}
		licenseStr := ""
		if t.License != "" {
			licenseStr = fmt.Sprintf(" [%s]", t.License)
		}
		fmt.Printf("  • %s%s%s%s (%s)\n", t.Name, versionStr, authorStr, licenseStr, t.Source)
		if t.Description != "" {
			lines := strings.Split(t.Description, "\n")
			fmt.Printf("    Description: %s\n", lines[0])
			for _, l := range lines[1:] {
				if strings.TrimSpace(l) != "" {
					fmt.Printf("                 %s\n", l)
				}
			}
		}
		if len(t.Keywords) > 0 {
			fmt.Printf("    Keywords: %v\n", t.Keywords)
		}
		fmt.Println("")
	}
}
