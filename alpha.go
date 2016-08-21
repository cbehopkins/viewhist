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
	"time"

	"github.com/cbehopkins/viewhist/common"

	"github.com/cbehopkins/viewhist/tumblr_man"
	"github.com/gorilla/mux"
)

// main performs the Tumblr OAuth1 user flow from the command line
var sub_map *common.SubscriptionMap

func main() {
	//test()

	filename := "alpha_sub_map.json"

	dat, err := ioutil.ReadFile(filename)
	if err == nil {
		sub_map = common.NewSubsctiptionMap()
		err := json.Unmarshal([]byte(dat), sub_map)
		check(err)
	} else {
		fmt.Println("Using default data set:", err)
		sub_map = common.GetTestData()
	}

	json_string := sub_map.GenJson()

	fmt.Println("Config used is:", json_string)
	go FileWriterConstant(filename)
	mux := mux.NewRouter()

	mux.HandleFunc("/login/", loginHandler)
	mux.HandleFunc("/logout/", logoutHandler)
	mux.HandleFunc("/view/", viewHandler)
	mux.HandleFunc("/add/", addHandler)
	mux.HandleFunc("/", handler)

	http.ListenAndServe(":8091", mux)
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
func logoutHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Running logout Handler")
	un_cookie := http.Cookie{Name: "username", Value: "", Path: "/"}
	un_cookie.MaxAge = -1
	http.SetCookie(w, &un_cookie)
	http.Redirect(w, r, "/", 301)
	return
}
func loginHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Running login handler")
	r.ParseForm()
	var canceled bool
	var oked bool
	for _, itm := range r.Form["res"] {
		fmt.Println("We have result:", itm)
		if itm == "cancel" {
			canceled = true
		}
		if itm == "ok" {
			oked = true
		}
	}
	if canceled {
		http.Redirect(w, r, "/", 301)
		return
	}
	if oked {
		fmt.Println("Examining", r.Form)
		for _, itm := range r.Form["username"] {
			if itm != "" {
				fmt.Println("We've been given username to login", itm)
				// REVISIT add an interface to the submap
				_, ok := sub_map.UsrInfo[itm]
				if ok {
					// Set up the cookie with our current status
					expiration := time.Now().Add(2 * 24 * time.Hour)
					un_cookie := http.Cookie{Name: "username", Value: itm, Expires: expiration, Path: "/"}
					http.SetCookie(w, &un_cookie)
					http.Redirect(w, r, "/", 301)
					return
				} else {
					// Username not found
					expiration := time.Now().Add(12 * time.Hour)
					un_cookie := http.Cookie{Name: "username", Value: itm, Expires: expiration, Path: "/"}
					http.SetCookie(w, &un_cookie)
					http.Redirect(w, r, "/", 301)

					return
				}
			}
		}
	}
	var HomePage common.Page
	HomePage = common.Page{Title: "Login"}
	// TBD this is misnamed
	template, err := template.ParseFiles("src/github.com/cbehopkins/flktst/login_user.html")
	if err != nil {
		cwd, _ := os.Getwd()
		fmt.Println("running in directory ", cwd)
		panic(err)
	}
	err = template.Execute(w, HomePage)
	check(err)

}
func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Running root handler")
	return_username := ""
	// We might already have a username set up
	for _, cookie := range r.Cookies() {
		fmt.Printf("Cookie is:\"%s\"\n", cookie.Name)
		if cookie.Name == "username" {
			return_username = cookie.Value
		}
	}

	var HomePage common.Page
	HomePage = common.Page{Title: "Content"}
	HomePage.Username = return_username
	r.ParseForm()
	for _, itm := range r.Form["action"] {
		if itm == "view" {
			fmt.Println("We've been told to view")
			http.Redirect(w, r, "view/", 301)
			return
		}
		if itm == "add" {
			fmt.Println("We've been told to add")
			http.Redirect(w, r, "add/", 301)
			return
		}
		if itm == "login" {
			fmt.Println("We've been told to login")
			http.Redirect(w, r, "login/", 301)
			return
		}
		if itm == "logout" {
			fmt.Println("We've been told to logout")
			http.Redirect(w, r, "logout/", 301)
			return
		}
	}

	template, err := template.ParseFiles("src/github.com/cbehopkins/flktst/test.html")
	if err != nil {
		cwd, _ := os.Getwd()
		fmt.Println("running in directory ", cwd)
		panic(err)
	}
	err = template.Execute(w, HomePage)
	check(err)
}

func addHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Running addHandler")
	var AddPage common.Page
	AddPage = common.Page{Title: "Adding"}
	r.ParseForm()
	var canceled bool
	var oked bool
	for _, itm := range r.Form["res"] {
		fmt.Println("We have result:", itm)
		if itm == "cancel" {
			canceled = true
		}
		if itm == "ok" {
			oked = true
		}
	}
	if canceled && oked {
		http.Redirect(w, r, "/", 301)
		return
	}
	//for _, itm := range r.Form["type"] {
	fmt.Println("We have type:", r.Form)
	//	}
	for _, itm := range r.Form["type"] {

		if itm == "tumblr" {
			AddPage.Body = "Tumblr submitted"
			// So someone has submitted this page and said they want tumbler style add
			// So add needs to return a webpage for us to render as it's add form
			// we will then write that to the writer and finish
			tumblr_man.RenderAddForm(w)
			return
		}
	}
	_, ok := r.Form["tumblradd"]
	if ok {
		// So the webpagethe tumblr suggested has been submitted
		// so we pass r.Form to the tumblr manager to do what it will with that form
		// If all goes well it will then make the changes it wants to and we are done
		//tumblr_man.InterpretAddForm(w, r)
		http.Redirect(w, r, "/", 301)

		return

	}
	//FIXME encapsulate the access to sub_map in an interface
	//sub_map.Lock()
	//defer sub_map.Unlock()

	template, err := template.ParseFiles("src/github.com/cbehopkins/flktst/add.html")
	if err != nil {
		cwd, _ := os.Getwd()
		fmt.Println("running in directory ", cwd)
		panic(err)
	}
	err = template.Execute(w, AddPage)
	check(err)
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Running viewHandler")
	service_username := ""
	for _, cookie := range r.Cookies() {
		//fmt.Printf("Cookie is:\"%s\"\n", cookie.Name)
		if cookie.Name == "username" {
			service_username = cookie.Value
		}

	}
	if service_username == "" {
		// TBD add new template here to redirect to home
		fmt.Fprintf(w, "<body><h1>Error Username not set, return to home</h1></body>")
		return
	}
	sub_map.Lock()

	defer sub_map.Unlock()
	ok := sub_map.UserConfigured(service_username)
	if !ok {
		fmt.Fprintf(w, "<body><h1>Error Username not configured, return to home</h1></body>")
		return
	}

	// Parse the form to work out how many we should fetch
	r.ParseForm()
	posted := r.Form
	//for _, itm := range r.Form["count"] {
	//fmt.Println("We have an item:", itm)
	//}
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
	var blog_array []common.Bgbody
	blog_array = make([]common.Bgbody, 0)

	for blg_2_get, blog_progress := range sub_map.UsrInfo[service_username].Subscribed {
		current_page_id := blog_progress.GetViewCount()
		// Update the view count to the number that the user has seen
		// TBD do this on the next page fetch in a series?
		blog_progress.SetViewCount(current_page_id + num_to_fetch)
		//fmt.Println("Got pageid as current_page_id, %x, next is %x,  %s\n", current_page_id, next_page_id, r.URL.Path)
		//fmt.Println("page title is", title_string)
		post_bodies := tumblr_man.GetPosts(num_to_fetch, current_page_id, blg_2_get, sub_map.UsrInfo[service_username], sub_map.AppCfg)
		blog_array = append(blog_array, post_bodies...)
		sub_map.Updated = true
	}

	title_string := "Username:" + service_username
	p := &common.Page{Title: title_string}

	p.Items = blog_array
	p.DefaultCountValue = num_to_fetch

	// Parse the template ready to render out the final page
	t, err := template.ParseFiles("src/github.com/cbehopkins/flktst/view.html")
	if err != nil {
		cwd, _ := os.Getwd()
		fmt.Println("running in directory ", cwd)
		panic(err)
	}

	err = t.Execute(w, p)
	check(err)
}
