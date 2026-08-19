package utils

type PBorders int

const (
	BorderAll PBorders = iota
	BorderBottom
	BorderNone
)

type PAlign int

const (
	AlignNone PAlign = iota
	AlignHorCenter
	AlignHorRight
	AlignHorLeft
	AlignVertCenter
)
