package main

import "fmt"

type Event struct {
	label string
	data  string
}

type EventManager struct {
	dispatch chan Event
	elmap    []*EventListenerDict
}

func NewEventManager() *EventManager {
	dispatch := make(chan Event)
	elmap := make([]*EventListenerDict, 0)
	return &EventManager{
		dispatch: dispatch,
		elmap:    elmap,
	}
}

func (em *EventManager) FindEventListenerByLabel(label string) (*EventListenerDict, error) {
	for _, dict := range em.elmap {
		if dict.label == label {
			return dict, nil
		}
	}
	return nil, fmt.Errorf("not found")
}

func (em *EventManager) AddEventListener(label string, fn func(*Event)) {
	dict, err := em.FindEventListenerByLabel(label)

	if err != nil {
		newDict := &EventListenerDict{
			label: label,
			funcs: []func(*Event){fn},
		}
		em.elmap = append(em.elmap, newDict)
		return
	}
	dict.AddListener(fn)
}

func (em *EventManager) GetDispatchStream() chan Event {
	return em.dispatch
}

type EventListenerDict struct {
	label string
	funcs []func(*Event)
}

func (d *EventListenerDict) AddListener(fn func(*Event)) {
	d.funcs = append(d.funcs, fn)
}

type EventDispatcher struct {
	dispatch chan Event
}

func (em *EventManager) Run() {
	for event := range em.dispatch {
		for _, dict := range em.elmap {
			if dict.label == event.label {
				for _, fn := range dict.funcs {
					fn(&event)
				}
			}
		}
	}
}
