package ui

import _ "embed" // To use the embed directive.

// TemplateContent is the content of page template.
//go:embed page.gohtml
var TemplateContent string
