package jobs

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"text/template"
	"time"
)

type TeamsNotifier struct {
	WebhookUrl        string
	NotificationTitle string
	NotificationBody  string
	NotificationUrl   string
	parsedTemplate    *template.Template
	httpClient        *http.Client
}

type TeamsNotificationCard struct {
	Title           string
	Description     string
	NotificationUrl string
	Items           []CertCheckNotification
}

const messageTemplate = `{
	"type": "message",
	"attachments": [{
		"contentType": "application/vnd.microsoft.card.adaptive",
		"content": {
			"type": "AdaptiveCard",
			"version": "1.5",
			"$schema": "http://adaptivecards.io/schemas/adaptive-card.json",
			"body": [
				{
					"type": "TextBlock",
					"text": "{{ .Title }}",
					"size": "large",
					"weight": "bolder",
					"wrap": true
				},
				{
					"type": "TextBlock",
					"text": "{{ .Description }}",
					"isSubtle": true,
					"wrap": true
				},
				{
					"type": "Table",
					"columns": [
						{
							"width": 2
						},
						{
							"width": 4
						}
					],
					"rows": [
						{{- if gt (len .Items) 0}}
						{{- $max := len (slice .Items 1)}}
						{{- range $i, $item := .Items}}
						{
							"type": "TableRow",
							"cells": [
								{
									"type": "TableCell",
									"items": [
										{
										"type": "TextBlock",
										"text": "{{if not $item.IsValid}}❌{{else if $item.ExpirationWarning}}⚠️{{else}}✔️{{end}}{{$item.Hostname}}"
										}
									]
								},
								{
									"type": "TableCell",
									"items": [
										{
										"type": "TextBlock",
										"text": "{{ range $index, $element := $item.Messages}}{{if $index}}, {{end}}{{$element}}{{end}}",
										"wrap": true
										}
									]
								}
							]
						}{{if lt $i $max}},{{end}}
						{{- end}}
						{{- end}}
					]
				}
			]{{if .NotificationUrl}},
			"actions": [
				{
					"type": "Action.OpenUrl",
					"title": "View Details",
					"url": "{{ .NotificationUrl }}"
				}
			]{{end}}
		}
	}]
}`

func (m *TeamsNotifier) Init(webhookUrl string, notificationTitle string, notificationBody string, notificationUrl string) {
	if notificationTitle == "" {
		notificationTitle = "Sharp Cert Manager Summary"
	}

	if notificationBody == "" {
		notificationBody = fmt.Sprintf("The following certificates were checked on %s", time.Now().Format("01/02/2006"))
	}

	m.NotificationTitle = notificationTitle
	m.NotificationBody = notificationBody
	m.NotificationUrl = notificationUrl
	m.WebhookUrl = webhookUrl
}

func (m *TeamsNotifier) Notify(result []CertCheckNotification) error {
	client := m.getClient()
	parsedTemplate := m.getTemplate()
	card := TeamsNotificationCard{
		Title:           m.NotificationTitle,
		Description:     m.NotificationBody,
		NotificationUrl: m.NotificationUrl,
		Items:           result,
	}

	var templateBody bytes.Buffer
	err := parsedTemplate.Execute(&templateBody, card)

	if err != nil {
		return err
	}

	response, err := client.Post(m.WebhookUrl, "application/json", bytes.NewReader(templateBody.Bytes()))

	if err != nil {
		return err
	}

	defer response.Body.Close()

	if response.StatusCode != 200 {
		return errors.New("error sending notification to Teams")
	}

	return nil
}

func (m *TeamsNotifier) getTemplate() *template.Template {
	if m.parsedTemplate == nil {
		m.parsedTemplate, _ = template.New("template").Parse(messageTemplate)
	}

	return m.parsedTemplate
}

func (m *TeamsNotifier) getClient() *http.Client {
	if m.httpClient == nil {
		m.httpClient = &http.Client{}
	}

	return m.httpClient
}
