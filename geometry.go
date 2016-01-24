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

func (r Rectangle) Center() (x, y int) {
	return r.X + r.W/2, r.Y + r.H/2
}

func (r Rectangle) AddMargin(margin int) Rectangle {
	return Rectangle{r.X - margin, r.Y - margin, r.W + 2*margin, r.H + 2*margin}
}

func (r Rectangle) Contains(o Rectangle) bool {
	return o.X >= r.X && o.Y >= r.Y && o.X+o.W <= r.X+r.W && o.Y+o.H <= r.Y+r.H
}

type Point struct {
	X, Y int
}
