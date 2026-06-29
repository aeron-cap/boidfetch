package main

import (
	"math"

	"github.com/gdamore/tcell/v2"
)

type BoidType int

const (
	NormalBoid BoidType = iota
	PredatorBoid
)

type Boid struct {
	Position     Vector2D
	Velocity     Vector2D
	Acceleration Vector2D
	Type         BoidType
	IsDead       bool
}

const (
	PerceptionRadius = 4.0
	MaxForce         = 0.4
	MaxSpeed         = 0.6

	PredatorPerceptionRadius = 8.0
	PredatorMaxForce         = 0.2
	PredatorMaxSpeed         = 1.0

	PanicRadius = 6.0
)

func (b *Boid) Update(screenWidth, screenHeight int, flock []*Boid) {
	separationForce := b.Separation(flock)
	alignmentForce := b.Alignment(flock)
	cohesionForce := b.Cohesion(flock)

	separationForce = separationForce.Multiply(1.2)
	alignmentForce = alignmentForce.Multiply(1.0)
	cohesionForce = cohesionForce.Multiply(1.0)

	b.Acceleration = b.Acceleration.Add(separationForce)
	b.Acceleration = b.Acceleration.Add(alignmentForce)
	b.Acceleration = b.Acceleration.Add(cohesionForce)

	if b.Type == PredatorBoid {
		chaseForce := Vector2D{X: 0, Y: 0}
		chaseForce = b.Chase(flock)
		chaseForce = chaseForce.Multiply(1.5)
		b.Acceleration = b.Acceleration.Add(chaseForce)

		b.Eat(flock)
	}

	if b.Type == NormalBoid {
		fleeForce := Vector2D{X: 0, Y: 0}
		fleeForce = b.Flee(flock)
		fleeForce = fleeForce.Multiply(1.7)
		b.Acceleration = b.Acceleration.Add(fleeForce)
	}

	b.Velocity = b.Velocity.Add(b.Acceleration)
	if b.Velocity.Magnitude() > MaxSpeed {
		b.Velocity = b.Velocity.Normalize().Multiply(MaxSpeed)
	}

	b.Position = b.Position.Add(b.Velocity)
	b.Acceleration = Vector2D{X: 0, Y: 0}

	w := float64(screenWidth)
	h := float64(screenHeight)

	if b.Position.X < 0 {
		b.Position.X += w
	} else if b.Position.X >= w {
		b.Position.X -= w
	}
	if b.Position.Y < 0 {
		b.Position.Y += h
	} else if b.Position.Y >= h {
		b.Position.Y -= h
	}
}

func (b *Boid) Draw(screen tcell.Screen) {
	x := int(b.Position.X)
	y := int(b.Position.Y)

	predatorStyle := tcell.StyleDefault.Foreground(tcell.ColorRed).Bold(true)
	normalStyle := tcell.StyleDefault.Foreground(tcell.ColorGreen).Bold(true)

	if b.Type == PredatorBoid {
		screen.SetContent(x, y, 'o', nil, predatorStyle)
	} else {
		screen.SetContent(x, y, DirectionNormalBoid(b.Velocity.X, b.Velocity.Y), nil, normalStyle)
	}
}

func DirectionNormalBoid(x, y float64) rune {
	absX := math.Abs(x)
	absY := math.Abs(y)

	if absX > absY {
		if x > 0 {
			return '>'
		}
		return '<'
	}

	if y > 0 {
		return 'v'
	}
	return '^'
}

func (b *Boid) Separation(flock []*Boid) Vector2D {
	var steering Vector2D
	total := 0

	for _, other := range flock {
		if b == other {
			continue
		}

		distance := b.Position.Subtract(other.Position).Magnitude()
		if distance > 0 && distance < PerceptionRadius {
			esc := b.Position.Subtract(other.Position)
			esc = esc.Normalize()
			esc = esc.Divide(distance)
			steering = steering.Add(esc)

			total++
		}
	}

	if total > 0 {
		steering = steering.Divide(float64(total))

		if steering.Magnitude() > MaxForce {
			steering = steering.Normalize().Multiply(MaxForce)
		}
	}

	return steering
}

func (b *Boid) Alignment(flock []*Boid) Vector2D {
	var averageVelocity Vector2D
	total := 0

	for _, other := range flock {
		if b == other {
			continue
		}

		distance := b.Position.Subtract(other.Position).Magnitude()

		if distance > 0 && distance < PerceptionRadius {
			averageVelocity = averageVelocity.Add(other.Velocity)
			total++
		}
	}

	var steering Vector2D
	if total > 0 {
		averageVelocity = averageVelocity.Divide(float64(total))
		averageVelocity = averageVelocity.Normalize().Multiply(MaxSpeed)
		steering = averageVelocity.Subtract(b.Velocity)

		if steering.Magnitude() > MaxForce {
			steering = steering.Normalize().Multiply(MaxForce)
		}
	}

	return steering
}

func (b *Boid) Cohesion(flock []*Boid) Vector2D {
	var centerOfMass Vector2D
	total := 0

	for _, other := range flock {
		if b == other {
			continue
		}

		distance := b.Position.Subtract(other.Position).Magnitude()

		if distance > 0 && distance < PerceptionRadius {
			centerOfMass = centerOfMass.Add(other.Position)
			total++
		}
	}

	var steering Vector2D
	if total > 0 {
		centerOfMass = centerOfMass.Divide(float64(total))
		desiredVelocity := centerOfMass.Subtract(b.Position)
		desiredVelocity = desiredVelocity.Normalize().Multiply(MaxSpeed)
		steering = desiredVelocity.Subtract(b.Velocity)

		if steering.Magnitude() > MaxForce {
			steering = steering.Normalize().Multiply(MaxForce)
		}
	}

	return steering
}

func (b *Boid) Chase(flock []*Boid) Vector2D {
	var steering Vector2D
	total := 0

	for _, other := range flock {
		if b == other || other.Type != NormalBoid {
			continue
		}

		distance := b.Position.Subtract(other.Position).Magnitude()
		if distance > 0 && distance < PredatorPerceptionRadius {
			desiredVelocity := other.Position.Subtract(b.Position)
			desiredVelocity = desiredVelocity.Normalize().Multiply(PredatorMaxSpeed)
			steering = steering.Add(desiredVelocity.Subtract(b.Velocity))

			total++
		}
	}

	if total > 0 {
		steering = steering.Divide(float64(total))

		if steering.Magnitude() > PredatorMaxForce {
			steering = steering.Normalize().Multiply(PredatorMaxForce)
		}
	}

	return steering
}

func (b *Boid) Flee(flock []*Boid) Vector2D {
	var steering Vector2D
	total := 0

	for _, other := range flock {
		if b == other || other.Type != PredatorBoid {
			continue
		}

		distance := b.Position.Subtract(other.Position).Magnitude()
		if distance > 0 && distance < PanicRadius {
			desiredVelocity := b.Position.Subtract(other.Position)
			desiredVelocity = desiredVelocity.Normalize().Multiply(MaxSpeed * 1.5)
			steering = steering.Add(desiredVelocity.Subtract(b.Velocity))

			total++
		}
	}

	if total > 0 {
		steering = steering.Divide(float64(total))

		if steering.Magnitude() > MaxForce {
			steering = steering.Normalize().Multiply(MaxForce)
		}
	}

	return steering
}

func (b *Boid) Eat(flock []*Boid) {
	for _, other := range flock {
		if b == other || other.Type != NormalBoid {
			continue
		}

		distance := b.Position.Subtract(other.Position).Magnitude()
		if distance < 1.0 {
			other.IsDead = true
			break
		}
	}
}
