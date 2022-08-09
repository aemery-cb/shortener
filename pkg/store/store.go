package store

import (
	"log"
	"time"

	"github.com/couchbase/gocb/v2"
)

type Store struct {
	cluster *gocb.Cluster
	bucket  *gocb.Bucket
}

type Shortened struct {
	Key        string `json:"key"`
	Content    string `json:"content"`
	Type       string `json:"type"`
	HitCounter int    `json:"hitCounter"`
	Owner      string
}

func New(cluster *gocb.Cluster) *Store {
	bucketName := "shortener"

	bucket := cluster.Bucket(bucketName)

	err := bucket.WaitUntilReady(5*time.Second, nil)
	if err != nil {
		log.Fatal(err)
	}

	return &Store{cluster: cluster, bucket: bucket}
}

func (store *Store) StoreUrl(key, url, userId string) error {
	col := store.bucket.DefaultCollection()

	_, err := col.Upsert("u:"+key, Shortened{
		Key:     key,
		Content: url,
		Type:    "url",
		Owner:   userId,
	}, nil)
	if err != nil {
		return err
	}

	return nil
}

func (store *Store) GetUrl(key string) string {
	if key == "" {
		return ""
	}
	col := store.bucket.DefaultCollection()

	getResult, err := col.Get("u:"+key, nil)
	if err != nil {
		return ""
	}

	var result Shortened
	err = getResult.Content(&result)

	if err != nil {
		return ""
	}

	if result.Type == "url" {
		return result.Content
	}
	return ""
}

func (store *Store) GetStats(key string) int {
	col := store.bucket.DefaultCollection()

	getResult, err := col.Get("u:"+key, nil)
	if err != nil {
		return 0
	}

	var result Shortened
	err = getResult.Content(&result)

	if err != nil {
		return 0
	}

	return result.HitCounter
}

func (store *Store) UpdateHitCounter(key string) int {
	if key == "" {
		return 0
	}
	col := store.bucket.DefaultCollection()

	getResult, err := col.Get("u:"+key, nil)
	if err != nil {
		return 0
	}

	var result Shortened
	err = getResult.Content(&result)

	if err != nil {
		return 0
	}

	result.HitCounter += 1
	_, err = col.Upsert("u:"+key, result, nil)

	if err != nil {
		return 0
	}

	return result.HitCounter
}
