package config

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
)

type Config struct {
	AppName       string
	TemplateName  string
	CustomPath    string
	Force         bool
	ShowVersion   bool
	ShowHelp      bool
	ListTemplates bool
}

func ParseConfig(args []string) (*Config, error) {
	flags := flag.NewFlagSet("create-harness-app", flag.ContinueOnError)

	var cfg Config
	flags.BoolVar(&cfg.Force, "force", false, "Overwrite existing files")
	flags.BoolVar(&cfg.Force, "f", false, "Overwrite existing files (short)")
	flags.StringVar(&cfg.TemplateName, "template", "default", "Blueprint template name")
	flags.StringVar(&cfg.CustomPath, "path", "", "Custom template dir path")
	flags.BoolVar(&cfg.ShowVersion, "version", false, "Show version")
	flags.BoolVar(&cfg.ShowVersion, "v", false, "Show version (short)")
	flags.BoolVar(&cfg.ShowHelp, "help", false, "Show help")
	flags.BoolVar(&cfg.ShowHelp, "h", false, "Show help (short)")
	flags.BoolVar(&cfg.ListTemplates, "list", false, "List available templates")
	flags.BoolVar(&cfg.ListTemplates, "l", false, "List available templates (short)")

	// Reorder args so that options/flags can appear after app name argument
	var flagArgs []string
	var nonFlags []string
	for i := 0; i < len(args); i++ {
		arg := args[i]
		if strings.HasPrefix(arg, "-") {
			flagArgs = append(flagArgs, arg)
			// Check if flag expects value (like --template <name> or -path <dir>)
			if (arg == "--template" || arg == "-template" || arg == "--path" || arg == "-path") && i+1 < len(args) {
				i++
				flagArgs = append(flagArgs, args[i])
			}
		} else {
			nonFlags = append(nonFlags, arg)
		}
	}

	reorderedArgs := append(flagArgs, nonFlags...)

	if err := flags.Parse(reorderedArgs); err != nil {
		return nil, err
	}

	parsedNonFlags := flags.Args()
	if len(parsedNonFlags) > 0 {
		cfg.AppName = parsedNonFlags[0]
	}

	// Interactive Wizard Mode if AppName is not provided and non-flag execution
	if cfg.AppName == "" && !cfg.ShowVersion && !cfg.ShowHelp && !cfg.ListTemplates {
		runWizard(&cfg)
	}

	if cfg.AppName == "" {
		cfg.AppName = "my-harness-app"
	}

	return &cfg, nil
}

func runWizard(cfg *Config) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("? 프로젝트 디렉토리명을 입력하세요 (기본값: my-harness-app): ")
	appName, _ := reader.ReadString('\n')
	appName = strings.TrimSpace(appName)
	if appName != "" {
		cfg.AppName = appName
	} else {
		cfg.AppName = "my-harness-app"
	}

	fmt.Println("? 사용할 하네스 블루프린트를 선택하세요:")
	fmt.Println("  1) default (기존 SDLC 4단계 모듈)")
	fmt.Println("  2) web-api (RESTful Web API: OpenAPI + DTO + TDD)")
	fmt.Println("  3) cli-app (Go CLI App)")
	fmt.Print("선택 번호 입력 [1]: ")

	choice, _ := reader.ReadString('\n')
	choice = strings.TrimSpace(choice)
	switch choice {
	case "2":
		cfg.TemplateName = "web-api"
	case "3":
		cfg.TemplateName = "cli-app"
	default:
		cfg.TemplateName = "default"
	}
}
