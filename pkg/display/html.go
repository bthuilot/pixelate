package display

import (
	"fmt"
	"html"
)

type HTMLAttributes interface {
	GenerateHTML() string
}

type HTMLLink struct {
	Text string
	Href string
}

func (l HTMLLink) GenerateHTML() string {
	return fmt.Sprintf("<a href=\"%s\">%s/a>", html.EscapeString(l.Href), html.EscapeString(l.Text))
}

type HTMLText struct {
	Content string
}

func (t HTMLText) GenerateHTML() string {
	return fmt.Sprintf("<p>%s</p>", html.EscapeString(t.Content))
}
