http://localhost:8081/users/login
http://localhost:8081/users/register
http://localhost:8080/products/1
{
    "name": "Test Product 22",
    "description": "A test product",
    "price": 30.99,
    "stock": 100
}


User                    User Service
  |                          |
  |--- Login Request ------->|
  |   (email + password)     |
  |                          |--- Verify Credentials
  |                          |--- Generate JWT Token
  |<-- Token Response -------|
  |   (JWT token)            |


  User                    API Gateway/Service
  |                          |
  |--- API Request --------->|
  | + JWT in Header          |
  |                          |--- Verify JWT Token
  |                          |--- Check Permissions
  |                          |--- Process Request
  |<-- API Response ---------|


  POST http://localhost:8081/users/login
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
Header: Authorization: Bearer eyJhbGciOiJIUzI1...


Client              AuthMiddleware           GetProfile Handler
   |                      |                        |
   |-- GET /profile ----->|                        |
   | + JWT token          |                        |
   |                      |-- Verify token ------->|
   |                      |-- Add user context --->|
   |                      |                        |-- Get user data
   |                      |                        |-- Return profile
   |<---- Profile Data ---|------------------------|