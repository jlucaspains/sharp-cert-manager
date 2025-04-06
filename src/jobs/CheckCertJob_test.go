package jobs

import (
	"testing"

	"github.com/jlucaspains/sharp-cert-manager/models"
	"github.com/stretchr/testify/assert"
)

type mockNotifier struct {
	executed bool
}

func (m *mockNotifier) Notify(result []CertCheckNotification) error {
	m.executed = true
	return nil
}

func (m *mockNotifier) IsReady() bool {
	return true
}

var certList = []models.CheckCertItem{
	{Name: "blog.lpains.net", Url: "https://blog.lpains.net", Type: models.CertCheckURL},
}

func TestJobInit(t *testing.T) {
	checkCertJob := &CheckCertJob{}

	checkCertJob.Init("* * * * *", "", 1, certList, &mockNotifier{})

	assert.Equal(t, "* * * * *", checkCertJob.cron)
	assert.Equal(t, "https://blog.lpains.net", checkCertJob.certList[0].Url)
	assert.Equal(t, "blog.lpains.net", checkCertJob.certList[0].Name)
	assert.Equal(t, models.CertCheckURL, checkCertJob.certList[0].Type)

	checkCertJob.ticker.Stop()
}

func TestJobInitDefaultWarningDays(t *testing.T) {
	checkCertJob := &CheckCertJob{}

	checkCertJob.Init("* * * * *", "", 0, certList, &mockNotifier{})

	assert.Equal(t, "* * * * *", checkCertJob.cron)
	assert.Equal(t, "blog.lpains.net", checkCertJob.certList[0].Name)
	assert.Equal(t, "https://blog.lpains.net", checkCertJob.certList[0].Url)
	assert.Equal(t, models.CertCheckURL, checkCertJob.certList[0].Type)
	assert.Equal(t, 30, checkCertJob.warningDays)

	checkCertJob.ticker.Stop()
}

func TestJobInitBadCron(t *testing.T) {
	checkCertJob := &CheckCertJob{}

	err := checkCertJob.Init("* * * *", "", 0, certList, &mockNotifier{})

	assert.Equal(t, "a valid cron schedule is required", err.Error())
}

func TestJobInitBadNotifier(t *testing.T) {
	checkCertJob := &CheckCertJob{}

	err := checkCertJob.Init("* * * * *", "", 0, certList, nil)

	assert.Equal(t, "a valid notifier is required", err.Error())
}

func TestJobStartStop(t *testing.T) {
	checkCertJob := &CheckCertJob{}

	err := checkCertJob.Init("* * * * *", "", 0, certList, &mockNotifier{})
	assert.Nil(t, err)
	checkCertJob.Start()
	assert.True(t, checkCertJob.running)
	checkCertJob.Stop()
	assert.False(t, checkCertJob.running)
}

func TestTryExecuteNotDue(t *testing.T) {
	checkCertJob := &CheckCertJob{}
	notifier := &mockNotifier{}
	checkCertJob.Init("0 0 1 1 1", "", 0, certList, &mockNotifier{})
	checkCertJob.notifier = notifier
	checkCertJob.tryExecute()

	assert.False(t, notifier.executed)
}

func TestTryExecuteDue(t *testing.T) {
	checkCertJob := &CheckCertJob{}
	notifier := &mockNotifier{}
	checkCertJob.Init("* * * * *", "", 0, certList, &mockNotifier{})
	checkCertJob.notifier = notifier
	checkCertJob.tryExecute()

	assert.True(t, notifier.executed)
}

func TestExecuteNow(t *testing.T) {
	checkCertJob := &CheckCertJob{}
	notifier := &mockNotifier{}
	checkCertJob.Init("* * * * *", "", 0, certList, &mockNotifier{})
	checkCertJob.notifier = notifier
	checkCertJob.RunNow()

	assert.True(t, notifier.executed)
}

func TestTryExecuteDueWarning(t *testing.T) {
	checkCertJob := &CheckCertJob{}
	notifier := &mockNotifier{}
	checkCertJob.Init("* * * * *", "", 10000, certList, &mockNotifier{})
	checkCertJob.notifier = notifier
	checkCertJob.tryExecute()

	assert.True(t, notifier.executed)
}
