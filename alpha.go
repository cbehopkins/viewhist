// alpha.go
package main

//import (
//	"fmt"

//	"github.com/dghubble/oauth1"
//	//	"github.com/MariaTerzieva/gotumblr"
//)

import (
	"fmt"
	"html/template"
	//"log"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/cbehopkins/viewhist/tumblr_man"
)

// so a user is subscribed to blogs
// Once we have the username we want to know what blogs they are subscribed to
type SubscriptionMap struct {
	sync.Mutex
	UsrInfo map[string]*tumblr_man.UserConfig `json:"blog_subscriptions"`
	AppCfg  tumblr_man.AppConfig              `json:"app_config"`
	Updated bool
}

func NewSubsctiptionMap() *SubscriptionMap {
	itm := new(SubscriptionMap)
	itm.UsrInfo = make(map[string]*tumblr_man.UserConfig)
	itm.AppCfg = *tumblr_man.NewAppConfig("", "")
	return itm

}

func (sm *SubscriptionMap) GenJson() string {

	output, err := json.MarshalIndent(sm, "", "    ")
	if err != nil {
		fmt.Printf("error: %v\n", err)
	}
	return string(output)
}

func (sm *SubscriptionMap) WriteJson(filename string) {
	sm.Lock()
	if sm.Updated {
		sm.Updated = false

		text_to_write := sm.GenJson()
		d1 := []byte(text_to_write)
		err := ioutil.WriteFile(filename, d1, 0644)
		check(err)
	}
	sm.Unlock()
}

func set_test_data() *SubscriptionMap {
	username := ""        // FIXME Fill these in with usable data
	blog_subscribed := "" // FIXME fill in
	sub_map := NewSubsctiptionMap()
	sub_map.UsrInfo[username] = tumblr_man.NewUserConfig()
	sub_map.UsrInfo[username].Subscribed[blog_subscribed] = &tumblr_man.BlogProgress{BlogType: tumblr_man.TumblrBlg, ViewCount: 0}
	sub_map.Updated = true
	return sub_map

}

// main performs the Tumblr OAuth1 user flow from the command line
var sub_map *SubscriptionMap

func main() {
	//test()

	filename := "alpha_sub_map.json"

	dat, err := ioutil.ReadFile(filename)
	if err == nil {
		sub_map = NewSubsctiptionMap()
		err := json.Unmarshal([]byte(dat), sub_map)
		check(err)
	} else {
		fmt.Println("Using default data set:", err)
		sub_map = set_test_data()
	}

	json_string := sub_map.GenJson()

	fmt.Println("Config used is:", json_string)
	go FileWriterConstant(filename)
	http.HandleFunc("/", handler)
	http.HandleFunc("/view/", viewHandler)
	http.HandleFunc("/next/", nextHandler)
	http.ListenAndServe(":8090", nil)
}

type Page struct {
	Title             string
	Body              string
	Items             []tumblr_man.Bgbody
	NextPageId        int
	DefaultCountValue int
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
func FileWriterConstant(filename string) {
	for {
		time.Sleep(time.Second * 60)
		sub_map.WriteJson(filename)

	}
}
func handler(w http.ResponseWriter, r *http.Request) {
	return_username := ""

	// We might already have a username set up
	for _, cookie := range r.Cookies() {
		fmt.Printf("Cookie is:\"%s\"\n", cookie.Name)
		if cookie.Name == "username" {
			return_username = cookie.Value
		}
	}

	// If ther is no username set in the cookies then
	// something is amis, set it to me for now
	// TBD set up form so you can enter a different username
	if return_username == "" {
		// Set up the cookie with our current status
		expiration := time.Now().Add(2 * 24 * time.Hour)
		un_cookie := http.Cookie{Name: "username", Value: "cbehopkins", Expires: expiration}
		http.SetCookie(w, &un_cookie)
	}

	t, err := template.ParseFiles("src/github.com/cbehopkins/viewhist/test.html")
	if err != nil {
		cwd, _ := os.Getwd()
		fmt.Println("running in directory ", cwd)
		panic(err)
	}
	p := Page{Title: "Content",
		Body: "<p>Hello</p>"}

	p.NextPageId = 3

	err = t.Execute(w, p)
	check(err)
}

func nextHandler(w http.ResponseWriter, r *http.Request) {
}

func viewHandler(w http.ResponseWriter, r *http.Request) {

	return_username := ""
	for _, cookie := range r.Cookies() {
		fmt.Printf("Cookie is:\"%s\"\n", cookie.Name)
		if cookie.Name == "username" {
			return_username = cookie.Value
		}

	}
	if return_username == "" {
		// TBD add new template here to redirect to home
		fmt.Fprintf(w, "<body><h1>Error Username not set, return to home</h1></body>")
		return
	}
	sub_map.Lock()

	defer sub_map.Unlock()
	subscriptions, ok := sub_map.UsrInfo[return_username]
	if !ok {
		fmt.Fprintf(w, "<body><h1>Error Username not configured, return to home</h1></body>")
		return
	}

	// Parse the form to work out how many we should fetch
	r.ParseForm()
	posted := r.Form
	for _, itm := range r.Form["count"] {
		fmt.Println("We have an item:", itm)
	}
	sub_cnt_array := posted["count"]
	var num_to_fetch int
	if len(sub_cnt_array) > 0 && (sub_cnt_array[0] != "") {
		var err error
		num_to_fetch, err = strconv.Atoi(sub_cnt_array[0])
		if err != nil {
			fmt.Println("Received submission of:", sub_cnt_array[0])
			num_to_fetch = 1
		}
	}
	if num_to_fetch < 1 {
		num_to_fetch = 1
	}

	//current_page_id, err := strconv.Atoi(r.URL.Path[6:])
	//check(err)
	var blog_array []tumblr_man.Bgbody
	blog_array = make([]tumblr_man.Bgbody, 0)
	for blg_2_get, current_subscription := range subscriptions.Subscribed {
		//current_subscription, ok := ["cpliso"]
		//if !ok {
		//	log.Fatal("cpliso not defined")
		//}
		current_page_id := current_subscription.ViewCount

		subscriptions.Subscribed[blg_2_get].ViewCount = current_page_id + num_to_fetch
		//fmt.Println("Got pageid as current_page_id, %x, next is %x,  %s\n", current_page_id, next_page_id, r.URL.Path)
		//fmt.Println("page title is", title_string)
		post_bodies := tumblr_man.GetPosts(num_to_fetch, current_page_id, blg_2_get, sub_map.UsrInfo[return_username], sub_map.AppCfg)
		blog_array = append(blog_array, post_bodies...)
		sub_map.Updated = true
	}

	title_string := "Username:" + return_username
	p := &Page{Title: title_string}

	p.Items = blog_array
	//p.NextPageId = next_page_id
	p.NextPageId = 0
	p.DefaultCountValue = num_to_fetch

	// Parse the template ready to render out the final page
	t, err := template.ParseFiles("src/github.com/cbehopkins/viewhist/view.html")
	if err != nil {
		cwd, _ := os.Getwd()
		fmt.Println("running in directory ", cwd)
		panic(err)
	}

	err = t.Execute(w, p)
	check(err)
}
