package imagemanager

import (
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
	// m := map[uint]string{
	// 	0: "a",
	// 	1: "b",
	// 	2: "c",
	// 	3: "d",
	// 	4: "e",
	// 	5: "f",
	// }
	//
	// v, ok := m[0]
	// s.Equal("", ok)
	// s.Equal("", v)

	// encoder := base64.StdEncoding
	// encoded := encoder.EncodeToString([]byte("ftp://host/path/to/file"))
	// s.Equal("", encoded)
	// decoded, _ := encoder.DecodeString(encoded)
	// s.Equal("", string(decoded))
}
