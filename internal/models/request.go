package models

type contextKey string

const RequestContextKey contextKey = "requestContext"

type RequestContext struct {
	UserInfo interface{}
	Token    string
}
