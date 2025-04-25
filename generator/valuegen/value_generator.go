package valuegen

import (
	"math"
	"math/rand/v2"
	"strings"
	"time"
)

type Option[O any] interface {
	apply(*O)
}

type optionFunc[O any] func(*O)

func (f optionFunc[O]) apply(opts *O) {
	f(opts)
}

type StringOption struct {
	MinLen  *int
	MaxLen  *int
	Charset string
}

func WithStrLimit(minLen int, maxLen int) Option[StringOption] {
	return optionFunc[StringOption](func(so *StringOption) {
		so.MinLen = &minLen
		so.MaxLen = &maxLen
	})
}

func WithCharset(charset string) Option[StringOption] {
	return optionFunc[StringOption](func(so *StringOption) {
		so.Charset = charset
	})
}

func GenString(opts ...Option[StringOption]) string {
	var minLen = 5
	var maxLen = 100

	option := &StringOption{
		MinLen:  &minLen,
		MaxLen:  &maxLen,
		Charset: "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789",
	}

	for _, opt := range opts {
		opt.apply(option)
	}

	minLen = *option.MinLen
	maxLen = *option.MaxLen

	length := rand.IntN(maxLen-minLen) + minLen
	charset := option.Charset

	var sb strings.Builder
	for range length {
		idx := rand.IntN(len(charset))
		sb.WriteByte(charset[idx])
	}
	return sb.String()
}

// 定义一个泛型类型约束
type Number interface {
	int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64 | float32 | float64
}

func GenNumber[T Number](minNum T, maxNum T) T {
	var result T

	// 判断具体类型并返回相应的最大值和最小值
	switch any(result).(type) {
	case uint, uint8, uint16, uint32, uint64:
		result = T(rand.Uint64N(uint64(maxNum-minNum)) + uint64(minNum))
	case int:
		if minNum < 0 && rand.IntN(2) == 0 {
			m := int(minNum)
			if m == math.MinInt {
				m = -math.MaxInt
			}
			negativeResult := T(rand.Int64N(int64(-m)))
			result = -negativeResult
		} else {
			result = T(rand.Int64N(int64(maxNum)))
		}
	case int8:
		if minNum < 0 && rand.IntN(2) == 0 {
			m := int8(minNum)
			if m == math.MinInt8 {
				m = -math.MaxInt8
			}
			negativeResult := T(rand.Int64N(int64(-m)))
			result = -negativeResult
		} else {
			result = T(rand.Int64N(int64(maxNum)))
		}

	case int16:
		if minNum < 0 && rand.IntN(2) == 0 {
			m := int16(minNum)
			if m == math.MinInt16 {
				m = -math.MaxInt16
			}
			negativeResult := T(rand.Int64N(int64(-m)))
			result = -negativeResult
		} else {
			result = T(rand.Int64N(int64(maxNum)))
		}

	case int32:
		if minNum < 0 && rand.IntN(2) == 0 {
			m := int32(minNum)
			if m == math.MinInt32 {
				m = -math.MaxInt32
			}
			negativeResult := T(rand.Int64N(int64(-m)))
			result = -negativeResult
		} else {
			result = T(rand.Int64N(int64(maxNum)))
		}

	case int64:
		if minNum < 0 && rand.IntN(2) == 0 {
			m := int64(minNum)
			if m == math.MinInt64 {
				m = -math.MaxInt64
			}
			negativeResult := T(rand.Int64N(int64(-m)))
			result = -negativeResult
		} else {
			result = T(rand.Int64N(int64(maxNum)))
		}

	case float32:
		if minNum < 0 && rand.IntN(2) == 0 {
			m := float32(minNum)
			result = T(rand.Float32() * m)
		} else {
			m := float32(maxNum)
			result = T(rand.Float32() * m)
		}

	case float64:
		if minNum < 0 && rand.IntN(2) == 0 {
			m := float64(minNum)
			result = T(rand.Float64() * m)
		} else {
			m := float64(maxNum)
			result = T(rand.Float64() * m)
		}
	}
	return result
}

func GenBool() bool {
	return rand.IntN(2) == 0
}

func GenTime(startTime *time.Time, endTime *time.Time) *time.Time {
	if startTime == nil {
		s := time.Now().AddDate(0, 0, -365)
		startTime = &s
	}

	if endTime == nil {
		e := time.Now()
		endTime = &e
	}

	diff := endTime.Sub(*startTime)
	randomDuration := time.Duration(rand.Int64N(int64(diff)))
	result := startTime.Add(randomDuration)
	return &result
}

func GenListValue[T any](list []T) T {
	idx := rand.IntN(len(list))
	return list[idx]
}
