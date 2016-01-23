package main

type Rectangle struct {
	X, Y, W, H int
}

func (r Rectangle) MoveBy(dx, dy int) Rectangle {
	return Rectangle{r.X + dx, r.Y + dy, r.W, r.H}
}

func (r Rectangle) MoveTo(x, y int) Rectangle {
	return Rectangle{x, y, r.W, r.H}
}

func (r Rectangle) Overlaps(o Rectangle) bool {
	return r.X+r.W > o.X && r.Y+r.H > o.Y &&
		o.X+o.W > r.X && o.Y+o.H > r.Y
}

type Point struct {
	X, Y int
}
