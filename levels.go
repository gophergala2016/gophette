package main

type LevelImage struct {
	ID   string
	X, Y int
}

type LevelObject struct {
	X, Y, W, H int
	Solid      bool
}

type Level struct {
	Objects []LevelObject
	Images  []LevelImage
}
