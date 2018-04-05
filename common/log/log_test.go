package log_test

import (
	"bytes"
	"fmt"
	"github.com/marcusva/docproc/common/log"
	"github.com/marcusva/docproc/common/testing/assert"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

func TestPackage(t *testing.T) {
	if log.Logger() == nil {
		t.Error("_log is nil, although a package initialization was done")
	}

	// None of those should cause a panic
	log.Alert("test")
	log.Alertf("test")
	log.Critical("test")
	log.Criticalf("test")
	log.Debug("test")
	log.Debugf("test")
	log.Info("test")
	log.Infof("test")
	log.Notice("test")
	log.Noticef("test")
	log.Error("test")
	log.Errorf("test")
	log.Warning("test")
	log.Warningf("test")
	log.Emergency("test")
	log.Emergencyf("test")
}

func TestLogger(t *testing.T) {
	logger := log.Logger()
	assert.NotNil(t, logger)

	var buf bytes.Buffer
	log.Init(&buf, log.LevelDebug, true)
	logger2 := log.Logger()
	assert.NotEqual(t, logger, logger2)
}

func TestInitFile(t *testing.T) {
	fp, err := ioutil.TempFile(os.TempDir(), "docproc-logtest")
	assert.NoErr(t, err)
	fname := fp.Name()
	fp.Close()

	err = log.InitFile(fname, log.LevelDebug, false)
	assert.NoErr(t, err)

	log.Init(os.Stdout, log.LevelDebug, false)
	assert.NoErr(t, os.Remove(fname))

	err = log.InitFile("", log.LevelDebug, false)
	assert.Err(t, err)

}

func TestGetLogLevel(t *testing.T) {
	levelsInt := []string{"0", "1", "2", "3", "4", "5", "6", "7"}
	levelsTxt := []string{
		"Emergency", "Alert", "Critical", "Error", "Warning", "Notice", "Info",
		"Debug",
	}

	for idx, v := range levelsInt {
		if v1, err := log.GetLogLevel(v); err != nil {
			t.Error(err)
		} else {
			if v2, err := log.GetLogLevel(levelsTxt[idx]); err != nil {
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
		if v1, err := log.GetLogLevel(""); err == nil || v1 != -1 {
			t.Errorf("invalid level '%s' was accepted", v)
		}
	}
}

func TestLog(t *testing.T) {
	callbacks := map[string]func(...interface{}){
		"DEBUG":     log.Debug,
		"INFO":      log.Info,
		"NOTICE":    log.Notice,
		"WARNING":   log.Warning,
		"ERROR":     log.Error,
		"CRITICAL":  log.Critical,
		"ALERT":     log.Alert,
		"EMERGENCY": log.Emergency,
	}

	var buf bytes.Buffer
	log.Init(&buf, log.LevelDebug, true)

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
		"DEBUG":     log.Debugf,
		"INFO":      log.Infof,
		"NOTICE":    log.Noticef,
		"WARNING":   log.Warningf,
		"ERROR":     log.Errorf,
		"CRITICAL":  log.Criticalf,
		"ALERT":     log.Alertf,
		"EMERGENCY": log.Emergencyf,
	}

	var buf bytes.Buffer
	log.Init(&buf, log.LevelDebug, true)

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
	levels := []log.Level{
		log.LevelDebug, log.LevelInfo, log.LevelNotice, log.LevelWarning,
		log.LevelError, log.LevelAlert, log.LevelCritical, log.LevelEmergency,
	}
	for _, level := range levels {
		var buf bytes.Buffer
		log.Init(&buf, level, true)
		assert.Equal(t, level, log.CurrentLevel())
	}
}
