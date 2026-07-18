package swallow

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// Read prints a stored log verbatim to stdout. The same-origin gate only
// admits paths pointing directly into the origin directory of the current
// working directory; everything else is refused without disclosing whether
// the path exists. The comparison is lexical — the gate scopes an agent to
// its own project's logs, it is not a security boundary (see ADR 0008).
func Read(path string) int {
	origin, err := originDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "swallow: %v\n", err)
		return 1
	}

	resolved := path
	if !filepath.IsAbs(resolved) {
		cwd, err := os.Getwd()
		if err != nil {
			fmt.Fprintf(os.Stderr, "swallow: %v\n", err)
			return 1
		}
		resolved = filepath.Join(cwd, resolved)
	}
	resolved = filepath.Clean(resolved)

	if filepath.Dir(resolved) != filepath.Clean(origin) {
		fmt.Fprintf(os.Stderr, "swallow: refusing to read %s: not a log of the current origin\n", path)
		return 1
	}

	file, err := os.Open(resolved)
	if err != nil {
		fmt.Fprintf(os.Stderr, "swallow: %v\n", err)
		return 1
	}
	defer func() { _ = file.Close() }()

	if _, err := io.Copy(os.Stdout, file); err != nil {
		fmt.Fprintf(os.Stderr, "swallow: %v\n", err)
		return 1
	}
	return 0
}
