package main

import (
	"encoding/json"
	"errors"
	"flag"
	"strconv"
	"strings"
	"syscall"
	"time"
)

func main() {

	username := flag.String("username", "", "Nerve Centre username")
	password := flag.String("password", "", "Nerve Centre password")
	webhookUrl := flag.String("webhook", "", "Slack webhook url")
	flag.Parse()

	if *username == "" || *password == "" || *webhookUrl == "" {
		flag.Usage()
		syscall.Exit(1)
	}

	err := Login(*username, *password)

	if err != nil {
		panic(err)
	}

	users := GetUsers()
	schedules := *GetSchedules()

	if len(*users) == 0 || len(schedules) == 0 {
		panic(errors.New("could not load users or schedules, check username and password"))
	}

	schedule := schedules[0]

	planningTime := time.Now()
	planningEnd := time.Now()
	planning := GetPlanning(schedule, planningTime)
	today := planning.GetActiveSlot(planningTime)
	var next *Slot

	var currentPlanningEnd time.Time
	var currentMembers []string

	if planning.HasMembers() {
		slot := planning.GetActiveSlot(planningTime)
		currentPlanningEnd = slot.End
		currentMembers = slot.GetMembers(users)
	}

	foundOther := false

	for planning.HasMembers() {

		for _, slot := range planning.BaseTimeSlots{
			planningEnd = slot.End

			if !foundOther {
				members := slot.GetMembers(users)

				if !Equal(currentMembers, members) {
					foundOther = true
					next = &slot
				} else {
					currentPlanningEnd = slot.End
				}
			}
		}

		planningTime = planningTime.Add(24 * time.Hour)
		planning = GetPlanning(schedule, planningTime)
	}

	todayMembersString := "<<geen>>"
	todayColor := "#ec0045"
	if len(currentMembers) > 0 {
		todayMembersString = strings.Join(currentMembers, ", ") + " tot " + currentPlanningEnd.Format("02-01-2006T15:04Z")
		todayColor = "#007a5a"
	}

	attachments := make([]Attachment, 0, 3)

	attachments = append(attachments, Attachment{
		Fallback: "Vandaag: " + todayMembersString,
		Color:    todayColor,
		Title:    "Vandaag",
		Text:     todayMembersString,
	})

	if foundOther && next != nil {
		nextMembers := next.GetMembers(users)
		nextMembersString := "<<geen>>"
		nextColor := "#ec0045"

		if len(currentMembers) > 0 {
			nextMembersString = strings.Join(nextMembers, ", ") + " om " + next.Start.Format("02-01-2006T15:04Z")
			nextColor = "#ffc917"
		}

		attachments = append(attachments, Attachment{
			Fallback: "Volgende: " + nextMembersString,
			Color:    nextColor,
			Title:    "Volgende",
			Text:     nextMembersString,
			Ts:       json.Number(strconv.FormatInt(next.Start.Unix(), 10)),
		})
	}

	if len(today.GetMembers(users)) > 0 {
		attachments = append(attachments, Attachment{
			Fallback: "Er is een rooster tot " + planningEnd.Format("02-01-2006T15:04Z"),
			Color:    "#ec0045",
			Title:    "Einde rooster",
			Text:     "Er is een rooster tot " + planningEnd.Format("02-01-2006T15:04Z"),
			Ts:       json.Number(strconv.FormatInt(planningEnd.Unix(), 10)),
		})
	}

	message := SlackPayload{
		Username:    "📞 Wachtdienst " + schedule.GroupName,
		Text:        "Een overzicht van de de huidige wachtdiensten die zijn ingeregeld voor " + schedule.GroupName + " in Nerve Centre",
		Attachments: attachments,
	}

	err = SendSlack(*webhookUrl, &message)

	if err != nil {
		panic(err)
	}
}
