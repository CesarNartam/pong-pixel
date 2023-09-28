package models

import "github.com/faiface/pixel"

type Ball struct {
	Body, Speed pixel.Vec
}

func NewBall(body pixel.Vec, speed pixel.Vec) *Ball {
	return &Ball{
		Body:  body,
		Speed: speed,
	}
}

func (b *Ball) Update(body pixel.Vec) {
	b.Body = body
}
