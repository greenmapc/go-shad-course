//go:build !solution

package ciletters

import (
	_ "embed"
	"strings"
	"text/template"
)

//go:embed resources/template.txt
var notificationTemplate string

func MakeLetter(n *Notification) (string, error) {
	tmpl, err := template.New("Notification Template").Funcs(template.FuncMap{
		"truncate": func(s string, length int) string {
			if len(s) > length {
				return s[:length]
			}
			return s
		},
		"cut": func(s string, rows int) []string {
			r := strings.Split(s, "\n")
			beginning := 0
			if len(r) > rows {
				beginning = len(r) - rows
			}

			return r[beginning:]
		},
	}).Parse(notificationTemplate)

	if err != nil {
		return "", err
	}

	var outputBuilder strings.Builder

	tmpl.Execute(&outputBuilder, &n)

	return outputBuilder.String(), nil
}
