package todos

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func createTestDB(t *testing.T) *gorm.DB {
	t.Helper()

	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open sqlite memory: %v", err)
	}

	if err := db.AutoMigrate(&Todo{}); err != nil {
		t.Fatalf("failed to migrate: %v", err)
	}

	t.Log("sqlite memory database created and migrated")

	return db
}

func TestRepository_CRUD(t *testing.T) {
	db := createTestDB(t)
	repo := NewTodoRepository(db)

	// Create
	todo := &Todo{Title: "A", Description: "d"}
	assert.NoError(t, repo.Create(todo))
	assert.NotZero(t, todo.ID)
	t.Logf("created todo with ID=%d", todo.ID)

	// GetByID
	got, err := repo.GetByID(todo.ID)
	assert.NoError(t, err)
	assert.Equal(t, "A", got.Title)
	t.Logf("fetched by id: %+v", got)

	// ExistsByTitle
	exists, err := repo.ExistsByTitle("A")
	assert.NoError(t, err)
	assert.True(t, exists)
	t.Logf("exists by title 'A': %v", exists)

	notExists, err := repo.ExistsByTitle("B")
	assert.NoError(t, err)
	assert.False(t, notExists)
	t.Logf("exists by title 'B': %v", notExists)

	// GetAll
	list, err := repo.GetAll()
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(list), 1)
	t.Logf("list size=%d", len(list))

	// Update
	got.Description = "new"
	assert.NoError(t, repo.Update(got))
	t.Log("updated description to 'new'")

	got2, err := repo.GetByID(todo.ID)
	assert.NoError(t, err)
	assert.Equal(t, "new", got2.Description)
	t.Logf("verified update: %+v", got2)

	// Delete
	assert.NoError(t, repo.Delete(todo.ID))
	_, err = repo.GetByID(todo.ID)
	assert.Error(t, err)
	t.Logf("deleted todo id=%d", todo.ID)
}
