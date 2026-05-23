package assemble

import "encoding/base64"

func BuildBasicAuth(user string, password string) string {
	raw := user + ":" + password
	encoded := base64.StdEncoding.EncodeToString([]byte(raw))
	return "Basic " + encoded

}
