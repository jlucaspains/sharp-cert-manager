package handlers

import (
	"fmt"
	"net/http"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/jlucaspains/sharp-cert-checker/models"
)

type Handlers struct {
	SiteList []string
}

func (h Handlers) GetConfigSites() []string {
	siteList := []string{}
	for i := 1; true; i++ {
		site, ok := os.LookupEnv(fmt.Sprintf("SITE_%d", i))
		if !ok {
			break
		}

		siteList = append(siteList, site)
	}

	return siteList
}

func (h Handlers) ErrorToHttpResult(err error) (int, *models.ErrorResult) {
	if vErrs, ok := err.(validator.ValidationErrors); ok {
		out := translateErrors(vErrs)
		return http.StatusBadRequest, &models.ErrorResult{Errors: out}
	}

	return http.StatusInternalServerError, &models.ErrorResult{Errors: []string{"Unknown error"}}
}

func translateErrors(err validator.ValidationErrors) []string {
	out := make([]string, len(err))
	for i, fe := range err {
		out[i] = getValidationErrorMsg(fe)
	}
	return out
}

func getValidationErrorMsg(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", fe.Field())
	case "lte":
		return fmt.Sprintf("%s should be less than or equal to %s", fe.Field(), fe.Param())
	case "lt":
		return fmt.Sprintf("%s should be less than %s", fe.Field(), fe.Param())
	case "gte":
		return fmt.Sprintf("%s should be greater than or equal to %s", fe.Field(), fe.Param())
	case "gt":
		return fmt.Sprintf("%s should be greater than %s", fe.Field(), fe.Param())
	case "min":
		return fmt.Sprintf("%s should have minimum length of %s", fe.Field(), fe.Param())
	case "max":
		return fmt.Sprintf("%s should have maximum length of %s", fe.Field(), fe.Param())
	case "alpha":
		return fmt.Sprintf("%s should contain alpha characters only", fe.Field())
	}
	return "Unknown error"
}
