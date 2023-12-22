package protoc_gen_ts

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestBuf(t *testing.T) {
	if err := os.RemoveAll(filepath.Join("build", "buf")); err != nil {
		t.Fatalf("failed to remove build directory: %v", err)
	}

	output := bytes.Buffer{}
	cmd := exec.Command("go", "run", "github.com/bufbuild/buf/cmd/buf@v1.28.1", "generate")
	cmd.Stderr = &output
	cmd.Stdout = &output
	cmd.Dir = "testdata"
	if err := cmd.Run(); err != nil {
		t.Fatalf("failed to run buf: %v\n%v", err, output.String())
	}

	for _, path := range []string{
		filepath.Join("build", "buf", "ts", "helloworld.ts"),
	} {
		if _, err := os.Stat(path); err != nil {
			t.Errorf("failed to stat %v: %v", path, err)
		}
	}
}

func TestProtoc(t *testing.T) {
	if _, err := exec.LookPath("protoc"); err != nil {
		t.Skip("protoc not found")
	}

	outDir := filepath.Join("build", "protoc")
	if err := os.RemoveAll(outDir); err != nil {
		t.Fatalf("failed to remove build directory: %v", err)
	}
	if err := os.RemoveAll(filepath.Join("build", "plugins")); err != nil {
		t.Fatalf("failed to remove build directory: %v", err)
	}

	plugin := "ts"
	output := bytes.Buffer{}
	cmd := exec.Command("go", "build", "-o", filepath.Join("build", "plugins", "protoc-gen-"+plugin), "./cmd/protoc-gen-"+plugin)
	cmd.Stderr = &output
	cmd.Stdout = &output
	if err := cmd.Run(); err != nil {
		t.Fatalf("failed to build plugin %v: %v\n%v", plugin, err, output.String())
	}

	if err := os.MkdirAll(filepath.Join(outDir, plugin), 0o755); err != nil {
		t.Fatalf("failed to create directory %v: %v", filepath.Join(outDir, plugin), err)
	}
	output = bytes.Buffer{}
	env := os.Environ()
	for i, val := range env {
		if strings.HasPrefix(val, "PATH=") {
			env[i] = "PATH=" + filepath.Join("build", "plugins") + string(os.PathListSeparator) + val[len("PATH="):]
		}
	}
	cmd = exec.Command(
		"protoc",
		"--ts_out="+filepath.Join(outDir, "ts"),
		"-I"+filepath.Join("testdata", "protos"),
		filepath.Join("testdata", "protos", "helloworld.proto"),
	)
	cmd.Stderr = &output
	cmd.Stdout = &output
	cmd.Env = env
	if err := cmd.Run(); err != nil {
		t.Fatalf("failed to run protoc: %v\n%v", err, output.String())
	}

	for _, path := range []string{
		filepath.Join("build", "protoc", "ts", "helloworld.ts"),
	} {
		if _, err := os.Stat(path); err != nil {
			t.Errorf("failed to stat %v: %v", path, err)
		}
	}
}
