package common

const (
	TumblrBlg = iota
)

// Each user is subscribed to a number of blogs
type UserConfig struct {
	// There is one of these per user

	// REVISIT currently since there is only 1 servce there should be another
	// level of indirection to the multiple services a single user could be subscribed to
	// And for each service the username on that service

	Type   int    `json:"type"`   // Service type e.g. Tumblr, URL etc
	Token  string `json:"token"`  // Tokens needed to access the service
	Secret string `json:"secret"` // secret to access the service

	// This maps from blog name to a structure saying how far we are through that blog
	Subscribed map[string]*BlogProgress `json:"blg_prg"`
}

func NewUserConfig() *UserConfig {
	itm := new(UserConfig)
	itm.Token = ""
	itm.Secret = ""
	itm.Subscribed = make(map[string]*BlogProgress)
	return itm
}
