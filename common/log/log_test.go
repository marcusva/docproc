package log

import (
	"bytes"
	"fmt"
	"github.com/marcusva/docproc/common/testing/assert"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

func TestPackage(t *testing.T) {
	if logger == nil {
		t.Error("_log is nil, although a package initialization was done")
	}

	// None of those should cause a panic
	Alert("test")
	Alertf("test")
	Critical("test")
	Criticalf("test")
	Debug("test")
	Debugf("test")
	Info("test")
	Infof("test")
	Notice("test")
	Noticef("test")
	Error("test")
	Errorf("test")
	Warning("test")
	Warningf("test")
	Emergency("test")
	Emergencyf("test")
}

func TestLogger(t *testing.T) {
	logger := Logger()
	assert.NotNil(t, logger)

	var buf bytes.Buffer
	Init(&buf, LevelDebug, true)
	logger2 := Logger()
	assert.NotEqual(t, logger, logger2)
}

func TestInitFile(t *testing.T) {
	fp, err := ioutil.TempFile(os.TempDir(), "docproc-logtest")
	assert.NoErr(t, err)
	fname := fp.Name()
	fp.Close()

	err = InitFile(fname, LevelDebug, false)
	assert.NoErr(t, err)

	Init(os.Stdout, LevelDebug, false)
	assert.NoErr(t, os.Remove(fname))

	err = InitFile("", LevelDebug, false)
	assert.Err(t, err)

}

func TestGetLogLevel(t *testing.T) {
	levelsInt := []string{"0", "1", "2", "3", "4", "5", "6", "7"}
	levelsTxt := []string{
		"Emergency", "Alert", "Critical", "Error", "Warning", "Notice", "Info",
		"Debug",
	}

	for idx, v := range levelsInt {
		if v1, err := GetLogLevel(v); err != nil {
			t.Error(err)
		} else {
			if v2, err := GetLogLevel(levelsTxt[idx]); err != nil {
				t.Error(err)
			} else {
				if v1 != v2 {
					t.Errorf("Log level mismatch: '%s' - '%s'",
						v, levelsTxt[idx])
				}
			}
		}
	}

	levelsInvalid := []string{"", "10", "SomeText"}
	for _, v := range levelsInvalid {
		if v1, err := GetLogLevel(""); err == nil || v1 != -1 {
			t.Errorf("invalid level '%s' was accepted", v)
		}
	}
}

func TestLog(t *testing.T) {
	callbacks := map[string]func(...interface{}){
		"DEBUG":     Debug,
		"INFO":      Info,
		"NOTICE":    Notice,
		"WARNING":   Warning,
		"ERROR":     Error,
		"CRITICAL":  Critical,
		"ALERT":     Alert,
		"EMERGENCY": Emergency,
	}

	var buf bytes.Buffer
	Init(&buf, LevelDebug, true)

	for prefix, cb := range callbacks {
		cb("Test")
		result := string(buf.Bytes())
		assert.FailIfNot(t, strings.Contains(result, prefix),
			"'%s' not found in %s", prefix, result)
		assert.FailIfNot(t, strings.Contains(result, "Test"))
		buf.Reset()
	}
}

func TestLogf(t *testing.T) {
	callbacks := map[string]func(f string, args ...interface{}){
		"DEBUG":     Debugf,
		"INFO":      Infof,
		"NOTICE":    Noticef,
		"WARNING":   Warningf,
		"ERROR":     Errorf,
		"CRITICAL":  Criticalf,
		"ALERT":     Alertf,
		"EMERGENCY": Emergencyf,
	}

	var buf bytes.Buffer
	Init(&buf, LevelDebug, true)

	fmtstring := "Formatted result: '%s'"
	for prefix, cb := range callbacks {
		fmtresult := fmt.Sprintf(fmtstring, "TestLogf")
		cb(fmtstring, "TestLogf")
		result := string(buf.Bytes())
		assert.FailIfNot(t, strings.Contains(result, prefix),
			"'%s' not found in %s", prefix, result)
		assert.FailIfNot(t, strings.Contains(result, fmtresult))
		buf.Reset()
	}
}

func TestLogLevel(t *testing.T) {
	levels := []Level{
		LevelDebug, LevelInfo, LevelNotice, LevelWarning,
		LevelError, LevelAlert, LevelCritical, LevelEmergency,
	}
	for _, level := range levels {
		var buf bytes.Buffer
		Init(&buf, level, true)
		assert.Equal(t, level, CurrentLevel())
	}
}
