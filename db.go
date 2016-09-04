package main

import (
	"encoding/json"
	"gopkg.in/redis.v4"
)

func setMatch(client *redis.Client, m match) {
	err := client.Set(m.toKey(), m.toJSON(), m.ttl()).Err()
	logErr(err)
}

func getMatch(client *redis.Client, key string) (*match, error) {
	b, err := client.Get(key).Bytes()
	if err != nil {
		return nil, err
	}

	m := match{}
	err = json.Unmarshal(b, &m)
	if err != nil {
		return nil, err
	}
	return &m, nil
}
