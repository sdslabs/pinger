// +build dev

package ui

import (
	"io/ioutil"
	"path"
)

// TemplateContent is contains the HTML template for standalone agent page.
var TemplateContent string

func init() {
	content, err := ioutil.ReadFile(path.Join("pkg", "components", "agent", "ui", "page.gohtml"))
	if err != nil {
		panic(err)
	}
	TemplateContent = string(content)
}
