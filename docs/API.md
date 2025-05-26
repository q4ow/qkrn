# qkrn API Documentation

## Overview
qkrn provides a simple HTTP REST API for key-value operations with optional API key authentication.

## Base URL
```
http://localhost:8080
```

## Authentication
qkrn supports optional API key authentication. When authentication is enabled, all endpoints except `/`, `/health` require a valid API key.

### Authentication Methods
1. **Bearer Token** (Recommended): Include in `Authorization` header
   ```
   Authorization: Bearer YOUR_API_KEY
   ```

2. **X-API-Key Header**: Include in custom header
   ```
   X-API-Key: YOUR_API_KEY
   ```

3. **Query Parameter**: Include as URL parameter (less secure)
   ```
   ?api_key=YOUR_API_KEY
   ```

### Enabling Authentication
Start the server with authentication enabled:
```bash
# With a specific API key
./qkrn --auth-enabled --api-key "your-secure-api-key"

# Auto-generate a secure API key
./qkrn --auth-enabled
```

### Error Responses
When authentication fails, you'll receive:
```json
{
  "success": false,
  "error": "Missing authentication token"
}
```
or
```json
{
  "success": false,
  "error": "Invalid authentication token"
}
```

## Endpoints

### Service Information

#### GET /
Returns basic service information including authentication status.

**Response:**
```json
{
  "service": "qkrn",
  "version": "0.1.0",
  "status": "running",
  "authentication": false
}
```

When authentication is enabled:
```json
{
  "service": "qkrn",
  "version": "0.1.0",
  "status": "running",
  "authentication": true
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

Store a value (with authentication):
```bash
curl -X PUT http://localhost:8080/kv/hello \
  -H "Authorization: Bearer YOUR_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{"value":"world"}'
```

Store a value (without authentication, if disabled):
```bash
curl -X PUT http://localhost:8080/kv/hello \
  -H "Content-Type: application/json" \
  -d '{"value":"world"}'
```

Retrieve a value:
```bash
curl -H "Authorization: Bearer YOUR_API_KEY" \
  http://localhost:8080/kv/hello
```

List all keys:
```bash
curl -H "Authorization: Bearer YOUR_API_KEY" \
  http://localhost:8080/keys
```

Delete a key:
```bash
curl -X DELETE \
  -H "Authorization: Bearer YOUR_API_KEY" \
  http://localhost:8080/kv/hello
```

Using X-API-Key header:
```bash
curl -X PUT http://localhost:8080/kv/hello \
  -H "X-API-Key: YOUR_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{"value":"world"}'
```

## Error Codes

- `200` - Success
- `201` - Created (for PUT operations)
- `400` - Bad Request (invalid JSON, empty key)
- `401` - Unauthorized (missing or invalid authentication token)
- `404` - Not Found (key doesn't exist)
- `405` - Method Not Allowed
- `500` - Internal Server Error

## Security Notes

- API keys should be kept secure and not exposed in logs or URLs when possible
- Use HTTPS in production environments
- The Bearer token method is preferred over query parameters
- API keys are generated using cryptographically secure random number generation
- Token validation uses constant-time comparison to prevent timing attacks
