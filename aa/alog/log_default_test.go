package alog

import "testing"

func TestNewDefaultLogIsIndependent(t *testing.T) {
	a := NewDefaultLog(LevelDebug).(*xlog)
	b := NewDefaultLog(LevelError).(*xlog)
	if a == b {
		t.Fatal("NewDefaultLog returned one shared logger")
	}
	if a.Level != LevelDebug || b.Level != LevelError {
		t.Fatalf("levels = %v, %v", a.Level, b.Level)
	}
}
