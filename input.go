package main

type InputEvent struct {
	Action  InputAction
	Pressed bool
}

type InputAction int

const (
	GoLeft InputAction = iota + 1
	GoRight
	Jump
	QuitGame
)

func (a InputAction) String() string {
	switch a {
	case GoLeft:
		return "GoLeft"
	case GoRight:
		return "GoRight"
	case Jump:
		return "Jump"
	case QuitGame:
		return "QuitGame"
	default:
		return "unknown input"
	}
}
