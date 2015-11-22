package engi

import (
	"reflect"
	"strings"
)

type Entity struct {
	id         string
	components map[string]Component
	requires   map[string]bool
	Pattern    string
}

func NewEntity(requires []string) *Entity {
	e := &Entity{
		id:         generateUUID(),
		requires:   make(map[string]bool),
		components: make(map[string]Component),
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
	e.components[component.Type()] = component
}

func (e *Entity) RemoveComponent(component Component) {
	delete(e.components, component.Type())
}

// GetComponent takes a double pointer to a Component,
// and populates it with the value of the right type.
func (e *Entity) Component(x interface{}) bool {
	v := reflect.ValueOf(x).Elem() // *T
	typeName := v.Type().String()
	dotIndex := strings.Index(typeName, ".")
	if dotIndex > 0 {
		typeName = typeName[dotIndex+1:]
	}
	c, ok := e.components[typeName]
	if !ok {
		return false
	}
	v.Set(reflect.ValueOf(c))
	return true
}

// ComponentFast returns the same object as GetComponent
// but without using reflect (and thus faster)
// Be sure to define the .Type() such that it takes a pointer receiver
func (e *Entity) ComponentFast(c Component) interface{} {
	return e.components[c.Type()]
}

func (e *Entity) ID() string {
	return e.id
}
