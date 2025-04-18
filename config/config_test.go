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

func (s *ConfigTestSuite) TestInitConfig()  {
	cfg := NewConfig(":8080")
	home := os.Getenv("HOME")

	s.Equal(filepath.Join(home, ".local/share", DATA_DIR), cfg.GetPrefixPath())
	s.Equal(filepath.Join(home, ".local/share", DATA_DIR, FILES_DIR), cfg.GetFilesPath())
}
