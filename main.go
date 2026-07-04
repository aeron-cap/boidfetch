package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/gdamore/tcell/v2"
)

const infoAreaWidth = 67

func main() {
	screen, err := tcell.NewScreen()
	if err != nil {
		log.Fatalf("Failed to create screen: %v", err)
	}
	if err := screen.Init(); err != nil {
		log.Fatalf("Failed to initialize screen: %v", err)
	}
	defer screen.Fini()

	screen.EnableMouse()
	info := GatherInfo()

	width, height := screen.Size()
	boidPaneWidth := width - infoAreaWidth 

	var flock []*Boid
	for range 100 {
		b := &Boid{
			Position: Vector2D{X: rand.Float64() * float64(boidPaneWidth), Y: rand.Float64() * float64(height)},
			Velocity: Vector2D{X: (rand.Float64() * 2) - 1, Y: (rand.Float64() * 2) - 1},
			Type:     NormalBoid,
		}
		flock = append(flock, b)
	}

	eventChan := make(chan tcell.Event)
	go func() {
		for {
			eventChan <- screen.PollEvent()
		}
	}()

	ticker := time.NewTicker(time.Second / 30)
	defer ticker.Stop()

	hasUi := true

	for {
		select {
		case ev := <-eventChan:
			switch ev := ev.(type) {
			case *tcell.EventKey:
				if ev.Key() == tcell.KeyEscape || ev.Key() == tcell.KeyCtrlC {
					return
				}

				if ev.Key() == tcell.KeyRune && ev.Rune() == 'r' {
					flock = nil
					for range 10 {
						b := &Boid{
							Position: Vector2D{X: rand.Float64() * float64(boidPaneWidth), Y: rand.Float64() * float64(height)},
							Velocity: Vector2D{X: (rand.Float64() * 2) - 1, Y: (rand.Float64() * 2) - 1},
							Type:     NormalBoid,
						}
						flock = append(flock, b)
					}
				}

				if ev.Key() == tcell.KeyRune && ev.Rune() == 'h' {
					hasUi = !hasUi
				}
				
			case *tcell.EventMouse:
				x, y := ev.Position()
				button := ev.Buttons()
				if x > boidPaneWidth{
					continue
				}

				if len(flock) >= 1000 {
					continue
				}
				
				relX := x
				if button == tcell.Button1 {
					newBoid := &Boid{
						Position: Vector2D{X: float64(relX), Y: float64(y)},
						Velocity: Vector2D{X: (rand.Float64() * 2) - 1, Y: (rand.Float64() * 2) - 1},
					}

					flock = append(flock, newBoid)
				}

				if button == tcell.Button2 {
					newBoid := &Boid{
						Position: Vector2D{X: float64(relX), Y: float64(y)},
						Velocity: Vector2D{X: (rand.Float64() * 2) - 1, Y: (rand.Float64() * 2) - 1},
						Type:     PredatorBoid,
					}

					flock = append(flock, newBoid)
				}
			}

		case <-ticker.C:
			screen.Clear()
			width, height = screen.Size()
			boidPaneWidth = width - infoAreaWidth
			for _, b := range flock {
				if b.IsDead {
					continue
				}
				
				b.Update(boidPaneWidth, height, flock)
				b.Draw(screen)
			}

			var survivors []*Boid
			for _, b := range flock {
				if !b.IsDead {
					survivors = append(survivors, b)
				}
			}
			flock = survivors

			drawInfo(screen, info, width - infoAreaWidth, height)

			if hasUi {
				uiStyle := tcell.StyleDefault.Foreground(tcell.ColorYellow)
				count := fmt.Sprintf("Normal Boids: %v | Predators: %v", getNormalCount(flock), len(flock)-getNormalCount(flock))
				DrawText(screen, 1, 0, uiStyle, count)
				controls := "Left Click: Spawn Boid | Right Click: Spawn Predator | Esc: Quit | r: Reset | h: hide ui"
				DrawText(screen, 1, height - 1, uiStyle, controls)
			}

			screen.Show()
		}
	}
}

func drawInfo(screen tcell.Screen, info []InfoLine, width, height int) {
	labelStyle := tcell.StyleDefault.Foreground(tcell.ColorYellow).Bold(true)
	valueStyle := tcell.StyleDefault.Foreground(tcell.ColorWhite)

	row := height/2 - len(info)/2
	if row < 0 {
		row = 0
	}

	for _, line := range info {
		DrawText(screen, width + 4, row, labelStyle, line.Label)
		DrawText(screen, width + 16, row, valueStyle, line.Value)
		row++
	}
}

func DrawText(screen tcell.Screen, x, y int, style tcell.Style, text string) {
	for i, char := range text {
		screen.SetContent(x+i, y, char, nil, style)
	}
}

func getNormalCount(flock []*Boid) (normalCount int) {
	for _, b := range flock {
		if b.Type == NormalBoid {
			normalCount++
		} 
	}
	return
}
