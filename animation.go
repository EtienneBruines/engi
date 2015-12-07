package engi

import (
	"log"

	"github.com/paked/engi/ecs"
)

// AnimationAction defines the sequence of indexes beloning to a single (graphical) action
type AnimationAction struct {
	Name   string
	Frames []int
}

// AnimationComponent is the `Component` that controls animation in rendering entities
type AnimationComponent struct {
	index            int              // What frame in the is being used
	Rate             float32          // How often frames should increment, in seconds.
	change           float32          // The time since the last incrementation
	Renderables      []Renderable     // Renderables
	Animations       map[string][]int // All possible animations
	CurrentAnimation []int            // The current animation
}

// NewAnimationComponent returns a pointer to a newly created `AnimationComponent`
func NewAnimationComponent(renderables []Renderable, rate float32) *AnimationComponent {
	return &AnimationComponent{
		Animations:  make(map[string][]int),
		Renderables: renderables,
		Rate:        rate,
	}
}

// SelectAnimationByName sets the `CurrentAnimation` to the `AnimationComponent`, with given name
func (ac *AnimationComponent) SelectAnimationByName(name string) {
	ac.CurrentAnimation = ac.Animations[name]
}

// SelectAnimationByAction sets the `CurrentAnimation` to the `AnimationComponent`, with the same name as `action`
func (ac *AnimationComponent) SelectAnimationByAction(action *AnimationAction) {
	ac.CurrentAnimation = ac.Animations[action.Name]
}

// AddAnimationAction sets the `Frames` of the `action` to the 'Animation' within the `AnimationComponent`
func (ac *AnimationComponent) AddAnimationAction(action *AnimationAction) {
	ac.Animations[action.Name] = action.Frames
}

// AddAnimationActions adds a number of `AnimationAction`s to the `AnimationComponent`
func (ac *AnimationComponent) AddAnimationActions(actions []*AnimationAction) {
	for _, action := range actions {
		ac.Animations[action.Name] = action.Frames
	}
}

// Cell returns the Renderable version of the `CurrentAnimation`
func (ac *AnimationComponent) Cell() Renderable {
	idx := ac.CurrentAnimation[ac.index]

	return ac.Renderables[idx]
}

// Type returns the (unique) string representation of `AnimationComponent`
func (*AnimationComponent) Type() string {
	return "AnimationComponent"
}

// AnimationSystem is a `System` that manages animations
type AnimationSystem struct {
	*ecs.System
}

// New initializes a new `AnimationSystem`
func (a *AnimationSystem) New(*ecs.World) {
	a.System = ecs.NewSystem()
}

// Type returns the (unique) string representation of `AnimationSystem`
func (AnimationSystem) Type() string {
	return "AnimationSystem"
}

// Update is called for each `ecs.Entity`, to allow `AnimationSystem` to update
func (a *AnimationSystem) Update(e *ecs.Entity, dt float32) {
	var (
		ac *AnimationComponent
		r  *RenderComponent
		ok bool
	)

	if ac, ok = e.ComponentFast(ac).(*AnimationComponent); !ok {
		return
	}
	if r, ok = e.ComponentFast(r).(*RenderComponent); !ok {
		return
	}

	ac.change += dt
	if ac.change >= ac.Rate {
		a.NextFrame(ac)
		r.Display = ac.Cell()
	}
}

// NextFrame updates the `CurrentAnimation` to move to the next frame
func (a *AnimationSystem) NextFrame(ac *AnimationComponent) {
	if len(ac.CurrentAnimation) == 0 {
		log.Println("No data for this animation")
		return
	}

	ac.index++
	if ac.index >= len(ac.CurrentAnimation) {
		ac.index = 0
	}
	ac.change = 0
}
