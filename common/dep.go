package common

type Dependency struct {
	// Location is the url/uri of the remote location where the dependency can
	// be found at. Required.
	Location string `yaml:"loc"`

	// Path is the path the dependency should be installed to. This should
	// correspond to whatever the dependency expects to be imported as. For
	// example: "code.google.com/p/protobuf". Default: Value of Location field
	Path string `yaml:"path"`

	// Type is what kind of project the dependency should be fetched as.
	// Available options are: goget, git, hg. Default: goget.
	Type string `yaml:"type"`

	// Reference is only valueable for version-control Types (e.g. git). It can
	// be any valid reference in that version control system (branch name, tag,
	// commit hash). Default depends on the Type, git is "master", hg is "tip".
	Reference string `yaml:"reference"`
}
