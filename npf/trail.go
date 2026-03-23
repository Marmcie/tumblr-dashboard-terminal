package npf
type TrailPost struct {
	Post struct {
		Id            string
		Timestamp     int64
		Is_commercial bool
	}
	Blog             Blog
	Content          []Content
	Layout           []Layout
	Broken_blog_name string
}
