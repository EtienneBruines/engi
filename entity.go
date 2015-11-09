package engi

import (
	"reflect"
)

type Entity struct {
	id         string
	components []componentHolder
	requires   map[string]bool
	Pattern    string
}

type componentHolder struct {
	c  Component
	id string
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
	e.components = append(e.components, componentHolder{component, reflect.TypeOf(component).String()})
}

func (e *Entity) RemoveComponent(component Component) {
	index := -1
	str := reflect.TypeOf(component).String()
	for i := 0; i < len(e.components); i++ {
		if e.components[i].id == str {
			index = i
			break
		}
	}
	if index >= 0 {
		if index == len(e.components)-1 {
			e.components = e.components[:index]
		} else {
			e.components = append(e.components[:index], e.components[index+1:]...)
		}
	}
}

// GetComponent takes a double pointer to a Component,
// and populates it with the value of the right type.
func (e *Entity) GetComponent(x interface{}) bool {
	v := reflect.ValueOf(x).Elem() // *T
	c := e.ComponentFast(v.Type().String())
	if c == nil {
		return false
	}
	v.Set(reflect.ValueOf(c))
	return true
}

// ComponentFast returns the same object as GetComponent
// but without using reflect (and thus faster)
func (e *Entity) ComponentFast(t string) interface{} {
	for i := 0; i < len(e.components); i++ {
		if e.components[i].id == t {
			return e.components[i].c
		}
	}
	return nil
}

func (e *Entity) ID() string {
	return e.id
}
