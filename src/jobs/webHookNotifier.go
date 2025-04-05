package jobs

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"text/template"
	"time"
)

type WebHookNotifier struct {
	NotifierType      NotifierType
	WebhookUrl        string
	NotificationTitle string
	NotificationBody  string
	NotificationUrl   string
	Mentions          []string
	parsedTemplate    *template.Template
	httpClient        *http.Client
}

type WebHookNotificationCard struct {
	Title           string
	Description     string
	NotificationUrl string
	Items           []CertCheckNotification
	Mentions        []string
}

func (m *WebHookNotifier) Init(notifierType NotifierType, webhookUrl string, notificationTitle string, notificationBody string, notificationUrl string, messageMentions string) {
	if notificationTitle == "" {
		notificationTitle = "Sharp Cert Manager Summary"
	}

	if notificationBody == "" {
		notificationBody = fmt.Sprintf("The following certificates were checked on %s", time.Now().Format("01/02/2006"))
	}

	m.NotifierType = notifierType
	m.NotificationTitle = notificationTitle
	m.NotificationBody = notificationBody
	m.NotificationUrl = notificationUrl
	m.WebhookUrl = webhookUrl
	m.Mentions = parseMentions(messageMentions)
}

func parseMentions(mentions string) []string {
	if mentions == "" {
		return []string{}
	}

	return strings.Split(mentions, ",")
}

func (m *WebHookNotifier) Notify(result []CertCheckNotification) error {
	client := m.getClient()
	parsedTemplate := m.getTemplate()
	card := WebHookNotificationCard{
		Title:           m.NotificationTitle,
		Description:     m.NotificationBody,
		NotificationUrl: m.NotificationUrl,
		Items:           result,
		Mentions:        m.Mentions,
	}

	var templateBody bytes.Buffer
	err := parsedTemplate.Execute(&templateBody, card)

	if err != nil {
		return err
	}

	stringBody := templateBody.String()
	fmt.Println(stringBody)

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

func (m *WebHookNotifier) getTemplate() *template.Template {
	if m.parsedTemplate == nil {
		m.parsedTemplate, _ = template.New("template").Funcs(template.FuncMap{
			"split": func(s, sep string) []string {
				return strings.Split(s, sep)
			},
		}).Parse(NotificationTemplates[m.NotifierType])
	}

	return m.parsedTemplate
}

func (m *WebHookNotifier) getClient() *http.Client {
	if m.httpClient == nil {
		m.httpClient = &http.Client{}
	}

	return m.httpClient
}
