package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"unicode"
)

func quoteSieveString(value string) (string, error) {
	for _, r := range value {
		if unicode.IsControl(r) {
			return "", fmt.Errorf("control characters are not supported in quoted strings")
		}
	}
	value = strings.ReplaceAll(value, "\\", "\\\\")
	value = strings.ReplaceAll(value, "\"", "\\\"")
	return "\"" + value + "\"", nil
}

// installNoReplace publishes source at target atomically, failing if target
// already exists. The temporary file is created in the target directory so
// os.Link uses the same filesystem.
func installNoReplace(source, target string) error {
	info, err := os.Stat(source)
	if err != nil {
		return fmt.Errorf("stat source: %w", err)
	}
	if !info.Mode().IsRegular() {
		return fmt.Errorf("source must be a regular file")
	}

	directory := filepath.Dir(target)
	dirInfo, err := os.Stat(directory)
	if err != nil {
		return fmt.Errorf("destination directory: %w", err)
	}
	if !dirInfo.IsDir() {
		return fmt.Errorf("destination parent is not a directory")
	}

	input, err := os.Open(source)
	if err != nil {
		return fmt.Errorf("open source: %w", err)
	}
	defer input.Close()

	temporary, err := os.CreateTemp(directory, ".vacation-manager.*")
	if err != nil {
		return fmt.Errorf("create temporary file: %w", err)
	}
	temporaryName := temporary.Name()
	defer os.Remove(temporaryName)

	if _, err := io.Copy(temporary, input); err != nil {
		temporary.Close()
		return fmt.Errorf("copy source: %w", err)
	}
	if err := temporary.Chmod(info.Mode().Perm()); err != nil {
		temporary.Close()
		return fmt.Errorf("set file mode: %w", err)
	}
	if err := temporary.Close(); err != nil {
		return fmt.Errorf("close temporary file: %w", err)
	}

	// Unlike a stat-then-rename sequence, link creation fails atomically if
	// target already exists (including when it is a symlink).
	if err := os.Link(temporaryName, target); err != nil {
		return fmt.Errorf("publish target without replacement: %w", err)
	}
	return nil
}

func run(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("usage: vacation-manager-prototype <quote VALUE | install SOURCE TARGET>")
	}

	switch args[0] {
	case "quote":
		flags := flag.NewFlagSet("quote", flag.ContinueOnError)
		if err := flags.Parse(args[1:]); err != nil {
			return err
		}
		if flags.NArg() != 1 {
			return fmt.Errorf("usage: vacation-manager-prototype quote VALUE")
		}
		quoted, err := quoteSieveString(flags.Arg(0))
		if err != nil {
			return err
		}
		fmt.Println(quoted)
		return nil
	case "install":
		flags := flag.NewFlagSet("install", flag.ContinueOnError)
		if err := flags.Parse(args[1:]); err != nil {
			return err
		}
		if flags.NArg() != 2 {
			return fmt.Errorf("usage: vacation-manager-prototype install SOURCE TARGET")
		}
		return installNoReplace(flags.Arg(0), flags.Arg(1))
	default:
		return fmt.Errorf("unknown command %q", args[0])
	}
}

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, "vacation-manager prototype:", err)
		os.Exit(1)
	}
}
