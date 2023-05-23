package util

// UnpackTarball unpacks a tarball to the specified target
func UnpackTarball(source, target string) error {
	// Unpack via tar
	return RunCmd("tar", "-xf", source, "-C", target)
}
