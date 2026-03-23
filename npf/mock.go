package npf

import (
	"encoding/json"
	"math/rand"
	"os"
	"strconv"
	"strings"
)

const sep string = string(os.PathSeparator)

var rootDir string = ""

func getRootDir() string {
	if len(rootDir) == 0 {
		d, _ := os.Getwd()

		parts := strings.Split(d, sep)
		for i := len(parts) - 1; i >= 0; i-- {
			dir := strings.Join(parts[:i], sep)
			_, err := os.ReadFile(dir + sep + "main.go")
			if err == nil {
				rootDir = dir
				break
			}
		}
	}
	return rootDir
}

func loadTestFile(path string) ([]byte, error) {
	b, err := os.ReadFile(getRootDir() + sep + "testdata" + sep + path)
	return b, err
}

func TestBlog() Blog {
	blogBytes, err := loadTestFile("blog.json")
	if err != nil {
		panic(err)
	}

	blog := Blog{}
	json.Unmarshal(blogBytes, &blog)
	return blog
}
func TestPosts(i int) []Post {
	res := []Post{}
	for range i {
		res = append(res, TestPost())
	}
	return res
}

func TestPost() Post {
	blogBytes, err := loadTestFile("post.json")
	if err != nil {
		panic(err)
	}

	post := Post{}
	json.Unmarshal(blogBytes, &post)
	post.Id_string += strconv.Itoa(rand.Int())
	post.Blog = TestBlog()
	for i := 0; i < rand.Int()%6; i++ {
		post.Content = append(post.Content, getRandomContent())
	}

	for i := 0; i < rand.Int()%6; i++ {
		post.Trail = append(post.Trail, TestTrail())
	}

	return post
}

func TestTrail() TrailPost {
	blogBytes, err := loadTestFile("trail.json")
	if err != nil {
		panic(err)
	}

	trail := TrailPost{}
	json.Unmarshal(blogBytes, &trail)
	trail.Blog = TestBlog()

	for i := 0; i < rand.Int()%6; i++ {
		trail.Content = append(trail.Content, getRandomContent())
	}

	return trail
}

var contentKeys = []string{
	"text",
	"image",
	"link",
}

func getRandomContent() Content {
	return TestContent(contentKeys[rand.Int()%len(contentKeys)])
}

func TestContent(contentType string) Content {
	blogBytes, err := loadTestFile("content" + sep + contentType + ".json")
	if err != nil {
		panic(err)
	}
	content := Content{}
	json.Unmarshal(blogBytes, &content)
	return content
}
func TestContentText() Content {
	return TestContent("text")
}
func TestContentImage() Content {
	return TestContent("image")
}
func TestContentLink() Content {
	return TestContent("link")
}
