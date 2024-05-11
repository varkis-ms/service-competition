package model

import "errors"

var (
	ErrCompNotFound   = errors.New("competition not found")
	ErrCompExists     = errors.New("competition with same title already exists")
	ErrNoAccessToComp = errors.New("no access to edit this competition")
)
