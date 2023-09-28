package models

import "github.com/faiface/pixel"

type Player struct {
	Body pixel.Rect
}

func NewPlayer(body pixel.Rect) *Player {
	return &Player{
		Body: body,
	}
}
