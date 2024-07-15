package io

import "errors"

var ErrReadHeader = Error{errors.New("error reading header")}

type Error struct {
	error
}
