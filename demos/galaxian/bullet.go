package main

import (
	"github.com/paked/engi"
	//"log"
)

type BulletSystem struct {
	*engi.System
}

func (BulletSystem) Type() string {
	return "BulletSystem"
}

func (c *BulletSystem) New() {
	c.System = &engi.System{}
}

func (c *BulletSystem) Update(bullet *engi.Entity, dt float32) {
	var space *engi.SpaceComponent

	if !bullet.GetComponent(&space) {
		return
	}

	//log.Println(len(game.Entities()))
	if space.Position.Y < 0 {
		defer game.RemoveEntity(bullet)
	}
}
