package common

import (
	"fmt"
	"log"
	"sort"
)

type Bgbody struct {
	Id       int64
	Title    string
	Line     []string
	ImageSrc string
	Time     int
}

func (bg Bgbody) String() string {
	var ret_string string
	if len(bg.Line) > 0 {
		ret_string += "\n\n"
		ret_string += fmt.Sprintf("Blog Title:%s\n", bg.Title)
		for _, str := range bg.Line {
			ret_string += str
		}
	} else {
		log.Fatal("Zero length post")
	}
	return ret_string
}

type BgbodyArray []Bgbody

func (bga BgbodyArray) String() string {
	var ret_string string
	//ret_string += "\n\n"

	for _, bg := range bga {
		ret_string += bg.String()
	}
	return ret_string
}
func (a BgbodyArray) Len() int           { return len(a) }
func (a BgbodyArray) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a BgbodyArray) Less(i, j int) bool { return a[i].Time < a[j].Time }

func (bga BgbodyArray) TimeSorted() BgbodyArray {
	sort.Sort(bga)
	return bga
}
