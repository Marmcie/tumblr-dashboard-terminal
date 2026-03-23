package component

import tea "charm.land/bubbletea/v2"

type GlobalValues struct {
	Msg             tea.Msg
	Time            int
	Elements        []Component
	EventDispatches map[string]map[string]map[string]func(tea.Msg, int)
	Command         tea.Cmd
	Logger          []func() string
}

var Global = &GlobalValues{
	Elements: []Component{},
}

func (g *GlobalValues) AddElement(c Component) int {
	g.Elements = append(g.Elements, c)
	return len(g.Elements) - 1
}
func (g *GlobalValues) DeleteElement(i int) {
	g.Elements[i] = g.Elements[len(g.Elements)-1]
	g.Elements[i].SetGlobalIndex(i)
	g.Elements = g.Elements[:len(g.Elements)-1]
}

func (g *GlobalValues) PrintLog() {
	f, _ := tea.LogToFile("debug.log", "debug")
	for _, cb := range g.Logger {
		f.WriteString(cb())
	}
	f.Close()
}

func (g *GlobalValues) SubscribeLogger(cb func() string) {
	g.Logger = append(g.Logger, cb)
}

func (g *GlobalValues) BlurAll() {
	for _, v := range g.Elements {
		v.Blur()
	}
}

func (g *GlobalValues) SetCmd(cmd tea.Cmd) {
	g.Command = cmd
}

func (g *GlobalValues) AddEventCallback(event string, uuid string, callbackUUID string, cb func(tea.Msg, int)) {
	if g.EventDispatches == nil {
		g.EventDispatches = map[string]map[string]map[string]func(tea.Msg, int){}
	}
	if g.EventDispatches[event] == nil {
		g.EventDispatches[event] = map[string]map[string]func(tea.Msg, int){}
	}

	if g.EventDispatches[event][uuid] == nil {
		g.EventDispatches[event][uuid] = map[string]func(tea.Msg, int){}
	}
	g.EventDispatches[event][uuid][callbackUUID] = cb
}

func (g *GlobalValues) CallEvents() {
	for _, v := range g.EventDispatches {
		for _, cbs := range v {
			for _, cb := range cbs {
				cb(g.Msg, g.Time)
			}
		}
	}
	g.EventDispatches = map[string]map[string]map[string]func(tea.Msg, int){}
}

func UpdateGlobalValues(msg tea.Msg, time int) {
	Global.Msg = msg
	Global.Time = time
}
