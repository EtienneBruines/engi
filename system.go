package engi

type Systemer interface {
	Update(entity *Entity, dt float32)
	Type() string
	Priority() int
	Pre()
	Post()
	New()
	Entities() []*Entity
	AddEntity(entity *Entity)
	RemoveEntity(entity *Entity)
	SkipOnHeadless() bool
	SetWorld(*World)
	// Push(message Message)
	// Receive(message Message)
	// Messages() []Message
	// Dismiss(i int)
}

type System struct {
	entities             map[string]*Entity
	messageQueue         []Message
	ShouldSkipOnHeadless bool
	World                *World
}

func NewSystem() *System {
	s := &System{}
	s.entities = make(map[string]*Entity)
	return s
}

func (s System) New()  {}
func (s System) Pre()  {}
func (s System) Post() {}

func (s System) Priority() int {
	return 0
}

func (s System) Entities() []*Entity {
	list := make([]*Entity, len(s.entities))
	i := 0
	for _, ent := range s.entities {
		list[i] = ent
		i++
	}
	return list
}

func (s *System) AddEntity(entity *Entity) {
	s.entities[entity.ID()] = entity
}

func (s *System) RemoveEntity(entity *Entity) {
	delete(s.entities, entity.ID())
}

func (s System) SkipOnHeadless() bool {
	return s.ShouldSkipOnHeadless
}

func (s *System) SetWorld(w *World) {
	s.World = w
}

type CollisionSystem struct {
	*System
}

func (cs *CollisionSystem) New() {
	cs.System = NewSystem()
}

func (cs *CollisionSystem) Update(entity *Entity, dt float32) {
	var (
		space     *SpaceComponent
		collision *CollisionComponent
		ok        bool
	)
	if space, ok = entity.ComponentFast("*engi.SpaceComponent").(*SpaceComponent); !ok {
		return
	}
	if collision, ok = entity.ComponentFast("*engi.CollisionComponent").(*CollisionComponent); !ok {
		return
	}

	if !collision.Main {
		return
	}

	var otherSpace *SpaceComponent
	var otherCollision *CollisionComponent

	for _, other := range cs.Entities() {
		if other.ID() != entity.ID() {
			if otherSpace, ok = other.ComponentFast("*engi.SpaceComponent").(*SpaceComponent); !ok {
				continue // with other entities
			}
			if otherCollision, ok = other.ComponentFast("*engi.CollisionComponent").(*CollisionComponent); !ok {
				continue // with other entities
			}

			entityAABB := space.AABB()
			offset := Point{collision.Extra.X / 2, collision.Extra.Y / 2}
			entityAABB.Min.X -= offset.X
			entityAABB.Min.Y -= offset.Y
			entityAABB.Max.X += offset.X
			entityAABB.Max.Y += offset.Y
			otherAABB := otherSpace.AABB()
			offset = Point{otherCollision.Extra.X / 2, otherCollision.Extra.Y / 2}
			otherAABB.Min.X -= offset.X
			otherAABB.Min.Y -= offset.Y
			otherAABB.Max.X += offset.X
			otherAABB.Max.Y += offset.Y
			if IsIntersecting(entityAABB, otherAABB) {
				if otherCollision.Solid && collision.Solid {
					mtd := MinimumTranslation(entityAABB, otherAABB)
					space.Position.X += mtd.X
					space.Position.Y += mtd.Y
				}

				Mailbox.Dispatch(CollisionMessage{Entity: entity, To: other})
			}
		}
	}
}

func (*CollisionSystem) Type() string {
	return "CollisionSystem"
}

type NilSystem struct {
	*System
}

func (ns *NilSystem) New() {
	ns.System = NewSystem()
}

func (*NilSystem) Update(*Entity, float32) {}

func (*NilSystem) Type() string {
	return "NilSystem"
}

type PriorityLevel int

const (
	// HighestGround is the highest PriorityLevel that will be rendered
	HighestGround PriorityLevel = 50
	// HUDGround is a PriorityLevel from which everything isn't being affected by the Camera
	HUDGround    PriorityLevel = 40
	Foreground   PriorityLevel = 30
	MiddleGround PriorityLevel = 20
	ScenicGround PriorityLevel = 10
	// Background is the lowest PriorityLevel that will be rendered
	Background PriorityLevel = 0
	// Hidden indicates that it should not be rendered by the RenderSystem
	Hidden PriorityLevel = -1
)

type RenderSystem struct {
	renders map[PriorityLevel][]*Entity
	changed bool
	*System
}

func (rs *RenderSystem) New() {
	rs.renders = make(map[PriorityLevel][]*Entity)
	rs.System = NewSystem()
	rs.ShouldSkipOnHeadless = true

	Mailbox.Listen("renderChangeMessage", func(m Message) {
		rs.changed = true
	})
}

func (rs *RenderSystem) AddEntity(e *Entity) {
	rs.changed = true
	rs.System.AddEntity(e)
}

func (rs *RenderSystem) RemoveEntity(e *Entity) {
	rs.changed = true
	rs.System.RemoveEntity(e)
}

func (rs *RenderSystem) Pre() {
	if !rs.changed {
		return
	}

	rs.renders = make(map[PriorityLevel][]*Entity)
}

type Renderable interface {
	Render(b *Batch, render *RenderComponent, space *SpaceComponent)
}

func (rs *RenderSystem) Post() {
	var currentBatch *Batch

	for i := Background; i <= HighestGround; i++ {
		if len(rs.renders[i]) == 0 {
			continue
		}

		// Retrieve a batch, may be the default one -- then call .Begin() if we arent already using it
		batch := world.batch(i)
		if batch != currentBatch {
			if currentBatch != nil {
				currentBatch.End()
			}
			batch.Begin()
			currentBatch = batch
		}
		// Then render everything for this level
		for _, entity := range rs.renders[i] {
			var (
				render *RenderComponent
				space  *SpaceComponent
				ok     bool
			)

			if render, ok = entity.ComponentFast("*engi.RenderComponent").(*RenderComponent); !ok {
				continue
			}
			if space, ok = entity.ComponentFast("*engi.SpaceComponent").(*SpaceComponent); !ok {
				continue
			}

			render.Display.Render(batch, render, space)
		}
	}

	if currentBatch != nil {
		currentBatch.End()
	}

	rs.changed = false
}

func (rs *RenderSystem) Update(entity *Entity, dt float32) {
	if !rs.changed {
		return
	}

	var render *RenderComponent
	var ok bool
	if render, ok = entity.ComponentFast("*engi.RenderComponent").(*RenderComponent); !ok {
		return
	}

	rs.renders[render.priority] = append(rs.renders[render.priority], entity)
}

func (*RenderSystem) Type() string {
	return "RenderSystem"
}

func (rs *RenderSystem) Priority() int {
	return 1
}
