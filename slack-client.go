package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type Attachment struct {
	Color         string      `json:"color,omitempty"`
	Fallback      string      `json:"fallback,omitempty"`
	CallbackID    string      `json:"callback_id,omitempty"`
	ID            int         `json:"id,omitempty"`
	AuthorID      string      `json:"author_id,omitempty"`
	AuthorName    string      `json:"author_name,omitempty"`
	AuthorSubname string      `json:"author_subname,omitempty"`
	AuthorLink    string      `json:"author_link,omitempty"`
	AuthorIcon    string      `json:"author_icon,omitempty"`
	Title         string      `json:"title,omitempty"`
	TitleLink     string      `json:"title_link,omitempty"`
	Pretext       string      `json:"pretext,omitempty"`
	Text          string      `json:"text,omitempty"`
	ImageURL      string      `json:"image_url,omitempty"`
	ThumbURL      string      `json:"thumb_url,omitempty"`
	MarkdownIn    []string    `json:"mrkdwn_in,omitempty"`
	Ts            json.Number `json:"ts,omitempty"`
}

type SlackPayload struct {
	Username    string       `json:"username,omitempty"`
	Channel     string       `json:"channel,omitempty"`
	Text        string       `json:"text,omitempty"`
	Attachments []Attachment `json:"attachments,omitempty"`
}

func SendSlack(webhook string, payload *SlackPayload) error {
	if len(webhook) == 0 {
		return fmt.Errorf("no webhook url was provided")
	}

	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", webhook, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, _ := client.Do(req)

	if resp.StatusCode != 200 {
		return fmt.Errorf("could not send slack notification, service returned %d", resp.StatusCode)
	}

	return nil
}