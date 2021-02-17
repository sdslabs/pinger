// +build !dev

package ui

// Enables the embed directive
import _ "embed"

// TemplateContent is contains the HTML template for standalone agent page.
//go:embed page.gohtml
var TemplateContent string
