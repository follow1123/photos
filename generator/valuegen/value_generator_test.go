package valuegen_test

import (
	"fmt"
	"math"
	"testing"
	"time"

	"github.com/follow1123/photos/generator/valuegen"
	"github.com/stretchr/testify/suite"
)

type ValueGeneratorTestSuite struct {
	suite.Suite
}

func TestValueGeneratorTestSuite(t *testing.T) {
	suite.Run(t, &ValueGeneratorTestSuite{})
}

func (s *ValueGeneratorTestSuite) testGenString() {

	for i := range 30 {
		_ = i
		fmt.Printf("random string: %s\n", valuegen.GenString())
	}
	for i := range 30 {
		_ = i
		var minLen = 3
		var maxLen = 10
		fmt.Printf("random string(%d, %d): %s\n", minLen, maxLen, valuegen.GenString(valuegen.WithStrLimit(minLen, maxLen)))
	}

	s.False(true)
}

func (s *ValueGeneratorTestSuite) testGenNumber() {

	for i := range 10 {
		_ = i
		fmt.Printf("random int: %d\n", valuegen.GenNumber[int](math.MinInt, math.MaxInt))
	}

	for i := range 10 {
		_ = i
		var min int = -9
		var max int = 10
		fmt.Printf("random int(%d, %d): %d\n", min, max, valuegen.GenNumber[int](min, max))
	}

	for i := range 10 {
		_ = i
		fmt.Printf("random int8: %d\n", valuegen.GenNumber[int8](math.MinInt8, math.MaxInt8))
	}

	for i := range 10 {
		_ = i
		var min int8 = -9
		var max int8 = 10
		fmt.Printf("random int8(%d, %d): %d\n", min, max, valuegen.GenNumber[int8](min, max))
	}

	for i := range 10 {
		_ = i
		fmt.Printf("random int16: %d\n", valuegen.GenNumber[int16](math.MinInt16, math.MaxInt16))
	}

	for i := range 10 {
		_ = i
		var min int16 = -9
		var max int16 = 10
		fmt.Printf("random int16(%d, %d): %d\n", min, max, valuegen.GenNumber[int16](min, max))
	}

	for i := range 10 {
		_ = i
		fmt.Printf("random int32: %d\n", valuegen.GenNumber[int32](math.MinInt32, math.MaxInt32))
	}

	for i := range 10 {
		_ = i
		var min int32 = -9
		var max int32 = 10
		fmt.Printf("random int32(%d, %d): %d\n", min, max, valuegen.GenNumber[int32](min, max))
	}

	for i := range 10 {
		_ = i
		fmt.Printf("random int64: %d\n", valuegen.GenNumber[int64](math.MinInt64, math.MaxInt64))
	}

	for i := range 10 {
		_ = i
		var min int64 = -9
		var max int64 = 10
		fmt.Printf("random int64(%d, %d): %d\n", min, max, valuegen.GenNumber[int64](min, max))
	}

	for i := range 10 {
		_ = i
		fmt.Printf("random uint: %d\n", valuegen.GenNumber[uint](0, math.MaxUint))
	}

	for i := range 10 {
		_ = i
		var min uint = 200
		var max uint = 32421
		fmt.Printf("random uint(%d, %d): %d\n", min, max, valuegen.GenNumber[uint](min, max))
	}

	for i := range 10 {
		_ = i
		fmt.Printf("random float32: %f\n", valuegen.GenNumber[float32](-math.MaxInt32, math.MaxFloat32))
	}

	for i := range 10 {
		_ = i
		var min float32 = -3.9
		var max float32 = 5.0
		fmt.Printf("random float32(%f, %f): %f\n", min, max, valuegen.GenNumber[float32](min, max))
	}

	s.False(true)
}

func (s *ValueGeneratorTestSuite) testGenBool() {
	for i := range 5 {
		_ = i
		fmt.Printf("random bool: %v\n", valuegen.GenBool())
	}

	s.False(true)
}

func (s *ValueGeneratorTestSuite) testGenTime() {
	for i := range 10 {
		_ = i
		fmt.Printf("random time: %v\n", valuegen.GenTime(nil, nil))
	}
	start, e := time.Parse(time.RFC3339, "2025-11-01T22:08:41+08:00")
	s.Nil(e)

	end, e := time.Parse(time.RFC3339, "2025-11-21T22:08:41+08:00")
	s.Nil(e)

	for i := range 10 {
		_ = i
		fmt.Printf("random time(%v, %v): %v\n", start, end, valuegen.GenTime(&start, &end))
	}

	s.False(true)
}
