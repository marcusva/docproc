package path_test

import (
	"github.com/marcusva/docproc/common/path"
	"github.com/marcusva/docproc/common/testing/assert"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestWritable(t *testing.T) {
	ok, err := path.Writable("invalid")
	assert.Err(t, err)

	ok, err = path.Writable(os.TempDir())
	assert.FailOnErr(t, err)
	assert.Equal(t, ok, true)
}

func TestCreateDir(t *testing.T) {
	dir, err := ioutil.TempDir(os.TempDir(), "test_create_dir")
	assert.FailOnErr(t, err)
	sub := filepath.Join(dir, "subdir")

	err = path.CreateDir(sub)
	if err == nil {
		// TODO: add test for creating on an existing file.
		assert.FailOnErr(t, os.Remove(sub))
	}
	assert.FailOnErr(t, os.Remove(dir))
	assert.FailOnErr(t, err)
}

func TestDirExists(t *testing.T) {
	ok, err := path.DirExists("/")
	assert.FailOnErr(t, err)
	assert.Equal(t, ok, true)

	ok, err = path.DirExists("invalid")
	assert.FailOnErr(t, err)
	assert.Equal(t, ok, false)
}
