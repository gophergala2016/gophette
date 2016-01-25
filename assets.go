package main

type Graphics interface {
	FillRect(rect Rectangle, r, g, b, a uint8)
	ClearScreen(r, g, b uint8)
}

type Image interface {
	DrawAt(x, y int)
	Size() (width, height int)
}

type Sound interface {
	PlayOnce()
}

type AssetLoader interface {
	LoadImage(id string) Image
	LoadSound(id string) Sound
}
