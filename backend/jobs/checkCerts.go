package jobs

import (
	"fmt"
	"log"
	"time"

	"github.com/adhocore/gronx"
	"github.com/jlucaspains/sharp-cert-manager/models"
	"github.com/jlucaspains/sharp-cert-manager/shared"
)

type Notifier interface {
	Notify(result []models.CertCheckResult) error
}

type CheckCertJob struct {
	cron     string
	ticker   *time.Ticker
	gron     gronx.Gronx
	siteList []string
	running  bool
	notifier Notifier
}

func (c *CheckCertJob) Init(schedule string, siteList []string, notifier Notifier) error {
	c.gron = gronx.New()

	if schedule == "" || !c.gron.IsValid(schedule) {
		log.Printf("A valid cron schedule is required in the format e.g.: * * * * *")
		return fmt.Errorf("a valid cron schedule is required")
	}

	if notifier == nil {
		log.Printf("A valid notifier is required")
		return fmt.Errorf("a valid notifier is required")
	}

	c.cron = schedule
	c.siteList = siteList
	c.ticker = time.NewTicker(time.Minute)
	c.notifier = notifier

	return nil
}

func (c *CheckCertJob) Start() {
	c.running = true
	go func() {
		for range c.ticker.C {
			c.tryExecute()
		}
	}()
}

func (c *CheckCertJob) Stop() {
	c.running = false
	c.ticker.Stop()
}

func (c *CheckCertJob) tryExecute() {
	due, _ := c.gron.IsDue(c.cron, time.Now().Truncate(time.Minute))

	log.Printf("tryExecute job, isDue: %t", due)

	if due {
		c.execute()
	}
}

func (c *CheckCertJob) execute() {
	result := []models.CertCheckResult{}
	for _, url := range c.siteList {
		params := models.CertCheckParams{Url: url}
		checkStatus, err := shared.CheckCertStatus(params)

		if err != nil {
			log.Printf("Error checking cert status: %s", err)
			continue
		}

		log.Printf("Cert status for %s: %t", url, checkStatus.IsValid)

		result = append(result, checkStatus)
	}

	c.notifier.Notify(result)
}
