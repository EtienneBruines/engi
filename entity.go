package engi

import (
	"reflect"
)

type Entity struct {
	id         string
	components []Component
	requires   map[string]bool
	Pattern    string
}

func NewEntity(requires []string) *Entity {
	e := &Entity{
		id:       generateUUID(),
		requires: make(map[string]bool),
	}
	for _, req := range requires {
		e.requires[req] = true
	}
	return e
}

func (e *Entity) DoesRequire(name string) bool {
	return e.requires[name]
}

func (e *Entity) AddComponent(component Component) {
	e.components = append(e.components, component)
}

func (e *Entity) RemoveComponent(component Component) {
	for index := range e.components {
		if e.components[index].Type() == component.Type() {
			e.components = append(e.components[:index], e.components[index+1:]...)
		}
	}
}

// GetComponent takes a double pointer to a Component,
// and populates it with the value of the right type.
func (e *Entity) Component(x interface{}) bool {
	v := reflect.ValueOf(x).Elem() // *T
	c := e.ComponentFast((v.Interface().(Component)))
	if c == nil {
		return false
	}
	v.Set(reflect.ValueOf(c))
	return true
}

// ComponentFast returns the same object as GetComponent
// but without using reflect (and thus faster)
// Be sure to define the .Type() such that it takes a pointer receiver
func (e *Entity) ComponentFast(c Component) interface{} {
	for _, comp := range e.components {
		if comp.Type() == c.Type() {
			return comp
		}
	}
	return nil
}

func (e *Entity) ID() string {
	return e.id
}
