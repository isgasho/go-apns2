package apns2

var statusCode = map[int]string{
	200: "Success",
	400: "Bad request",
	403: "There was an error with the certificate.",
	405: "The request used a bad :method value. Only POST requests are supported.",
	410: "The device token is no longer active for the topic.",
	413: "The notification payload was too large.",
	429: "The server received too many requests for the same device token.",
	500: "Internal server error",
	503: "The server is shutting down and unavailable.",
}
