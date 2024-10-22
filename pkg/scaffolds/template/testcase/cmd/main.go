package cmd

import (
	"github.com/huanghj78/jepsenFuzz/pkg/scaffolds/file"
)

// Cmd uses for cmd/main.go
type Cmd struct {
	file.TemplateMixin
	CaseName string
}

// GetIfExistsAction ...
func (c *Cmd) GetIfExistsAction() file.IfExistsAction {
	return file.IfExistsActionError
}

// Validate ...
func (c *Cmd) Validate() error {
	return c.TemplateMixin.Validate()
}

// SetTemplateDefaults ...
func (c *Cmd) SetTemplateDefaults() error {
	c.TemplateBody = cmdTemplate
	return nil
}

const cmdTemplate = `
package main

import (
	"context"
	"flag"

	"github.com/huanghj78/jepsenFuzz/cmd/util"

	testcase "github.com/huanghj78/jepsenFuzz/testcase/{{ .CaseName }}"
)

func main() {
	flag.Parse()
	suit := util.Suit{}
	suit.Run(context.Background())
}
`
