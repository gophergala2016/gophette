package main

type Game struct {
	running bool
	hero    *Hero

	leftDown  bool
	rightDown bool
}

func NewGame(assets AssetLoader) *Game {
	hero := NewHero(assets)
	hero.X, hero.Y = 100, 800
	hero.Direction = RightDirectionIndex

	return &Game{
		running: true,
		hero:    hero,
	}
}

func (g *Game) HandleInput(event InputEvent) {
	if event.Action == GoLeft {
		g.leftDown = event.Pressed
	}
	if event.Action == GoRight {
		g.rightDown = event.Pressed
	}

	if event.Action == QuitGame {
		g.running = false
	}
}

func (g *Game) Update() {
	// decelerate the hero to 0
	if g.hero.SpeedX > 0 {
		g.hero.SpeedX -= HeroDecelerationX
		if g.hero.SpeedX < 0 {
			g.hero.SpeedX = 0
		}
	}
	if g.hero.SpeedX < 0 {
		g.hero.SpeedX += HeroDecelerationX
		if g.hero.SpeedX > 0 {
			g.hero.SpeedX = 0
		}
	}

	// accelerate the hero if pressing left or right (exclusively)
	if g.leftDown && !g.rightDown {
		g.hero.SpeedX -= HeroAccelerationX
		if g.hero.SpeedX < -HeroMaxSpeedX {
			g.hero.SpeedX = -HeroMaxSpeedX
		}
	}
	if g.rightDown && !g.leftDown {
		g.hero.SpeedX += HeroAccelerationX
		if g.hero.SpeedX > HeroMaxSpeedX {
			g.hero.SpeedX = HeroMaxSpeedX
		}
	}

	g.hero.Update()
}

func (g *Game) Running() bool {
	return g.running
}

func (g *Game) Render() {
	g.hero.Render()
}
