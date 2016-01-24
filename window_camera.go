package main

type windowCamera struct {
	position Rectangle
	bounds   Rectangle
}

func newWindowCamera(windowW, windowH int) *windowCamera {
	cam := &windowCamera{
		// initially set no bounds (big integers)
		bounds: Rectangle{-999999, -999999, 2 * 999999, 2 * 999999},
	}
	cam.setWindowSize(windowW, windowH)
	return cam
}

func (cam *windowCamera) setWindowSize(w, h int) {
	cx, cy := cam.position.Center()
	cam.position.W, cam.position.H = w, h
	cam.CenterAround(cx, cy)
}

func (cam *windowCamera) CenterAround(x, y int) {
	cam.position.X = x - cam.position.W/2
	cam.position.Y = y - cam.position.H/2

	// keep the camera in the bounds
	if cam.position.X < cam.bounds.X {
		cam.position.X = cam.bounds.X
	}
	if cam.position.Y < cam.bounds.Y {
		cam.position.Y = cam.bounds.Y
	}
	if cam.position.X+cam.position.W > cam.bounds.X+cam.bounds.W {
		cam.position.X = cam.bounds.X + cam.bounds.W - cam.position.W
	}
	if cam.position.Y+cam.position.H > cam.bounds.Y+cam.bounds.H {
		cam.position.Y = cam.bounds.Y + cam.bounds.H - cam.position.H
	}
}

func (cam *windowCamera) SetBounds(bounds Rectangle) {
	cam.bounds = bounds
}

func (cam *windowCamera) offset() (dx, dy int) {
	return -cam.position.X, -cam.position.Y
}
