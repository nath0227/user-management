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

2. Create a `.env` file in the root directory to customize configuration:
   ```env
   MONGO_URI=mongodb://mongo:27017/user-management
   JWT_SECRET=your_jwt_secret
   APP_PORT=8080
   ```

   - `MONGO_URI`: MongoDB connection string.
   - `JWT_SECRET`: Secret key for signing JWT tokens.
   - `APP_PORT`: Port on which the Go application will run.

3. Start the application using Docker Compose:
   ```bash
   docker-compose up
   ```

   This will start:
   - The Go application on `http://localhost:8080`
   - MongoDB on `localhost:27017`

### Using the Application

The application provides a REST API for managing users. Below is a guide for using JWT tokens and sample API requests.

#### JWT Token Usage

1. Obtain a JWT token by logging in or registering a user.
2. Include the token in the `Authorization` header for protected endpoints:
   ```
   Authorization: Bearer <your_token>
   ```

#### Sample API Requests/Responses

Below are examples of API requests using `curl` commands from `./tool/execute.curl`:

1. **Register a User**
   - **Command**:
     ```bash
     curl -X POST http://localhost:8080/api/register \
     -H "Content-Type: application/json" \
     -d '{"username": "testuser", "password": "password123"}'
     ```
   - **Response**:
     ```json
     {
       "message": "User registered successfully"
     }
     ```

2. **Login**
   - **Command**:
     ```bash
     curl -X POST http://localhost:8080/api/login \
     -H "Content-Type: application/json" \
     -d '{"username": "testuser", "password": "password123"}'
     ```
   - **Response**:
     ```json
     {
       "token": "your_jwt_token"
     }
     ```

3. **Get User Details (Protected)**
   - **Command**:
     ```bash
     curl -X GET http://localhost:8080/api/user \
     -H "Authorization: Bearer your_jwt_token"
     ```
   - **Response**:
     ```json
     {
       "id": "12345",
       "username": "testuser"
     }
     ```

4. **Update User (Protected)**
   - **Command**:
     ```bash
     curl -X PUT http://localhost:8080/api/user \
     -H "Authorization: Bearer your_jwt_token" \
     -H "Content-Type: application/json" \
     -d '{"username": "updateduser"}'
     ```
   - **Response**:
     ```json
     {
       "message": "User updated successfully"
     }
     ```

5. **Delete User (Protected)**
   - **Command**:
     ```bash
     curl -X DELETE http://localhost:8080/api/user \
     -H "Authorization: Bearer your_jwt_token"
     ```
   - **Response**:
     ```json
     {
       "message": "User deleted successfully"
     }
     ```

### Stopping the Application

To stop the backend and MongoDB, press `Ctrl+C` in the terminal running `docker-compose up`, then run:
```bash
docker-compose down
```

### Assumptions and Decisions

1. The application uses MongoDB as the database.
2. JWT tokens are used for authentication and must be included in the `Authorization` header for protected endpoints.
3. The `.env` file is used for configuration, and sensitive information like `JWT_SECRET` should not be hardcoded.
4. The application runs on port `8080` by default, but this can be customized in the `.env` file.

### Contributing

Feel free to fork the repository and submit pull requests for improvements or bug fixes.

### License

This project is licensed under the MIT License.