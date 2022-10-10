package agents

import (
	"fmt"
	"html"
)

type Attribute interface {
	GetHTML() string
}

// Button
type Button struct {
	Name string
	Link string
}

func (b Button) GetHTML() string {
	return fmt.Sprintf("<a href='%s'>%s</a>", b.Link, b.Name)
}

// Text
type Text struct {
	Content string
}

func (t Text) GetHTML() string {
	return html.EscapeString(t.Content)
}
