//go:generate go run -mod=mod github.com/golang/mock/mockgen -package pathlib -destination=./mock_path.go -source=pathlib.go -build_flags=-mod=mod
package pathlib

import (
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/pkg/errors"
)

// Path to a directory or file
type Path interface {
	Parts() []string
	Absolute() (*PathImpl, error)
	Cwd() (*Path, error)
	Parent() (*Path, error)
	Touch() error
	RmDir() error
	Unlink() error
	MkDir(mode os.FileMode, parents bool) (err error)
	Open() ([]byte, error)
	Chmod(mode os.FileMode) error
	JoinPath(elem ...string) *Path
	Exists() bool
	IsAbs() bool
	IsFile() bool
	IsDir() bool
}

// PathImpl
type PathImpl struct {
	Path string
}

// NewPathImpl Returns a new path.
func NewPathImpl(path string) *PathImpl {
	p := new(PathImpl)
	p.Path = path
	return p
}

// FromParts Reconstitute a path string from a list/slice
func FromParts(value []string) *PathImpl {
	return NewPathImpl(filepath.Join(value...))
}

// Absolute Returns an absolute representation of path.
func (p *PathImpl) Absolute() (*PathImpl, error) {
	pth, err := filepath.Abs(p.Path)
	if err != nil {
		return nil, errors.Wrap(err, "get absolute failed")
	}
	newP := NewPathImpl(pth)
	return newP, nil
}

// Cwd Return a new path pointing to the current working directory.
func (p *PathImpl) Cwd() (*PathImpl, error) {
	pth, err := os.Getwd()
	if err != nil {
		return nil, errors.Wrap(err, "get cwd failed")
	}
	newP := NewPathImpl(pth)
	return newP, nil
}

// Parent Return a new path for current path parent.
func (p *PathImpl) Parent() (*PathImpl, error) {
	pth, err := p.Absolute()
	if err != nil {
		return nil, errors.Wrap(err, "get parent failed")
	}
	dir := filepath.Dir(pth.Path)
	newP := NewPathImpl(dir)
	return newP, nil
}

// Touch Create creates the named file with mode 0666 (before umask), regardless of whether it exists.
func (p *PathImpl) Touch() error {
	f, err := os.Create(p.Path)
	f.Close()
	return err
}

// Unlink Remove this file or link.
func (p *PathImpl) Unlink() error {
	err := syscall.Unlink(p.Path)
	return err
}

// RmDir Remove this directory. The directory must be empty.
func (p *PathImpl) RmDir() error {
	err := os.Remove(p.Path)
	return err
}

// MkDir Create a new directory at this given path.
func (p *PathImpl) MkDir(mode os.FileMode, parents bool) (err error) {
	if parents {
		err = os.MkdirAll(p.Path, mode)
	} else {
		err = os.Mkdir(p.Path, mode)
	}
	return
}

// Open Reads the file named by filename and returns the contents.
func (p *PathImpl) Open() ([]byte, error) {
	buf, err := ioutil.ReadFile(p.Path)
	if err != nil {
		return nil, err
	}
	return buf, nil
}

// Chmod changes the mode of the named file to mode.
func (p *PathImpl) Chmod(mode os.FileMode) error {
	return os.Chmod(p.Path, mode)
}

// JoinPath Returns a new path, Combine current path with one or several arguments
func (p *PathImpl) JoinPath(elem ...string) *PathImpl {
	temp := []string{p.Path}
	elem = append(temp, elem[0:]...)
	newP := NewPathImpl(path.Join(elem...))
	return newP
}

// Exists reports current path parent exists.
func (p *PathImpl) Exists() bool {
	_, err := os.Stat(p.Path)
	return err == nil || os.IsExist(err)
}

// IsDir reports Whether this path is a directory.
func (p *PathImpl) IsDir() bool {
	f, err := os.Stat(p.Path)
	if err != nil {
		return false
	}
	return f.IsDir()
}

// IsFile reports Whether this path is a regular file.
func (p *PathImpl) IsFile() bool {
	f, e := os.Stat(p.Path)
	if e != nil {
		return false
	}
	return !f.IsDir()
}

// IsAbs reports whether the path is absolute.
func (p *PathImpl) IsAbs() bool {
	return filepath.IsAbs(p.Path)
}

// from https://github.com/golang/go/issues/33393
func RemoveEmpty(slice *[]string) {
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
func (p *PathImpl) Parts() []string {
	parts := strings.Split(p.Path, string(os.PathSeparator))
	RemoveEmpty(&parts)
	return parts
}
