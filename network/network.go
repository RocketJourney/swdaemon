package network

import (
	"bytes"
	l4g "code.google.com/p/log4go"
	"net/http"
	"strconv"
)

type Network struct {
	Server string
}

func (n *Network) SendCheck(way bool, club_id int, user_id int) {
	l4g.Info("Sending check")
	client := &http.Client{}
	l4g.Info(n.Server)
	path := n.Server + "/api/v1/clubs/" + strconv.Itoa(club_id) + "/users/" + strconv.Itoa(user_id)
	request, _ := http.NewRequest("POST", path, bytes.NewBuffer([]byte("{\"check_in\":true}")))
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Authorization", "g3QAAAACZAAEZGF0YXQAAAABZAACaWRhAWQABnNpZ25lZG4GAPD6pvVOAQ")
	client.Do(request)
}
