package npf

type Layout struct {
	Type           string
	Display        []map[string][]int64
	Truncate_after int64
	Blocks         []int64

	Attribution *struct {
		Type string
		Url  string
		Blog Blog
	}
}
