package appui

import "strings"

func buildWizardArgs(direction, iface, from, to, port, proto, comment string) []string {
	direction = strings.TrimSpace(direction)
	iface = strings.TrimSpace(iface)
	from = strings.TrimSpace(from)
	to = strings.TrimSpace(to)
	port = strings.TrimSpace(port)
	proto = strings.TrimSpace(proto)
	comment = strings.TrimSpace(comment)

	args := []string{}
	if direction != "" {
		args = append(args, direction)
	}
	if iface != "" {
		args = append(args, "on", iface)
	}
	if from != "" {
		args = append(args, "from", from)
	}
	if to != "" {
		args = append(args, "to", to)
	} else if from != "" || port != "" || proto != "" {
		args = append(args, "to", "any")
	}
	if port != "" {
		p, pproto := splitPortProto(port)
		args = append(args, "port", p)
		if proto == "" {
			proto = pproto
		}
	}
	if proto != "" {
		args = append(args, "proto", proto)
	}
	if comment != "" {
		args = append(args, "comment", comment)
	}
	return args
}

func splitPortProto(s string) (port string, proto string) {
	s = strings.TrimSpace(s)
	if s == "" {
		return "", ""
	}
	if strings.Contains(s, "/") {
		parts := strings.SplitN(s, "/", 2)
		return strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])
	}
	return s, ""
}

