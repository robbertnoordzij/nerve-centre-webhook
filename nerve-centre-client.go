package main

import (
	"4d63.com/tz"
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

var nerveCentreBaseUrl = "https://portal.ncaas.nl/2020-2"

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

var nerveCentreHttpClient *http.Client

func init() {
	jar, _ := cookiejar.New(nil)
	nerveCentreHttpClient = &http.Client{
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
	nerveCentreHttpClient.Do(req)

	form := url.Values{}
	form.Add("username", username)
	form.Add("password", password)

	req, _ = http.NewRequest("POST", nerveCentreBaseUrl+"/login.cshtml?ReturnUrl=~%2f", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, _ := nerveCentreHttpClient.Do(req)

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to login, Nerve Centre returned %d", resp.StatusCode)
	}

	return nil
}

func GetUsers() *[]User {
	req, _ := http.NewRequest("GET", nerveCentreBaseUrl+"/um/controller/1.0/users", nil)
	req.Header.Set("Accept", "application/json, text/plain, */*")

	resp, _ := nerveCentreHttpClient.Do(req)

	body, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	var users []User

	json.Unmarshal(body, &users)

	return &users
}

func GetSchedules() *[]Schedule {
	req, _ := http.NewRequest("GET", nerveCentreBaseUrl+"/reachability/controller/1.0/groups/config/schedules", nil)
	req.Header.Set("Accept", "application/json, text/plain, */*")

	resp, _ := nerveCentreHttpClient.Do(req)

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

	resp, _ := nerveCentreHttpClient.Do(req)

	body, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	var planning Planning

	json.Unmarshal(body, &planning)

	fixTimeZoneForPlanning(&planning)

	return &planning
}

func fixTimeZoneForPlanning(planning *Planning) {
	for i, _ := range planning.BaseTimeSlots {
		slot := &planning.BaseTimeSlots[i]
		slot.Start = fixTimeZone(slot.Start)
		slot.End = fixTimeZone(slot.End)
	}
	for i, _ := range planning.PrimaryTimeSlots {
		slot := &planning.PrimaryTimeSlots[i]
		slot.Start = fixTimeZone(slot.Start)
		slot.End = fixTimeZone(slot.End)
	}
}

func fixTimeZone(toFix time.Time) time.Time {
	loc, _ := tz.LoadLocation("Europe/Amsterdam")

	// Nerve Centre has the nerve to communicate local times as if they were Zulu, so we have to update te location
	return time.Date(
		toFix.Year(),
		toFix.Month(),
		toFix.Day(),
		toFix.Hour(),
		toFix.Minute(),
		toFix.Second(),
		0,
		loc,
	)
}


func (planning *Planning) HasMembers() bool {
	for _, slot := range planning.BaseTimeSlots {
		if len(slot.Members) > 0 {
			return true
		}
	}

	return false
}

func (planning *Planning) GetActiveSlot(time time.Time) *Slot {
	for _, slot := range planning.BaseTimeSlots {
		if (slot.Start.Before(time) || slot.Start.Equal(time)) && slot.End.After(time) {
			return &slot
		}
	}

	return nil
}

func (slot *Slot) GetMembers(users *[]User) []string {
	index := make(map[string]string)

	for _, user := range *users {
		index[user.Id] = user.Name
	}

	members := make(map[string]struct{})

	for _, member := range slot.Members {
		members[index[member]] = struct{}{}
	}

	keys := make([]string, 0, len(members))

	for key, _ := range members {
		keys = append(keys, key)
	}

	return keys
}