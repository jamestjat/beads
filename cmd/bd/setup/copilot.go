package setup

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/steveyegge/beads/internal/templates/agents"
)

const (
	copilotInstructionsFile = ".github/copilot-instructions.md"
	copilotPromptsDir       = ".github/prompts"
)

//go:embed copilot_prompts/*.prompt.md
var copilotPromptFiles embed.FS

// beadsPromptNames lists the prompt files shipped with the integration.
var beadsPromptNames = []string{
	"beads-ready.prompt.md",
	"beads-create.prompt.md",
	"beads-workflow.prompt.md",
	"plan-to-beads.prompt.md",
}

var copilotIntegration = agentsIntegration{
	name:         "GitHub Copilot",
	setupCommand: "bd setup copilot",
	readHint:     "GitHub Copilot reads .github/copilot-instructions.md automatically in VS Code.",
	docsURL:      "https://github.com/steveyegge/beads/blob/main/docs/COPILOT_INTEGRATION.md",
	profile:      agents.ProfileFull,
}

// InstallCopilot installs GitHub Copilot integration.
// When installPrompts is true, reusable prompt files are also copied to .github/prompts/.
func InstallCopilot(installPrompts bool) {
	fmt.Println("Installing GitHub Copilot integration...")

	// 1. Ensure .github directory exists
	if err := EnsureDir(".github", 0o755); err != nil {
		FatalError("%v", err)
	}

	// 2. Inject beads section into .github/copilot-instructions.md
	env := agentsEnv{
		agentsPath: copilotInstructionsFile,
		stdout:     os.Stdout,
		stderr:     os.Stderr,
	}
	if err := installAgents(env, copilotIntegration); err != nil {
		return
	}

	// 3. Optionally install reusable prompt files
	if installPrompts {
		installCopilotPrompts()
	}

	fmt.Println()
	fmt.Println("Reload VS Code for changes to take effect.")
	if !installPrompts {
		fmt.Println()
		fmt.Println("Tip: use --prompts to also install reusable prompt files to .github/prompts/")
	}
	fmt.Println()
	fmt.Println("For Copilot CLI, pre-approve bd commands:")
	fmt.Println("  copilot --allow-tool='shell(bd:*)'")
}

// installCopilotPrompts copies embedded prompt files to .github/prompts/.
func installCopilotPrompts() {
	if err := EnsureDir(copilotPromptsDir, 0o755); err != nil {
		FatalError("%v", err)
	}

	for _, name := range beadsPromptNames {
		data, err := copilotPromptFiles.ReadFile("copilot_prompts/" + name)
		if err != nil {
			FatalError("read embedded prompt %s: %v", name, err)
		}
		dest := filepath.Join(copilotPromptsDir, name)
		if err := atomicWriteFile(dest, data); err != nil {
			FatalError("write prompt %s: %v", dest, err)
		}
		fmt.Printf("✓ Prompt: %s\n", dest)
	}
}

// CheckCopilot reports whether the GitHub Copilot integration is installed.
func CheckCopilot() {
	instructionsExists := FileExists(copilotInstructionsFile)

	if instructionsExists {
		fmt.Println("✓ GitHub Copilot integration installed")
		fmt.Printf("  Instructions: %s\n", copilotInstructionsFile)

		// Report prompt files if present
		promptCount := 0
		for _, name := range beadsPromptNames {
			if FileExists(filepath.Join(copilotPromptsDir, name)) {
				promptCount++
			}
		}
		if promptCount > 0 {
			fmt.Printf("  Prompts:      %s (%d files)\n", copilotPromptsDir, promptCount)
		}
		return
	}

	FatalErrorWithHint("GitHub Copilot integration not installed", "Run: bd setup copilot")
}

// RemoveCopilot removes GitHub Copilot integration files.
// When removePrompts is true, beads prompt files are also removed from .github/prompts/.
func RemoveCopilot(removePrompts bool) {
	fmt.Println("Removing GitHub Copilot integration...")

	removed := false

	// Remove beads section (or the whole file if it only contained beads content)
	if FileExists(copilotInstructionsFile) {
		data, err := os.ReadFile(copilotInstructionsFile) //nolint:gosec // project-relative constant
		if err == nil {
			content := string(data)
			if containsBeadsMarker(content) {
				newContent := removeBeadsSection(content)
				if strings.TrimSpace(newContent) == "" {
					if err := os.Remove(copilotInstructionsFile); err != nil {
						FatalError("remove %s: %v", copilotInstructionsFile, err)
					}
					fmt.Printf("✓ Removed %s (beads section was the only content)\n", copilotInstructionsFile)
				} else {
					if err := atomicWriteFile(copilotInstructionsFile, []byte(newContent)); err != nil {
						FatalError("update %s: %v", copilotInstructionsFile, err)
					}
					fmt.Printf("✓ Removed beads section from %s\n", copilotInstructionsFile)
				}
				removed = true
			} else {
				fmt.Printf("  %s exists but has no beads section – skipping\n", copilotInstructionsFile)
			}
		}
	}

	// Remove beads prompt files from .github/prompts/
	if removePrompts {
		for _, name := range beadsPromptNames {
			p := filepath.Join(copilotPromptsDir, name)
			if FileExists(p) {
				if err := os.Remove(p); err != nil {
					FatalError("remove %s: %v", p, err)
				}
				fmt.Printf("✓ Removed %s\n", p)
				removed = true
			}
		}
		// Remove prompts dir if now empty
		_ = os.Remove(copilotPromptsDir)
	}

	if !removed {
		fmt.Println("No GitHub Copilot integration found")
	}
}
