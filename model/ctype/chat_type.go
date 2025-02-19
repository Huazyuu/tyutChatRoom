package ctype

type MsgType int

const (
	SystemMsg MsgType = iota + 1
	InRoomMsg
	OutRoomMsg

	TextMsg
	FileMsg
	FileProgressMsg
	ImageMsg
)
