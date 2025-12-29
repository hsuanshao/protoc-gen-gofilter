package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestEndToEndGeneration(t *testing.T) {
	// 1. Build the plugin binary
	tempDir := t.TempDir()
	pluginPath := filepath.Join(tempDir, "protoc-gen-gofilter")

	// Assuming we are in pkgs/filter directory, main.go is here.
	buildCmd := exec.Command("go", "build", "-o", pluginPath, ".")
	if out, err := buildCmd.CombinedOutput(); err != nil {
		t.Fatalf("Failed to build plugin: %v\nOutput: %s", err, out)
	}

	// 2. Prepare protoc command
	// We need to point to the project root for imports to work correctly
	projectRoot, err := filepath.Abs("../..")
	if err != nil {
		t.Fatalf("Failed to resolve project root: %v", err)
	}

	// Ensure protoc-gen-go is available or install it?
	// We assume the environment has it since it's a dev environment.
	// But just in case, we might need to rely on what's available.

	// Create output directory for generated files
	outDir := filepath.Join(tempDir, "out")
	if err := os.MkdirAll(outDir, 0755); err != nil {
		t.Fatalf("Failed to create output dir: %v", err)
	}

	// 3. Run protoc
	testProto := "testdata/test.proto"

	// We need to pass --plugin to use our built binary
	// And we need --gofilter_out to trigger our plugin
	// We also need --go_out because our plugin expects standard go struct to be generated?
	// Actually our plugin just generates the side file _filter.pb.go.
	// But if we want to compile the result we need --go_out too. For now we just check generated file content.

	protocCmd := exec.Command("protoc",
		"--plugin=protoc-gen-gofilter="+pluginPath,
		"--gofilter_out="+outDir,
		"--gofilter_opt=paths=source_relative",
		"--proto_path="+projectRoot,
		"--proto_path="+filepath.Join(projectRoot, "cmd/protoc-gen-gofilter"),
		testProto,
	)

	if out, err := protocCmd.CombinedOutput(); err != nil {
		t.Logf("Project Root: %s", projectRoot)
		t.Fatalf("Failed to run protoc: %v\nOutput: %s", err, out)
	}

	// 4. Verify generated file exists and contains expected content
	genFile := filepath.Join(outDir, "testdata/test_filter.pb.go")
	contentBytes, err := os.ReadFile(genFile)
	if err != nil {
		t.Fatalf("Failed to read generated file: %v", err)
	}
	content := string(contentBytes)

	// Check for expected snippets
	expectedSnippets := []string{
		`func init() {`,
		`_PermIdx_TestMessage_PrivateField = `,
		`.Register("test.private")`,
		`func (x *TestMessage) FilterFields(mask `,
		`if !mask.Has(_PermIdx_TestMessage_PrivateField) {`,
		`x.PrivateField = ""`,
		`x.SecretNumber = 0`,
	}

	for _, snippet := range expectedSnippets {
		if !strings.Contains(content, snippet) {
			t.Errorf("Generated code missing snippet: %q", snippet)
		}
	}

	// Verify that Test2Message (which has no filter options) does NOT have specific code generated
	unexpectedSnippets := []string{
		`func (x *Test2Message) FilterFields`,
	}
	for _, snippet := range unexpectedSnippets {
		if strings.Contains(content, snippet) {
			t.Errorf("Generated code contains unexpected snippet: %q", snippet)
		}
	}
}

func TestOptionalFieldGeneration(t *testing.T) {
	// 1. Build the plugin binary
	tempDir := t.TempDir()
	pluginPath := filepath.Join(tempDir, "protoc-gen-gofilter")

	buildCmd := exec.Command("go", "build", "-o", pluginPath, ".")
	if out, err := buildCmd.CombinedOutput(); err != nil {
		t.Fatalf("Failed to build plugin: %v\nOutput: %s", err, out)
	}

	// 2. Prepare protoc command
	projectRoot, err := filepath.Abs("../..")
	if err != nil {
		t.Fatalf("Failed to resolve project root: %v", err)
	}

	outDir := filepath.Join(tempDir, "out")
	if err := os.MkdirAll(outDir, 0755); err != nil {
		t.Fatalf("Failed to create output dir: %v", err)
	}

	// 3. Run protoc with optional.proto
	testProto := "testdata/optional.proto"

	protocCmd := exec.Command("protoc",
		"--plugin=protoc-gen-gofilter="+pluginPath,
		"--gofilter_out="+outDir,
		"--gofilter_opt=paths=source_relative",
		"--proto_path="+projectRoot,
		"--proto_path="+filepath.Join(projectRoot, "cmd/protoc-gen-gofilter"),
		testProto,
	)

	if out, err := protocCmd.CombinedOutput(); err != nil {
		// This is where we expect it to fail if optional is not supported
		t.Logf("Project Root: %s", projectRoot)
		t.Fatalf("Failed to run protoc: %v\nOutput: %s", err, out)
	}

	// 4. Verify generated file exists
	genFile := filepath.Join(outDir, "testdata/optional_filter.pb.go")
	contentBytes, err := os.ReadFile(genFile)
	if err != nil {
		t.Fatalf("Failed to read generated file: %v", err)
	}
	content := string(contentBytes)

	// Check for proper nil assignment for optional field
	expectedSnippets := []string{
		`func (x *OptionalMessage) FilterFields(mask `,
		`x.SecretOptional = nil`,
	}

	for _, snippet := range expectedSnippets {
		if !strings.Contains(content, snippet) {
			t.Errorf("Generated code missing snippet: %q", snippet)
		}
	}
}
