package typed

// Options for generating coke
type Options struct {
	AppName    string
	ModulePath string
	TypeName   string
}

// Validate that options are usuable
func (opts *Options) Validate() error {
	return nil
}
