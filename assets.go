package main

type Image interface{}

type AssetLoader interface {
	LoadImage(id string) Image
}
