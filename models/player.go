package models

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

type Player struct {
	Body   pixel.Rect
	Speed  float64
	MoveCh chan float64
	Score  int
}

func NewPlayer(body pixel.Rect, speed float64) *Player {
	return &Player{
		Body:   body,
		Speed:  speed,
		MoveCh: make(chan float64),
		Score:  0,
	}
}

func (p *Player) Move(win *pixelgl.Window, gameOver bool, windowHeight float64, playerNumber int) {
	for !gameOver {
		if win.Pressed(pixelgl.KeyW) && p.Body.Max.Y < windowHeight && playerNumber == 1 {
			p.MoveCh <- p.Speed
		}
		if win.Pressed(pixelgl.KeyS) && p.Body.Min.Y > 0 && playerNumber == 1 {
			p.MoveCh <- -p.Speed
		}

		if win.Pressed(pixelgl.KeyUp) && p.Body.Max.Y < windowHeight && playerNumber == 2 {
			p.MoveCh <- p.Speed
		}
		if win.Pressed(pixelgl.KeyDown) && p.Body.Min.Y > 0 && playerNumber == 2 {
			p.MoveCh <- -p.Speed
		}
	}
}

func (p *Player) UpdateScore(score int) {
	p.Score = score
}
