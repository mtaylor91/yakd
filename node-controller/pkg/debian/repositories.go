package debian

import (
	"os"
	"path"
	"text/template"
)

const AptSourcesTemplate = `
deb http://{{.Mirror}} {{.Release}} main contrib non-free
deb http://{{.Mirror}} {{.Release}}-updates main contrib non-free
deb http://security.debian.org/debian-security {{.Release}}-security main contrib non-free
`

// ConfigureRepositories configures the apt repositories for the target OS
func (c *BootstrapConfig) ConfigureRepositories() error {
	// Format the template
	tmpl, err := template.New("sources.list").Parse(AptSourcesTemplate)
	if err != nil {
		return nil
	}

	// Create the sources.list file
	f, err := os.Create(path.Join(c.Target, "etc", "apt", "sources.list"))
	if err != nil {
		return nil
	}

	// Execute the template
	if err := tmpl.Execute(f, c); err != nil {
		return nil
	}

	return nil
}
