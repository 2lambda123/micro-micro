// Package config contains helper methods for
// client side config management (`~/.micro/config.json` file).
// It uses the `JSONValues` helper
package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/juju/fslock"
	conf "github.com/micro/go-micro/v3/config"
	"github.com/micro/micro/v3/internal/user"
)

var (

	// lock in single process
	mtx sync.Mutex

	// file for global micro config
	file = filepath.Join(user.Dir, "config.json")

	// full path to file
	path, _ = filePath()

	// a global lock for the config
	lock = fslock.New(path)
)

// SetConfig sets the config file
func SetConfig(f string) {
	mtx.Lock()
	defer mtx.Unlock()

	// path is the full path
	path = f
	// the name of the file
	file = filepath.Base(f)
	// new lock for the file
	lock = fslock.New(path)
}

// config is a singleton which is required to ensure
// each function call doesn't load the .micro file
// from disk

// Get a value from the .micro file
func Get(path string) (string, error) {
	mtx.Lock()
	defer mtx.Unlock()

	config, err := newConfig()
	if err != nil {
		return "", err
	}

	// acquire lock
	if err := lock.Lock(); err != nil {
		return "", err
	}
	defer lock.Unlock()

	val := config.Get(path)
	v := strings.TrimSpace(val.String(""))
	if len(v) > 0 {
		return v, nil
	}

	// try as bytes
	v = string(val.Bytes())
	v = strings.TrimSpace(v)

	// don't return nil decoded value
	if v == "null" {
		return "", nil
	}

	return v, nil
}

// Set a value in the .micro file
func Set(value string, p ...string) error {
	mtx.Lock()
	defer mtx.Unlock()

	config, err := newConfig()
	if err != nil {
		return err
	}
	// acquire lock
	if err := lock.Lock(); err != nil {
		return err
	}
	defer lock.Unlock()

	// set the value
	config.Set(value, p...)

	// write to the file
	return ioutil.WriteFile(path, config.Bytes(), 0644)
}

func filePath() (string, error) {
	return file, nil
}

func moveConfig(from, to string) error {
	// read the config
	b, err := ioutil.ReadFile(from)
	if err != nil {
		return fmt.Errorf("Failed to read config file %s: %v", from, err)
	}
	// remove the file
	os.Remove(from)

	// create new directory
	dir := filepath.Dir(to)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("Failed to create dir %s: %v", dir, err)
	}
	// write the file to new location
	return ioutil.WriteFile(to, b, 0644)
}

// newConfig returns a loaded config
func newConfig() (*conf.JSONValues, error) {
	// check if the directory exists, otherwise create it
	dir := filepath.Dir(path)

	// for legacy purposes check if .micro is a file or directory
	if f, err := os.Stat(dir); err != nil {
		// check the error to see if the directory exists
		if os.IsNotExist(err) {
			if err := os.MkdirAll(dir, 0755); err != nil {
				return nil, fmt.Errorf("Failed to create dir %s: %v", dir, err)
			}
		} else {
			return nil, fmt.Errorf("Failed to create config dir %s: %v", dir, err)
		}
	} else {
		// if not a directory, copy and move the config
		if !f.IsDir() {
			if err := moveConfig(dir, path); err != nil {
				return nil, fmt.Errorf("Failed to move config from %s to %s: %v", dir, path, err)
			}
		}
	}

	// now write the file if it does not exist
	if _, err := os.Stat(path); os.IsNotExist(err) {
		ioutil.WriteFile(path, []byte(`{"env":"local"}`), 0644)
	} else if err != nil {
		return nil, fmt.Errorf("Failed to write config file %s: %v", path, err)
	}

	contents, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	c, err := conf.NewJSONValues(contents)
	if err != nil {
		return nil, err
	}

	// return the conf
	return c, nil
}

func Path(paths ...string) string {
	return strings.Join(paths, ".")
}
