package main

type Hero struct {
	runFrames  [DirectionCount][3]Image
	jumpFrames [DirectionCount]Image
}

const (
	DirectionCount      = 2
	LeftDirectionIndex  = 0
	RightDirectionIndex = 1
)

func NewHero(assets AssetLoader) *Hero {
	return &Hero{
		runFrames: [DirectionCount][3]Image{
			[3]Image{
				assets.LoadImage("gophette_left_run1"),
				assets.LoadImage("gophette_left_run2"),
				assets.LoadImage("gophette_left_run3"),
			},
			[3]Image{
				assets.LoadImage("gophette_right_run1"),
				assets.LoadImage("gophette_right_run2"),
				assets.LoadImage("gophette_right_run3"),
			},
		},
		jumpFrames: [DirectionCount]Image{
			assets.LoadImage("gophette_left_jump"),
			assets.LoadImage("gophette_right_jump"),
		},
	}
}
