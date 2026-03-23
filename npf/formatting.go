package npf
type Formatting struct {
	Start int64
	End   int64
	Type  string
	Url   string
	Blog  struct {
		Uuid string
		Name string
		Url  string
	}

	Hex string
}
