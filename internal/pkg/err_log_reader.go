package pkg

import (
	"errors"
	"strings"
)

func ErrReader(err error) (function string, e error) {
	str := strings.Split(err.Error(), ":")
	function = str[0]
	e = errors.New(str[1])
	return
}
