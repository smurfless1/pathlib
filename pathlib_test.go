package pathlib

import (
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPathInOut(t *testing.T) {
	// system testing a bit
	pp := PathImpl{Path: "/tmp"}
	assert.True(t, pp.Exists())
	assert.True(t, pp.IsDir())
	assert.True(t, pp.IsAbsolute())
	assert.Equal(t, "/tmp", pp.Path)
	parts := pp.Parts()
	assert.Equal(t, []string{"tmp"}, parts)

	// relative when extracted
	made := FromParts(parts)
	assert.Equal(t, "tmp", made.Path)
	absolute, err := made.Absolute(); if err != nil { log.Println(err)}
	assert.Equal(t, "/tmp", absolute.Path)

	made = FromParts([]string{"foo"})
	absolute, err = made.Absolute()
	assert.EqualError(t, err, "unable to resolve path to file")

	made = FromParts([]string{"pathlib/pathlib.go"})
	absolute, err = made.Absolute()
	assert.True(t, strings.Contains(made.Path,"pathlib/pathlib.go"))

	// making it absolute
	parts = append([]string{"/"}, parts...)
	made = FromParts(parts)
	assert.Equal(t, "/tmp", made.Path)
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
