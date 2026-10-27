package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestQuoteSieveStringEscapesQuotesAndBackslashes(t *testing.T) {
	got, err := quoteSieveString("Hello \"world\" \\ mailbox")
	if err != nil {
		t.Fatal(err)
	}
	want := "\"Hello \\\"world\\\" \\\\ mailbox\""
	if got != want {
		t.Fatalf("quoteSieveString() = %q, want %q", got, want)
	}
}

func TestQuoteSieveStringRejectsControlCharacters(t *testing.T) {
	if _, err := quoteSieveString("line one\nline two"); err == nil {
		t.Fatal("expected control character to be rejected")
	}
}

func TestQuoteSieveStringAllowsUnicode(t *testing.T) {
	got, err := quoteSieveString("Grüße 世界")
	if err != nil {
		t.Fatal(err)
	}
	if want := "\"Grüße 世界\""; got != want {
		t.Fatalf("quoteSieveString() = %q, want %q", got, want)
	}
}

func newInstallFixture(t *testing.T) (string, string, string) {
	t.Helper()
	root := t.TempDir()
	source := filepath.Join(root, "source.sieve")
	target := filepath.Join(root, "vacation.sieve")
	if err := os.WriteFile(source, []byte("require [\"vacation\"];\n"), 0o640); err != nil {
		t.Fatal(err)
	}
	if err := os.Chmod(source, 0o640); err != nil {
		t.Fatal(err)
	}
	return root, source, target
}

func TestInstallNoReplaceCopiesContentAndMode(t *testing.T) {
	_, source, target := newInstallFixture(t)
	if err := installNoReplace(source, target); err != nil {
		t.Fatal(err)
	}
	sourceBytes, err := os.ReadFile(source)
	if err != nil {
		t.Fatal(err)
	}
	targetBytes, err := os.ReadFile(target)
	if err != nil {
		t.Fatal(err)
	}
	if string(targetBytes) != string(sourceBytes) {
		t.Fatalf("installed content = %q, want %q", targetBytes, sourceBytes)
	}
	info, err := os.Stat(target)
	if err != nil {
		t.Fatal(err)
	}
	if got := info.Mode().Perm(); got != 0o640 {
		t.Fatalf("installed mode = %04o, want 0640", got)
	}
}

func TestInstallNoReplacePreservesExistingTarget(t *testing.T) {
	_, source, target := newInstallFixture(t)
	if err := os.WriteFile(target, []byte("keep me\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := installNoReplace(source, target); err == nil {
		t.Fatal("expected existing target to cause an error")
	}
	got, err := os.ReadFile(target)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != "keep me\n" {
		t.Fatalf("existing target changed: %q", got)
	}
}

func TestInstallNoReplacePreservesSymlink(t *testing.T) {
	root, source, target := newInstallFixture(t)
	other := filepath.Join(root, "other")
	if err := os.WriteFile(other, []byte("keep me\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(other, target); err != nil {
		t.Skipf("symlinks unavailable: %v", err)
	}
	if err := installNoReplace(source, target); err == nil {
		t.Fatal("expected symlink destination to cause an error")
	}
	got, err := os.ReadFile(other)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != "keep me\n" {
		t.Fatalf("symlink target changed: %q", got)
	}
}

func TestInstallNoReplaceMissingDirectory(t *testing.T) {
	root, source, _ := newInstallFixture(t)
	target := filepath.Join(root, "missing", "vacation.sieve")
	if err := installNoReplace(source, target); err == nil {
		t.Fatal("expected missing destination directory to cause an error")
	}
}
