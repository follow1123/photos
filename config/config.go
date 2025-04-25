package config

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/follow1123/photos/common"
)

const (
	DATA_DIR  = "photos"
	FILES_DIR = "files"
)

func WithAddress(addr string) common.Option[Config] {
	return common.OptionFunc[Config](func(c *Config) {
		c.address = addr
	})
}

func WithPath(path string) common.Option[Config] {
	return common.OptionFunc[Config](func(c *Config) {
		c.prefixPath = path
	})
}

type Config struct {
	address    string
	prefixPath string
}

func NewConfig(opts ...common.Option[Config]) *Config {
	conf := &Config{}
	for _, opt := range opts {
		opt.Apply(conf)
	}
	if conf.address == "" {
		conf.address = ":8080"
	}

	if conf.prefixPath == "" {
		conf.prefixPath = initPath()
	}
	return conf
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
	return dataHome
}

func (c *Config) CreatePath() error {
	err := os.MkdirAll(c.GetFilesPath(), 0755)
	if err != nil {
		return err
	}
	return nil
}

func (c *Config) DeletePath() error {
	err := os.RemoveAll(c.GetPrefixPath())
	if err != nil {
		return err
	}
	return nil
}
