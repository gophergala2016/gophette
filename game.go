package main

type Game struct {
	graphics Graphics
	camera   Camera

	running bool
	hero    *Hero

	leftDown          bool
	rightDown         bool
	jumpDown          bool
	mustJumpThisFrame bool

	objects []CollisionObject
}

type Camera interface {
	CenterAround(x, y int)
}

func NewGame(assets AssetLoader, graphics Graphics, cam Camera) *Game {
	hero := NewHero(assets)
	hero.SetBottomCenterTo(380, -900)
	hero.Direction = RightDirectionIndex

	objects := []CollisionObject{
		{Rectangle{10, -10000, 1, 20000}},   // left wall
		{Rectangle{1998, -10000, 1, 20000}}, // right wall
		{Rectangle{0, 800, 2000, 50}},       // floor

		{Rectangle{400, 610, 200, 50}},
		{Rectangle{800, 420, 200, 50}},
		{Rectangle{380, 230, 200, 50}},
		{Rectangle{820, 40, 200, 50}},
		{Rectangle{360, -150, 200, 50}},
		{Rectangle{840, -340, 200, 50}},
		{Rectangle{340, -530, 200, 50}},
		{Rectangle{100, -783, 200, 50}}, // max jump height is 253
		{Rectangle{300, -783 - 253, 200, 50}},
	}

	return &Game{
		running:  true,
		graphics: graphics,
		hero:     hero,
		objects:  objects,
		camera:   cam,
	}
}

func (g *Game) HandleInput(event InputEvent) {
	if event.Action == GoLeft {
		g.leftDown = event.Pressed
	}
	if event.Action == GoRight {
		g.rightDown = event.Pressed
	}
	if event.Action == Jump {
		g.mustJumpThisFrame = event.Pressed
		g.jumpDown = event.Pressed
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

	// mustJumpThisFrame is for avoiding jumping again after a jump is over.
	// If you press jump and keep holding it until you land, you should not
	// launch into the next jump right away. Only when you release the jump
	// button and press it again will you launch another jump
	if g.mustJumpThisFrame && !g.hero.InAir {
		g.hero.SpeedY = HeroInitialJumpSpeedY
	}
	g.mustJumpThisFrame = false

	goingUp := g.hero.SpeedY < 0
	if goingUp && g.jumpDown {
		// make her jump higher if holding jump while going up
		g.hero.SpeedY += HeroLowGravity
	} else {
		g.hero.SpeedY += HeroHighGravity
	}
	if g.hero.SpeedY > HeroMaxSpeedY {
		g.hero.SpeedY = HeroMaxSpeedY
	}

	g.hero.Update(g)

	g.camera.CenterAround(g.hero.Position.Center())
}

func (g *Game) MoveInX(bounds Rectangle, dx int) (newBounds Rectangle, collided bool) {
	newBounds = bounds.MoveBy(dx, 0)
	// create a rectangle that occupies all space from current to new
	// position and then check if it overlaps any object
	if dx < 0 {
		moveSpace := bounds
		moveSpace.X += dx
		moveSpace.W -= dx // make it wider, dx is negative
		for i := range g.objects {
			if g.objects[i].Bounds.Overlaps(moveSpace) {
				collided = true
				overlap := g.objects[i].Bounds.X + g.objects[i].Bounds.W - moveSpace.X
				moveSpace.X += overlap
				moveSpace.W -= overlap
			}
		}
		newBounds = bounds.MoveTo(moveSpace.X, moveSpace.Y)
	}
	if dx > 0 {
		moveSpace := bounds
		moveSpace.W += dx
		for i := range g.objects {
			if g.objects[i].Bounds.Overlaps(moveSpace) {
				collided = true
				overlap := moveSpace.X + moveSpace.W - g.objects[i].Bounds.X
				moveSpace.W -= overlap
			}
		}
		newBounds = bounds.MoveTo(moveSpace.X+moveSpace.W-bounds.W, moveSpace.Y)
	}
	return
}

func (g *Game) MoveInY(bounds Rectangle, dy int) (newBounds Rectangle, collided bool) {
	// this works analogous to MoveInX
	newBounds = bounds.MoveBy(0, dy)
	if dy < 0 {
		moveSpace := bounds
		moveSpace.Y += dy
		moveSpace.H -= dy // make it higher, dy is negative
		for i := range g.objects {
			if g.objects[i].Bounds.Overlaps(moveSpace) {
				collided = true
				overlap := g.objects[i].Bounds.Y + g.objects[i].Bounds.H - moveSpace.Y
				moveSpace.Y += overlap
				moveSpace.H -= overlap
			}
		}
		newBounds = bounds.MoveTo(moveSpace.X, moveSpace.Y)
	}
	if dy > 0 {
		moveSpace := bounds
		moveSpace.H += dy
		for i := range g.objects {
			if g.objects[i].Bounds.Overlaps(moveSpace) {
				collided = true
				overlap := moveSpace.Y + moveSpace.H - g.objects[i].Bounds.Y
				moveSpace.H -= overlap
			}
		}
		newBounds = bounds.MoveTo(moveSpace.X, moveSpace.Y+moveSpace.H-bounds.H)
	}
	return
}

func (g *Game) Running() bool {
	return g.running
}

func (g *Game) Render() {
	for i := range g.objects {
		g.graphics.FillRect(g.objects[i].Bounds, 255, 0, 0, 255)
	}

	g.hero.Render()
}
