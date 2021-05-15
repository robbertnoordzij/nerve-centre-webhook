package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"time"
)

const nerveCentreBaseUrl = "https://portal.ncaas.nl/2020-2"

type User struct {
	Id   string
	Name string
}

type Schedule struct {
	GroupId     string
	ParameterId string
	GroupName   string
}

type Slot struct {
	Start   time.Time
	End     time.Time
	Members []string
}

type Planning struct {
	BaseTimeSlots    []Slot
	PrimaryTimeSlots []Slot
}

var client *http.Client

func InitializeClient() {
	jar, _ := cookiejar.New(nil)
	client = &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
		Timeout: 10 * time.Second,
		Jar:     jar,
	}
}

func Login(username string, password string) error {

	if len(username) == 0 || len(password) == 0 {
		return fmt.Errorf("username or password is not provided")
	}

	req, _ := http.NewRequest("GET", nerveCentreBaseUrl+"/login.cshtml", nil)
	client.Do(req)

	form := url.Values{}
	form.Add("username", username)
	form.Add("password", password)

	req, _ = http.NewRequest("POST", nerveCentreBaseUrl+"/login.cshtml?ReturnUrl=~%2f", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, _ := client.Do(req)

	if resp.StatusCode != 200 {
		return fmt.Errorf("failed to login, Nerve Centre returned %d", resp.StatusCode)
	}

	return nil
}

func GetUsers() *[]User {
	req, _ := http.NewRequest("GET", nerveCentreBaseUrl+"/um/controller/1.0/users", nil)
	req.Header.Set("Accept", "application/json, text/plain, */*")

	resp, _ := client.Do(req)

	body, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	var users []User

	json.Unmarshal(body, &users)

	return &users
}

func GetSchedules() *[]Schedule {
	req, _ := http.NewRequest("GET", nerveCentreBaseUrl+"/reachability/controller/1.0/groups/config/schedules", nil)
	req.Header.Set("Accept", "application/json, text/plain, */*")

	resp, _ := client.Do(req)

	body, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	var schedules []Schedule

	json.Unmarshal(body, &schedules)

	return &schedules
}

func GetPlanning(schedule Schedule, date time.Time) *Planning {
	dateString := date.Format("2006-01-02")

	req, _ := http.NewRequest("GET", nerveCentreBaseUrl+"/reachability/controller/1.0/groups/"+schedule.GroupId+"/config/"+schedule.ParameterId+"/schedule/"+dateString, nil)
	req.Header.Set("Accept", "application/json, text/plain, */*")

	resp, _ := client.Do(req)

	body, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	var planning Planning

	json.Unmarshal(body, &planning)

	return &planning
}


func (planning *Planning) HasMembers() bool {

	for _, slot := range planning.BaseTimeSlots {
		if len(slot.Members) > 0 {
			return true
		}
	}

	return false
}

func (planning *Planning) GetEnd() time.Time {
	var end time.Time

	for _, slot := range planning.BaseTimeSlots {
		if len(slot.Members) > 0 {
			end = slot.End
		}
	}

	return end
}

func (planning *Planning) GetStart() time.Time {
	return planning.BaseTimeSlots[0].Start
}

func (planning *Planning) GetMembers(users *[]User) []string {
	index := make(map[string]string)

	for _, user := range *users {
		index[user.Id] = user.Name
	}

	members := make(map[string]struct{})

	for _, slot := range planning.BaseTimeSlots {

		for _, member := range slot.Members {
			members[index[member]] = struct{}{}
		}
	}

	keys := make([]string, 0, len(members))

	for key, _ := range members {
		keys = append(keys, key)
	}

	return keys
}