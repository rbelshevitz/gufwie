package ufw

type Status struct {
	Active bool
}

type Rule struct {
	Number int
	To     string
	Action string
	From   string
	V6     bool
	Raw    string
}

type NumberedStatus struct {
	Status Status
	Rules  []Rule
}

