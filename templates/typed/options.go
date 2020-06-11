package typed

// Options ...
type Options struct {
	AppName    string
	ModulePath string
	TypeName   string
	Fields     []string
}

// Validate that options are usuable
func (opts *Options) Validate() error {
	return nil
}
