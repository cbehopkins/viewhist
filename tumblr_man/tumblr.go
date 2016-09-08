// tumblr.go
package tumblr_man

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/cbehopkins/viewhist/common"

	"github.com/MariaTerzieva/gotumblr"
	"github.com/dghubble/oauth1"
	"github.com/dghubble/oauth1/tumblr"
)

var config oauth1.Config

//func main() {

//	token := ""
//	token_secret := ""
//	callback_url := ""
//	client := gotumblr.NewTumblrRestClient(consumer_key, secret_key, token, token_secret, callback_url, "http://api.tumblr.com")

//	fmt.Println("Hello World!")
//}
// Each blog has a progress through it

type BlogStruct struct {
	// The aim here is to keep the produced JSON as small as possible
	BlogName  string   `json:"blog_name,attr"`
	PostUrl   string   `json:"post_url,attr"`
	Body      string   `json:"body,omitempty"`
	Id        int      `json:"id,attr"`
	CType     string   `json:"type,attr":`
	Date      string   `json:"date,attr":`
	Timestamp int      `json:"timestamp,attr"`
	Format    string   `json:"format,attr"`
	Tags      []string `json:"tags,attr"`
	Summary   string   `json:"summary,attr"`
	Title     string   `json:"title,attr"`
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
func UnMarshalJson(input []byte) (v *BlogStruct, err error) {
	v = new(BlogStruct)
	err = json.Unmarshal(input, v)
	//txt = v.Body
	return
}

func GetPosts(count, offset int, blg_2_get string, cfg *common.UserConfig, app common.AppConfig) []common.Bgbody {
	fmt.Println("Running GetPosts")
	ret_posts := make([]common.Bgbody, count)
	blog_array, err := GetTumblr(count, offset, blg_2_get, cfg, app)
	check(err)
	for i, blpost := range blog_array {
		html_post := blpost.Body

		r := strings.NewReplacer("<p>", "", "</p>", "", "<br/>", "", "<b>", "", "</b>", "\n", "<span>", "", "</span>", "", "<strong>", "", "</strong>", "")
		tt0 := strings.Split(r.Replace(html_post), "\n")
		var ttp common.Bgbody

		ttp.Line = make([]string, len(tt0))
		copy(ttp.Line, tt0)
		//ret_posts[i] = Bgbody{Line: []string{html_post}}
		ret_posts[i].Line = tt0
		ret_posts[i].Id = offset + i
		ret_posts[i].Title = blpost.Title
	}

	return ret_posts
}

func getUsrStrings(cfg *common.UserConfig) (token, secret string) {

	//get_key()
	token = cfg.Token
	secret = cfg.Secret
	return
}

func GetTumblr(count, offset int, blg_2_get string, cfg *common.UserConfig, app common.AppConfig) (blog_array []*BlogStruct, err error) {
	consumer_key := app.ConsumerKey
	secret_key := app.SecretKey
	token, secret := getUsrStrings(cfg)
	callback_url := "http://www.doubleudoubleudoubleu.co.uk/tum"
	client := gotumblr.NewTumblrRestClient(consumer_key, secret_key, token, secret, callback_url, "http://api.tumblr.com")

	//info := client.Info()
	//fmt.Println(info.User.Name)

	//following := client.Following(map[string]string{})
	//fmt.Println(following.Total_blogs)
	//for _, ky := range following.Blogs {
	//	fmt.Println(ky)
	//}

	// First run a fetch just to get the number of posts in totoal
	// undet the category of text
	post_options := make(map[string]string)
	post_options["limit"] = strconv.Itoa(1)
	post_options["offset"] = strconv.Itoa(0)
	blg := client.Posts(blg_2_get, "text", post_options)

	// Now ask for the oldest post
	reversed_offset := blg.Total_posts - int64(offset)
	post_options["limit"] = strconv.Itoa(count)
	post_options["offset"] = strconv.Itoa(int(reversed_offset))
	blg = client.Posts(blg_2_get, "text", post_options)

	num_of_blogs := len(blg.Posts)
	//fmt.Printf("Requested offset %d, Reversed offset %d, Num_posts %d, total %d \n", offset, reversed_offset, num_of_blogs, blg.Total_posts)

	//steve := blg.Posts[1]
	//_, err = UnMarshalJson(steve)
	//fmt.Printf("Its:%s\n", steve)
	if err == nil {
		//pt_txt := pt_tzt.BlogName
		//fmt.Println("Its:", pt_txt)
	} else {
		fmt.Println("Error is:", err)
	}

	blog_array = make([]*BlogStruct, num_of_blogs)
	for i, blpost := range blg.Posts {
		//fmt.Printf("%s", blpost)

		//blog_array[i], err = UnMarshalJson(blpost)
		tmp, err := UnMarshalJson(blpost)
		check(err)
		//r_txt := tmp.Body
		//fmt.Printf("%s", r_txt)
		blog_array[(num_of_blogs-i)-1] = tmp
	}

	return
}

type APageT struct {
	Title string
	Body  string
}

func RenderAddForm(w http.ResponseWriter) {
	AddPage := APageT{Title: "Add a Tumblr", Body: "Test"}
	template, err := template.ParseFiles("src/github.com/cbehopkins/flktst/add_tumblr.html")
	if err != nil {
		panic(err)
	}
	err = template.Execute(w, AddPage)
	check(err)
}
func InterpretAddForm(w http.ResponseWriter, r *http.Request) {
	// add code here to:
	// Get the current username from the cookie
	// Extract the required blog name from the submitted tumblradd text field
	// Add this to the required data structure
	// Return to root

	http.Redirect(w, r, "/", 301)

}
func get_key(app common.AppConfig) {
	// read credentials from constants
	//consumerKey, consumerSecret := getAppStrings()
	consumer_key := app.ConsumerKey
	secret_key := app.SecretKey

	if consumer_key == "" || secret_key == "" {
		log.Fatal("Required variables missing.")
	}

	config = oauth1.Config{
		ConsumerKey:    consumer_key,
		ConsumerSecret: secret_key,
		// Tumblr does not support oob, uses consumer registered callback
		CallbackURL: "",
		Endpoint:    tumblr.Endpoint,
	}

	requestToken, requestSecret, err := login()
	if err != nil {
		log.Fatalf("Request Token Phase: %s", err.Error())
	}
	accessToken, err := receivePIN(requestToken, requestSecret)
	if err != nil {
		log.Fatalf("Access Token Phase: %s", err.Error())
	}

	fmt.Println("Consumer was granted an access token to act on behalf of a user.")
	fmt.Printf("token: %s\nsecret: %s\n", accessToken.Token, accessToken.TokenSecret)
}

func login() (requestToken, requestSecret string, err error) {
	requestToken, requestSecret, err = config.RequestToken()
	if err != nil {
		return "", "", err
	}
	authorizationURL, err := config.AuthorizationURL(requestToken)
	if err != nil {
		return "", "", err
	}
	fmt.Printf("Open this URL in your browser:\n%s\n", authorizationURL.String())
	return requestToken, requestSecret, err
}

func receivePIN(requestToken, requestSecret string) (*oauth1.Token, error) {
	fmt.Printf("Choose whether to grant the application access.\nPaste " +
		"the oauth_verifier parameter (excluding trailing #_=_) from the " +
		"address bar: ")
	var verifier string
	_, err := fmt.Scanf("%s", &verifier)
	accessToken, accessSecret, err := config.AccessToken(requestToken, requestSecret, verifier)
	if err != nil {
		return nil, err
	}
	return oauth1.NewToken(accessToken, accessSecret), err
}
