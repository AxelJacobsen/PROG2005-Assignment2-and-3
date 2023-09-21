package endpoints_test

import (
	vars "cloudAss2/code/constants"
	public "cloudAss2/code/endpoints/public"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

/*
	Tests original public endpoint
*/
func TestPublic(t *testing.T) {
	// Send get request to each different type of public EP
	req := httptest.NewRequest(http.MethodGet, vars.LOCAL_HOST_PRE, nil)
	//req.URL.Path = test

	w := httptest.NewRecorder()
	public.PublicEntry(w, req)

	// Read result of query
	res := w.Result()
	defer res.Body.Close()

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Error("Error reading body of Public endpoint Test result", err)
	}
	strDat := string(data)
	strDat = strDat[:len(strDat)-1] // Remove last character as it's a Line Feed

	if strDat != vars.PUBLIC_SHORT {
		t.Errorf("\nExpected '%s'\n but got '%v'", vars.PUBLIC_SHORT, strDat)
	}
}

/*
	Tests the Endpoint giving "illegal" endpoint reply
*/
func TestIncorrect(t *testing.T) {
	// Send get request to each different type of public EP
	req := httptest.NewRequest(http.MethodGet, vars.LOCAL_HOST_URL_TOT, nil)
	//req.URL.Path = test

	w := httptest.NewRecorder()
	public.IncorrectURL(w, req)

	// Read result of query
	res := w.Result()
	defer res.Body.Close()

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Error("Error reading body of IncorrectURL endpoint Test result", err)
	}
	strDat := string(data)
	strDat = strDat[:len(strDat)-1] // Remove last character as it's a Line Feed

	if strDat != vars.PUBLIC_ILLEGAL {
		t.Errorf("\nExpected '%s'\n but got '%v'", vars.PUBLIC_ILLEGAL, strDat)
	}
}
