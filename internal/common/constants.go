package common

const (
	USER_KEY         = "user"    // Key to get the user from the context locals.
	USER_CONTEXT_KEY = "userJWT" // Context key to store user information from the token into context.
	USER_COOKIE_KEY  = "user"    // Key to get the user cookie from the cookie string.
	USER_CLAIM_KEY   = "userId"  // Key to get the user id from the JWT claims.
)
