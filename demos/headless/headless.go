package main

import (
	"fmt"
	"github.com/paked/engi"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"runtime/pprof"
	"sync"
)

type PongGame struct{}

var (
	basicFont *engi.Font
)

func (pong *PongGame) Preload() {
	engi.Files.Add("assets/ball.png", "assets/paddle.png")

	basicFont = (&engi.Font{URL: "assets/Roboto-Regular.ttf", Size: 32, FG: engi.Color{255, 255, 255, 255}})
	if err := basicFont.Create(); err != nil {
		log.Fatalln("Could not load font:", err)
	}
}

func (pong *PongGame) Setup(w *engi.World) {
	engi.SetBg(0x2d3739)
	w.AddSystem(&engi.RenderSystem{})
	w.AddSystem(&engi.CollisionSystem{})
	w.AddSystem(&SpeedSystem{})
	w.AddSystem(&ControlSystem{})
	w.AddSystem(&BallSystem{})
	w.AddSystem(&ScoreSystem{})

	ball := engi.NewEntity([]string{"RenderSystem", "CollisionSystem", "SpeedSystem", "BallSystem"})
	ballTexture := engi.Files.Image("ball.png")
	ballRender := engi.NewRenderComponent(ballTexture, engi.Point{2, 2}, "ball")
	ballSpace := &engi.SpaceComponent{engi.Point{(engi.Width() - ballTexture.Width()) / 2, (engi.Height() - ballTexture.Height()) / 2}, ballTexture.Width() * ballRender.Scale.X, ballTexture.Height() * ballRender.Scale.Y}
	ballCollision := &engi.CollisionComponent{Main: true, Solid: true}
	ballSpeed := &SpeedComponent{}
	ballSpeed.Point = engi.Point{300, 100}
	ball.AddComponent(ballRender)
	ball.AddComponent(ballSpace)
	ball.AddComponent(ballCollision)
	ball.AddComponent(ballSpeed)
	w.AddEntity(ball)

	score := engi.NewEntity([]string{"RenderSystem", "ScoreSystem"})

	scoreRender := engi.NewRenderComponent(basicFont.Render(" "), engi.Point{1, 1}, "YOLO <3")
	scoreSpace := &engi.SpaceComponent{engi.Point{100, 100}, 100, 100}
	score.AddComponent(scoreRender)
	score.AddComponent(scoreSpace)
	w.AddEntity(score)

	schemes := []string{"WASD", ""}
	for i := 0; i < 2; i++ {
		paddle := engi.NewEntity([]string{"RenderSystem", "CollisionSystem", "ControlSystem"})
		paddleTexture := engi.Files.Image("paddle.png")
		paddleRender := engi.NewRenderComponent(paddleTexture, engi.Point{2, 2}, "paddle")
		x := float32(0)
		if i != 0 {
			x = 800 - 16
		}
		paddleSpace := &engi.SpaceComponent{engi.Point{x, (engi.Height() - paddleTexture.Height()) / 2}, paddleRender.Scale.X * paddleTexture.Width(), paddleRender.Scale.Y * paddleTexture.Height()}
		paddleControl := &ControlComponent{schemes[i]}
		paddleCollision := &engi.CollisionComponent{Main: false, Solid: true}
		paddle.AddComponent(paddleRender)
		paddle.AddComponent(paddleSpace)
		paddle.AddComponent(paddleControl)
		paddle.AddComponent(paddleCollision)
		w.AddEntity(paddle)
	}
}

type SpeedSystem struct {
	*engi.System
}

func (ms *SpeedSystem) New() {
	ms.System = engi.NewSystem()
	engi.Mailbox.Listen("CollisionMessage", func(message engi.Message) {
		collision, isCollision := message.(engi.CollisionMessage)
		if isCollision {
			var speed *SpeedComponent
			var ok bool
			if speed, ok = collision.Entity.ComponentFast("*SpeedComponent").(*SpeedComponent); !ok {
				return
			}

			speed.X *= -1
		}
	})
}

func (*SpeedSystem) Type() string {
	return "SpeedSystem"
}

func (ms *SpeedSystem) Update(entity *engi.Entity, dt float32) {
	var speed *SpeedComponent
	var space *engi.SpaceComponent
	var ok bool

	if speed, ok = entity.ComponentFast("*SpeedComponent").(*SpeedComponent); !ok {
		return
	}
	if space, ok = entity.ComponentFast("*engi.SpaceComponent").(*engi.SpaceComponent); !ok {
		return
	}

	space.Position.X += speed.X * dt
	space.Position.Y += speed.Y * dt
}

func (ms *SpeedSystem) Receive(message engi.Message) {}

type SpeedComponent struct {
	engi.Point
}

func (SpeedComponent) Type() string {
	return "SpeedComponent"
}

type BallSystem struct {
	*engi.System
}

func (bs *BallSystem) New() {
	bs.System = engi.NewSystem()
}

func (BallSystem) Type() string {
	return "BallSystem"
}

func (bs *BallSystem) Update(entity *engi.Entity, dt float32) {
	var space *engi.SpaceComponent
	var speed *SpeedComponent
	var ok bool

	if space, ok = entity.ComponentFast("*engi.SpaceComponent").(*engi.SpaceComponent); !ok {
		return
	}
	if speed, ok = entity.ComponentFast("*SpeedComponent").(*SpeedComponent); !ok {
		return
	}

	if space.Position.X < 0 {
		engi.Mailbox.Dispatch(ScoreMessage{1})

		space.Position.X = 400 - 16
		space.Position.Y = 400 - 16
		speed.X = 800 * rand.Float32()
		speed.Y = 800 * rand.Float32()
	}

	if space.Position.Y < 0 {
		space.Position.Y = 0
		speed.Y *= -1
	}

	if space.Position.X > (800 - 16) {
		engi.Mailbox.Dispatch(ScoreMessage{2})

		space.Position.X = 400 - 16
		space.Position.Y = 400 - 16
		speed.X = 800 * rand.Float32()
		speed.Y = 800 * rand.Float32()
	}

	if space.Position.Y > (800 - 16) {
		space.Position.Y = 800 - 16
		speed.Y *= -1
	}
}

type ControlSystem struct {
	*engi.System
}

func (ControlSystem) Type() string {
	return "ControlSystem"
}
func (c *ControlSystem) New() {
	c.System = engi.NewSystem()
}

func (c *ControlSystem) Update(entity *engi.Entity, dt float32) {
	//Check scheme
	// -Move entity based on that
	var control *ControlComponent
	var space *engi.SpaceComponent
	var ok bool

	if control, ok = entity.ComponentFast("*ControlComponent").(*ControlComponent); !ok {
		return
	}
	if space, ok = entity.ComponentFast("*engi.SpaceComponent").(*engi.SpaceComponent); !ok {
		return
	}

	up := false
	down := false
	if control.Scheme == "WASD" {
		up = engi.Keys.KEY_W.Down()
		down = engi.Keys.KEY_S.Down()
	} else {
		up = engi.Keys.KEY_UP.Down()
		down = engi.Keys.KEY_DOWN.Down()
	}

	if up {
		space.Position.Y -= 800 * dt
	}

	if down {
		space.Position.Y += 800 * dt
	}

}

type ControlComponent struct {
	Scheme string
}

func (ControlComponent) Type() string {
	return "ControlComponent"
}

type ScoreSystem struct {
	*engi.System
	PlayerOneScore, PlayerTwoScore int
	upToDate                       bool
	scoreLock                      sync.RWMutex
}

func (ScoreSystem) Type() string {
	return "ScoreSystem"
}

func (sc *ScoreSystem) New() {
	sc.upToDate = true
	sc.System = engi.NewSystem()
	engi.Mailbox.Listen("ScoreMessage", func(message engi.Message) {
		scoreMessage, isScore := message.(ScoreMessage)
		if !isScore {
			return
		}

		sc.scoreLock.Lock()
		if scoreMessage.Player != 1 {
			sc.PlayerOneScore += 1
		} else {
			sc.PlayerTwoScore += 1
		}
		log.Println("The score is now", sc.PlayerOneScore, "vs", sc.PlayerTwoScore)
		sc.upToDate = false
		sc.scoreLock.Unlock()
	})
}

var sum float32

func (c *ScoreSystem) Update(entity *engi.Entity, dt float32) {
	var render *engi.RenderComponent
	var space *engi.SpaceComponent
	var ok bool

	sum += dt

	if sum > 1 {
		sum = 0
		fmt.Println(engi.Time.Fps())
	}

	if render, ok = entity.ComponentFast("*engi.RenderComponent").(*engi.RenderComponent); !ok {
		return
	}
	if space, ok = entity.ComponentFast("*engi.SpaceComponent").(*engi.SpaceComponent); !ok {
		return
	}

	if !c.upToDate {
		c.scoreLock.RLock()
		render.Label = fmt.Sprintf("%v vs %v", c.PlayerOneScore, c.PlayerTwoScore)
		c.upToDate = true
		c.scoreLock.RUnlock()

		render.Display = basicFont.Render(render.Label)
		width := len(render.Label) * 20

		space.Position.X = float32(400 - (width / 2))
	}
}

type ScoreMessage struct {
	Player int
}

func (ScoreMessage) Type() string {
	return "ScoreMessage"
}

const dir = "/tmp/cpu.out"

func main() {
	if dir != "" {
		f, err := os.Create(dir)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for sig := range c {
			fmt.Sprint(sig)
			fmt.Println("STOPPING")
			pprof.StopCPUProfile()
			fmt.Println("STOPPING")
			os.Exit(0)
		}
	}()

	defer func() {
		for sig := range c {
			fmt.Sprint(sig)
			pprof.StopCPUProfile()
			os.Exit(0)
		}
	}()

	engi.SetFPSLimit(0)
	engi.OpenHeadless(&PongGame{})
}
