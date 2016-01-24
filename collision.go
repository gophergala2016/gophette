package main

type CollisionObject struct {
	Bounds    Rectangle
	Solidness Solidness
}

type Solidness int

const (
	// Solid means not walkable from any side, the character can never overlap
	// the object
	Solid Solidness = iota
	// TopSolid means only when jumping on the object from above will it stop
	// you, you can walk through it sideways and jump through it from below.
	TopSolid
)
