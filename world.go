package engi

import (
	"runtime"
	"sync"
)

type World struct {
	Game
	entities map[string]*Entity
	systems  []Systemer

	defaultBatch *Batch
	hudBatch     *Batch

	isSetup bool
	paused  bool

	entitiesMu sync.RWMutex
}

func (w *World) New() {
	if !w.isSetup {
		w.entities = make(map[string]*Entity)
		if !headless {
			w.defaultBatch = NewBatch(Width(), Height(), batchVert, batchFrag)
			w.hudBatch = NewBatch(Width(), Height(), hudVert, hudFrag)
		}

		// Default WorldBounds values
		if WorldBounds.Max.X == 0 && WorldBounds.Max.Y == 0 {
			WorldBounds.Max = Point{Width(), Height()}
		}

		// Initialize cameraSystem
		cam = &cameraSystem{}
		w.AddSystem(cam)

		w.isSetup = true
	}
}

func (w *World) AddEntity(entity *Entity) {
	w.entities[entity.ID()] = entity

	for _, system := range w.systems {
		if entity.DoesRequire(system.Type()) {
			system.AddEntity(entity)
		}
	}
}

func (w *World) RemoveEntity(entity *Entity) {
	w.entitiesMu.Lock()
	delete(w.entities, entity.ID())
	w.entitiesMu.Unlock()
}

func (w *World) AddSystem(system Systemer) {
	system.New()
	w.entitiesMu.Lock()
	w.systems = append(w.systems, system)
	w.entitiesMu.Unlock()
}

func (w *World) Entities() []*Entity {
	entities := make([]*Entity, len(w.entities))
	w.entitiesMu.RLock()
	for _, v := range w.entities {
		entities = append(entities, v)
	}
	w.entitiesMu.RUnlock()
	return entities
}

func (w *World) Systems() []Systemer {
	return w.systems
}

func (w *World) Pre() {
	if !headless {
		Gl.Clear(Gl.COLOR_BUFFER_BIT)
	}
}

func (w *World) Post() {}

func (w *World) Update(dt float32) {
	w.Pre()

	var unp *UnpauseComponent

	for _, system := range w.Systems() {
		if headless && system.SkipOnHeadless() {
			continue // so skip it
		}

		system.Pre()

		entities := system.Entities()

		// It's not always faster to multithread; so in this case we're not going to
		if len(entities) < 2*runtime.NumCPU() {
			for _, ent := range entities {
				system.Update(ent, dt)
			}
			system.Post()
			continue // with other Systems
		}

		entityChannel := make(chan *Entity)
		wg := sync.WaitGroup{}

		// Launch workers
		for i := 0; i < runtime.NumCPU(); i++ {
			go func() {
				for ent := range entityChannel {
					system.Update(ent, dt)
					wg.Done()
				}
			}()
		}

		// Give them something to do
		for _, entity := range entities {
			if w.paused {
				ok := entity.GetComponent(&unp)
				if !ok {
					continue // so skip it
				}
			}
			if entity.Exists {
				wg.Add(1)
				entityChannel <- entity
			}
		}

		// Wait until they're done, before continuing to other Systems
		wg.Wait()
		close(entityChannel)

		system.Post()
	}

	if Keys.KEY_ESCAPE.JustPressed() {
		Exit()
	}

	w.Post()
}

func (w *World) Batch(prio PriorityLevel) *Batch {
	if prio >= HUDGround {
		return w.hudBatch
	} else {
		return w.defaultBatch
	}
}
