package main

type ImageObject struct {
	image Image
	X, Y  int
}

func (img *ImageObject) Render() {
	img.image.DrawAt(img.X, img.Y)
}
