package config

import (
	"github.com/marcusva/docproc/common/testing/assert"
	"strings"
	"testing"
)

var (
	_cfg = `
[log]
# log to a specific file instead of stdout
# file=<path/to/the/file>
# Emergency,Alert,Critical,Error,Warning,Notice,Info,Debug
level = Debug

[section]
k = v
intval = 1234
boolval = t

[arrays]
ar = this,is,    an, array     , with	some values, yay
	`
)

func TestLoad(t *testing.T) {
	_, err := Load(strings.NewReader(_cfg), NoValidate)
	assert.FailOnErr(t, err)

	_, err = Load(strings.NewReader(_cfg), nil)
	assert.FailOnErr(t, err)

	_broken := `
	[log]
	log to a specific file instead of stdout
	# file=<path/to/the/file>
	# Emergency,Alert,Critical,Error,Warning,Notice,Info,Debug
	level = Debug`
	_, err = Load(strings.NewReader(_broken), NoValidate)
	assert.Err(t, err)

	_broken2 := `
	nosecval = 1234
	[log]
	level = Debug`
	_, err = Load(strings.NewReader(_broken2), NoValidate)
	assert.Err(t, err)

	_broken3 := `
	[log]
	level = Debug
	[log]
	level = Info`
	_, err = Load(strings.NewReader(_broken3), NoValidate)
	assert.Err(t, err)

	_broken4 := `
	[ ]
	some = value`
	_, err = Load(strings.NewReader(_broken4), NoValidate)
	assert.Err(t, err)
}

func TestGet(t *testing.T) {
	cfg, err := Load(strings.NewReader(_cfg), NoValidate)
	assert.FailOnErr(t, err)

	val, err := cfg.Get("log", "level")
	assert.FailOnErr(t, err)
	assert.Equal(t, val, "Debug")

	val, err = cfg.Get("section", "k")
	assert.FailOnErr(t, err)
	assert.Equal(t, val, "v")

	val, err = cfg.Get("section", "intval")
	assert.FailOnErr(t, err)
	assert.Equal(t, val, "1234")

	_, err = cfg.Get("section", "invalidkey")
	assert.Err(t, err)

	_, err = cfg.Get("invalidsection", "invalidkey")
	assert.Err(t, err)

}

func TestGetOrPanic(t *testing.T) {
	cfg, err := Load(strings.NewReader(_cfg), NoValidate)
	assert.FailOnErr(t, err)

	v1 := cfg.GetOrPanic("section", "k")
	assert.Equal(t, v1, "v")

	assert.Panics(t, func() { cfg.GetOrPanic("section", "invalidkey") })
}

func TestGetDefault(t *testing.T) {
	cfg, err := Load(strings.NewReader(_cfg), NoValidate)
	assert.FailOnErr(t, err)

	val := cfg.GetDefault("log", "level", "Info")
	assert.Equal(t, val, "Debug")

	val2 := cfg.GetDefault("log", "invalid", "Info")
	assert.Equal(t, val2, "Info")
}

func TestInt(t *testing.T) {
	cfg, err := Load(strings.NewReader(_cfg), NoValidate)
	assert.FailOnErr(t, err)

	intval, err := cfg.Int("section", "intval")
	assert.FailOnErr(t, err)
	assert.Equal(t, intval, 1234)

	_, err = cfg.Int("section", "k")
	assert.Err(t, err)

	_, err = cfg.Int("section", "invalidkey")
	assert.Err(t, err)
}
func TestBool(t *testing.T) {
	cfg, err := Load(strings.NewReader(_cfg), NoValidate)
	assert.FailOnErr(t, err)

	bv, err := cfg.Bool("section", "boolval")
	assert.FailOnErr(t, err)
	assert.Equal(t, bv, true)

	_, err = cfg.Bool("section", "k")
	assert.Err(t, err)

	_, err = cfg.Bool("section", "invalidkey")
	assert.Err(t, err)

}

func TestHasSection(t *testing.T) {
	cfg, err := Load(strings.NewReader(_cfg), NoValidate)
	assert.FailOnErr(t, err)

	assert.Equal(t, cfg.HasSection("section"), true)
	assert.Equal(t, cfg.HasSection("invalidsection"), false)
}

func TestAllFor(t *testing.T) {
	cfg, err := Load(strings.NewReader(_cfg), NoValidate)
	assert.FailOnErr(t, err)

	sec, err := cfg.AllFor("section")
	assert.FailOnErr(t, err)
	assert.Equal(t, sec["k"], "v")
	assert.Equal(t, sec["intval"], "1234")
	assert.Equal(t, sec["boolval"], "t")

	_, err = cfg.AllFor("invalidsection")
	assert.Err(t, err)
}

func TestArray(t *testing.T) {
	cfg, err := Load(strings.NewReader(_cfg), NoValidate)
	assert.FailOnErr(t, err)

	ar, err := cfg.Array("arrays", "ar")
	assert.Equal(t, ar[0], "this")
	assert.Equal(t, ar[1], "is")
	assert.Equal(t, ar[2], "an")
	assert.Equal(t, ar[3], "array")
	assert.Equal(t, ar[4], "with	some values")
	assert.Equal(t, ar[5], "yay")

	_, err = cfg.Array("arrays", "invalid")
	assert.Err(t, err)
}

func TestLoadFile(t *testing.T) {
	cfg, err := LoadFile("test/test.ini", NoValidate)
	assert.FailOnErr(t, err)

	val, err := cfg.Get("log", "level")
	assert.FailOnErr(t, err)
	assert.Equal(t, val, "Debug")

	_, err = LoadFile("invalid.ini", NoValidate)
	assert.Err(t, err)
}
