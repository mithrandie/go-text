package color

import "testing"

func TestError(t *testing.T) {
	s := "message"
	expect := "\033[31;1mmessage\033[0m"
	result := Error(s)
	if result != expect {
		t.Errorf("result = %q, want %q", result, expect)
	}
}

func TestWarn(t *testing.T) {
	s := "message"
	expect := "\033[33;1mmessage\033[0m"
	result := Warn(s)
	if result != expect {
		t.Errorf("result = %q, want %q", result, expect)
	}
}

func TestNotice(t *testing.T) {
	s := "message"
	expect := "\033[32mmessage\033[0m"
	result := Notice(s)
	if result != expect {
		t.Errorf("result = %q, want %q", result, expect)
	}
}

func TestInfo(t *testing.T) {
	s := "message"
	expect := "message"
	result := Info(s)
	if result != expect {
		t.Errorf("result = %q, want %q", result, expect)
	}
}
