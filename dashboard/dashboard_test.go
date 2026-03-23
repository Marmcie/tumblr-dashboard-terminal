package dashboard_test

import (
	"testing"
	"time"
	"tumblr-dt/dashboard"
	"tumblr-dt/modules"

	tea "charm.land/bubbletea/v2"
)

func TestDashboardLoad(t *testing.T) {
	config := modules.Config{}
	config.Testing = true
	config.Initialized = true
	ch := make(chan *dashboard.Dashboard)
	go func(ch chan *dashboard.Dashboard) {
		ch <- dashboard.NewDashboard(config)
	}(ch)
	dashboard := <-ch
	if dashboard == nil {
		t.Errorf("a")
	}
}

func TestDashboardDisplay(t *testing.T) {
	config := modules.Config{}
	config.Testing = true
	config.Initialized = true

	ch := make(chan *dashboard.Dashboard)
	go func(ch chan *dashboard.Dashboard) {
		db := dashboard.NewDashboard(config)
		ch <- db
	}(ch)
	dashboard := <-ch
	for dashboard.IsLoading {
		time.Sleep(time.Second / 2)
	}
	dashboard.DisplayPost(dashboard.GetSelectedPost(), true)
}

func BenchmarkDashboardLoad(b *testing.B) {
	config := modules.Config{}
	config.Testing = true
	config.Initialized = true
	ch := make(chan *dashboard.Dashboard)
	go func(ch chan *dashboard.Dashboard) {
		db := dashboard.NewDashboard(config)
		ch <- db
	}(ch)
	dashboard := <-ch
	for dashboard.IsLoading {
		time.Sleep(time.Second / 2)
	}
	for b.Loop() {
		dashboard.DisplayPost(dashboard.GetSelectedPost(), true)
	}
}
func BenchmarkDashboardUpdate(b *testing.B) {
	config := modules.Config{}
	config.Testing = true
	config.Initialized = true

	ch := make(chan *dashboard.Dashboard)
	go func(ch chan *dashboard.Dashboard) {
		ch <- dashboard.NewDashboard(config)
	}(ch)
	dashboard := <-ch

	msg := tea.KeyPressMsg{}
	for b.Loop() {
		dashboard.GetRootModel().Update(msg)
	}
}
