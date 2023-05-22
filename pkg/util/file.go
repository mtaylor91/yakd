package util

import "os"

// WriteFile writes a file with the given contents and permissions 0644.
func WriteFile(path string, contents string) error {
	return WriteFileWithPermissions(path, contents, 0644)
}

// WriteFileWithPermissions writes a file with the given contents and permissions.
func WriteFileWithPermissions(path string, contents string, perm os.FileMode) error {
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, perm)
	if err != nil {
		return err
	}

	defer f.Close()

	_, err = f.WriteString(contents)
	if err != nil {
		return err
	}

	return nil
}
