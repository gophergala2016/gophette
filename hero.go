package main

// maximum jump height is currently 253

const (
	HeroAccelerationX     = 2
	HeroDecelerationX     = 1
	HeroMaxSpeedX         = 10
	HeroMaxSpeedY         = 32
	HeroInitialJumpSpeedY = -23
	HeroHighGravity       = 2
	HeroLowGravity        = 1

	HeroRunFrameDelay = 3
)

type Hero struct {
	Direction int
	Position  Rectangle
	SpeedX    int
	SpeedY    int

	InAir bool

	runFrames     [DirectionCount][3]Image
	jumpFrames    [DirectionCount]Image
	runFrameIndex int
	nextRunFrame  int
}

func NewHero(assets AssetLoader) *Hero {
	return &Hero{
		Position: HeroCollisionRect,
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

func (h *Hero) SetBottomCenterTo(x, y int) {
	h.Position.X = x - h.Position.W/2
	h.Position.Y = y - h.Position.H
}

func (h *Hero) Render() {
	var frame Image
	if h.InAir {
		frame = h.jumpFrames[h.Direction]
	} else {
		// the order of animation images is 0,1,0,2  0,1,0,2  0,1 ...
		frameIndex := h.runFrameIndex
		if frameIndex == 2 {
			frameIndex = 0
		}
		if frameIndex == 3 {
			frameIndex = 2
		}
		frame = h.runFrames[h.Direction][frameIndex]
	}

	// the position is that of the collision rectangle, the image does not have
	// the same size as the collision rectangle so it must be offset relative
	// to the collision rectangle's top-left corner for drawing
	frame.DrawAt(
		h.Position.X-HeroCollisionRect.X,
		h.Position.Y-HeroCollisionRect.Y,
	)
}

type Collider interface {
	MoveInX(bounds Rectangle, dx int) (newBounds Rectangle, collided bool)
	MoveInY(bounds Rectangle, dy int) (newBounds Rectangle, collided bool)
}

func (h *Hero) Update(collider Collider) {
	// face the way of the horizonal speed
	if h.SpeedX < 0 {
		h.Direction = LeftDirectionIndex
	}
	if h.SpeedX > 0 {
		h.Direction = RightDirectionIndex
	}

	// NOTE putting this code AFTER the collision detection will not have her
	// run while stuck in a wall; BEFORE the collision detection will keep the
	// running animation while hitting the wall
	if h.SpeedX == 0 {
		h.runFrameIndex = 0
		h.nextRunFrame = 0
	} else {
		h.nextRunFrame--
		if h.nextRunFrame <= 0 {
			h.runFrameIndex = (h.runFrameIndex + 1) % 4
			h.nextRunFrame = HeroRunFrameDelay
		}
	}

	// Move in Y first, this assures that you land on a platform even if it is
	// at the maximum jump height; in this case you move up above the platform,
	// then you move in X in the next step and land on the platform.
	// Were it the other way round would mean moving in X, hitting the platform,
	// then moving in Y above the platform but to the side of it
	var collided bool
	h.InAir = true // assume this until proven otherwise
	h.Position, collided = collider.MoveInY(h.Position, h.SpeedY)
	if collided {
		if h.SpeedY > 0 {
			// if she was going down, she now landed on the ground
			h.InAir = false
		}
		h.SpeedY = 0
	}

	// move in X
	h.Position, collided = collider.MoveInX(h.Position, h.SpeedX)
	if collided {
		h.SpeedX = 0
	}

	if h.Position.Y+h.Position.H > 800 {
		h.Position.Y = 800 - h.Position.H
		h.SpeedY = 0
		h.InAir = false
	}
}
