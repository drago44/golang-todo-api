// Package todos contains benchmark tests for the todos module.
package todos

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// BenchmarkCreateAndList simulates steady QPS create + list with minimal locking
func BenchmarkCreateAndList(b *testing.B) {
	gin.SetMode(gin.ReleaseMode)

	file := filepath.Join(b.TempDir(), "bench.db")
	dsn := file + "?_journal_mode=WAL&_synchronous=NORMAL&_busy_timeout=5000&_foreign_keys=ON&_cache_size=-20000"

	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{PrepareStmt: true, SkipDefaultTransaction: true})
	if err != nil {
		b.Fatalf("open sqlite: %v", err)
	}

	if err := db.AutoMigrate(&Todo{}); err != nil {
		b.Fatalf("migrate: %v", err)
	}

	if sqlDB, err2 := db.DB(); err2 == nil {
		// Few conns for WAL; it's a single-process benchmark
		sqlDB.SetMaxOpenConns(2)
		sqlDB.SetMaxIdleConns(2)
	}

	repo := NewTodoRepository(db)
	svc := NewTodoService(repo)
	h := NewTodoHandler(svc)
	r := gin.New()
	rg := r.Group("/")
	h.RegisterTodoRoutes(rg)

	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		// Create
		body := CreateTodoRequest{Title: "t-" + itoa(i+1), Description: "d"}
		buf, _ := json.Marshal(body)
		req := httptest.NewRequest(http.MethodPost, "/todos", bytes.NewReader(buf))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusCreated {
			b.Fatalf("create status=%d", w.Code)
		}
		// List
		req2 := httptest.NewRequest(http.MethodGet, "/todos", nil)
		w2 := httptest.NewRecorder()
		r.ServeHTTP(w2, req2)

		if w2.Code != http.StatusOK {
			b.Fatalf("list status=%d", w2.Code)
		}
	}
}

// micro itoa without fmt to reduce allocs
func itoa(n int) string {
	if n == 0 {
		return "0"
	}

	var a [20]byte

	i := len(a)
	for n > 0 {
		i--
		a[i] = byte('0' + n%10)
		n /= 10
	}

	return string(a[i:])
}
