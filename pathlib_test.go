package pathlib

import (
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"runtime"
	"strings"
	"testing"
)

func TestPathInOut(t *testing.T) {
	// system testing a bit
	pp := PathImpl{Value: "/tmp"}
	assert.True(t, pp.Exists())
	assert.True(t, pp.IsDir())
	assert.True(t, pp.IsAbsolute())
	assert.Equal(t, "/tmp", pp.String())
	parts := pp.Parts()
	assert.Equal(t, []string{"tmp"}, parts)

	// relative when extracted
	made := fromParts(parts)
	assert.Equal(t, "tmp", made.String())
	absolute, err := made.Absolute()
	if err != nil {
		log.Println(err)
	}
	assert.Equal(t, "/tmp", absolute.String())

	made = fromParts([]string{"foo"})
	absolute, err = made.Absolute()
	assert.EqualError(t, err, "unable to resolve path to file")

	made = fromParts([]string{"pathlib/pathlib.go"})
	absolute, err = made.Absolute()
	assert.True(t, strings.Contains(made.String(), "pathlib/pathlib.go"))

	// making it absolute
	parts = append([]string{"/"}, parts...)
	made = fromParts(parts)
	assert.Equal(t, "/tmp", made.String())
}

func TestExistsIsMockable(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()
	MockInterface := NewMockPath(controller)
	MockInterface.EXPECT().Exists().Return(false)

	// now this can be injected into something else
	result := MockInterface.Exists()
	assert.Equal(t, false, result, "a should be removed")
}

func TestExpandUser(t *testing.T) {
	// ~/something can be expanded correctly
	pp := New("~/tmp")
	assert.Equal(t, "~/tmp", pp.String())

	expanded := pp.ExpandUser()

	if runtime.GOOS == "windows" {
		assert.True(t, strings.HasPrefix(expanded.String(), "C:"))
	} else if runtime.GOOS == "darwin" {
		assert.True(t, strings.HasPrefix(expanded.String(), "/Users"))
	} else if runtime.GOOS == "linux" {
		assert.True(t, strings.HasPrefix(expanded.String(), "/"))
		assert.True(t, strings.Contains(expanded.String(), "home")) // probably?
	}
	assert.True(t, strings.HasSuffix(expanded.String(), "tmp"))
	assert.True(t, strings.Contains(expanded.String(), os.Getenv("USER")))
	assert.False(t, strings.HasSuffix(expanded.String(), "/"))
}

func TestExpandJustUser(t *testing.T) {
	// Just the tilde expands correctly
	pp := New("~")
	assert.Equal(t, "~", pp.String())

	expanded := pp.ExpandUser()

	if runtime.GOOS == "windows" {
		assert.True(t, strings.HasPrefix(expanded.String(), "C:"))
	} else if runtime.GOOS == "darwin" {
		assert.True(t, strings.HasPrefix(expanded.String(), "/Users"))
	} else if runtime.GOOS == "linux" {
		assert.True(t, strings.HasPrefix(expanded.String(), "/"))
		assert.True(t, strings.Contains(expanded.String(), "home")) // probably?
	}
	assert.True(t, strings.Contains(expanded.String(), os.Getenv("USER")))
	assert.False(t, strings.HasSuffix(expanded.String(), "/"))
}

func TestReturnedCopyTypes(t *testing.T) {
	// if you return a copy (absolute, expanduser, etc.) I'm having trouble
	// thinking out the correct return types. still learning.
	logs := New("/tmp")
	joined := logs.JoinPath("foo")
	if !joined.Exists() {
		logs = joined.JoinPath("bar")
	}
}
