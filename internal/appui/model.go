package appui

import (
	"strings"

	"github.com/rbelshevitz/gufwie/internal/ufw"
)

type model struct {
	status ufw.Status
	rules  []ufw.Rule

	filter string

	// indexMap maps visible table row (1-based, excluding header) to rules index.
	indexMap []int
}

func (m *model) setStatus(ns ufw.NumberedStatus) {
	m.status = ns.Status
	m.rules = ns.Rules
}

func (m *model) applyFilter() {
	f := strings.TrimSpace(strings.ToLower(m.filter))
	m.indexMap = m.indexMap[:0]
	for idx, r := range m.rules {
		if f == "" {
			m.indexMap = append(m.indexMap, idx)
			continue
		}
		hay := strings.ToLower(r.Raw + " " + r.To + " " + r.From + " " + r.Action)
		if strings.Contains(hay, f) {
			m.indexMap = append(m.indexMap, idx)
		}
	}
}

func (m *model) ruleForRow(row int) (ufw.Rule, bool) {
	if row <= 0 {
		return ufw.Rule{}, false
	}
	i := row - 1
	if i < 0 || i >= len(m.indexMap) {
		return ufw.Rule{}, false
	}
	ridx := m.indexMap[i]
	if ridx < 0 || ridx >= len(m.rules) {
		return ufw.Rule{}, false
	}
	return m.rules[ridx], true
}

