package gentoo

import (
	"fmt"
	"path"
	"strings"

	"github.com/mtaylor91/yakd/pkg/system"
	"github.com/mtaylor91/yakd/pkg/util"
)

// acceptKeywords unmasks a package with the given keywords
func acceptKeywords(
	target, section, pkg string, priority int, keywords ...string,
) error {
	filename := fmt.Sprintf("%02d-%s", priority, pkg)
	return util.AppendFile(
		path.Join(target, "etc", "portage", "package.accept_keywords", filename),
		strings.Join(append(
			[]string{fmt.Sprintf("%s/%s", section, pkg)}, keywords...,
		), " "),
	)
}

// installPackages installs the given packages
func installPackages(sys system.System, pkgs ...string) error {
	for _, pkg := range pkgs {
		err := sys.RunCommand("emerge", "--usepkg", pkg)
		if err != nil {
			return err
		}
	}

	return nil
}
