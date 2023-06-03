package debian

import (
	"context"
	"os"
	"path"
	"text/template"

	"github.com/mtaylor91/yakd/pkg/util/http"
)

const AptSourcesTemplate = `
deb {{.Mirror}} {{.Suite}} main contrib non-free
deb {{.Mirror}} {{.Suite}}-updates main contrib non-free
deb http://security.debian.org/debian-security {{.Suite}}-security main contrib non-free
`

const SignedAptSourcesTemplate = `
deb [signed-by={{.Keyring}}] {{.Url}} {{.Components}}
`

const crioVersion = "1.24"
const debianVersion = "Debian_11"
const libcontainers = "https://download.opensuse.org/repositories/devel:/kubic:/libcontainers:/stable"

type AptSource struct {
	Keyring    string
	Url        string
	Components string
}

// configureRepositories configures the apt repositories for the target OS
func (b *BootstrapInstaller) configureRepositories(ctx context.Context) error {
	keyrings := "usr/share/keyrings"

	// Setup sources.list
	if err := writeTemplateToFile("sources.list", AptSourcesTemplate,
		path.Join(b.Target, "etc", "apt", "sources.list"), b); err != nil {
		return err
	}

	// Setup kubernetes repository
	keyring := path.Join(keyrings, "kubernetes-archive-keyring.gpg")
	keyringDownload := http.NewDownload(
		"https://packages.cloud.google.com/apt/doc/apt-key.gpg",
		path.Join(b.Target, keyring))
	// Download keyring
	if err := keyringDownload.DownloadAndDearmorGPG(ctx); err != nil {
		return err
	}
	// Write template to apt source file
	if err := writeTemplateToFile("kubernetes.list", SignedAptSourcesTemplate,
		path.Join(b.Target, "etc", "apt", "sources.list.d", "kubernetes.list"),
		AptSource{
			path.Join("/", keyring),
			"https://apt.kubernetes.io/",
			"kubernetes-xenial main",
		},
	); err != nil {
		return err
	}

	// Setup libcontainers repository
	keyring = path.Join(keyrings, "libcontainers-archive-keyring.gpg")
	releaseKeyringUrl := libcontainersUrl(debianVersion) + "Release.key"
	keyringDownload = http.NewDownload(
		releaseKeyringUrl, path.Join(b.Target, keyring))
	// Download keyring
	if err := keyringDownload.DownloadAndDearmorGPG(ctx); err != nil {
		return err
	}
	// Write template to apt source file
	if err := writeTemplateToFile("libcontainers.list", SignedAptSourcesTemplate,
		path.Join(
			b.Target, "etc", "apt", "sources.list.d", "libcontainers.list",
		),
		AptSource{
			path.Join("/", keyring),
			libcontainersUrl(debianVersion),
			"/",
		},
	); err != nil {
		return err
	}

	// Setup libcontainers crio repository
	keyring = path.Join(keyrings, "libcontainers-crio-archive-keyring.gpg")
	releaseKeyringUrl = crioArchiveUrl(crioVersion, debianVersion) + "Release.key"
	keyringDownload = http.NewDownload(
		releaseKeyringUrl, path.Join(b.Target, keyring))
	// Download keyring
	if err := keyringDownload.DownloadAndDearmorGPG(ctx); err != nil {
		return err
	}
	// Write template to apt source file
	if err := writeTemplateToFile(
		"libcontainers-crio.list",
		SignedAptSourcesTemplate,
		path.Join(
			b.Target, "etc", "apt", "sources.list.d",
			"libcontainers-crio.list",
		),
		AptSource{
			path.Join("/", keyring),
			crioArchiveUrl(crioVersion, debianVersion),
			"/",
		},
	); err != nil {
		return err
	}

	return nil
}

func writeTemplateToFile(name, src, dest string, data interface{}) error {
	tmpl, err := template.New(name).Parse(src)

	f, err := os.Create(dest)
	if err != nil {
		return err
	}

	defer f.Close()

	if err := tmpl.Execute(f, data); err != nil {
		return err
	}

	return nil
}

func libcontainersUrl(os string) string {
	return libcontainers + "/" + os + "/"
}

func crioArchiveUrl(version, os string) string {
	return libcontainers + ":/cri-o:/" + version + "/" + os + "/"
}
