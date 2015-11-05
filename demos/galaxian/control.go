package main

import (
	"github.com/paked/engi"
)

var (
	bulletCount uint64
)

type ControlSystem struct {
	*engi.System
}

func (ControlSystem) Type() string {
	return "ControlSystem"
}

func (c *ControlSystem) New() {
	c.System = &engi.System{}
}

func (c *ControlSystem) Update(entity *engi.Entity, dt float32) {
	left := engi.Keys.KEY_LEFT.Down()
	right := engi.Keys.KEY_RIGHT.Down()
	fire := engi.Keys.KEY_SPACE.JustReleased() || engi.Keys.KEY_SPACE.Down()

	if left {
		c.moveLeft(entity, dt)
	}

	if right {
		c.moveRight(entity, dt)
	}

	if fire {
		c.fire(entity)
	}

}

func (c *ControlSystem) moveLeft(ship *engi.Entity, dt float32) {
	var space *engi.SpaceComponent

	if !ship.GetComponent(&space) {
		return
	}

	newX := space.Position.X - 800*dt
	if newX >= 0 {
		space.Position.X = newX
	}
}

func (c *ControlSystem) moveRight(ship *engi.Entity, dt float32) {
	var space *engi.SpaceComponent

	if !ship.GetComponent(&space) {
		return
	}

	newX := space.Position.X + 800*dt
	if newX <= engi.Width() {
		space.Position.X = newX
	}
}

var counter = 0

func (c *ControlSystem) fire(ship *engi.Entity) {
	var space *engi.SpaceComponent

	if !ship.GetComponent(&space) {
		return
	}

	//fmt.Println(counter, len(game.World.Entities()))

	counter++

	bulletPoint := engi.Point{space.Position.X + 7*scale.X, space.Position.Y - 1*scale.Y}
	bullet := engi.NewEntity([]string{"RenderSystem", "SpeedSystem", "BulletSystem"})
	bulletCount += 1
	bulletRender := engi.NewRenderComponent(bulletSprite, scale, "bullet"+string(bulletCount))
	bulletSpace := &engi.SpaceComponent{bulletPoint, scale.X * bulletSprite.Width(), scale.Y * bulletSprite.Height()}
	bulletSpeed := &SpeedComponent{Point: engi.Point{X: 0, Y: scale.Y * -100}}

	bullet.AddComponent(bulletRender)
	bullet.AddComponent(bulletSpace)
	bullet.AddComponent(bulletSpeed)

	game.AddEntity(bullet)
}
