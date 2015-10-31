package main

import (
	"github.com/paked/engi"
	"log"
)

type GalaxianGame struct {
	engi.World
}

var (
	game         *GalaxianGame
	scale        engi.Point
	basicFont    *engi.Font
	shipSprite   *engi.Region
	bulletSprite *engi.Region
)

func (galaxianGame *GalaxianGame) Preload() {
	galaxianGame.New()
	engi.Files.Add("assets/galaxian.png")

	basicFont = (&engi.Font{URL: "assets/Roboto-Regular.ttf", Size: 32, FG: engi.Color{255, 255, 255, 255}})
	if err := basicFont.Create(); err != nil {
		log.Fatalln("Could not load font:", err)
	}
}

func (galaxianGame *GalaxianGame) Setup() {
	engi.SetBg(0x000000)
	scale = engi.Point{10, 10}

	galaxianGame.AddSystem(&engi.RenderSystem{})
	galaxianGame.AddSystem(&ControlSystem{})
	galaxianGame.AddSystem(&BulletSystem{})
	galaxianGame.AddSystem(&SpeedSystem{})

	galaxianSprite := engi.Files.Image("galaxian.png")

	shipSprite = engi.NewRegion(galaxianSprite, 3, 70, 15, 15)
	bulletSprite = engi.NewRegion(galaxianSprite, 70, 65, 1, 3)

	shipPoint := engi.Point{engi.Width() / 2, engi.Height()}
	ship := engi.NewEntity([]string{"RenderSystem", "ControlSystem"})
	shipRender := engi.NewRenderComponent(shipSprite, scale, "ship")
	shipSpace := &engi.SpaceComponent{shipPoint, scale.X * shipSprite.Width(), scale.Y * shipSprite.Height()}

	ship.AddComponent(shipRender)
	ship.AddComponent(shipSpace)

	galaxianGame.AddEntity(ship)
	game = galaxianGame
}

func main() {
	engi.Open("GalaxianGame", 1024, 768, false, &GalaxianGame{})
}
