package auth

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

type Token struct {
	Sub       string
	Event_Id  string
	Token_use string
	Scope     string
	Auth_time int
	Iss       string
	Exp       int
	Iat       int
	Client_id string
	Username  string
}

// This function vaidate the token. If the token is not valid
// The token expected must be access_token because this contains the Username attribute.
func Validate(token string) (
	bool,
	error,
	string,
) {

	// First part contains the payload
	// Second part contains the headers
	// Third part contains the sign
	tokenItems := strings.Split(token, ".")

	if len(tokenItems) != 3 {
		message := "The token is not valid."
		fmt.Println(message)
		return false, nil, message
	}

	//userInfo, err := base64.StdEncoding.DecodeString(tokenItems[1])
	userInfo, err := base64.RawStdEncoding.DecodeString(tokenItems[1])

	if err != nil {
		errorMessage := err.Error()
		fmt.Println("It was an error to decode (base64) the token:", errorMessage)
		return false, err, errorMessage
	}

	var t Token
	err = json.Unmarshal(userInfo, &t)

	if err != nil {
		errorMessage := err.Error()
		fmt.Println("It was possible to create an Token struct. It was an error:", errorMessage)

		return false, err, errorMessage
	}

	// Check if the token is not expired
	now := time.Now()
	tm := time.Unix(int64(t.Exp), 0)

	// Token expired if tm is before now
	if tm.Before(now) {
		errorMessage := "Token has expired."
		fmt.Printf("%s Date token is: %s and today: %s", errorMessage, tm.String(), now.String())

		return false, err, errorMessage
	}

	// Get Username from access_token.
	return true, nil, string(t.Username)

}
