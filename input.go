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
