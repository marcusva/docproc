// Package config provides a simple access to INI-style configuration files.
package config

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

// Config is a simple configuration store.
// It consists of unique sections, which contain key-value pairs.
type Config struct {
	Sections map[string]map[string]string
}

// Validator allows a Config to be checked for invalid configuration settings.
type Validator func(cfg *Config) error

// NoValidate is a no-op validation function for Config instances.
func NoValidate(cfg *Config) error { return nil }

// Get gets a value for a certain key within the specified section.
func (cfg *Config) Get(section, key string) (string, error) {
	opts, ok := cfg.Sections[section]
	if !ok {
		return "", fmt.Errorf("section '%s' does not exist", section)
	}
	if v, ok := opts[key]; ok {
		return v, nil
	}
	return "", fmt.Errorf("key '%s' not found in section '%s'", key, section)
}

// GetOrPanic gets a value for a certain key or panics.
func (cfg *Config) GetOrPanic(section, key string) string {
	val, err := cfg.Get(section, key)
	if err != nil {
		panic(err)
	}
	return val
}

// GetDefault gets a value for a certain key. If the section or key could not
// be found, the provided default value will be returned.
func (cfg *Config) GetDefault(section, key, def string) string {
	if val, err := cfg.Get(section, key); err == nil {
		return val
	}
	return def
}

// Int gets a value for a certain key within the specified section as int.
func (cfg *Config) Int(section, key string) (int, error) {
	val, err := cfg.Get(section, key)
	if err != nil {
		return -1, err
	}
	return strconv.Atoi(val)
}

// Bool gets a value for a certain key within the specified section as bool.
func (cfg *Config) Bool(section, key string) (bool, error) {
	val, err := cfg.Get(section, key)
	if err != nil {
		return false, err
	}
	return strconv.ParseBool(val)
}

// Array transforms a comma-separated value into a string array and returns
// it.
func (cfg *Config) Array(section, key string) ([]string, error) {
	val, err := cfg.Get(section, key)
	if err != nil {
		return nil, err
	}
	result := []string{}
	values := strings.Split(val, ",")
	for _, v := range values {
		result = append(result, strings.TrimSpace(v))
	}
	return result, nil
}

// HasSection checks, if the specified section exists within the Config.
func (cfg *Config) HasSection(section string) bool {
	_, ok := cfg.Sections[section]
	return ok
}

// AllFor retrieves a map containing all options for the specified section.
func (cfg *Config) AllFor(section string) (map[string]string, error) {
	opts, ok := cfg.Sections[section]
	if !ok {
		return nil, fmt.Errorf("section '%s' does not exist", section)
	}
	result := make(map[string]string)
	for k, v := range opts {
		result[k] = v
	}
	return result, nil
}

// LoadFile loads the configuration from the passed file. Files have to follow
// the INI file configuration layout.
//
//   # Comments can only be declared on a seperate line
//   # Any line starting with a semicolon (;) or number sign (hash - #)
//   ; is recognized as a comment
//   #
//   # Empty lines are skipped and do not have any effect.
//        # Whitespace are removed from the beginning and end of each line
//
//   # square brackets declare a new section
//   [section]
//   key = value
func LoadFile(filename string, validator Validator) (*Config, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return Load(file, validator)
}

// Load loads the configuration from a io.Reader.
func Load(r io.Reader, validate Validator) (*Config, error) {
	cfg := &Config{
		Sections: make(map[string]map[string]string),
	}

	offset := 0
	scanner := bufio.NewScanner(r)
	var cursection string
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		llen := len(line)
		offset++
		switch {
		case llen == 0, line[0] == '#', line[0] == ';':
			// Skip empty lines and comments
			continue
		case line[0] == '[' && line[llen-1] == ']':
			cursection = strings.TrimSpace(line[1 : llen-1])
			if len(cursection) == 0 {
				return nil, fmt.Errorf("line %d: invalid, empty section name", offset)
			}
			if _, ok := cfg.Sections[cursection]; ok {
				// Section already exists, fail
				return nil, fmt.Errorf("line %d: section '%s' was defined before", offset, cursection)
			}
			cfg.Sections[cursection] = make(map[string]string)
		default:
			if cursection == "" {
				return nil, fmt.Errorf("line %d: key-value definition without section", offset)
			}
			kv := strings.SplitN(line, "=", 2)
			if len(kv) < 2 {
				return nil, fmt.Errorf("line %d: key-value definition misses assignment", offset)
			}
			key := strings.TrimSpace(kv[0])
			val := strings.TrimSpace(kv[1])
			cfg.Sections[cursection][key] = val
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	if validate != nil {
		return cfg, validate(cfg)
	}
	return cfg, nil
}
