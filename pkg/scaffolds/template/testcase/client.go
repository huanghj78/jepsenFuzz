package testcase

import "github.com/huanghj78/jepsenFuzz/pkg/scaffolds/file"

// Client uses for client.go
type Client struct {
	file.TemplateMixin
	CaseName string
}

// GetIfExistsAction ...
func (c *Client) GetIfExistsAction() file.IfExistsAction {
	return file.IfExistsActionError
}

// Validate ...
func (c *Client) Validate() error {
	return c.TemplateMixin.Validate()
}

// SetTemplateDefaults ...
func (c *Client) SetTemplateDefaults() error {
	c.TemplateBody = clientTemplate
	return nil
}

const clientTemplate = `
package testcase


`
