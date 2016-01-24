package main

type LevelObject struct {
	X, Y, W, H int
	Solid      bool
}

var LevelObjects = []LevelObject{
	{217, 533, 1492, 42, true},
	{596, 364, 142, 34, false},
	{907, 261, 160, 36, false},
	{1731, 689, 178, 49, false},
}
