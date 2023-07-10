package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/fatih/color"
	consul "github.com/hashicorp/consul/api"
	"github.com/servehub/utils/gabs"
)

func ConsulClient(consulAddress string) (*consul.Client, error) {
	conf := consul.DefaultConfig()
	conf.Address = consulAddress
	return consul.NewClient(conf)
}

var PutConsulKv = func(client *consul.Client, key string, value string) error {
	log.Printf("consul put `%s`: %s", key, value)
	_, err := client.KV().Put(&consul.KVPair{Key: strings.TrimPrefix(key, "/"), Value: []byte(value)}, nil)
	return err
}

func ListConsulKv(client *consul.Client, prefix string, q *consul.QueryOptions) (consul.KVPairs, error) {
	log.Printf("consul list `%s`", prefix)
	list, _, err := client.KV().List(prefix, q)
	return list, err
}

var DelConsulKv = func(client *consul.Client, key string) error {
	log.Printf("consul delete `%s`", key)
	_, err := client.KV().Delete(strings.TrimPrefix(key, "/"), nil)
	return err
}

var RegisterPluginData = func(plugin string, packageName string, data string, consulAddress string) error {
	consulApi, err := ConsulClient(consulAddress)
	if err != nil {
		return err
	}

	return PutConsulKv(consulApi, "services/data/"+packageName+"/"+plugin, data)
}

var DeletePluginData = func(plugin string, packageName string, consulAddress string) error {
	log.Println(color.YellowString("Delete %s for %s package in consul", plugin, packageName))
	consulApi, err := ConsulClient(consulAddress)
	if err != nil {
		return err
	}

	return DelConsulKv(consulApi, "services/data/"+packageName+"/"+plugin)
}

func MarkAsOutdated(client *consul.Client, name string, delay time.Duration) error {
	log.Printf("Mark service `%s` as outdated\n", name)
	json := fmt.Sprintf(`{"endOfLife":"%s"}`, time.Now().Add(delay).Format(time.RFC3339))
	return PutConsulKv(client, "services/outdated/"+name, json)
}

type CachedObject struct {
	ExpiredAt int64
	Obj       *gabs.Container
}

func ConsulCache(client *consul.Client, key string, ttlSeconds int, f func() *gabs.Container) *gabs.Container {
	if pair, _, err := client.KV().Get("services/cache/"+key, nil); err == nil && pair != nil {
		var cachedObject CachedObject
		json.Unmarshal(pair.Value, &cachedObject)

		if cachedObject.ExpiredAt > time.Now().Unix() {
			return cachedObject.Obj
		}
	}

	obj := f()
	client.KV().Put(&consul.KVPair{
		Key:   "services/cache/" + key,
		Value: []byte(fmt.Sprintf(`{"expiredAt": %d, "obj": %s}`, time.Now().Unix()+int64(ttlSeconds), obj.String())),
	}, nil)

	return obj
}
