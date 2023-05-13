package assert

import (
	"errors"

	tassert "github.com/stretchr/testify/assert"
)

var (
	ErrMock = errors.New("an error")
)

func Error(t tassert.TestingT, err error, msgAndArgs ...interface{}) bool {
	return tassert.Error(t, err, msgAndArgs...)
}

func NoError(t tassert.TestingT, err error, msgAndArgs ...interface{}) bool {
	return tassert.NoError(t, err, msgAndArgs...)
}

func WantError(t tassert.TestingT, want bool, err error, msgAndArgs ...interface{}) bool {
	if want {
		return tassert.Error(t, err, msgAndArgs...)
	} else {
		return tassert.NoError(t, err, msgAndArgs...)
	}
}

func True(t tassert.TestingT, value bool, msgAndArgs ...interface{}) bool {
	return tassert.True(t, value, msgAndArgs...)
}

func Equal(t tassert.TestingT, expected, actual interface{}, msgAndArgs ...interface{}) bool {
	return tassert.Equal(t, expected, actual, msgAndArgs...)
}

func NotEqual(t tassert.TestingT, expected, actual interface{}, msgAndArgs ...interface{}) bool {
	return tassert.NotEqual(t, expected, actual, msgAndArgs...)
}

func NotNil(t tassert.TestingT, object interface{}, msgAndArgs ...interface{}) bool {
	return tassert.NotNil(t, object, msgAndArgs...)
}

func Panic(t tassert.TestingT, f func()) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("the code did not panic")
		}
	}()
	f()
}

func NoPanic(t tassert.TestingT, f func()) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("the code panic")
		}
	}()
	f()
}
