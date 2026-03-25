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
	copilotHooksDir         = ".github/hooks"
	copilotSkillsDir        = ".github/skills"
	copilotAgentsDir        = ".github/agents"
	copilotModularInsDir    = ".github/instructions"
)

//go:embed copilot_prompts/*.prompt.md
var copilotPromptFiles embed.FS

//go:embed copilot_cli/hooks.json
var copilotHooksFile embed.FS

//go:embed copilot_cli/skills/beads/SKILL.md
var copilotSkillFile embed.FS

//go:embed copilot_cli/agents/beads.agent.md
var copilotAgentFile embed.FS

//go:embed copilot_cli/instructions/beads.instructions.md
var copilotInstructionsModularFile embed.FS

// beadsPromptNames lists the prompt files shipped with the integration.
var beadsPromptNames = []string{
	"beads-ready.prompt.md",
	"beads-create.prompt.md",
	"beads-workflow.prompt.md",
	"beads-fleet.prompt.md",
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
// When installCLI is true, Copilot CLI hooks, skills, and agents are installed.
func InstallCopilot(installPrompts bool, installCLI bool) {
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

	// 4. Install Copilot CLI features (hooks, skills, agents)
	if installCLI {
		installCopilotCLIFeatures()
	}

	fmt.Println()
	if installCLI {
		fmt.Println("Copilot CLI integration installed.")
		fmt.Println("  Hooks auto-prime beads context on session start and push on session end.")
		fmt.Println("  The beads skill and agent are available via /skills and /agent.")
	} else {
		fmt.Println("Reload VS Code for changes to take effect.")
	}
	if !installPrompts && !installCLI {
		fmt.Println()
		fmt.Println("Tip: use --prompts to also install reusable prompt files to .github/prompts/")
	}
	if !installCLI {
		fmt.Println()
		fmt.Println("For Copilot CLI, install hooks, skills, and agents:")
		fmt.Println("  bd setup copilot --cli")
	}
	fmt.Println()
	fmt.Println("Pre-approve bd commands in Copilot CLI:")
	fmt.Println("  copilot --allow-tool='shell(bd:*)'")
	fmt.Println("  Or use /allow-all in a session")
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

// installCopilotCLIFeatures installs Copilot CLI hooks, skills, and agents.
func installCopilotCLIFeatures() {
	// Hooks
	if err := EnsureDir(copilotHooksDir, 0o755); err != nil {
		FatalError("%v", err)
	}
	hooksData, err := copilotHooksFile.ReadFile("copilot_cli/hooks.json")
	if err != nil {
		FatalError("read embedded hooks: %v", err)
	}
	hooksDest := filepath.Join(copilotHooksDir, "beads.json")
	if err := atomicWriteFile(hooksDest, hooksData); err != nil {
		FatalError("write hooks %s: %v", hooksDest, err)
	}
	fmt.Printf("✓ Hook:   %s\n", hooksDest)

	// Skill
	skillDir := filepath.Join(copilotSkillsDir, "beads")
	if err := EnsureDir(skillDir, 0o755); err != nil {
		FatalError("%v", err)
	}
	skillData, err := copilotSkillFile.ReadFile("copilot_cli/skills/beads/SKILL.md")
	if err != nil {
		FatalError("read embedded skill: %v", err)
	}
	skillDest := filepath.Join(skillDir, "SKILL.md")
	if err := atomicWriteFile(skillDest, skillData); err != nil {
		FatalError("write skill %s: %v", skillDest, err)
	}
	fmt.Printf("✓ Skill:  %s\n", skillDest)

	// Agent
	if err := EnsureDir(copilotAgentsDir, 0o755); err != nil {
		FatalError("%v", err)
	}
	agentData, err := copilotAgentFile.ReadFile("copilot_cli/agents/beads.agent.md")
	if err != nil {
		FatalError("read embedded agent: %v", err)
	}
	agentDest := filepath.Join(copilotAgentsDir, "beads.agent.md")
	if err := atomicWriteFile(agentDest, agentData); err != nil {
		FatalError("write agent %s: %v", agentDest, err)
	}
	fmt.Printf("✓ Agent:  %s\n", agentDest)

	// Modular instructions (Copilot CLI v1.0.11+ auto-discovers .github/instructions/)
	if err := EnsureDir(copilotModularInsDir, 0o755); err != nil {
		FatalError("%v", err)
	}
	insData, err := copilotInstructionsModularFile.ReadFile("copilot_cli/instructions/beads.instructions.md")
	if err != nil {
		FatalError("read embedded instructions: %v", err)
	}
	insDest := filepath.Join(copilotModularInsDir, "beads.instructions.md")
	if err := atomicWriteFile(insDest, insData); err != nil {
		FatalError("write instructions %s: %v", insDest, err)
	}
	fmt.Printf("✓ Instructions: %s\n", insDest)
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

		// Report CLI features
		if FileExists(filepath.Join(copilotHooksDir, "beads.json")) {
			fmt.Printf("  Hooks:        %s/beads.json\n", copilotHooksDir)
		}
		if FileExists(filepath.Join(copilotSkillsDir, "beads", "SKILL.md")) {
			fmt.Printf("  Skill:        %s/beads/SKILL.md\n", copilotSkillsDir)
		}
		if FileExists(filepath.Join(copilotAgentsDir, "beads.agent.md")) {
			fmt.Printf("  Agent:        %s/beads.agent.md\n", copilotAgentsDir)
		}
		if FileExists(filepath.Join(copilotModularInsDir, "beads.instructions.md")) {
			fmt.Printf("  Modular:      %s/beads.instructions.md\n", copilotModularInsDir)
		}
		return
	}

	FatalErrorWithHint("GitHub Copilot integration not installed", "Run: bd setup copilot")
}

// RemoveCopilot removes GitHub Copilot integration files.
// When removePrompts is true, beads prompt files are also removed from .github/prompts/.
// When removeCLI is true, Copilot CLI hooks, skills, and agents are also removed.
func RemoveCopilot(removePrompts bool, removeCLI bool) {
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

	// Remove CLI features
	if removeCLI {
		removed = removeCopilotCLIFile(filepath.Join(copilotHooksDir, "beads.json")) || removed
		removed = removeCopilotCLIFile(filepath.Join(copilotSkillsDir, "beads", "SKILL.md")) || removed
		removed = removeCopilotCLIFile(filepath.Join(copilotAgentsDir, "beads.agent.md")) || removed
		removed = removeCopilotCLIFile(filepath.Join(copilotModularInsDir, "beads.instructions.md")) || removed
		// Clean up empty directories
		_ = os.Remove(filepath.Join(copilotSkillsDir, "beads"))
		_ = os.Remove(copilotSkillsDir)
		_ = os.Remove(copilotAgentsDir)
		_ = os.Remove(copilotHooksDir)
		_ = os.Remove(copilotModularInsDir)
	}

	if !removed {
		fmt.Println("No GitHub Copilot integration found")
	}
}

// removeCopilotCLIFile removes a single file and prints status. Returns true if removed.
func removeCopilotCLIFile(path string) bool {
	if FileExists(path) {
		if err := os.Remove(path); err != nil {
			FatalError("remove %s: %v", path, err)
		}
		fmt.Printf("✓ Removed %s\n", path)
		return true
	}
	return false
}
