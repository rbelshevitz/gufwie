package appui

import "strings"

func freeformHelpText() string {
	return strings.Join([]string{
		"Freeform add:",
		"  You type arguments exactly as you would after `ufw <action>`.",
		"  Example: if you run `ufw allow 22/tcp`, enter just `22/tcp` here.",
		"",
		"Examples:",
		"  22/tcp",
		"  in 80/tcp",
		"  in on eth0 to any port 443 proto tcp",
		"  from 192.168.1.0/24 to any port 22 proto tcp",
		"  from 10.0.0.5 to any port 5432 proto tcp",
		"  proto tcp from any to any port 80 comment 'web'",
		"",
		"Tips:",
		"  - Use Service/Profile mode for known services (OpenSSH, Nginx, …).",
		"  - Use Wizard mode if you prefer structured fields.",
	}, "\n")
}

func wizardHelpText() string {
	return strings.Join([]string{
		"Wizard add:",
		"  Builds a rule from fields and shows a preview command before applying.",
		"",
		"Notes:",
		"  - Direction: in/out (empty means default).",
		"  - From/To can be IP/CIDR or 'any' (leave blank to omit).",
		"  - Port/service can be '22' or '22/tcp'.",
		"  - Proto optional; if you type '22/tcp' it will be inferred when Proto is empty.",
		"  - Comment is optional (stored as a rule comment).",
	}, "\n")
}
