# API Reference

REST API for managing Todo items. The API uses JSON for both requests and responses.

## Base URL
```
http://localhost:8080/api/v1
```

## Authentication
Currently, no authentication is required for API access.

## Endpoints

### Todos

#### Create Todo
**POST** `/todos`

Creates a new todo item.

**Request Body:**
```json
{
    "title": "Task title",
    "description": "Task description (optional)"
}
```

**Request Fields:**
- `title` (string, required) - The title of the todo
- `description` (string, optional) - The description of the todo

**Response:**
```json
{
    "id": 1,
    "title": "Task title",
    "description": "Task description",
    "completed": false,
    "created_at": "2023-01-01T12:00:00Z",
    "updated_at": "2023-01-01T12:00:00Z"
}
```

**Status Codes:**
- `201 Created` - Todo successfully created
- `400 Bad Request` - Invalid request data
- `409 Conflict` - Todo with this title already exists
- `500 Internal Server Error` - Server error

---

#### Get All Todos
**GET** `/todos`

Returns a list of all todo items.

**Response:**
```json
[
    {
        "id": 1,
        "title": "First task",
        "description": "Description of first task",
        "completed": false,
        "created_at": "2023-01-01T12:00:00Z",
        "updated_at": "2023-01-01T12:00:00Z"
    },
    {
        "id": 2,
        "title": "Second task",
        "description": "Description of second task",
        "completed": true,
        "created_at": "2023-01-01T13:00:00Z",
        "updated_at": "2023-01-01T14:00:00Z"
    }
]
```

**Status Codes:**
- `200 OK` - List of todos successfully retrieved
- `500 Internal Server Error` - Server error

---

#### Get Todo by ID
**GET** `/todos/{id}`

Returns a specific todo by its ID.

**Parameters:**
- `id` (integer) - Todo identifier

**Response:**
```json
{
    "id": 1,
    "title": "Task title",
    "description": "Task description",
    "completed": false,
    "created_at": "2023-01-01T12:00:00Z",
    "updated_at": "2023-01-01T12:00:00Z"
}
```

**Status Codes:**
- `200 OK` - Todo found
- `400 Bad Request` - Invalid ID
- `404 Not Found` - Todo not found
- `500 Internal Server Error` - Server error

---

#### Update Todo
**PUT** `/todos/{id}`

Updates an existing todo item.

**Parameters:**
- `id` (integer) - Todo identifier

**Request Body:**
```json
{
    "title": "Updated task title",
    "description": "Updated task description",
    "completed": true
}
```

**Request Fields:**
- `title` (string, optional) - New title for the todo
- `description` (string, optional) - New description for the todo
- `completed` (boolean, optional) - Completion status of the todo

**Response:**
```json
{
    "id": 1,
    "title": "Updated task title",
    "description": "Updated task description",
    "completed": true,
    "created_at": "2023-01-01T12:00:00Z",
    "updated_at": "2023-01-01T15:00:00Z"
}
```

**Status Codes:**
- `200 OK` - Todo successfully updated
- `400 Bad Request` - Invalid data or ID
- `404 Not Found` - Todo not found
- `409 Conflict` - Todo with this title already exists
- `500 Internal Server Error` - Server error

---

#### Delete Todo
**DELETE** `/todos/{id}`

Deletes a todo by ID.

**Parameters:**
- `id` (integer) - Todo identifier

**Response:**
```json
{
    "message": "Todo deleted successfully"
}
```

**Status Codes:**
- `200 OK` - Todo successfully deleted
- `400 Bad Request` - Invalid ID
- `404 Not Found` - Todo not found
- `500 Internal Server Error` - Server error

## Data Structures

### Todo
```json
{
    "id": 1,
    "title": "Task title",
    "description": "Task description",
    "completed": false,
    "created_at": "2023-01-01T12:00:00Z",
    "updated_at": "2023-01-01T12:00:00Z"
}
```

**Fields:**
- `id` (integer) - Unique identifier
- `title` (string) - Todo title
- `description` (string) - Todo description
- `completed` (boolean) - Completion status (default: false)
- `created_at` (datetime) - Creation timestamp
- `updated_at` (datetime) - Last update timestamp

### ErrorResponse
```json
{
    "error": "Error description"
}
```

### MessageResponse
```json
{
    "message": "Success message"
}
```

## Usage Examples

### Creating a Todo
```bash
curl -X POST http://localhost:8080/api/v1/todos \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Learn Go",
    "description": "Complete Go programming course"
  }'
```

### Getting All Todos
```bash
curl -X GET http://localhost:8080/api/v1/todos
```

### Updating a Todo
```bash
curl -X PUT http://localhost:8080/api/v1/todos/1 \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Learn Go (completed)",
    "completed": true
  }'
```

### Deleting a Todo
```bash
curl -X DELETE http://localhost:8080/api/v1/todos/1
```

## Error Codes

| Code | Description |
|------|-------------|
| 400  | Bad Request - Invalid request data |
| 404  | Not Found - Resource not found |
| 409  | Conflict - Conflict (e.g., duplicate title) |
| 500  | Internal Server Error - Server error |

## Business Rules

1. **Unique Titles**: Todo titles must be unique across all non-deleted todos
2. **Required Fields**: Title is required when creating a todo
3. **Soft Delete**: Todos are soft-deleted (marked as deleted but not removed from database)
4. **Timestamps**: All todos have creation and update timestamps

## Swagger/OpenAPI

For interactive API documentation, you can use Swagger UI:

1. Enable Swagger in `.env`:
   ```
   ENABLE_SWAGGER=true
   ```

2. Visit: http://localhost:8080/swagger/index.html

The Swagger documentation is auto-generated from code annotations.

## Rate Limiting

Rate limiting can be enabled via configuration. When enabled, it applies globally to all endpoints.

## CORS

CORS is configurable and can be set to allow specific origins or credentials as needed.