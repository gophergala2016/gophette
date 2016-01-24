package main

type LevelObject struct {
	X, Y, W, H int
	Solid      bool
}

var LevelObjects = []LevelObject{
	{214, 537, 1106, 42, true},
	{596, 364, 142, 34, false},
	{907, 261, 160, 36, false},
	{1360, 688, 178, 49, false},
}
