# qkrn API Documentation

## Overview
qkrn provides a simple HTTP REST API for key-value operations.

## Base URL
```
http://localhost:8080
```

## Endpoints

### Service Information

#### GET /
Returns basic service information.

**Response:**
```json
{
  "service": "qkrn",
  "version": "0.1.0",
  "status": "running"
}
```

#### GET /health
Health check endpoint.

**Response:**
```json
{
  "status": "healthy"
}
```

### Key-Value Operations

#### GET /kv/{key}
Retrieve a value by key.

**Parameters:**
- `key` (path): The key to retrieve

**Response:**
```json
{
  "success": true,
  "value": "stored_value"
}
```

**Error Response (404):**
```json
{
  "success": false,
  "error": "Key not found"
}
```

#### PUT /kv/{key}
Store a key-value pair.

**Parameters:**
- `key` (path): The key to store

**Request Body:**
```json
{
  "value": "new_value"
}
```

**Response (201):**
```json
{
  "success": true
}
```

#### DELETE /kv/{key}
Delete a key-value pair.

**Parameters:**
- `key` (path): The key to delete

**Response:**
```json
{
  "success": true
}
```

**Error Response (404):**
```json
{
  "success": false,
  "error": "Key not found"
}
```

#### GET /keys
List all keys in the store.

**Response:**
```json
{
  "keys": ["key1", "key2", "key3"]
}
```

## Examples

### Using curl

Store a value:
```bash
curl -X PUT http://localhost:8080/kv/hello \
  -H "Content-Type: application/json" \
  -d '{"value":"world"}'
```

Retrieve a value:
```bash
curl http://localhost:8080/kv/hello
```

List all keys:
```bash
curl http://localhost:8080/keys
```

Delete a key:
```bash
curl -X DELETE http://localhost:8080/kv/hello
```

## Error Codes

- `200` - Success
- `201` - Created (for PUT operations)
- `400` - Bad Request (invalid JSON, empty key)
- `404` - Not Found (key doesn't exist)
- `405` - Method Not Allowed
- `500` - Internal Server Error
