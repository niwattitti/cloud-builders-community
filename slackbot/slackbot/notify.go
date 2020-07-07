package slackbot

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	cloudbuild "google.golang.org/api/cloudbuild/v1"
)

// Notify posts a notification to Slack that the build is complete.
func Notify(b *cloudbuild.Build, title string, icon string, tag string, webhook string) {
	burl := fmt.Sprintf("https://console.cloud.google.com/cloud-build/builds/%s?project=%s", b.Id, b.ProjectId)
	query := fmt.Sprintf("tags=\"%s\"", tag)
	params := url.Values{}
	params.Add("query", query)
	params.Add("project", b.ProjectId)
	turl := fmt.Sprintf("https://console.cloud.google.com/cloud-build/builds?%s", params.Encode())

	var i string
	var c string
	switch b.Status {
	case "SUCCESS":
		i = ":white_check_mark:"
		c = "#9CCC65"
	case "FAILURE":
		i = ":x:"
		c = "#FF5252"
	case "CANCELLED":
		i = ":wastebasket:"
		c = "#CCD1D9"
	case "TIMEOUT":
		i = ":hourglass:"
		c = "#FF5252"
	case "STATUS_UNKNOWN", "INTERNAL_ERROR":
		i = ":interrobang:"
		c = "#FF5252"
	default:
		i = ":question:"
		c = "#FF5252"
	}

	startTime, err := time.Parse(time.RFC3339, b.StartTime)
	if err != nil {
		log.Fatalf("Failed to parse Build.StartTime: %v", err)
	}
	finishTime, err := time.Parse(time.RFC3339, b.FinishTime)
	if err != nil {
		log.Fatalf("Failed to parse Build.FinishTime: %v", err)
	}
	buildDuration := finishTime.Sub(startTime).Truncate(time.Second)
	text := fmt.Sprintf("[%s] %s. %s Id: %s\nBuildDuration:%s\nStartTime: %s\nFinishTime: %s", b.Status, title, i, b.Id, buildDuration, b.StartTime, b.FinishTime)

	msgFmt := `{
		"icon_emoji": "%s",
		"username": "Cloud Build/%s",
		"attachments": [{
				"color": "%s",
				"text": "%s",
				"actions": [{
						"type": "button",
						"text": "Details",
						"url": "%s"
					},
					{
						"type": "button",
						"text": "Results using %s tag",
						"url": "%s"
				}]
		}]
	}`

	j := fmt.Sprintf(msgFmt, icon, b.ProjectId, c, text, burl, tag, turl)

	r := strings.NewReader(j)
	resp, err := http.Post(webhook, "application/json", r)
	if err != nil {
		log.Fatalf("Failed to post to Slack: %v", err)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	log.Printf("Posted message to Slack: [%v], got response [%s]", j, body)
}
