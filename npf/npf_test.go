package npf_test

import (
	"testing"
	"tumblr-dt/npf"
)

func TestLinks(t *testing.T) {
	posts := npf.TestPosts(1)
	posts[0].Render()
	if len(posts[0].GetLinks()) == 0 {
		t.Errorf("Failed to load links from post")
	}

}
