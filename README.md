# user-management

## Getting Started

Follow these instructions to set up and run the application on your local machine.

### Prerequisites

Ensure you have the following installed:
- [Docker](https://www.docker.com/) and [Docker Compose](https://docs.docker.com/compose/)

### Project Setup

1. Clone the repository:
   ```bash
   git clone https://github.com/nath0227/user-management.git
   cd user-management
   ```

2. Create a `.env` file in the root directory to customize the configuration:
   - `HTTP_SERVER_PORT`: Port for the HTTP server.
   - `GRPC_SERVER_PORT`: Port for the gRPC server.
   - `CRYPTO_JWT_KEY`: Secret key for signing JWT tokens.
   - `CRYPTO_JWT_EXPIRE_DURATION`: Duration before the JWT token expires.
   - `MONGO_CONFIG_URI`: MongoDB connection string.
   - `MONGO_CONFIG_USERNAME`: MongoDB username (ensure it matches the configuration in `mongo-init/init.js`).
   - `MONGO_CONFIG_PASSWORD`: MongoDB password (ensure it matches the configuration in `mongo-init/init.js`).
   - `MONGO_CONFIG_DATABASE`: MongoDB database name (ensure it matches the configuration in `mongo-init/init.js`).
   - `MONGO_CONFIG_USER_COLLECTION`: MongoDB collection name for users (ensure it matches the configuration in `mongo-init/init.js`).
   - `USER_COUNT_INTERVAL`: Interval duration for logging the user count.

3. Start the application using Docker Compose:
   ```bash
   docker-compose up
   ```

   This will start:
   - The Go application on `http://localhost:8080` (REST API) and `http://localhost:50051` (gRPC).
   - MongoDB on `localhost:27017`.

### Using the Application

The application provides both REST API and gRPC endpoints for managing users.

#### REST API

The REST API allows you to perform user management operations such as registering, logging in, and managing user details. The API documentation is available in Swagger format.

- **Swagger Documentation**: Once the application is running, you can access the Swagger UI at:
  ```
  http://localhost:8081/
  ```

- **Sample REST API Requests**:
  Below are examples of REST API requests using `curl` commands:

  1. **Login**
     ```bash
     curl -X POST http://localhost:8080/login \
     -H "Content-Type: application/json" \
     -d '{"email": "admin@example.com", "password": "passwordstring"}'
     ```

  2. **Register a User**
     ```bash
     curl -X POST http://localhost:8080/register \
     -H "Content-Type: application/json" \
     -H "Authorization: Bearer <your_jwt_token>" \
     -d '{"name": "testuser", "email": "test@example.com", "password": "password123"}'
     ```

  3. **Get All Users (Protected)**
     ```bash
     curl -X GET http://localhost:8080/users \
     -H "Authorization: Bearer <your_jwt_token>"
     ```

  4. **Get User by ID (Protected)**
     ```bash
     curl -X GET http://localhost:8080/users/{id} \
     -H "Authorization: Bearer <your_jwt_token>"
     ```

  5. **Update User (Protected)**
     ```bash
     curl -X PUT http://localhost:8080/users/{id} \
     -H "Authorization: Bearer <your_jwt_token>" \
     -H "Content-Type: application/json" \
     -d '{"name": "updateduser", "email": "updated@example.com"}'
     ```

  6. **Delete User (Protected)**
     ```bash
     curl -X DELETE http://localhost:8080/users/{id} \
     -H "Authorization: Bearer <your_jwt_token>"
     ```

#### gRPC

The application also provides gRPC endpoints for user management. Currently, only the following methods are available:
- `CreateUser`
- `GetUser`

These endpoints are defined in the `.proto` files located in:
```
app/user/grpc/proto
```

- **gRPC Server**:
  - Runs on port `50051` by default.
  - Postman is used for testing the gRPC endpoints. Ensure you have the Postman gRPC client set up.

- **Using Postman for gRPC**:
  1. Open Postman, click the **New** button, and select **gRPC**.
  2. Enter the gRPC server address:
     ```
     localhost:50051
     ```
  3. Import the `.proto` file:
     - Go to the **Service definition** tab in the request.
     - Select the `.proto` file from the `app/user/grpc/proto` directory.
  4. Select the desired method (e.g., `CreateUser` or `GetUserById`) from the dropdown.
  5. Provide the request body in JSON format (see examples below).
  6. Click **Invoke** to execute the request and view the response.

- **Sample gRPC Requests**:
  1. **CreateUser**:
     - Method: `CreateUser`
     - Request Body:
       ```json
       {
         "name": "John Doe",
         "email": "john.doe@example.com",
         "password": "password123"
       }
       ```

  2. **GetUser**:
     - Method: `GetUser`
     - Request Body:
       ```json
       {
         "id": "60d5ec49f1f1c939b4f2f0c2"
       }
       ```

  Refer to the `.proto` files for the full request/response structure.

### Stopping the Application

To stop the backend and MongoDB, press `Ctrl+C` in the terminal running `docker-compose up`, then run:
```bash
docker-compose down
```

### Assumptions and Decisions

1. The application uses MongoDB as the database.
2. JWT tokens are used for authentication and must be included in the `Authorization` header for protected endpoints.
3. The `.env` file is used for configuration, and sensitive information like `CRYPTO_JWT_KEY` should not be hardcoded.
4. The application runs on port `8080` for REST API and `50051` for gRPC by default, but these can be customized in the `.env` file.
