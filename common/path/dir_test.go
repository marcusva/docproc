package path

import (
	"github.com/marcusva/docproc/common/testing/assert"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestWritable(t *testing.T) {
	ok, err := Writable("invalid")
	assert.Err(t, err)

	ok, err = Writable(os.TempDir())
	assert.FailOnErr(t, err)
	assert.Equal(t, ok, true)
}

func TestCreateDir(t *testing.T) {
	dir, err := ioutil.TempDir(os.TempDir(), "test_create_dir")
	assert.FailOnErr(t, err)
	sub := filepath.Join(dir, "subdir")

	err = CreateDir(sub)
	if err == nil {
		// TODO: add test for creating on an existing file.
		assert.FailOnErr(t, os.Remove(sub))
	}
	assert.FailOnErr(t, os.Remove(dir))
	assert.FailOnErr(t, err)
}

func TestDirExists(t *testing.T) {
	ok, err := DirExists("/")
	assert.FailOnErr(t, err)
	assert.Equal(t, ok, true)

	ok, err = DirExists("invalid")
	assert.FailOnErr(t, err)
	assert.Equal(t, ok, false)
}
