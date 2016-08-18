// tumblr.go
package tumblr_man

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"

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
const (
	TumblrBlg = iota
)

type BlogProgress struct {
	BlogType  int
	ViewCount int
}

// Each user is subscribed to a number of blogs
type UserConfig struct {
	// There is one of these per user
	// TBD add things like the login paraneters to the services into here
	Token  string `json:"token"`
	Secret string `json:"secret"`

	Subscribed map[string]*BlogProgress `json:"blg_prg"`
}
type AppConfig struct {
	ConsumerKey string `json:"consumer_key"`
	SecretKey   string `json:"secret_key"`
}

func NewAppConfig(consumer_key, secret_key string) *AppConfig {
	itm := new(AppConfig)
	itm.ConsumerKey = consumer_key
	itm.SecretKey = secret_key
	return itm
}
func NewUserConfig() *UserConfig {
	itm := new(UserConfig)
	itm.Token = ""
	itm.Secret = ""
	itm.Subscribed = make(map[string]*BlogProgress)
	return itm
}

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

type Bgbody struct {
	Id    int
	Title string
	Line  []string
}

func GetPosts(count, offset int, blg_2_get string, cfg *UserConfig, app AppConfig) []Bgbody {
	ret_posts := make([]Bgbody, count)
	blog_array, err := GetTumblr(count, offset, blg_2_get, cfg, app)
	check(err)
	for i, blpost := range blog_array {
		html_post := blpost.Body

		r := strings.NewReplacer("<p>", "", "</p>", "", "<br/>", "", "<b>", "", "</b>", "\n")
		tt0 := strings.Split(r.Replace(html_post), "\n")
		var ttp Bgbody

		ttp.Line = make([]string, len(tt0))
		copy(ttp.Line, tt0)
		//ret_posts[i] = Bgbody{Line: []string{html_post}}
		ret_posts[i].Line = tt0
		ret_posts[i].Id = offset + i
		ret_posts[i].Title = blpost.Title
	}

	return ret_posts
}

func getUsrStrings(cfg *UserConfig) (token, secret string) {

	//get_key()
	token = cfg.Token
	secret = cfg.Secret
	return
}

func GetTumblr(count, offset int, blg_2_get string, cfg *UserConfig, app AppConfig) (blog_array []*BlogStruct, err error) {
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
	post_options := make(map[string]string)
	post_options["limit"] = strconv.Itoa(count)
	post_options["offset"] = strconv.Itoa(offset)

	blg := client.Posts(blg_2_get, "text", post_options)

	num_of_blogs := len(blg.Posts)
	fmt.Println("There are ", num_of_blogs, "Posts")

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
		blog_array[i] = tmp
	}

	return
}
func get_key(app AppConfig) {
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
