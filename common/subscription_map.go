package common

import (
	"fmt"
	//"log"
	"encoding/json"
	"io/ioutil"
	"sync"
)

type UCType interface {
	ViewCount() int // Number of pages viewed so far
	SetViewCount(int)
	GetPosts() []Bgbody
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

type SubscriptionMap struct {
	sync.Mutex
	UsrInfo map[string]*UserConfig `json:"blog_subscriptions"`
	AppCfg  AppConfig              `json:"app_config"`
	Updated bool
}

func NewSubsctiptionMap() *SubscriptionMap {
	itm := new(SubscriptionMap)
	itm.UsrInfo = make(map[string]*UserConfig)
	itm.AppCfg = *NewAppConfig("", "")
	return itm

}
func (sm *SubscriptionMap) UserConfigured(username string) bool {
	_, ok := sm.UsrInfo[username]
	return ok

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

func GetTestData() *SubscriptionMap {
	//	username := ""        // FIXME Fill these in with usable data
	//	blog_subscribed := "" // FIXME fill in
	sub_map := NewSubsctiptionMap()
	//sub_map.UsrInfo[username] = NewUserConfig()
	//sub_map.UsrInfo[username].Subscribed[blog_subscribed] = &BlogProgress{BlogType: TumblrBlg, ViewCount: 0}
	sub_map.Updated = true
	return sub_map
}
