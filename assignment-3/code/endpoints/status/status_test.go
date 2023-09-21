package endpoints_test

import (
	vars "cloudAss2/code/constants"
	status "cloudAss2/code/endpoints/status"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestStatus(t *testing.T) {
	methods := []string{http.MethodGet, http.MethodPut}

	testNames := []string{"Expected query", "Missing Name query", "Illegal name query", "Illegal method"}

	statusRet := vars.Status{WebhookA: len(vars.WebhookRefs), Version: vars.CURRENT_VER, Uptime: vars.MaxCacheTime}

	result1 := statusRet
	result1.CasesApi = "200 OK"
	result1.PolicyApi = "200 OK"

	results := []string{`{"cases_api":"200 OK","policy_api":"200 OK","webhooks":0,"version":"v1","uptime":172800}`, vars.ERR_ILLEGAL_METHOD}

	// Test subtests
	for i, methods := range methods {
		t.Run(testNames[i], func(t *testing.T) {
			//Update mock in status file
			status.Mock = true

			// Send request
			req := httptest.NewRequest(methods, vars.LOCAL_HOST_URL_TOT+vars.COV_STATUS, nil)

			w := httptest.NewRecorder()
			status.StatusEntry(w, req)
			res := w.Result()
			defer res.Body.Close()

			//Read respons data
			data, err := ioutil.ReadAll(res.Body)
			if err != nil {
				t.Error("Error reading body of Status test result", err)
			}
			strDat := string(data)
			strDat = strDat[:len(strDat)-1]

			if strDat != results[i] {
				t.Errorf("Expected '%s' but got '%v'", results[i], strDat)
			}
		})
	}
	status.Mock = false
}
