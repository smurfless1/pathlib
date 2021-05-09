//go:generate go run -mod=mod github.com/golang/mock/mockgen -package pathlib -destination=./mock_path.go -source=pathlib.go -build_flags=-mod=mod
package pathlib

import (
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"

	"github.com/pkg/errors"
)

// Path to a directory or file
type Path interface {
	Parts() []string
	Absolute() (Path, error)
	Cwd() (Path, error)
	Parent() (Path, error)
	Touch() error
	RmDir() error
	Unlink() error
	MkDir(mode os.FileMode, parents bool) (err error)
	Open() ([]byte, error)
	Chmod(mode os.FileMode) error
	JoinPath(elem ...string) Path
	Exists() bool
	IsAbsolute() bool
	IsFile() bool
	IsDir() bool
	ExpandUser() Path
	String() string
}

// PathImpl is the real implementation of interface Path over os/filepath and fs and so on
type PathImpl struct {
	Path string
}

// New Returns a new path.
func New(path string) PathImpl {
	return PathImpl{Path: path}
}

// FromParts Reconstitute a path string from a list/slice
func FromParts(value []string) PathImpl {
	return New(filepath.Join(value...))
}

// Absolute Returns an absolute representation of path.
func (p PathImpl) Absolute() (Path, error) {
	pth, err := filepath.Abs(p.Path)
	if err != nil {
		return nil, errors.Wrap(err, "get absolute failed")
	}
	newP := New(pth)
	if !newP.Exists() {
		parts := p.Parts()
		parts = append([]string{"/"}, parts...)
		newP = FromParts(parts)
		if !newP.Exists() {
			return nil, errors.New("unable to resolve path to file")
		}
	}

	return newP, nil
}

// Cwd Return a new path pointing to the current working directory.
func (p PathImpl) Cwd() (Path, error) {
	pth, err := os.Getwd()
	if err != nil {
		return nil, errors.Wrap(err, "get cwd failed")
	}
	newP := New(pth)
	return newP, nil
}

// Parent Return a new path for current path parent.
func (p PathImpl) Parent() (Path, error) {
	pth, err := p.Absolute()
	if err != nil {
		return nil, errors.Wrap(err, "get parent failed")
	}
	dir := filepath.Dir(pth.String())
	newP := New(dir)
	return newP, nil
}

// Touch Create creates the named file with mode 0666 (before umask), regardless of whether it exists.
func (p PathImpl) Touch() error {
	f, err := os.Create(p.Path)
	if err != nil {
		return err
	}
	return f.Close()
}

// Unlink Remove this file or link.
func (p PathImpl) Unlink() error {
	err := syscall.Unlink(p.Path)
	return err
}

// RmDir Remove this directory. The directory must be empty.
func (p PathImpl) RmDir() error {
	err := os.Remove(p.Path)
	return err
}

// MkDir Create a new directory at this given path.
func (p PathImpl) MkDir(mode os.FileMode, parents bool) (err error) {
	if parents {
		err = os.MkdirAll(p.Path, mode)
	} else {
		err = os.Mkdir(p.Path, mode)
	}
	return
}

// Open Reads the file named by filename and returns the contents.
func (p PathImpl) Open() ([]byte, error) {
	buf, err := ioutil.ReadFile(p.Path)
	if err != nil {
		return nil, err
	}
	return buf, nil
}

// Chmod changes the mode of the named file to mode.
func (p PathImpl) Chmod(mode os.FileMode) error {
	return os.Chmod(p.Path, mode)
}

// JoinPath Returns a new path, Combine current path with one or several arguments
func (p PathImpl) JoinPath(elem ...string) Path {
	temp := []string{p.Path}
	elem = append(temp, elem[0:]...)
	newP := New(path.Join(elem...))
	return newP
}

// Exists reports current path parent exists.
func (p PathImpl) Exists() bool {
	_, err := os.Stat(p.Path)
	return err == nil || os.IsExist(err)
}

// IsDir reports Whether this path is a directory.
func (p PathImpl) IsDir() bool {
	f, err := os.Stat(p.Path)
	if err != nil {
		return false
	}
	return f.IsDir()
}

// IsFile reports Whether this path is a regular file.
func (p PathImpl) IsFile() bool {
	f, e := os.Stat(p.Path)
	if e != nil {
		return false
	}
	return !f.IsDir()
}

// IsAbsolute reports whether the path is absolute.
func (p PathImpl) IsAbsolute() bool {
	return filepath.IsAbs(p.Path)
}

// from https://github.com/golang/go/issues/33393
// removeEmpty removes empty string elements from a slice
func removeEmpty(slice *[]string) {
	i := 0
	p := *slice
	for _, entry := range p {
		if strings.Trim(entry, " ") != "" {
			p[i] = entry
			i++
		}
	}
	*slice = p[0:i]
}

// Parts get the list of path components
func (p PathImpl) Parts() []string {
	parts := strings.Split(p.Path, string(os.PathSeparator))
	removeEmpty(&parts)
	return parts
}

// userHomeDir ...
func userHomeDir() string {
	if runtime.GOOS == "windows" {
		home := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
		if home == "" {
			home = os.Getenv("USERPROFILE")
		}
		return home
	} else if runtime.GOOS == "linux" {
		home := os.Getenv("XDG_CONFIG_HOME")
		if home != "" {
			return home
		}
	}
	return os.Getenv("HOME")
}

// normalizePath ...
func normalizePath(path string) string {
	// expand tilde
	if strings.HasPrefix(path, "~/") {
		path = filepath.Join(userHomeDir(), path[2:])
	} else if strings.HasPrefix(path, "~") {
		path = userHomeDir()
	}

	return path
}

// ExpandUser returns a copy of this path with ~ expanded
func (p PathImpl) ExpandUser() Path {
	return New(normalizePath(p.Path))
}

// String conversion
func (p PathImpl) String() string {
	return p.Path
}
