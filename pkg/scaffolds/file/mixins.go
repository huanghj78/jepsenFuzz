package file

import "fmt"

// TemplateMixin is the mixin that should be embedded in Template builders
type TemplateMixin struct {
	// Path is which the template should locate
	Path string
	// TemplateBody is the template body to execute
	TemplateBody string
}

// GetPath is the Path field getter
func (t *TemplateMixin) GetPath() string {
	return t.Path
}

// GetBody is the TemplateBody field getter
func (t *TemplateMixin) GetBody() string {
	return t.TemplateBody
}

// Validate ...
func (t *TemplateMixin) Validate() error {
	if t.GetPath() == "" {
		return fmt.Errorf("template's Path cannot be empty")
	}
	return nil
}

// InserterMixin is the mixin that should be embedded in Inserter builders
type InserterMixin struct {
	// Path is which the template should locate
	Path string
}

// GetPath is the Path field getter
func (t *InserterMixin) GetPath() string {
	return t.Path
}

// Validate ...
func (t *InserterMixin) Validate() error {
	if t.GetPath() == "" {
		return fmt.Errorf("inserter's Path cannot be empty")
	}
	return nil
}
