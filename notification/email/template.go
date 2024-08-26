package email

import (
	"bytes"
	"html/template"
)

func buildHtml(templateName string, data map[string]any) (string, error) {
	templatePath := "./templates/" + templateName + ".html.gohtml"
	t, err := template.New("email-verification-html").ParseFiles(templatePath)
	if err != nil {
		return "", err
	}

	var b bytes.Buffer
	err = t.ExecuteTemplate(&b, "body", data)
	if err != nil {
		return "", err
	}

	return b.String(), nil
}
