package cli

import "fmt"

func HasHelpFlag(args []string) bool {
	for _, arg := range args {
		if arg == "-h" || arg == "--help" || arg == "help" {
			return true
		}
	}
	return false
}

func PrintSubCommandHelp(cmd string) {
	switch cmd {
	case "status":
		fmt.Println("Usage: create-harness-app status")
		fmt.Println("Description: Displays full SDLC harness execution progress and node statuses in JSON format.")

	case "next":
		fmt.Println("Usage: create-harness-app next")
		fmt.Println("Description: Analyzes harness state and infers the next pending action node, inputs, and outputs.")

	case "state":
		fmt.Println("Usage: create-harness-app state update --node <node_id> --status <pending|completed|blocked>")
		fmt.Println("Description: Deterministically updates node status after output file validation.")

	case "extract":
		fmt.Println("Usage: create-harness-app extract <sample_app_dir> [--template <name>]")
		fmt.Println("Description: Extracts optimized .agents/ and docs/ assets from a sample app back into harness templates.")

	case "hub":
		fmt.Println("Usage: create-harness-app hub <list|info|pull|push> [args]")
		fmt.Println("       create-harness-app hub list [--agent <agent>] [--query <q>]")
		fmt.Println("       create-harness-app hub info <name> [version]")
		fmt.Println("       create-harness-app hub pull <name> [version]")
		fmt.Println("       create-harness-app hub push <name> [version] [--description <desc>] [--agent <agent>]")
		fmt.Println("Description: Interacts with HarnessHub backend registry for searching, inspecting, pulling, and pushing templates.")

	case "templates", "template", "list":
		fmt.Println("Usage: create-harness-app templates [list]")
		fmt.Println("       create-harness-app templates remove <name>")
		fmt.Println("       create-harness-app -list")
		fmt.Println("Description: Lists available templates or removes downloaded custom template from local cache.")
	}
}
