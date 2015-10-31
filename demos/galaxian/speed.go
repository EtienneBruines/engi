package main

import (
	"github.com/paked/engi"
)

type SpeedSystem struct {
	*engi.System
}

type SpeedComponent struct {
	engi.Point
}

func (SpeedComponent) Type() string {
	return "SpeedComponent"
}

func (ms *SpeedSystem) New() {
	ms.System = &engi.System{}
}

func (*SpeedSystem) Type() string {
	return "SpeedSystem"
}

func (ms *SpeedSystem) Update(entity *engi.Entity, dt float32) {
	var speed *SpeedComponent
	var space *engi.SpaceComponent
	if !entity.GetComponent(&speed) || !entity.GetComponent(&space) {
		return
	}
	space.Position.X += speed.X * dt
	space.Position.Y += speed.Y * dt
}

func (ms *SpeedSystem) Receive(message engi.Message) {
}
