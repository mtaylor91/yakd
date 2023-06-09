package util

import (
	"io"
	"os"
)

// AppendFile appends the given contents to a file.
func AppendFile(path string, contents string) error {
	return AppendFileWithPermissions(path, contents, 0644)
}

func AppendFileWithPermissions(path string, contents string, perm os.FileMode) error {
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, perm)
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

// CopyFile copies a file from src to dst.
func CopyFile(src string, dst string) error {
	return CopyFileWithPermissions(src, dst, 0644)
}

// CopyFileWithPermissions copies a file from src to dst with the given permissions.
func CopyFileWithPermissions(src string, dst string, perm os.FileMode) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}

	defer srcFile.Close()

	dstFile, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, perm)
	if err != nil {
		return err
	}

	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return err
	}

	return nil
}
