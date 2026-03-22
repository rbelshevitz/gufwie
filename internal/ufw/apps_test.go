package ufw

import "testing"

func TestParseAppList(t *testing.T) {
	out := `
Available applications:
  OpenSSH
  Nginx Full
  Nginx HTTP
`
	apps, err := ParseAppList(out)
	if err != nil {
		t.Fatalf("ParseAppList: %v", err)
	}
	if len(apps) != 3 {
		t.Fatalf("expected 3, got %d: %#v", len(apps), apps)
	}
	if apps[0] != "OpenSSH" || apps[1] != "Nginx Full" {
		t.Fatalf("unexpected: %#v", apps)
	}
}

