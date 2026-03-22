package ufw

import (
	"bufio"
	"fmt"
	"regexp"
	"strings"
)

var reRule = regexp.MustCompile(`^\[\s*(\d+)\]\s+(.+?)\s{2,}(.+?)\s{2,}(.+?)\s*$`)

func ParseNumberedStatus(out string) (NumberedStatus, error) {
	s, err := ParseStatus(out)
	if err != nil {
		return NumberedStatus{}, err
	}
	ns := NumberedStatus{Status: s}

	sc := bufio.NewScanner(strings.NewReader(out))
	for sc.Scan() {
		line := strings.TrimRight(sc.Text(), "\r\n")
		m := reRule.FindStringSubmatch(line)
		if len(m) == 0 {
			continue
		}
		num, convErr := atoi(m[1])
		if convErr != nil {
			continue
		}
		to := strings.TrimSpace(m[2])
		action := strings.TrimSpace(m[3])
		from := strings.TrimSpace(m[4])

		v6 := false
		if strings.HasSuffix(to, "(v6)") {
			v6 = true
			to = strings.TrimSpace(strings.TrimSuffix(to, "(v6)"))
		}
		if strings.HasSuffix(from, "(v6)") {
			v6 = true
			from = strings.TrimSpace(strings.TrimSuffix(from, "(v6)"))
		}
		ns.Rules = append(ns.Rules, Rule{
			Number: num,
			To:     to,
			Action: action,
			From:   from,
			V6:     v6,
			Raw:    strings.TrimSpace(line),
		})
	}
	if err := sc.Err(); err != nil {
		return NumberedStatus{}, err
	}
	return ns, nil
}

func atoi(s string) (int, error) {
	n := 0
	for _, r := range s {
		if r < '0' || r > '9' {
			return 0, fmt.Errorf("not int")
		}
		n = n*10 + int(r-'0')
	}
	return n, nil
}
