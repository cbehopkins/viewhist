package common

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
