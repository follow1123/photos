package modelgen_test

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/follow1123/photos/config"
	"github.com/follow1123/photos/generator/appgen"
	"github.com/follow1123/photos/generator/imagegen"
	"github.com/follow1123/photos/generator/valuegen"
	"github.com/follow1123/photos/imagemanager"
	"github.com/follow1123/photos/model/dto"
	"github.com/follow1123/photos/service"
	"github.com/stretchr/testify/suite"
)

type ModelGeneratorTestSuite struct {
	suite.Suite
}

func TestModelGeneratorTestSuite(t *testing.T) {
	suite.Run(t, &ModelGeneratorTestSuite{})
}

// 生成photos表数据，包括图片，方法名第一个改大写执行
func (s *ModelGeneratorTestSuite) testCreatePhoto() {
	appComponents := &appgen.AppComponents{
		Config: config.NewConfig(),
	}
	ctx, err := appgen.GenAppContext(appComponents)
	s.Nil(err)
	db, err := appgen.GenDatabase(appComponents)
	s.Nil(err)
	migrator, err := appgen.GenDBMigrator(appComponents)
	s.Nil(err)
	migrator.InitOrMigrate()

	serv := service.NewPhotoService(ctx, db)

	var recordCount = 100
	params := make([]dto.CreatePhotoParam, 0, recordCount)

	for i := range recordCount {
		_ = i

		var param = dto.CreatePhotoParam{
			Desc:      valuegen.GenString(valuegen.WithStrLimit(5, 30)),
			PhotoDate: *valuegen.GenTime(nil, nil),
		}

		buf := new(bytes.Buffer)
		_, err := imagegen.GenImage(buf)
		s.Nil(err)
		param.ImageSource = imagemanager.NewReaderSource(bytes.NewReader(buf.Bytes()), param.Desc)
		params = append(params, param)
		fmt.Printf("param: %v\n", param)
	}
	failureResults := serv.CreatePhoto(params)
	s.True(len(failureResults) == 0)
	s.True(false)
}
