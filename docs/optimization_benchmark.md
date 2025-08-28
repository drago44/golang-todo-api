### 📊 Benchmark Comparison: BEFORE vs AFTER Optimization

#### BEFORE Optimization (previous version):

```
goos: darwin
goarch: amd64
cpu: Intel(R) Core(TM) i9-9880H CPU @ 2.30GHz
BenchmarkCreateTodo-16    	     177	   6748422 ns/op	 3177838 B/op	   43072 allocs/op
```

#### AFTER Optimization:

```
goos: darwin
goarch: amd64
cpu: Intel(R) Core(TM) i9-9880H CPU @ 2.30GHz
BenchmarkCreateTodo-16    	     261	   4586188 ns/op	 1974065 B/op	   28613 allocs/op
```

---

### 🚀 Optimization Results:

| Metric                      | BEFORE    | AFTER     | Improvement                     |
| --------------------------- | --------- | --------- | ------------------------------- |
| **Time (ns/op)**            | 6,748,422 | 4,586,188 | **~1.47x faster** ⚡            |
| **Memory (B/op)**           | 3,177,838 | 1,974,065 | **~1.61x less memory** 💾       |
| **Allocations (allocs/op)** | 43,072    | 28,613    | **~1.51x fewer allocations** 📉 |

---

### 🔍 Detailed Improvement Analysis:

1.  **Speed (Time):**

    - **BEFORE:** 6.75ms per operation
    - **AFTER:** 4.59ms per operation
    - **Gain:** ~32% faster (reduction of ~2.16ms)

2.  **Memory:**

    - **BEFORE:** ~3.18MB per operation
    - **AFTER:** ~1.97MB per operation
    - **Gain:** ~38% less memory (saved ~1.2MB)

3.  **Allocations:**
    - **BEFORE:** 43,072 allocations per operation
    - **AFTER:** 28,613 allocations per operation
    - **Gain:** ~34% fewer allocations (saved ~14,459 allocations)

---

### 💡 What had the biggest impact:

- **SQLite PRAGMAs + prepared statements** - faster DB queries
- **sync.Pool for DTO** - reusing structs instead of creating new ones
- **Static JSON responses** - less reflection and allocations
- **Gin Release Mode** - disabling debug overhead
- **Optimized SQL queries** - `SELECT id LIMIT 1` instead of `COUNT()`
