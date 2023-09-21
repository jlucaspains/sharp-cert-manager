package jobs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockNotifier struct {
	executed bool
}

func (m *mockNotifier) Notify(result []CertCheckNotification) error {
	m.executed = true
	return nil
}

func TestJobInit(t *testing.T) {
	checkCertJob := &CheckCertJob{}

	checkCertJob.Init("* * * * *", "", 1, []string{"https://blog.lpains.net"}, &mockNotifier{})

	assert.Equal(t, "* * * * *", checkCertJob.cron)
	assert.Equal(t, "https://blog.lpains.net", checkCertJob.siteList[0])

	checkCertJob.ticker.Stop()
}

func TestJobInitDefaultWarningDays(t *testing.T) {
	checkCertJob := &CheckCertJob{}

	checkCertJob.Init("* * * * *", "", 0, []string{"https://blog.lpains.net"}, &mockNotifier{})

	assert.Equal(t, "* * * * *", checkCertJob.cron)
	assert.Equal(t, "https://blog.lpains.net", checkCertJob.siteList[0])
	assert.Equal(t, 30, checkCertJob.warningDays)

	checkCertJob.ticker.Stop()
}

func TestJobInitBadCron(t *testing.T) {
	checkCertJob := &CheckCertJob{}

	err := checkCertJob.Init("* * * *", "", 0, []string{"https://blog.lpains.net"}, &mockNotifier{})

	assert.Equal(t, "a valid cron schedule is required", err.Error())
}

func TestJobInitBadNotifier(t *testing.T) {
	checkCertJob := &CheckCertJob{}

	err := checkCertJob.Init("* * * * *", "", 0, []string{"https://blog.lpains.net"}, nil)

	assert.Equal(t, "a valid notifier is required", err.Error())
}

func TestJobStartStop(t *testing.T) {
	checkCertJob := &CheckCertJob{}

	err := checkCertJob.Init("* * * * *", "", 0, []string{"https://blog.lpains.net"}, &mockNotifier{})
	assert.Nil(t, err)
	checkCertJob.Start()
	assert.True(t, checkCertJob.running)
	checkCertJob.Stop()
	assert.False(t, checkCertJob.running)
}

func TestTryExecuteNotDue(t *testing.T) {
	checkCertJob := &CheckCertJob{}
	notifier := &mockNotifier{}
	checkCertJob.Init("0 0 1 1 1", "", 0, []string{"https://blog.lpains.net"}, &mockNotifier{})
	checkCertJob.notifier = notifier
	checkCertJob.tryExecute()

	assert.False(t, notifier.executed)
}

func TestTryExecuteDue(t *testing.T) {
	checkCertJob := &CheckCertJob{}
	notifier := &mockNotifier{}
	checkCertJob.Init("* * * * *", "", 0, []string{"https://blog.lpains.net"}, &mockNotifier{})
	checkCertJob.notifier = notifier
	checkCertJob.tryExecute()

	assert.True(t, notifier.executed)
}

func TestTryExecuteDueWarning(t *testing.T) {
	checkCertJob := &CheckCertJob{}
	notifier := &mockNotifier{}
	checkCertJob.Init("* * * * *", "", 10000, []string{"https://blog.lpains.net"}, &mockNotifier{})
	checkCertJob.notifier = notifier
	checkCertJob.tryExecute()

	assert.True(t, notifier.executed)
}
