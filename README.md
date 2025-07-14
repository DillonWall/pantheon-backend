# Pantheon Lab Programming Assignment: Backend
A robust GraphQL backend for user authentication and multi-source image search, built as part of the Pantheon Lab Programming Assignment.


## Features
✅ Concurrent Image Search from Multiple APIs  
✅ Secure Password Storage (bcrypt hashing)  
✅ Configurable Timeouts via Environment Variables  
✅ GraphQL Playground for Interactive Testing  
### Optional Tasks
✅ [Level: Intermediate] Implement an authentication feature with Registration, Login, and token-based access control for the image API (only registered users with valid tokens can access the API).   
✅ [Level: Hard] Use GraphQL instead of REST for the web services layer.   
✅ [Level: Challenging] Build and package the application as a Docker image via a CI/CD pipeline (e.g., GitHub Actions).   
❌ [Level: Very Challenging] Create proper infrastructure and deploy the application to a cloud provider (e.g., AWS, Google Cloud, Azure) with documentation on the implementation process.   
❌ [Level: Insane] Optimize the application to handle high traffic (> 1000 requests/second) without performance degradation, including load test methods and results.   

#### Note:
To achieve the last 2 optional tasks, I originally planned to use AWS Lambda functions and DynamoDB to host the api and deploy automatically with Serverless. However, it seems no one has really done this before and it would take too long to get working for the scope of this assignment :( . I did tinker with this a bit and got somewhat close, so please let me know if you want to see the code for that.

## Prerequisites
- Go v1.24.5 or equivalent installed (https://golang.org/doc/install)  
- Environment variables configured (see below)  


## Quick Start

### 1. Clone the Repository
```
git clone https://github.com/DillonWall/pantheon-backend.git
cd pantheon-backend
```

### 2. Install Dependencies
```
go mod download
```

### 3. Configure Environment Variables
Create a `.env` file (or set these in your shell) based on `example.env`

### 4. Run the Server
```
go run server.go
```

### 5. Access the GraphQL Playground
Open your browser and navigate to:  
`http://localhost:8080` (or the port configured in your server)  

## GraphQL API Reference

### Authentication

#### Register a New User
```graphql
mutation RegisterUser {
  register(input: {
    username: "newuser",
    password: "securePassword123"
  }) {
    token  # JWT token for subsequent requests. Use this token for protected endpoints
  }
}
```

#### Login
```graphql
mutation LoginUser {
  login(input: {
    username: "newuser",
    password: "securePassword123"
  }) {
    token  # JWT token for subsequent requests. Use this token for protected endpoints
  }
}
```

#### Verify a Token
```graphql
mutation VerifyToken {
  verify(token: "YOUR_JWT_TOKEN")  # Returns true if valid
}
```


### Image Search
```graphql
query SearchImages {
  searchImages(
    token: "YOUR_JWT_TOKEN",  # From login/register
    query: "yellow flower"    # Search term
  ) {
    image_ID       # Unique image ID
    title          # Image title/description
    thumbnails     # URL to thumbnail image
    preview        # URL to preview image
    source         # Provider (i.e., "UNSPLASH", "PIXABAY", "STORYBLOCKS")
    tags           # Related tags
  }
}
```


## Project Structure
```
pantheon-backend/
├── graph/                        # GraphQL schema and resolvers
│   ├── model/                    # GraphQL type definitions
│   └── schema.resolvers.go       # Resolver implementations (auth, image search)
├── pkg/
│   ├── auth/                     # Authentication logic (JWT, password hashing)
│   └── imageapi/                 # Image API clients (Unsplash, Pixabay, Storyblocks)
├── scripts/                      # Useful scripts (Generate JWT secret)
├── server.go                     # Entry point: starts the GraphQL server
└── go.mod                        # Dependencies
```


## Troubleshooting
- **Invalid token**: Ensure the token is from a recent login and hasn’t expired.
- **No images returned**: Check image provider API keys and rate limits.
- **Server startup fails**: Verify all required environment variables are set.
- **Timeouts**: Increase `IMAGE_API_TIMEOUT_SEC` for slower connections.


## License

This project is part of a Pantheon Lab Programming Assignment and is intended for educational purposes only.

