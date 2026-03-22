package ufw

import (
	"bufio"
	"fmt"
	"regexp"
	"strings"
)

var reStatus = regexp.MustCompile(`(?i)^Status:\s*(active|inactive)\s*$`)

func ParseStatus(out string) (Status, error) {
	s := Status{}
	sc := bufio.NewScanner(strings.NewReader(out))
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" {
			continue
		}
		m := reStatus.FindStringSubmatch(line)
		if len(m) == 0 {
			continue
		}
		switch strings.ToLower(m[1]) {
		case "active":
			s.Active = true
		case "inactive":
			s.Active = false
		}
		return s, nil
	}
	if err := sc.Err(); err != nil {
		return Status{}, err
	}
	return Status{}, fmt.Errorf("unable to parse ufw status")
}

