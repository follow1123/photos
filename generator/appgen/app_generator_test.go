package appgen_test

import (
	"os"
	"testing"

	"github.com/follow1123/photos/generator/appgen"
	"github.com/stretchr/testify/suite"
)

type AppGeneratorTestSuite struct {
	suite.Suite
}

func TestAppGeneratorTestSuite(t *testing.T) {
	suite.Run(t, &AppGeneratorTestSuite{})
}

func (s *AppGeneratorTestSuite) TestGenConfig() {
	appComponents := &appgen.AppComponents{}
	conf, err := appgen.GenConfig(appComponents)
	s.Nil(err)
	s.Equal(":8088", conf.GetAddr())
	err = conf.CreatePath()
	s.Nil(err)
	err = conf.DeletePath()
	s.Nil(err)
}

func (s *AppGeneratorTestSuite) TestGenLogger() {
	appComponents := &appgen.AppComponents{}
	_, err := appgen.GenGinLogger(appComponents)
	s.Nil(err)
	s.NotNil(appComponents.BaseLogger)
}

func (s *AppGeneratorTestSuite) TestGenDatabase() {
	appComponents := &appgen.AppComponents{}
	db, err := appgen.GenDatabase(appComponents)
	s.Nil(err)
	defer db.Config.DeletePath()
	s.NotNil(appComponents.BaseLogger)
	s.NotNil(appComponents.GormLogger)
	s.NotNil(appComponents.Config)
	info, err := os.Stat(db.DBFile)
	s.Nil(err)
	s.Equal(info.Name(), "photos.db")
}

func (s *AppGeneratorTestSuite) TestGenDBMigrator() {
	appComponents := &appgen.AppComponents{}
	migrator, err := appgen.GenDBMigrator(appComponents)
	s.Nil(err)
	defer migrator.DB.Config.DeletePath()
	s.NotNil(appComponents.BaseLogger)
	s.NotNil(appComponents.GormLogger)
	s.NotNil(appComponents.Config)
	s.NotNil(appComponents.DB)
	s.NotNil(migrator)
}
