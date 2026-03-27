package ufw

import (
	"context"
	"errors"
	"testing"
)

type fakeRunner struct {
	stdout string
	stderr string
	code   int
	err    error
}

func (r fakeRunner) Run(context.Context, string, ...string) (string, string, int, error) {
	return r.stdout, r.stderr, r.code, r.err
}

func TestIsNotPrivilegedMessage(t *testing.T) {
	cases := []struct {
		in   string
		want bool
	}{
		{"ERROR: You need to be root to run this program", true},
		{"you need to be root to run this program", true},
		{"Permission denied", true},
		{"some other error", false},
		{"", false},
	}

	for _, tc := range cases {
		if got := isNotPrivilegedMessage(tc.in); got != tc.want {
			t.Fatalf("isNotPrivilegedMessage(%q)=%v, want %v", tc.in, got, tc.want)
		}
	}
}

func TestClientRun_NotPrivilegedError(t *testing.T) {
	c := NewClient(fakeRunner{
		stderr: "ERROR: You need to be root to run this program",
		code:   1,
	}, Config{})

	_, err := c.run(context.Background(), "status", "numbered")
	if err == nil {
		t.Fatalf("expected error")
	}
	var npe *NotPrivilegedError
	if !errors.As(err, &npe) {
		t.Fatalf("expected NotPrivilegedError, got %T: %v", err, err)
	}
	if npe.Detail == "" {
		t.Fatalf("expected detail to be set")
	}
}

