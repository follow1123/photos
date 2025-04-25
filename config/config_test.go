package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/suite"
)

type ConfigTestSuite struct {
	suite.Suite
}

func TestConfigTestSuite(t *testing.T) {
	suite.Run(t, &ConfigTestSuite{})
}

func (s *ConfigTestSuite) TestInitConfig() {
	conf := NewConfig()
	home := os.Getenv("HOME")

	s.Equal(":8080", conf.GetAddr())
	s.Equal(filepath.Join(home, ".local/share", DATA_DIR), conf.GetPrefixPath())
	s.Equal(filepath.Join(home, ".local/share", DATA_DIR, FILES_DIR), conf.GetFilesPath())
}
