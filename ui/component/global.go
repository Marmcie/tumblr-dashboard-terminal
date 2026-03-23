package component

import (
	"time"
	"tumblr-dt/ui/helper"

	tea "charm.land/bubbletea/v2"
)

// This class is there to keep track of the state information that should be accessible globally.
type GlobalValues struct {
	Msg             tea.Msg
	Elements        []Component
	EventDispatches map[string]map[string]map[string]func(tea.Msg)
	Command         tea.Cmd
	Logger          []func() string
	TickInterval    time.Duration
	IsSmall         bool
}

var Global = &GlobalValues{
	Elements:     []Component{},
	TickInterval: time.Second,
	IsSmall:      false,
}

// Add a component to the global list of all components
func (g *GlobalValues) AddElement(c Component) int {
	g.Elements = append(g.Elements, c)
	return len(g.Elements) - 1
}

// Remove a component from the global list of all components
func (g *GlobalValues) DeleteElement(i int) {
	g.Elements[i] = g.Elements[len(g.Elements)-1]
	g.Elements[i].SetGlobalIndex(i)
	g.Elements = g.Elements[:len(g.Elements)-1]
}

// Print all logs from callbacks subscribed to the global logger.
func (g *GlobalValues) PrintLog() {
	for _, cb := range g.Logger {
		helper.Log(cb())
	}
}

// Subscribe a callbacks to be outputted to the log file on command.
func (g *GlobalValues) SubscribeLogger(cb func() string) {
	g.Logger = append(g.Logger, cb)
}

// Remove focus from all components
func (g *GlobalValues) BlurAll() {
	for _, v := range g.Elements {
		v.Blur()
	}
}

// Set the tea.Cmd to be returned to the tea application
func (g *GlobalValues) SetCmd(cmd tea.Cmd) {
	g.Command = cmd
}

// Add a callback to a queue to be executed later in batch.
func (g *GlobalValues) AddEventCallback(event string, uuid string, callbackUUID string, cb func(tea.Msg)) {
	if g.EventDispatches == nil {
		g.EventDispatches = map[string]map[string]map[string]func(tea.Msg){}
	}
	if g.EventDispatches[event] == nil {
		g.EventDispatches[event] = map[string]map[string]func(tea.Msg){}
	}

	if g.EventDispatches[event][uuid] == nil {
		g.EventDispatches[event][uuid] = map[string]func(tea.Msg){}
	}
	g.EventDispatches[event][uuid][callbackUUID] = cb
}

// Execute all event callbacks
func (g *GlobalValues) CallEvents() {
	for _, v := range g.EventDispatches {
		for _, cbs := range v {
			for _, cb := range cbs {
				cb(g.Msg)
			}
		}
	}
	g.EventDispatches = map[string]map[string]map[string]func(tea.Msg){}
}

// Update globally accessible values
func UpdateGlobalValues(msg tea.Msg) {
	Global.Msg = msg
}
