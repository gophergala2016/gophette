package main

type LevelObject struct {
	X, Y, W, H int
	Solid      bool
}

var LevelObjects = []LevelObject{
	{214, 537, 1106, 42, true},
	{596, 364, 142, 34, false},
	{904, 263, 145, 36, false},
	{1379, 731, 178, 49, false},
	{387, -90, 537, 39, true},
	{175, -608, 29, 1192, true},
}
