package jobs

import "github.com/jlucaspains/sharp-cert-manager/models"

type EmptyNotifier struct{}

func (m *EmptyNotifier) Notify(result []models.CertCheckResult) error {
	return nil
}
