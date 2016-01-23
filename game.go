package main

type Game struct {
	running bool
	hero    *Hero
}

func NewGame(assets AssetLoader) *Game {
	return &Game{
		running: true,
		hero:    NewHero(assets),
	}
}

func (g *Game) HandleInput(event InputEvent) {
	// TODO make her run and jump

	if event.Action == QuitGame {
		g.running = false
	}
}

func (g *Game) Running() bool {
	return g.running
}

func (g *Game) Render() {
	g.hero.Render()
}
