package config

import (
	"os"
	"path/filepath"
	"strings"
)

const (
	DATA_DIR  = "photos"
	FILES_DIR = "files"
)

type Config struct {
	address    string
	prefixPath string
}

func NewConfig(addr string) *Config {
	return &Config{address: addr, prefixPath: initPath()}
}

func (c *Config) GetFilesPath() string {
	return filepath.Join(c.prefixPath, FILES_DIR)
}

func (c *Config) GetAddr() string {
	return c.address
}

func (c *Config) GetPrefixPath() string {
	return c.prefixPath
}

func initPath() string {
	xdgDatahome := strings.TrimSpace(os.Getenv("XDG_DATA_HOME"))
	if xdgDatahome == "" {
		xdgDatahome = filepath.Join(os.Getenv("HOME"), ".local/share")
	}
	dataHome := filepath.Join(xdgDatahome, DATA_DIR)

	filesHome := filepath.Join(dataHome, FILES_DIR)
	err := os.MkdirAll(filesHome, 0755)
	if err != nil {
		panic(err)
	}

	return dataHome
}
