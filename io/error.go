package io

import "errors"

var (
	ErrReadHeader      = Error{errors.New("error reading header")}
	ErrIndexOutOfRange = Error{errors.New("index out of range")}
)

type Error struct {
	error
}
