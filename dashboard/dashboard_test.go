package dashboard_test

import (
	"testing"
	"tumblr-dt/dashboard"
	"tumblr-dt/modules"

	tea "charm.land/bubbletea/v2"
)

func TestDashboardLoad(t *testing.T) {
	config := modules.Config{}
	config.Testing = true
	dash := dashboard.NewDashboard(config)
	if dash == nil {
		t.Errorf("a")
	}
}

func TestDashboardDisplay(t *testing.T) {
	config := modules.Config{}
	config.Testing = true
	dash := dashboard.NewDashboard(config)
	dash.DisplayPost(dash.GetSelectedPost(),true)
}

func BenchmarkDashboardLoad(b *testing.B) {
	config := modules.Config{}
	config.Testing = true
	dashboard := dashboard.NewDashboard(config)
	for b.Loop() {
		dashboard.DisplayPost(dashboard.GetSelectedPost(),true)
	}
}
func BenchmarkDashboardUpdate(b *testing.B) {
	config := modules.Config{}
	config.Testing = true
	dashboard := dashboard.NewDashboard(config)
	msg := tea.KeyPressMsg{}
	for b.Loop() {
		dashboard.GetRootModel().Update(msg)
	}
}
