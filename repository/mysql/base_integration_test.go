package mysql

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/mysql"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"

	"github.com/tusmasoma/connectHub-backend/repository"
)

type Item struct {
	ID        string    `db:"id"`
	UserID    string    `db:"user_id"`
	Text      string    `db:"text"`
	Count     int       `db:"count"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func Test_Base(t *testing.T) {
	dialect := goqu.Dialect("mysql")
	ctx := context.Background()
	userID := uuid.NewString()
	items := []Item{
		{ID: uuid.NewString(), UserID: userID, Text: "bar", Count: 1, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: uuid.NewString(), UserID: "bat", Text: "baz", Count: 2, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: uuid.NewString(), UserID: "qux", Text: "quux", Count: 3, CreatedAt: time.Now(), UpdatedAt: time.Now()},
	}
	repo := newBase[Item](db, &dialect, "TestItems")

	// create
	err := repo.Create(ctx, items[0])
	ValidateErr(t, err, nil)

	// create: already exists
	err = repo.Create(ctx, items[0])
	wantErrMsg := "Error 1062 (23000): Duplicate entry"
	if err == nil || !strings.Contains(err.Error(), wantErrMsg) {
		t.Errorf("Expected error containing '%s', got %v", wantErrMsg, err)
	}

	// batch create
	err = repo.BatchCreate(ctx, []Item{items[1], items[2]})
	ValidateErr(t, err, nil)

	// get
	item, err := repo.Get(ctx, items[0].ID)
	ValidateErr(t, err, nil)
	if d := cmp.Diff(*item, items[0], cmpopts.IgnoreFields(Item{}, "CreatedAt", "UpdatedAt")); len(d) != 0 {
		t.Errorf("Get()differs: (-got +want)\n%s", d)
	}

	// list
	qcs := []repository.QueryCondition{
		{
			Field: "Text",
			Value: "baz",
		},
		{
			Field: "Count",
			Value: 2,
		},
	}
	listItems, err := repo.List(ctx, qcs)
	ValidateErr(t, err, nil)
	if d := cmp.Diff(listItems[0], items[1], cmpopts.IgnoreFields(Item{}, "CreatedAt", "UpdatedAt")); len(d) != 0 {
		t.Errorf("List()differs: (-got +want)\n%s", d)
	}

	// update
	updatedItem := items[0]
	updatedItem.Text = "updated text"
	err = repo.Update(ctx, updatedItem.ID, updatedItem)
	ValidateErr(t, err, nil)
	gotItem, err := repo.Get(ctx, updatedItem.ID)
	ValidateErr(t, err, nil)
	if d := cmp.Diff(*gotItem, updatedItem, cmpopts.IgnoreFields(Item{}, "CreatedAt", "UpdatedAt")); len(d) != 0 {
		t.Errorf("Update() differs: (-got +want)\n%s", d)
	}

	// delete
	err = repo.Delete(ctx, updatedItem.ID)
	ValidateErr(t, err, nil)
	_, err = repo.Get(ctx, updatedItem.ID)
	if err == nil {
		t.Errorf("Expected error for deleted item, got nil")
	}
}
