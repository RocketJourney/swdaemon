package network

import (
	"bytes"
	l4g "code.google.com/p/log4go"
	"net/http"
)

func SendCheck(bool) {
	l4g.Info("Sending check")
	client := &http.Client{}
	request, _ := http.NewRequest("POST", "http://localhost:4000/api/v1/clubs/1/users/855893", bytes.NewBuffer([]byte("{\"check_in\":true}")))
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Authorization", "g3QAAAACZAAEZGF0YXQAAAABZAACaWRhAWQABnNpZ25lZG4GAPD6pvVOAQ")
	client.Do(request)
}
