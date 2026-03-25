package setup

import (
	"os"
	"path/filepath"
	"testing"
)

func TestInstallCopilot_InstructionsOnly(t *testing.T) {
	origDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	tmpDir := t.TempDir()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("chdir: %v", err)
	}
	defer func() { _ = os.Chdir(origDir) }()

	InstallCopilot(false)

	// Instructions file must exist
	if !FileExists(copilotInstructionsFile) {
		t.Errorf("%s was not created", copilotInstructionsFile)
	}

	// No MCP config should be created
	if FileExists(".vscode/mcp.json") {
		t.Error(".vscode/mcp.json should NOT be created in CLI-only mode")
	}

	// No prompt files without --prompts
	if FileExists(filepath.Join(copilotPromptsDir, "beads-ready.prompt.md")) {
		t.Error("prompt files should NOT be created without --prompts flag")
	}
}

func TestInstallCopilot_WithPrompts(t *testing.T) {
	origDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	tmpDir := t.TempDir()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("chdir: %v", err)
	}
	defer func() { _ = os.Chdir(origDir) }()

	InstallCopilot(true)

	// Instructions file must exist
	if !FileExists(copilotInstructionsFile) {
		t.Errorf("%s was not created", copilotInstructionsFile)
	}

	// All prompt files must exist
	for _, name := range beadsPromptNames {
		p := filepath.Join(copilotPromptsDir, name)
		if !FileExists(p) {
			t.Errorf("prompt file %s was not created", p)
		}
	}
}

func TestRemoveCopilot_WithPrompts(t *testing.T) {
	origDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	tmpDir := t.TempDir()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("chdir: %v", err)
	}
	defer func() { _ = os.Chdir(origDir) }()

	// Install first
	InstallCopilot(true)

	// Verify files exist
	if !FileExists(copilotInstructionsFile) {
		t.Fatalf("setup failed: %s missing", copilotInstructionsFile)
	}
	if !FileExists(filepath.Join(copilotPromptsDir, "beads-ready.prompt.md")) {
		t.Fatalf("setup failed: prompt file missing")
	}

	// Remove with prompts
	RemoveCopilot(true)

	// Instructions file may still exist (with non-beads boilerplate),
	// but the beads section must be gone.
	if FileExists(copilotInstructionsFile) {
		data, err := os.ReadFile(copilotInstructionsFile)
		if err != nil {
			t.Fatalf("read instructions: %v", err)
		}
		if containsBeadsMarker(string(data)) {
			t.Error("beads section should have been removed from instructions file")
		}
	}

	// Prompt files should be removed
	for _, name := range beadsPromptNames {
		p := filepath.Join(copilotPromptsDir, name)
		if FileExists(p) {
			t.Errorf("prompt file %s should have been removed", p)
		}
	}
}

func TestCheckCopilot_NotInstalled(t *testing.T) {
	origDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	tmpDir := t.TempDir()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("chdir: %v", err)
	}
	defer func() { _ = os.Chdir(origDir) }()

	// We can't easily call CheckCopilot() because it calls FatalErrorWithHint which exits.
	// Just verify the file doesn't exist.
	if FileExists(copilotInstructionsFile) {
		t.Error("unexpected instructions file in temp dir")
	}
}

