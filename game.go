package main

type Game struct {
	graphics Graphics
	camera   Camera

	running          bool
	characters       [2]*Character
	inputStates      [2]inputState
	primaryCharIndex int

	objects      []CollisionObject
	imageObjects []ImageObject
}

type Camera interface {
	CenterAround(x, y int)
	SetBounds(Rectangle)
}

type inputState struct {
	leftDown          bool
	rightDown         bool
	jumpDown          bool
	mustJumpThisFrame bool
}

func NewGame(
	assets AssetLoader,
	graphics Graphics,
	cam Camera,
	cameraFocusCharIndex int,
) *Game {
	hero := NewHero(assets)
	hero.SetBottomCenterTo(100, 800)
	hero.Direction = RightDirectionIndex

	barney := NewBarney(assets)
	barney.SetBottomCenterTo(100, 800)
	barney.Direction = RightDirectionIndex

	objects := []CollisionObject{
		{Rectangle{0, -10000, 1, 20000}, Solid},    // left wall
		{Rectangle{2500, -10000, 1, 20000}, Solid}, // right wall
		{Rectangle{0, 800, 2000, 50}, Solid},       // floor

		{Rectangle{385, 610, 253, 50}, TopSolid},
		{Rectangle{800, 420, 200, 50}, TopSolid},
		{Rectangle{380, 230, 200, 50}, TopSolid},
		{Rectangle{820, 40, 200, 50}, TopSolid},
		{Rectangle{360, -150, 200, 50}, TopSolid},
		{Rectangle{840, -340, 1050, 50}, Solid},
		{Rectangle{340, -530, 200, 50}, TopSolid},
		{Rectangle{100, -783, 200, 50}, TopSolid},
		{Rectangle{300, -1036, 1000, 50}, TopSolid},
		{Rectangle{1700, -1036, 200, 50}, TopSolid},
		{Rectangle{1840, -840, 50, 1000}, Solid},
	}

	y := 600
	imageObjects := []ImageObject{
		{assets.LoadImage("grass left"), 370, y},
		{assets.LoadImage("grass center 1"), 422, y},
		{assets.LoadImage("grass center 2"), 475, y},
		{assets.LoadImage("grass center 3"), 528, y},
		{assets.LoadImage("grass right"), 583, y},
	}

	cam.SetBounds(Rectangle{0, -1399, 2500, 2200})

	return &Game{
		running:          true,
		graphics:         graphics,
		characters:       [2]*Character{hero, barney},
		primaryCharIndex: cameraFocusCharIndex,
		objects:          objects,
		imageObjects:     imageObjects,
		camera:           cam,
	}
}

func (g *Game) HandleInput(event InputEvent) {
	recordInput(event)

	inputState := &g.inputStates[event.CharacterIndex]

	if event.Action == GoLeft {
		inputState.leftDown = event.Pressed
	}
	if event.Action == GoRight {
		inputState.rightDown = event.Pressed
	}
	if event.Action == Jump {
		inputState.mustJumpThisFrame = event.Pressed
		inputState.jumpDown = event.Pressed
	}

	if event.Action == QuitGame {
		g.running = false
		if recordingInput {
			saveRecordedInputs()
		}
	}
}

func (g *Game) Update() {
	for len(recordedInputs) > 0 && recordedInputs[0].frame == frame {
		if recordedInputs[0].event.Action != QuitGame {
			recordedInputs[0].event.CharacterIndex = 1
			g.HandleInput(recordedInputs[0].event)
		}
		recordedInputs = recordedInputs[1:]
	}

	frame++

	g.updateCharacter(0)
	g.updateCharacter(1)

	g.camera.CenterAround(g.characters[g.primaryCharIndex].Position.Center())
}

func (g *Game) updateCharacter(charIndex int) {
	char := g.characters[charIndex]
	inputState := &g.inputStates[charIndex]

	// decelerate to 0
	if char.SpeedX > 0 {
		char.SpeedX -= char.Params.DecelerationX
		if char.SpeedX < 0 {
			char.SpeedX = 0
		}
	}
	if char.SpeedX < 0 {
		char.SpeedX += char.Params.DecelerationX
		if char.SpeedX > 0 {
			char.SpeedX = 0
		}
	}

	// accelerate the character if pressing left or right (exclusively)
	if inputState.leftDown && !inputState.rightDown {
		char.SpeedX -= char.Params.AccelerationX
		if char.SpeedX < -char.Params.MaxSpeedX {
			char.SpeedX = -char.Params.MaxSpeedX
		}
	}
	if inputState.rightDown && !inputState.leftDown {
		char.SpeedX += char.Params.AccelerationX
		if char.SpeedX > char.Params.MaxSpeedX {
			char.SpeedX = char.Params.MaxSpeedX
		}
	}

	// mustJumpThisFrame is for avoiding jumping again after a jump is over.
	// If you press jump and keep holding it until you land, you should not
	// launch into the next jump right away. Only when you release the jump
	// button and press it again will you launch another jump
	if inputState.mustJumpThisFrame && !char.InAir {
		char.SpeedY = char.Params.InitialJumpSpeedY
	}
	inputState.mustJumpThisFrame = false

	goingUp := char.SpeedY < 0
	if goingUp && inputState.jumpDown {
		// make her jump higher if holding jump while going up
		char.SpeedY += char.Params.LowGravity
	} else {
		char.SpeedY += char.Params.HighGravity
	}
	if char.SpeedY > char.Params.MaxSpeedY {
		char.SpeedY = char.Params.MaxSpeedY
	}

	char.Update(g)
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
			if g.objects[i].Solidness == Solid &&
				g.objects[i].Bounds.Overlaps(moveSpace) {
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
			if g.objects[i].Solidness == Solid &&
				g.objects[i].Bounds.Overlaps(moveSpace) {
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
	newBounds = bounds.MoveBy(0, dy)
	if dy < 0 {
		moveSpace := bounds
		moveSpace.Y += dy
		moveSpace.H -= dy // make it wider, dy is negative
		for i := range g.objects {
			if g.objects[i].Solidness == Solid &&
				g.objects[i].Bounds.Overlaps(moveSpace) {
				collided = true
				overlap := g.objects[i].Bounds.Y + g.objects[i].Bounds.H - moveSpace.Y
				moveSpace.Y += overlap
				moveSpace.H -= overlap
			}
		}
		newBounds = bounds.MoveTo(moveSpace.X, moveSpace.Y)
	}
	// when jumping up you are allowed to go through TopSolid objects from the
	// bottom when you land on an object (going down) you come to a halt and
	// stand on it; this means the only going down needs to be considered for
	// collision detection
	if dy > 0 {
		moveSpace := bounds
		moveSpace.Y += bounds.H
		moveSpace.H = dy
		for i := range g.objects {
			objBounds := g.objects[i].Bounds
			objBounds.H = 1
			if objBounds.Overlaps(moveSpace) {
				collided = true
				overlap := moveSpace.Y + moveSpace.H - g.objects[i].Bounds.Y
				moveSpace.H -= overlap
			}
		}
		newBounds = bounds.MoveBy(0, moveSpace.H)
	}
	return
}

func (g *Game) Running() bool {
	return g.running
}

func (g *Game) Render() {
	for i := range g.objects {
		if g.objects[i].Solidness == Solid {
			g.graphics.FillRect(g.objects[i].Bounds, 255, 0, 0, 255)
		} else {
			g.graphics.FillRect(g.objects[i].Bounds, 133, 98, 98, 255)
			//g.graphics.FillRect(g.objects[i].Bounds, 0, 192, 0, 255)
		}
	}

	for i := range g.imageObjects {
		g.imageObjects[i].Render()
	}

	g.characters[1].Render()
	g.characters[0].Render()
}
