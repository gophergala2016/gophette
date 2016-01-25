package main

const (
	PrePlayFrameDelay    = 100
	PlayerDyingDelay     = 100
	LosingSoundDelay     = 90
	BarneyWinDelay       = 100
	PlayerWinDelay       = 80
	WhistleSoundDuration = 30
	IntroDuration        = 930
)

type Game struct {
	graphics Graphics
	camera   Camera

	state                GameState
	prePlayCountDown     int
	playerDyingCountDown int
	dieBounds            Rectangle
	aiInputs             []inputRecord
	goalBounds           Rectangle
	losingSoundCountDown int
	barneyWinCountDown   int
	playerWinCountDown   int
	introCountUp         int

	running          bool
	characters       [2]*Character
	inputStates      [2]inputState
	primaryCharIndex int

	objects      []CollisionObject
	imageObjects []ImageObject

	winningSound         Sound
	losingSound          Sound
	fallingSound         Sound
	barneyWinSound       Sound
	whistleSound         Sound
	barneyIntroTextSound Sound
	introInstructions    Sound

	introPC1            Image
	introPC2            Image
	introGophette       Image
	currentIntroPCImage int
	introBarneyTalking  bool
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

type GameState int

const (
	Playing GameState = iota + 1
	PrePlaying
	PlayerDying
	PlayerWinning
	PlayerRealizingLoss
	CameraShowsBarneyWinning
	IntroPCScene
)

func NewGame(
	assets AssetLoader,
	graphics Graphics,
	cam Camera,
	cameraFocusCharIndex int,
) *Game {
	hero := NewHero(assets)
	hero.SetBottomCenterTo(500, 537)
	//hero.SetBottomCenterTo(8800, -800) // TODO
	hero.Direction = RightDirectionIndex

	barney := NewBarney(assets)
	//barney.SetBottomCenterTo(8500, -800)
	barney.SetBottomCenterTo(300, 537) // TODO
	barney.Direction = RightDirectionIndex

	cameraBounds := Rectangle{200, -1399, 9150, 2100}
	cam.SetBounds(cameraBounds)

	game := &Game{
		running:              true,
		graphics:             graphics,
		characters:           [2]*Character{hero, barney},
		primaryCharIndex:     cameraFocusCharIndex,
		camera:               cam,
		dieBounds:            cameraBounds.AddMargin(200),
		aiInputs:             recordedInputs,
		goalBounds:           Rectangle{9200, -1000, 1000, 350},
		winningSound:         assets.LoadSound("win"),
		losingSound:          assets.LoadSound("lose"),
		fallingSound:         assets.LoadSound("fall"),
		barneyWinSound:       assets.LoadSound("barney wins"),
		whistleSound:         assets.LoadSound("whistle"),
		barneyIntroTextSound: assets.LoadSound("barney intro text"),
		introInstructions:    assets.LoadSound("instructions"),
		introPC1:             assets.LoadImage("intro pc 1"),
		introPC2:             assets.LoadImage("intro pc 2"),
		introGophette:        assets.LoadImage("intro gophette"),
	}
	game.loadLevel(assets, &level1)
	game.state = IntroPCScene
	return game
}

func (g *Game) loadLevel(assets AssetLoader, level *Level) {
	g.imageObjects = make([]ImageObject, len(level.Images))
	for i := range level.Images {
		img := &level.Images[i]
		g.imageObjects[i].image = assets.LoadImage(img.ID)
		g.imageObjects[i].X = img.X
		g.imageObjects[i].Y = img.Y
	}

	g.objects = make([]CollisionObject, len(level.Objects))
	for i := range level.Objects {
		g.objects[i].Solidness = TopSolid
		if level.Objects[i].Solid {
			g.objects[i].Solidness = Solid
		}
		g.objects[i].Bounds = Rectangle{
			level.Objects[i].X,
			level.Objects[i].Y,
			level.Objects[i].W,
			level.Objects[i].H,
		}
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
	if g.state == IntroPCScene {
		g.introCountUp++

		if g.introCountUp == 100 {
			g.introBarneyTalking = true
			g.barneyIntroTextSound.PlayOnce()
		}
		if g.currentIntroPCImage < 2 && g.introCountUp%10 == 0 {
			g.currentIntroPCImage = 1 - g.currentIntroPCImage
		}
		if g.introCountUp == 580 {
			g.currentIntroPCImage = 2
		}
		if g.introCountUp == 620 {
			g.losingSound.PlayOnce()
		}
		if g.introCountUp == 700 {
			g.introInstructions.PlayOnce()
		}

		if g.introCountUp >= IntroDuration {
			g.prePlayCountDown = PrePlayFrameDelay
			g.state = PrePlaying
		}
	} else if g.state == Playing {
		for len(g.aiInputs) > 0 && g.aiInputs[0].frame == frame {
			if g.aiInputs[0].event.Action != QuitGame {
				g.aiInputs[0].event.CharacterIndex = 1
				g.HandleInput(g.aiInputs[0].event)
			}
			g.aiInputs = g.aiInputs[1:]
		}

		frame++

		g.updateCharacter(0)
		g.updateCharacter(1)

		if !g.dieBounds.Overlaps(g.characters[0].Position) {
			g.fallingSound.PlayOnce()
			g.state = PlayerDying
			g.playerDyingCountDown = PlayerDyingDelay
		}

		if g.goalBounds.Contains(g.characters[0].Position) {
			g.winningSound.PlayOnce()
			g.playerWinCountDown = PlayerWinDelay
			g.state = PlayerWinning
		} else if g.goalBounds.Contains(g.characters[1].Position) {
			g.losingSound.PlayOnce()
			g.state = PlayerRealizingLoss
			g.losingSoundCountDown = LosingSoundDelay
		}

		g.camera.CenterAround(g.characters[g.primaryCharIndex].Position.Center())
	} else if g.state == PrePlaying {
		g.camera.CenterAround(g.characters[g.primaryCharIndex].Position.Center())
		g.prePlayCountDown--
		if g.prePlayCountDown == WhistleSoundDuration {
			g.whistleSound.PlayOnce()
		}
		if g.prePlayCountDown <= 0 {
			g.state = Playing
		}
	} else if g.state == PlayerDying {
		g.playerDyingCountDown--
		if g.playerDyingCountDown <= 0 {
			g.resetLevel()
		}
	} else if g.state == PlayerWinning {
		g.characters[0].Reset(LeftDirectionIndex)
		g.characters[1].Reset(RightDirectionIndex)
		g.playerWinCountDown--
		if g.playerWinCountDown < 0 {
			// TODO go to end cut-scene
		}
	} else if g.state == PlayerRealizingLoss {
		g.losingSoundCountDown--
		if g.losingSoundCountDown <= 0 {
			g.state = CameraShowsBarneyWinning
			g.barneyWinCountDown = BarneyWinDelay
			g.barneyWinSound.PlayOnce()
			g.characters[1].Reset(LeftDirectionIndex)
		}
	} else if g.state == CameraShowsBarneyWinning {
		g.camera.CenterAround(g.characters[1].Position.Center())
		g.barneyWinCountDown--
		if g.barneyWinCountDown <= 0 {
			g.resetLevel()
		}
	}
}

func (g *Game) resetLevel() {
	g.characters[0].SetBottomCenterTo(500, 537)
	g.characters[0].Reset(RightDirectionIndex)

	g.characters[1].SetBottomCenterTo(300, 537)
	g.characters[1].Reset(RightDirectionIndex)

	g.aiInputs = make([]inputRecord, len(recordedInputs))
	copy(g.aiInputs, recordedInputs)
	frame = 0

	g.state = PrePlaying
	g.prePlayCountDown = PrePlayFrameDelay
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
	if g.state == IntroPCScene {
		x, y := 1000, 0
		g.camera.CenterAround(x, y)
		g.graphics.ClearScreen(0, 0, 0)

		img := g.introPC1
		if g.introBarneyTalking && g.currentIntroPCImage == 1 {
			img = g.introPC2
		}
		if g.currentIntroPCImage == 2 {
			img = g.introGophette
		}

		w, h := img.Size()
		img.DrawAt(x-w/2, y-h/2)
	} else {
		for i := range g.imageObjects {
			g.imageObjects[i].Render()
		}

		g.characters[1].Render()
		g.characters[0].Render()
	}
}
