package swallow

import (
	"crypto/rand"
	"encoding/hex"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// timestampFormat sorts lexically and avoids colons for filesystem friendliness.
const timestampFormat = "2006-01-02T15-04-05"

func swallowDir() (string, error) {
	if dir := os.Getenv("SWALLOW_DIR"); dir != "" {
		return dir, nil
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".swallow"), nil
}

// originDir resolves the log directory for the current origin, the working
// directory swallow was invoked from.
func originDir() (string, error) {
	dir, err := swallowDir()
	if err != nil {
		return "", err
	}

	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	return filepath.Join(dir, slug(cwd)), nil
}

func createLog(argv0 string) (string, *os.File, error) {
	origin, err := originDir()
	if err != nil {
		return "", nil, err
	}
	if err := os.MkdirAll(origin, 0o755); err != nil {
		return "", nil, err
	}

	suffix, err := randomSuffix()
	if err != nil {
		return "", nil, err
	}

	// The command component is slugged so the log name is always shell-safe:
	// the read hint printed after a run must be runnable verbatim.
	name := time.Now().Format(timestampFormat) + "-" + slug(filepath.Base(argv0)) + "-" + suffix + ".log"
	path := filepath.Join(origin, name)
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND|os.O_EXCL, 0o644)
	if err != nil {
		return "", nil, err
	}
	return path, file, nil
}

func slug(path string) string {
	var b strings.Builder
	lastDash := true
	for _, r := range path {
		alnum := (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9')
		if alnum {
			b.WriteRune(r)
			lastDash = false
		} else if !lastDash {
			b.WriteByte('-')
			lastDash = true
		}
	}

	result := strings.TrimSuffix(b.String(), "-")
	if result == "" {
		return "root"
	}
	return result
}

func randomSuffix() (string, error) {
	bytes := make([]byte, 3)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
