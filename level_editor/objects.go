package main

type LevelObject struct {
	X, Y, W, H int
	Solid      bool
}

var LevelObjects = []LevelObject{
	{175, -608, 29, 1192, true},
	{204, 537, 5293, 47, true},
}
