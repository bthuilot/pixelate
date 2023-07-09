package display

import (
	"fmt"
	"html"
	"html/template"
)

// HTMLAttribute represents an HTML element to be rendered
// on the dashboard
type HTMLAttribute interface {
	GenerateHTML() template.HTML
}

// HTMLLink represents an `<a>` link in HTML
type HTMLLink struct {
	// Text is the text the link should display
	Text string
	// Href is the HRef of the link
	Href string
}

func (l HTMLLink) GenerateHTML() template.HTML {
	return template.HTML(fmt.Sprintf("<a href=\"%s\">%s</a>", html.EscapeString(l.Href), html.EscapeString(l.Text)))
}

// HTMLText renders plain HTML `<p>` to the screen
type HTMLText struct {
	Content string
}

func (t HTMLText) GenerateHTML() template.HTML {
	return template.HTML(fmt.Sprintf("<p>%s</p>", html.EscapeString(t.Content)))
}
