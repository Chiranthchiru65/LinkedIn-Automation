package models

// Lead represents a potential connection target
type Lead struct {
	Name       string
	ProfileURL string
	Headline   string // Useful for filtering (e.g., skip "Student")
	Location   string
}

// String returns a pretty print of the lead
func (l Lead) String() string {
	return l.Name + " (" + l.Headline + ")"
}