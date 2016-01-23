package main

const (
	HeroAccelerationX     = 2
	HeroDecelerationX     = 1
	HeroMaxSpeedX         = 10
	HeroInitialJumpSpeedY = -22
	HeroHighGravity       = 2
	HeroLowGravity        = 1

	HeroRunFrameDelay = 4
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
}

func (h *Hero) Update(collider Collider) {
	if h.SpeedX < 0 {
		h.Direction = LeftDirectionIndex
	}
	if h.SpeedX > 0 {
		h.Direction = RightDirectionIndex
	}

	// NOTE putting this code AFTER the collision detection will not have her
	// run while stuck in a wall
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

	var collided bool
	h.Position, collided = collider.MoveInX(h.Position, h.SpeedX)
	if collided {
		h.SpeedX = 0
	}

	// TODO collide correctly in Y as well
	h.InAir = true // assume this until proven otherwise
	h.Position.Y += h.SpeedY

	if h.Position.Y+h.Position.H > 800 {
		h.Position.Y = 800 - h.Position.H
		h.SpeedY = 0
		h.InAir = false
	}
}
