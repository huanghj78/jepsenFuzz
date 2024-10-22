package file

// Builder defines the basic methods that any file builder must implement
type Builder interface {
	// GetPath returns the path to the file location
	GetPath() string
	Validate() error
	GetIfExistsAction() IfExistsAction
}

// Template is the file builder based on a template file
type Template interface {
	Builder
	GetBody() string
	SetTemplateDefaults() error
}

// Inserter is a file builder that inserts code fragments in marked positions
type Inserter interface {
	Builder
	// GetCodeFragments returns a map that binds markers to code fragments
	GetCodeFragments() map[Marker]CodeFragment
}
