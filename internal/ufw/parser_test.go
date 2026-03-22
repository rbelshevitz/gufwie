package ufw

import "testing"

func TestParseNumberedStatus_ActiveWithRules(t *testing.T) {
	out := `
Status: active

To                         Action      From
--                         ------      ----
[ 1] 22/tcp                ALLOW IN    Anywhere
[ 2] 22/tcp (v6)           ALLOW IN    Anywhere (v6)
[ 3] 80/tcp                DENY IN     10.0.0.0/8
`
	ns, err := ParseNumberedStatus(out)
	if err != nil {
		t.Fatalf("ParseNumberedStatus: %v", err)
	}
	if !ns.Status.Active {
		t.Fatalf("expected active")
	}
	if len(ns.Rules) != 3 {
		t.Fatalf("expected 3 rules, got %d", len(ns.Rules))
	}
	if ns.Rules[0].Number != 1 || ns.Rules[0].To != "22/tcp" || ns.Rules[0].From != "Anywhere" || ns.Rules[0].V6 {
		t.Fatalf("bad rule[0]: %+v", ns.Rules[0])
	}
	if !ns.Rules[1].V6 || ns.Rules[1].From != "Anywhere" {
		t.Fatalf("bad rule[1]: %+v", ns.Rules[1])
	}
}

func TestSplitArgs(t *testing.T) {
	args, err := SplitArgs(`from "192.168.1.0/24" to any port 22 proto tcp`)
	if err != nil {
		t.Fatalf("SplitArgs: %v", err)
	}
	if len(args) < 2 || args[1] != "192.168.1.0/24" {
		t.Fatalf("unexpected args: %#v", args)
	}
}

