## AWS Cognito Authentication in Golang

Read the blog for more implementation details: [https://aashishkoshti.in/blog/aws-cognito-golang](https://aashishkoshti.in/blog/aws-cognito-golang)

### Run this project

1. Setup [AWS Cognito User Pool](https://aashishkoshti.in/blog/aws-cognito-golang#aws-cognito-user-pool-setup)
2. Create the `.env` file and fill all the variables, any problem refer the blog.

```sh
ENV=development
PORT=:8080

AWS_COGNITO_USER_POOL_ID=<user_pool_id>
AWS_COGNITO_CLIENT_ID=<application_client_id>
AWS_COGNITO_CLIENT_SECRET=<application_client_secret>
AWS_COGNITO_TOKEN_URL=<user_pool_token_signing_url>
AWS_COGNITO_JWT_ISSUER_URL=<user_pool_issuer_url>
```

3. Install go dependencies using `go mod tidy`.
4. Run the project using `make run` or `go run cmd/main.go`.
