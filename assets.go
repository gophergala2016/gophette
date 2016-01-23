package main

type Graphics interface {
	FillRect(rect Rectangle, r, g, b, a uint8)
}

type Image interface {
	DrawAt(x, y int)
	Size() (width, height int)
}

type AssetLoader interface {
	LoadImage(id string) Image
}
