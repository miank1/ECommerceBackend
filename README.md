# Authentication & API Flow Example

## User Service Endpoints

- **Login**:  
  `POST http://localhost:8081/users/login`

- **Register**:  
  `POST http://localhost:8081/users/register`

- **Products Service**:  
  `GET http://localhost:8080/products/1`

---

## Example Product Payload

```json
{
  "name": "Test Product 22",
  "description": "A test product",
  "price": 30.99,
  "stock": 100
}

USER LOGIN FLOW

User                    User Service
  |                          |
  |--- Login Request ------->|
  |   (email + password)     |
  |                          |--- Verify Credentials
  |                          |--- Generate JWT Token
  |<-- Token Response -------|
  |   (JWT token)            |

API GATEWAY

User                    API Gateway/Service
  |                          |
  |--- API Request --------->|
  | + JWT in Header          |
  |                          |--- Verify JWT Token
  |                          |--- Check Permissions
  |                          |--- Process Request
  |<-- API Response ---------|

Example Login request

POST http://localhost:8081/users/login
Content-Type: application/json

{
    "email": "user@example.com",
    "password": "securepassword123"
}
{
    "message": "Login successful",
    "user_id": "1",
    "token": "eyJhbGciOiJIUzI1..."
}
GET http://localhost:8080/products
Authorization: Bearer eyJhbGciOiJIUzI1...

Profile Request with middleware

Client              AuthMiddleware           GetProfile Handler
   |                      |                        |
   |-- GET /profile ----->|                        |
   | + JWT token          |                        |
   |                      |-- Verify token ------->|
   |                      |-- Add user context --->|
   |                      |                        |-- Get user data
   |                      |                        |-- Return profile
   |<---- Profile Data ---|------------------------|

