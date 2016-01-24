package main

type LevelImage struct {
	ID   string
	X, Y int
}

var LevelImages = []LevelImage{
	{"grass left", 0, 0},
	{"grass right", 0, 50},
	{"grass center 1", 0, 100},
	{"grass center 2", 0, 150},
	{"grass center 3", 0, 200},
	{"grass right", 283, 174},
	{"grass right", 381, 209},
	{"grass right", 492, 301},

}
