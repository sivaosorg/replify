package test

import (
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"
	"time"

	"github.com/sivaosorg/replify/pkg/sysx"
)

func TestSysx_Move(t *testing.T) {
	dir, err := os.MkdirTemp("", "sysx_move_test_*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	src := filepath.Join(dir, "src.txt")
	dst := filepath.Join(dir, "dst.txt")
	content := "move me"

	if err := os.WriteFile(src, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	if err := sysx.Move(src, dst); err != nil {
		t.Fatalf("Move failed: %v", err)
	}

	if sysx.FileExists(src) {
		t.Error("source file still exists after Move")
	}

	got, err := os.ReadFile(dst)
	if err != nil {
		t.Fatalf("could not read destination file: %v", err)
	}
	if string(got) != content {
		t.Errorf("destination content = %q; want %q", string(got), content)
	}
}

func TestSysx_Touch(t *testing.T) {
	dir, err := os.MkdirTemp("", "sysx_touch_test_*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	path := filepath.Join(dir, "new_file.txt")

	// Test creation
	if err := sysx.Touch(path); err != nil {
		t.Fatalf("Touch creation failed: %v", err)
	}
	if !sysx.FileExists(path) {
		t.Error("Touch did not create file")
	}

	// Test update
	oldFi, _ := os.Stat(path)
	time.Sleep(10 * time.Millisecond) // Ensure time difference
	if err := sysx.Touch(path); err != nil {
		t.Fatalf("Touch update failed: %v", err)
	}
	newFi, _ := os.Stat(path)

	if !newFi.ModTime().After(oldFi.ModTime()) {
		t.Errorf("Touch did not update mod time: %v -> %v", oldFi.ModTime(), newFi.ModTime())
	}
}

func TestSysx_IsBinary(t *testing.T) {
	dir, err := os.MkdirTemp("", "sysx_binary_test_*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	txtPath := filepath.Join(dir, "text.txt")
	binPath := filepath.Join(dir, "bin.dat")

	if err := os.WriteFile(txtPath, []byte("hello world"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(binPath, []byte{0x48, 0x00, 0x4c, 0x4c, 0x4f}, 0644); err != nil {
		t.Fatal(err)
	}

	isBin, err := sysx.IsBinary(txtPath)
	if err != nil {
		t.Fatal(err)
	}
	if isBin {
		t.Errorf("IsBinary(%q) = true; want false", txtPath)
	}

	isBin, err = sysx.IsBinary(binPath)
	if err != nil {
		t.Fatal(err)
	}
	if !isBin {
		t.Errorf("IsBinary(%q) = false; want true", binPath)
	}
}

func TestSysx_Hashing(t *testing.T) {
	dir, err := os.MkdirTemp("", "sysx_hash_test_*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	path := filepath.Join(dir, "data.txt")
	content := "replify"
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	// MD5 of "replify" is 5f7be9021c5932d1e511859f201c86d8
	wantMD5 := "5f7be9021c5932d1e511859f201c86d8"
	gotMD5, err := sysx.FileMD5(path)
	if err != nil {
		t.Fatal(err)
	}
	if gotMD5 != wantMD5 {
		t.Errorf("MD5 = %q; want %q", gotMD5, wantMD5)
	}

	// SHA256 of "replify" is 89ac3e4273d8ce6ccab52871f2bbdffcdee4399617b0eff0f2b94aadfafe4f0e
	wantSHA256 := "89ac3e4273d8ce6ccab52871f2bbdffcdee4399617b0eff0f2b94aadfafe4f0e"
	gotSHA256, err := sysx.FileSHA256(path)
	if err != nil {
		t.Fatal(err)
	}
	if gotSHA256 != wantSHA256 {
		t.Errorf("SHA256 = %q; want %q", gotSHA256, wantSHA256)
	}
}

func TestSysx_CopyDir(t *testing.T) {
	parent, err := os.MkdirTemp("", "sysx_copydir_test_*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(parent)

	src := filepath.Join(parent, "src")
	dst := filepath.Join(parent, "dst")

	// Create structure:
	// src/
	//   file1.txt
	//   subdir/
	//     file2.txt
	if err := os.MkdirAll(filepath.Join(src, "subdir"), 0755); err != nil {
		t.Fatal(err)
	}
	os.WriteFile(filepath.Join(src, "file1.txt"), []byte("one"), 0644)
	os.WriteFile(filepath.Join(src, "subdir", "file2.txt"), []byte("two"), 0644)

	// Test symlink replication (on non-Windows)
	if runtime.GOOS != "windows" {
		os.Symlink("file1.txt", filepath.Join(src, "link.txt"))
	}

	if err := sysx.CopyDir(src, dst); err != nil {
		t.Fatalf("CopyDir failed: %v", err)
	}

	// Verify
	files := []string{
		filepath.Join(dst, "file1.txt"),
		filepath.Join(dst, "subdir", "file2.txt"),
	}
	for _, f := range files {
		if !sysx.FileExists(f) {
			t.Errorf("copied file missing: %s", f)
		}
	}

	if runtime.GOOS != "windows" {
		link := filepath.Join(dst, "link.txt")
		if !sysx.IsSymlink(link) {
			t.Errorf("symlink was not replicated: %s", link)
		}
	}
}

func TestSysx_ClearDir(t *testing.T) {
	dir, err := os.MkdirTemp("", "sysx_cleardir_test_*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	os.WriteFile(filepath.Join(dir, "f1.txt"), []byte("data"), 0644)
	os.Mkdir(filepath.Join(dir, "d1"), 0755)

	if err := sysx.ClearDir(dir); err != nil {
		t.Fatalf("ClearDir failed: %v", err)
	}

	entries, _ := os.ReadDir(dir)
	if len(entries) != 0 {
		t.Errorf("ClearDir did not empty directory; remaining entries: %d", len(entries))
	}
	if !sysx.DirExists(dir) {
		t.Error("ClearDir removed the directory itself")
	}
}

func TestSysx_IsEmpty(t *testing.T) {
	dir, err := os.MkdirTemp("", "sysx_isempty_test_*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	empty, _ := sysx.IsDirEmpty(dir)
	if !empty {
		t.Errorf("IsEmpty(%q) = false; want true", dir)
	}

	os.WriteFile(filepath.Join(dir, "not_empty.txt"), []byte("x"), 0644)
	empty, _ = sysx.IsDirEmpty(dir)
	if empty {
		t.Errorf("IsEmpty(%q) = true; want false", dir)
	}
}

func TestSysx_ReadTail(t *testing.T) {
	dir, err := os.MkdirTemp("", "sysx_tail_test_*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	path := filepath.Join(dir, "lines.txt")
	lines := []string{"line1", "line2", "line3", "line4", "line5"}
	if err := sysx.WriteLines(path, lines); err != nil {
		t.Fatal(err)
	}

	// Test tail 2
	got, err := sysx.Tail(path, 2)
	if err != nil {
		t.Fatalf("ReadTail failed: %v", err)
	}
	want := []string{"line4", "line5"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("ReadTail(2) = %v; want %v", got, want)
	}

	// Test tail more than exists
	got, err = sysx.Tail(path, 10)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, lines) {
		t.Errorf("ReadTail(10) = %v; want %v", got, lines)
	}
}

func TestSysx_CountLines_Optimized(t *testing.T) {
	dir, err := os.MkdirTemp("", "sysx_count_test_*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	path := filepath.Join(dir, "many_lines.txt")
	numLines := 1000
	lines := make([]string, numLines)
	for i := 0; i < numLines; i++ {
		lines[i] = "this is a line"
	}
	if err := sysx.WriteLines(path, lines); err != nil {
		t.Fatal(err)
	}

	got, err := sysx.CountLines(path)
	if err != nil {
		t.Fatalf("CountLines failed: %v", err)
	}
	if got != numLines {
		t.Errorf("CountLines = %d; want %d", got, numLines)
	}
}
