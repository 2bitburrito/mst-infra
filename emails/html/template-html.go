package html

import (
	"bytes"
	"html/template"
	"io"
)

type GenericEmailData struct {
	HighlightWord  *string
	FirstName      *string
	MainMessage    *string
	CtaText        *string
	CtaLink        *string
	SecondaryLink  *string
	ClosingMessage *string
	PreferencesUrl *string
}

func TemplateEmail[T any](reader io.Reader, data T) (*bytes.Buffer, error) {
	html, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	tmpl, err := template.New("html").Parse(string(html))
	if err != nil {
		return nil, err
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return nil, err
	}
	return &buf, nil
}
