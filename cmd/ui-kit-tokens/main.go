// Command ui-kit-tokens writes css/tokens.css from tokens/tokens.json. Run it from the repository root.
package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/ElderlyGate/ui-kit/tokens"
)

func main() {
	if err := run("tokens/tokens.json", "css/tokens.css"); err != nil {
		fmt.Fprintln(os.Stderr, "ui-kit-tokens:", err)
		os.Exit(1)
	}
}

func run(in, out string) error {
	data, err := os.ReadFile(in)
	if err != nil {
		return err
	}
	set, err := tokens.Parse(data)
	if err != nil {
		return err
	}
	return writeWhole(out, set.CSS())
}

// writeWhole writes into a temporary file beside path, then renames it into place (G-CODE-5).
func writeWhole(path string, data []byte) error {
	tmp, err := os.CreateTemp(filepath.Dir(path), ".tokens-*.css")
	if err != nil {
		return err
	}
	defer os.Remove(tmp.Name())
	if _, err := tmp.Write(data); err != nil {
		tmp.Close()
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	if err := os.Chmod(tmp.Name(), 0o644); err != nil {
		return err
	}
	return os.Rename(tmp.Name(), path)
}
