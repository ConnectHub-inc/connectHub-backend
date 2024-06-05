package redis

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
)

type Item struct {
	ID        string    `db:"id" json:"id"`
	UserID    string    `db:"user_id" json:"user_id"`
	Text      string    `db:"text" json:"text"`
	Count     int       `db:"count" json:"count"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

func Test_Base(t *testing.T) {
	ctx := context.Background()
	userID := uuid.NewString()
	items := []Item{
		{ID: uuid.NewString(), UserID: userID, Text: "bar", Count: 1, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: uuid.NewString(), UserID: "bat", Text: "baz", Count: 2, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: uuid.NewString(), UserID: "qux", Text: "quux", Count: 3, CreatedAt: time.Now(), UpdatedAt: time.Now()},
	}
	repo := newBase[Item](client)

	// set
	err := repo.Set(ctx, "item0", items[0])
	ValidateErr(t, err, nil)

	// get
	getItem, err := repo.Get(ctx, "item0")
	ValidateErr(t, err, nil)
	if d := cmp.Diff(getItem, &items[0], cmpopts.IgnoreFields(Item{}, "CreatedAt", "UpdatedAt")); len(d) != 0 {
		t.Errorf("Get()differs: (-got +want)\n%s", d)
	}

	// get: dont exists key
	_, err = repo.Get(ctx, "item1")
	ValidateErr(t, err, ErrCacheMiss)

	// exists
	exists := repo.Exists(ctx, "item0")
	if !reflect.DeepEqual(exists, true) {
		t.Errorf("Exists() \n got = %v,\n want = %v", exists, true)
	}

	// scan
	keys, err := repo.Scan(ctx, "item*")
	ValidateErr(t, err, nil)
	if !reflect.DeepEqual(keys, []string{"item0"}) {
		t.Errorf("Scan() \n got = %v,\n want = %v", keys, []string{"item0"})
	}

	// delete
	err = repo.Delete(ctx, "item0")
	ValidateErr(t, err, nil)
	exists = repo.Exists(ctx, "item0")
	if !reflect.DeepEqual(exists, false) {
		t.Errorf("Delete() \n got = %v,\n want = %v", exists, false)
	}

	// delete: dont exists key
	err = repo.Delete(ctx, "item1")
	ValidateErr(t, err, nil)
}
