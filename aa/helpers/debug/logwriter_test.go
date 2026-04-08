package debug_test

import (
	"fmt"
	"os"
	"path"
	"runtime"
	"testing"
	"time"

	"github.com/aarioai/airis/aa/helpers/debug"
)

func testLogWriterBuffer(t *testing.T, bufferSize int) {
	// not support windows yet
	if runtime.GOOS == "windows" {
		fmt.Println("Skipping TestLogWriterBuffer on Windows")
		return
	}

	dir := t.TempDir()
	symlink := path.Join(t.TempDir(), "logwriter_symlink_test.log")

	lw, err := debug.NewLogWriter(dir, 0777, bufferSize, symlink)
	if err != nil {
		t.Fatal(err.Error())
	}

	now := time.Now().Format("2006-01-02 15:04:05")
	testData := []byte(fmt.Sprintf(now+" buffer:%d\n", bufferSize))

	n := 0
	if n, err = lw.Write(testData); err != nil {
		t.Errorf("failed to write log: %v", err)
	}
	if n != len(testData) {
		t.Errorf("written bytes mismatch, expected %d, got %d", len(testData), n)
	}

	if bufferSize > 0 {
		debug.Console(lw.Flush())
	}

	// verify written data
	content, err := os.ReadFile(symlink)
	if err != nil {
		t.Errorf("failed to read log file: %v", err)
	}
	if string(content) != string(testData) {
		t.Errorf("log content mismatch, expected %q, got %q", testData, content)
	}

	fmt.Println("------------>>>>>>")
}

func TestLogWriter(t *testing.T) {
	testLogWriterBuffer(t, 0)
	testLogWriterBuffer(t, 0)
}

func TestLogWriterBuffer(t *testing.T) {
	testLogWriterBuffer(t, 2048)
}
