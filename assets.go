package main

type Image interface {
	DrawAt(x, y int)
}

type AssetLoader interface {
	LoadImage(id string) Image
}
