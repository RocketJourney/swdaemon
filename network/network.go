package network

import (
	"bytes"
	l4g "code.google.com/p/log4go"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

type Network struct {
	Server      string
	AccessToken string
}

func (n *Network) SendCheck(way int, club_id int, user_id int) {
	l4g.Info("Sending check")
	client := &http.Client{}
	l4g.Info(n.Server)
	check := "false"
	if way == 1 {
		check = "true"
	}
	path := n.Server + "/api/v1/clubs/" + strconv.Itoa(club_id) + "/users/" + strconv.Itoa(user_id)
	l4g.Info(path)
	request, _ := http.NewRequest("POST", path, bytes.NewBuffer([]byte("{\"check_in\":"+check+"}")))
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Authorization", n.AccessToken)
	client.Do(request)
}

func (n *Network) ReportAlive(pid string, club_id string) {
	l4g.Info("Reporting Alive")
	client := &http.Client{}

	path := n.Server + "/api/v1/swdaemon/" + club_id + "/" + pid
	l4g.Info(path)
	request, _ := http.NewRequest("POST", path, nil)
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Authorization", "g3QAAAACZAAEZGF0YXQAAAABZAACaWRhAWQABnNpZ25lZG4GAPD6pvVOAQ")
	client.Do(request)
}

func (n *Network) GetUpdateFile(config_file string) (*[]byte, error) {
	response, err := http.Get(config_file)

	if err != nil {
		l4g.Trace(fmt.Sprintf("Could not check for new swdameon version: %s", config_file))
		return nil, err
	}
	defer response.Body.Close()

	var conf_data []byte
	conf_data, err = ioutil.ReadAll(response.Body)
	if err != nil {
		l4g.Trace(fmt.Sprintf("Unable to get update file: %s error: %s ", err.Error()))
		return nil, err
	}

	responseCode := response.StatusCode
	if responseCode != http.StatusOK {
		errorMessage := fmt.Sprintf("Could not check for new swdameon version: %s, getting file: %s", response.Status, config_file)
		l4g.Trace(errorMessage)
		return nil, errors.New(errorMessage)
	}

	l4g.Trace("Configuration file downloaded")

	return &conf_data, nil

}
