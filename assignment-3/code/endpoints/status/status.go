package endpoints

import (
	"bytes"
	vars "cloudAss2/code/constants"
	inner "cloudAss2/code/endpoints/sharedFuncs/inner"
	"encoding/json"
	"net/http"
	"time"
)

var Mock = false

/*
	Makes sure the URL is correct and handels errors
*/
func StatusEntry(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		err1 := inner.CheckForLegalURL(r.URL.Path, vars.COV_STATUS)
		if len(err1) != 0 {
			http.Error(w, err1, http.StatusBadRequest)
			return
		}
		http.Header.Add(w.Header(), "content-type", "application/json")
		//Calls the actual status function
		writeData, err := statusInner()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		} else {
			json.NewEncoder(w).Encode(writeData)
		}
	default:
		http.Error(w, vars.ERR_ILLEGAL_METHOD, http.StatusBadRequest)
	}
}

/*
	Main Status function, handles all funcitonality outside of writing
*/
func statusInner() (vars.Status, error) {
	statusRet := vars.Status{WebhookA: len(vars.WebhookRefs), Version: vars.CURRENT_VER, Uptime: int(time.Now().Unix()) - vars.StartTime}
	if Mock {
		statusRet.Uptime = vars.MaxCacheTime
	}

	//Checks for status of the given website-

	dummyRequest := map[string]string{
		"query": `{
			country(name: " ") {
			name
			}
		}`,
	}

	json_data, err1 := json.Marshal(dummyRequest)
	if err1 != nil {
		return vars.Status{}, err1
	}
	if !Mock {
		//Sends dummy request to grab status code
		request, err2 := http.NewRequest("POST", vars.QUE_GRAPHQL, bytes.NewBuffer(json_data))
		if err2 != nil {
			return vars.Status{}, err2
		}

		client := &http.Client{Timeout: time.Second * 20}
		response, err3 := client.Do(request)
		if err3 != nil {
			return vars.Status{}, err3
		}

		//Adds status to Struct
		statusRet.CasesApi = response.Status
		defer response.Body.Close()

		//Checks for status of the given website
		resp, err2 := http.Get(vars.QUE_OXFORD_STATUS)
		if err2 != nil {
			return vars.Status{}, err2
		}
		statusRet.PolicyApi = resp.Status
	} else {
		statusRet.CasesApi = "200 OK"
		statusRet.PolicyApi = "200 OK"
	}

	return statusRet, nil
}
