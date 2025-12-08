# API Endpoints - cURL Examples

Base URL: `http://localhost:8080`

## 1. Health Check

### GET /health
Check if the service is running and database is connected.

```bash
curl -X GET http://localhost:8080/health
```

**Expected Response:**
```json
{
  "status": "healthy"
}
```

---

## 2. Create User

### POST /api/v1/users
Create a new user account.

```bash
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "username": "johndoe",
    "password": "password123"
  }'
```

**Expected Response (201 Created):**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "email": "user@example.com",
  "username": "johndoe",
  "created_at": "2025-12-07 20:58:10"
}
```

**Error Responses:**
- `400 Bad Request` - Invalid input (email format, username too short, password too short)
- `409 Conflict` - Email already exists

---

## 3. Login

### POST /api/v1/login
Authenticate user and get JWT token.

```bash
curl -X POST http://localhost:8080/api/v1/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password123"
  }'
```

**Expected Response (200 OK):**
```json
{
  "token": "jwt-token-example",
  "user": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "email": "user@example.com",
    "username": "johndoe",
    "created_at": "2025-12-07 20:58:10"
  }
}
```

**Error Responses:**
- `400 Bad Request` - Invalid input
- `401 Unauthorized` - Invalid credentials

---

## 4. Get User by ID

### GET /api/v1/users/:id
Get user information by UUID.

```bash
curl -X GET http://localhost:8080/api/v1/users/550e8400-e29b-41d4-a716-446655440000
```

**Expected Response (200 OK):**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "email": "user@example.com",
  "username": "johndoe",
  "created_at": "2025-12-07 20:58:10"
}
```

**Error Responses:**
- `400 Bad Request` - Invalid UUID format
- `404 Not Found` - User not found

---

## 5. Wallet Operations

### POST /api/v1/wallet
Process a wallet operation (DEPOSIT or WITHDRAW).

**Request Body:**
```json
{
  "walletId": "550e8400-e29b-41d4-a716-446655440000",
  "operationType": "DEPOSIT",
  "amount": 1000
}
```

**cURL (Bash/Linux/Mac):**
```bash
curl -X POST http://localhost:8080/api/v1/wallet \
  -H "Content-Type: application/json" \
  -d '{
    "walletId": "550e8400-e29b-41d4-a716-446655440000",
    "operationType": "DEPOSIT",
    "amount": 1000
  }'
```

**Expected Response (200 OK):**
```json
{
  "message": "operation completed successfully",
  "walletId": "550e8400-e29b-41d4-a716-446655440000",
  "operationType": "DEPOSIT",
  "amount": 1000
}
```

**Error Responses:**
- `400 Bad Request` - Invalid input, insufficient funds, or invalid operation type
- `404 Not Found` - Wallet not found
- `500 Internal Server Error` - Server error

---

## 6. Get Wallet Balance

### GET /api/v1/wallets/:walletId
Get wallet information and balance by wallet UUID.

**cURL:**
```bash
curl -X GET http://localhost:8080/api/v1/wallets/550e8400-e29b-41d4-a716-446655440000
```

**Expected Response (200 OK):**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "user_id": "660e8400-e29b-41d4-a716-446655440000",
  "balance": 1000.00,
  "created_at": "2025-12-07 20:58:10",
  "updated_at": "2025-12-07 21:30:00"
}
```

**Error Responses:**
- `400 Bad Request` - Invalid UUID format
- `404 Not Found` - Wallet not found

---

## Complete Example Workflow

### Step 1: Check health
```bash
curl -X GET http://localhost:8080/health
```

### Step 2: Create a user
```bash
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "username": "testuser",
    "password": "testpass123"
  }'
```

**Save the user ID from the response for the next steps.**

### Step 3: Login
```bash
curl -X POST http://localhost:8080/api/v1/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "testpass123"
  }'
```

**Save the token from the response if you need it for authenticated endpoints.**

### Step 4: Get user by ID
```bash
# Replace USER_ID with the actual UUID from Step 2
curl -X GET http://localhost:8080/api/v1/users/USER_ID
```

---

## Pretty Print JSON Responses

Add `| jq` to format JSON responses (requires jq installed):

```bash
curl -X GET http://localhost:8080/health | jq
```

Or use Python for pretty printing:

```bash
curl -X GET http://localhost:8080/health | python -m json.tool
```

---

## Windows PowerShell Examples

### Recommended: Use Invoke-RestMethod (Native PowerShell)

PowerShell has issues with curl.exe arguments. Use `Invoke-RestMethod` for better compatibility:

#### Health Check
```powershell
Invoke-RestMethod -Uri "http://localhost:8080/health" -Method GET
```

#### Wallet Operation (DEPOSIT/WITHDRAW)
```powershell
$body = @{
    walletId = "550e8400-e29b-41d4-a716-446655440000"
    operationType = "DEPOSIT"
    amount = 1000
} | ConvertTo-Json

Invoke-RestMethod -Uri "http://localhost:8080/api/v1/wallet" -Method POST -ContentType "application/json" -Body $body
```

#### Get Wallet Balance
```powershell
Invoke-RestMethod -Uri "http://localhost:8080/api/v1/wallets/550e8400-e29b-41d4-a716-446655440000" -Method GET
```

#### Create User
```powershell
$body = @{
    email = "user@example.com"
    username = "johndoe"
    password = "password123"
} | ConvertTo-Json

Invoke-RestMethod -Uri "http://localhost:8080/api/v1/users" -Method POST -ContentType "application/json" -Body $body
```

Or as a one-liner:
```powershell
Invoke-RestMethod -Uri "http://localhost:8080/api/v1/users" -Method POST -ContentType "application/json" -Body '{"email":"user@example.com","username":"johndoe","password":"password123"}'
```

#### Login
```powershell
Invoke-RestMethod -Uri "http://localhost:8080/api/v1/login" -Method POST -ContentType "application/json" -Body '{"email":"user@example.com","password":"password123"}'
```

#### Get User by ID
```powershell
Invoke-RestMethod -Uri "http://localhost:8080/api/v1/users/550e8400-e29b-41d4-a716-446655440000" -Method GET
```

---

### Alternative: Using curl.exe (Works in CMD, Git Bash, or WSL)

For simple GET requests in PowerShell:
```powershell
curl.exe http://localhost:8080/health
```

For POST requests, use a JSON file or use `Invoke-RestMethod` instead. PowerShell has issues parsing curl.exe arguments with JSON data.

**Create a JSON file (request.json):**
```json
{
  "email": "user@example.com",
  "username": "johndoe",
  "password": "password123"
}
```

Then use:
```powershell
curl.exe -X POST http://localhost:8080/api/v1/users -H "Content-Type: application/json" -d "@request.json"
```

---

### Using Invoke-RestMethod (Native PowerShell)

#### Create User
```powershell
Invoke-RestMethod -Uri "http://localhost:8080/api/v1/users" `
  -Method POST `
  -ContentType "application/json" `
  -Body '{"email":"user@example.com","username":"johndoe","password":"password123"}'
```

#### Login
```powershell
Invoke-RestMethod -Uri "http://localhost:8080/api/v1/login" `
  -Method POST `
  -ContentType "application/json" `
  -Body '{"email":"user@example.com","password":"password123"}'
```

#### Get User
```powershell
Invoke-RestMethod -Uri "http://localhost:8080/api/v1/users/550e8400-e29b-41d4-a716-446655440000" `
  -Method GET
```

#### Health Check
```powershell
Invoke-RestMethod -Uri "http://localhost:8080/health" -Method GET
```

#### Wallet Operation (DEPOSIT/WITHDRAW)
```powershell
$body = @{
    walletId = "550e8400-e29b-41d4-a716-446655440000"
    operationType = "DEPOSIT"
    amount = 1000
} | ConvertTo-Json

Invoke-RestMethod -Uri "http://localhost:8080/api/v1/wallet" -Method POST -ContentType "application/json" -Body $body
```

#### Get Wallet Balance
```powershell
Invoke-RestMethod -Uri "http://localhost:8080/api/v1/wallets/550e8400-e29b-41d4-a716-446655440000" -Method GET
```

