# Database Schema Documentation

This document describes the database schema, design decisions, and data management strategies used in the Todo API.

## Database Technology

### SQLite
- **File-based** database stored in `data/app.db`
- **ACID compliant** with transaction support
- **Zero configuration** - no separate database server required
- **CGO enabled** for better performance
- **In-memory** databases for testing

### GORM ORM
- **Auto-migration** - schema updates applied automatically
- **Soft delete** support with `deleted_at` timestamps
- **Relationship** management
- **Query optimization** with prepared statements

## Schema Overview

### Tables

Currently, the application has one main table:

#### `todos` Table

## Table Structure

### todos

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| `id` | INTEGER | PRIMARY KEY, AUTO_INCREMENT | Unique identifier |
| `title` | TEXT | NOT NULL, UNIQUE (where not deleted) | Todo title |
| `description` | TEXT | | Optional description |
| `completed` | BOOLEAN | DEFAULT FALSE | Completion status |
| `created_at` | DATETIME | NOT NULL | Record creation timestamp |
| `updated_at` | DATETIME | NOT NULL | Last update timestamp |
| `deleted_at` | DATETIME | NULL | Soft delete timestamp |

### GORM Entity Definition

```go
type Todo struct {
    ID          uint           `json:"id" gorm:"primaryKey"`
    Title       string         `json:"title" gorm:"type:text;uniqueIndex:idx_todos_title_not_deleted,where:deleted_at IS NULL;not null"`
    Description string         `json:"description"`
    Completed   bool           `json:"completed" gorm:"default:false"`
    CreatedAt   time.Time      `json:"created_at"`
    UpdatedAt   time.Time      `json:"updated_at"`
    DeletedAt   gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index" swaggerignore:"true"`
}
```

## Indexes

### Primary Index
- **`PRIMARY`** on `id` column
- Automatically created for primary key
- Ensures unique identification of records

### Unique Indexes
- **`idx_todos_title_not_deleted`** - Partial unique index
  - Columns: `title`
  - Condition: `WHERE deleted_at IS NULL`
  - Purpose: Ensures unique titles among active todos

### Regular Indexes
- **`idx_todos_deleted_at`** on `deleted_at` column
- Optimizes soft delete queries
- Helps with filtering out deleted records

## Database Features

### Soft Delete Implementation

The application implements soft delete using GORM's built-in functionality:

```sql
-- When deleting a todo (soft delete)
UPDATE todos SET deleted_at = '2023-01-01 12:00:00' WHERE id = 1;

-- Queries automatically exclude soft-deleted records
SELECT * FROM todos WHERE deleted_at IS NULL;
```

**Benefits:**
- **Data Recovery** - Accidentally deleted records can be restored
- **Audit Trail** - Historical data is preserved
- **Referential Integrity** - Related data remains intact

### Unique Constraints

#### Title Uniqueness
- Titles must be unique among **non-deleted** todos only
- Uses partial unique index with `WHERE deleted_at IS NULL`
- Allows the same title to be reused after deletion

```sql
-- This index ensures the constraint
CREATE UNIQUE INDEX idx_todos_title_not_deleted 
ON todos(title) WHERE deleted_at IS NULL;
```

### Timestamps

#### Automatic Timestamp Management
- **`created_at`** - Set once when record is created
- **`updated_at`** - Updated automatically on every change
- **`deleted_at`** - Set when record is soft-deleted

```go
// GORM automatically manages these timestamps
todo.CreatedAt = time.Now()  // On create
todo.UpdatedAt = time.Now()  // On update
todo.DeletedAt = time.Now()  // On delete
```

## Data Types and Validation

### Field Specifications

#### ID Field
- **Type**: `uint` (unsigned integer)
- **Database**: `INTEGER PRIMARY KEY`
- **Auto-increment**: Yes
- **Range**: 1 to 4,294,967,295

#### Title Field
- **Type**: `string`
- **Database**: `TEXT`
- **Constraints**: NOT NULL, Unique among active records
- **Validation**: Required at application level
- **Max Length**: SQLite TEXT has no practical limit

#### Description Field
- **Type**: `string`
- **Database**: `TEXT`
- **Constraints**: None
- **Validation**: Optional
- **Max Length**: SQLite TEXT has no practical limit

#### Completed Field
- **Type**: `bool`
- **Database**: `BOOLEAN`
- **Default**: `false`
- **Values**: `true` (1) or `false` (0)

#### Timestamp Fields
- **Type**: `time.Time`
- **Database**: `DATETIME`
- **Format**: ISO 8601 format
- **Timezone**: UTC (recommended)

## Database Configuration

### SQLite PRAGMA Settings

The application configures SQLite with performance optimizations:

```sql
PRAGMA journal_mode = WAL;          -- Write-Ahead Logging
PRAGMA synchronous = NORMAL;        -- Balance performance/safety
PRAGMA cache_size = 1000;          -- Cache size in pages
PRAGMA temp_store = memory;        -- Temp tables in memory
PRAGMA mmap_size = 268435456;      -- Memory-mapped I/O
```

### Connection Pool Settings

```go
// GORM connection pool configuration
sqlDB, _ := db.DB()
sqlDB.SetMaxIdleConns(10)
sqlDB.SetMaxOpenConns(100)
sqlDB.SetConnMaxLifetime(time.Hour)
```

## Query Patterns

### Common Queries

#### Create Todo
```sql
INSERT INTO todos (title, description, completed, created_at, updated_at) 
VALUES (?, ?, ?, ?, ?);
```

#### Get All Active Todos
```sql
SELECT id, title, description, completed, created_at, updated_at 
FROM todos 
WHERE deleted_at IS NULL 
ORDER BY created_at DESC;
```

#### Get Todo by ID
```sql
SELECT id, title, description, completed, created_at, updated_at 
FROM todos 
WHERE id = ? AND deleted_at IS NULL;
```

#### Update Todo
```sql
UPDATE todos 
SET title = ?, description = ?, completed = ?, updated_at = ? 
WHERE id = ? AND deleted_at IS NULL;
```

#### Soft Delete Todo
```sql
UPDATE todos 
SET deleted_at = ? 
WHERE id = ? AND deleted_at IS NULL;
```

#### Check Title Exists
```sql
SELECT id 
FROM todos 
WHERE title = ? AND deleted_at IS NULL 
LIMIT 1;
```

## Performance Considerations

### Query Optimization
- **Prepared Statements** - All queries use prepared statements
- **Index Usage** - Queries designed to use available indexes
- **Limit Results** - Pagination support ready for future implementation

### Memory Optimization
- **Connection Pooling** - Reuse database connections
- **Statement Caching** - GORM caches prepared statements
- **Minimal Columns** - Select only needed columns when possible

## Data Migration

### Auto-Migration
GORM handles schema changes automatically:

```go
// Auto-migrate on application startup
db.AutoMigrate(&Todo{})
```

### Migration Process
1. **Compare** current schema with entity definitions
2. **Add** new columns if needed
3. **Create** new indexes if needed
4. **Preserve** existing data

### Manual Migration
For complex changes, manual migration might be needed:

```sql
-- Example: Adding a new column
ALTER TABLE todos ADD COLUMN priority INTEGER DEFAULT 0;

-- Example: Creating a new index
CREATE INDEX idx_todos_priority ON todos(priority);
```

## Backup and Recovery

### Database Backup
```bash
# Create backup
sqlite3 data/app.db ".backup backup_$(date +%Y%m%d_%H%M%S).db"

# Or using file copy (when application is stopped)
cp data/app.db backup_$(date +%Y%m%d_%H%M%S).db
```

### Database Recovery
```bash
# Restore from backup
sqlite3 data/app.db ".restore backup_20230101_120000.db"

# Or using file copy
cp backup_20230101_120000.db data/app.db
```

### Data Export
```bash
# Export to SQL
sqlite3 data/app.db .dump > todos_export.sql

# Export to CSV
sqlite3 data/app.db <<< ".mode csv
.output todos_export.csv
SELECT * FROM todos WHERE deleted_at IS NULL;"
```

## Security Considerations

### SQL Injection Prevention
- **Prepared Statements** - All queries use parameterized statements
- **GORM ORM** - Provides additional protection
- **Input Validation** - Application-level validation prevents malicious input

### Data Integrity
- **ACID Properties** - SQLite provides ACID compliance
- **Constraints** - Database-level constraints enforce data integrity
- **Transactions** - Multi-step operations wrapped in transactions

### Access Control
- **File Permissions** - Database file protected by OS permissions
- **Application Layer** - Access control implemented in application

## Monitoring and Maintenance

### Database Statistics
```sql
-- Check table sizes
SELECT 
    name,
    COUNT(*) as total_rows,
    COUNT(CASE WHEN deleted_at IS NULL THEN 1 END) as active_rows,
    COUNT(CASE WHEN deleted_at IS NOT NULL THEN 1 END) as deleted_rows
FROM todos;

-- Check index usage
PRAGMA index_info(idx_todos_title_not_deleted);
```

### Maintenance Tasks
```bash
# Analyze database (update statistics)
sqlite3 data/app.db "ANALYZE;"

# Vacuum database (reclaim space)
sqlite3 data/app.db "VACUUM;"

# Check integrity
sqlite3 data/app.db "PRAGMA integrity_check;"
```

## Future Considerations

### Potential Schema Changes
- **User Management** - Add user_id foreign key to todos
- **Categories** - Add category support with separate table
- **Priority Levels** - Add priority field with enum values
- **Due Dates** - Add due_date timestamp field
- **Tags** - Many-to-many relationship with tags table

### Scaling Considerations
- **Pagination** - Add offset/limit support for large datasets
- **Archival** - Move old deleted records to archive tables
- **Partitioning** - Consider date-based partitioning for large datasets
- **Read Replicas** - Consider read replicas for high-load scenarios

### Migration to Other Databases
The current design is database-agnostic enough to migrate to:
- **PostgreSQL** - For better concurrent access
- **MySQL** - For familiarity and ecosystem
- **MongoDB** - For document-based storage needs

Database migration would primarily require:
1. Connection string changes
2. Driver changes in GORM
3. Minor SQL dialect adjustments