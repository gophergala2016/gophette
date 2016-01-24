package main

import "github.com/gophergala2016/gophette/resource"

type CharacterParams struct {
	AccelerationX     int
	DecelerationX     int
	MaxSpeedX         int
	MaxSpeedY         int
	InitialJumpSpeedY int
	HighGravity       int
	LowGravity        int
	RunFrameDelay     int
}

var HeroParams = CharacterParams{
	AccelerationX:     2,
	DecelerationX:     1,
	MaxSpeedX:         10,
	MaxSpeedY:         32,
	InitialJumpSpeedY: -23,
	HighGravity:       2,
	LowGravity:        1,
	RunFrameDelay:     3,
}

var BarneyParams = CharacterParams{
	AccelerationX:     2,
	DecelerationX:     1,
	MaxSpeedX:         12,
	MaxSpeedY:         32,
	InitialJumpSpeedY: -25,
	HighGravity:       2,
	LowGravity:        1,
	RunFrameDelay:     5,
}

type Character struct {
	Direction int
	Position  Rectangle
	SpeedX    int
	SpeedY    int

	InAir bool

	Params        CharacterParams
	collisionRect Rectangle

	runFrames   [DirectionCount][]Image
	standFrames [DirectionCount]Image
	jumpFrames  [DirectionCount]Image

	runFrameIndex int
	nextRunFrame  int
}

func toRect(r resource.Rectangle) Rectangle {
	return Rectangle{r.X, r.Y, r.W, r.H}
}

func NewHero(assets AssetLoader) *Character {
	return &Character{
		Position:      toRect(resource.HeroCollisionRect),
		collisionRect: toRect(resource.HeroCollisionRect),
		Params:        HeroParams,
		runFrames: [DirectionCount][]Image{
			[]Image{
				assets.LoadImage("gophette_left_run1"),
				assets.LoadImage("gophette_left_run2"),
				assets.LoadImage("gophette_left_run1"),
				assets.LoadImage("gophette_left_run3"),
			},
			[]Image{
				assets.LoadImage("gophette_right_run1"),
				assets.LoadImage("gophette_right_run2"),
				assets.LoadImage("gophette_right_run1"),
				assets.LoadImage("gophette_right_run3"),
			},
		},
		standFrames: [DirectionCount]Image{
			assets.LoadImage("gophette_left_run1"),
			assets.LoadImage("gophette_right_run1"),
		},
		jumpFrames: [DirectionCount]Image{
			assets.LoadImage("gophette_left_jump"),
			assets.LoadImage("gophette_right_jump"),
		},
	}
}

func NewBarney(assets AssetLoader) *Character {
	return &Character{
		Position:      toRect(resource.BarneyCollisionRect),
		collisionRect: toRect(resource.BarneyCollisionRect),
		Params:        BarneyParams,
		runFrames: [DirectionCount][]Image{
			[]Image{
				assets.LoadImage("barney_left_run1"),
				assets.LoadImage("barney_left_run2"),
				assets.LoadImage("barney_left_run3"),
				assets.LoadImage("barney_left_run4"),
				assets.LoadImage("barney_left_run5"),
				assets.LoadImage("barney_left_run6"),
			},
			[]Image{
				assets.LoadImage("barney_right_run1"),
				assets.LoadImage("barney_right_run2"),
				assets.LoadImage("barney_right_run3"),
				assets.LoadImage("barney_right_run4"),
				assets.LoadImage("barney_right_run5"),
				assets.LoadImage("barney_right_run6"),
			},
		},
		standFrames: [DirectionCount]Image{
			assets.LoadImage("barney_left_stand"),
			assets.LoadImage("barney_right_stand"),
		},
		jumpFrames: [DirectionCount]Image{
			assets.LoadImage("barney_left_jump"),
			assets.LoadImage("barney_right_jump"),
		},
	}
}

func (c *Character) SetBottomCenterTo(x, y int) {
	c.Position.X = x - c.Position.W/2
	c.Position.Y = y - c.Position.H
}

func (c *Character) Render() {
	var frame Image
	if c.InAir {
		frame = c.jumpFrames[c.Direction]
	} else if c.SpeedX == 0 {
		frame = c.standFrames[c.Direction]
	} else {
		frame = c.runFrames[c.Direction][c.runFrameIndex]
	}

	// the position is that of the collision rectangle, the image does not have
	// the same size as the collision rectangle so it must be offset relative
	// to the collision rectangle's top-left corner for drawing
	frame.DrawAt(
		c.Position.X-c.collisionRect.X,
		c.Position.Y-c.collisionRect.Y,
	)
}

type Collider interface {
	MoveInX(bounds Rectangle, dx int) (newBounds Rectangle, collided bool)
	MoveInY(bounds Rectangle, dy int) (newBounds Rectangle, collided bool)
}

func (c *Character) Update(collider Collider) {
	// face the way of the horizonal speed
	if c.SpeedX < 0 {
		c.Direction = LeftDirectionIndex
	}
	if c.SpeedX > 0 {
		c.Direction = RightDirectionIndex
	}

	// NOTE putting this code AFTER the collision detection will not have the
	// character run while stuck in a wall; BEFORE the collision detection will
	// keep the running animation while hitting the wall
	if c.SpeedX == 0 {
		c.runFrameIndex = 0
		c.nextRunFrame = 0
	} else {
		c.nextRunFrame--
		if c.nextRunFrame <= 0 {
			c.runFrameIndex = (c.runFrameIndex + 1) % len(c.runFrames[c.Direction])
			c.nextRunFrame = c.Params.RunFrameDelay
		}
	}

	// Move in Y first, this assures that you land on a platform even if it is
	// at the maximum jump height; in this case you move up above the platform,
	// then you move in X in the next step and land on the platform.
	// Were it the other way round would mean moving in X, hitting the platform,
	// then moving in Y above the platform but to the side of it
	var collided bool
	c.InAir = true // assume this until proven otherwise
	c.Position, collided = collider.MoveInY(c.Position, c.SpeedY)
	if collided {
		if c.SpeedY > 0 {
			// if she was going down, she now landed on the ground
			c.InAir = false
		}
		c.SpeedY = 0
	}

	// move in X
	c.Position, collided = collider.MoveInX(c.Position, c.SpeedX)
	if collided {
		c.SpeedX = 0
	}

	if c.Position.Y+c.Position.H > 800 {
		c.Position.Y = 800 - c.Position.H
		c.SpeedY = 0
		c.InAir = false
	}
}
