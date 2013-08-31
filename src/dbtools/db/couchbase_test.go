package db

import (
	"testing"
)

func TestGetBucket(t *testing.T) {
	b := getLocalBucket("test")
	if b == nil {
		t.Errorf("Problem retreiving local bucket %s", "test")
	}
}

func TestView(t *testing.T) {
	params := map[string]interface{}{}
	b := getLocalBucket("test")
	vres, _ := b.View("dev_all", "all", params)
	if vres.TotalRows == 0 {
		t.Errorf("Problem accessing documents from view %s", "dev_all/all")
	}
}
