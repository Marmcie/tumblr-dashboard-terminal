package npf
type Media struct {
	Type   string
	Url    string
	Width  int64
	Height int64
	Hd     bool
}
func (m *Media) Render() string {
	var str = ""
	str += "![Image]("
	str += m.Url
	str += ")"
	return str
}
