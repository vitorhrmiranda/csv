package io

import "errors"

var (
	ErrReadHeader      = Error{errors.New("error reading header")}
	ErrReadLine        = Error{errors.New("error reading line")}
	ErrIndexOutOfRange = Error{errors.New("index out of range")}
	ErrMapRows         = Error{errors.New("error mapping rows")}
)

type Error struct {
	error
}
