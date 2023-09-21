package endpoints

import (
	"bytes"
	inner "cloudAss2/endpoints/inner"
	vars "cloudAss2/variables"
	"encoding/json"
	"net/http"
	"strings"
	"time"
)

func StatusEntry(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		parts := strings.Split(r.URL.Path, "/")
		if len(parts) < 4 || !inner.CompUrlString(r.URL.Path, vars.COV_URL+vars.COV_STATUS) {
			http.Error(w, "Malformed URL", http.StatusBadRequest)
			return
		} else {
			http.Header.Add(w.Header(), "content-type", "application/json")
			writeData, err := statusInner()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			} else {
				json.NewEncoder(w).Encode(writeData)
			}
		}
	default:
		http.Error(w, vars.ERR_ILLEGAL_METHOD, http.StatusBadRequest)
	}
}

func statusInner() (vars.Status, error) {
	statusRet := vars.Status{WebhookA: len(vars.WebhookRefs), Version: vars.CURRENT_VER, Uptime: int(time.Now().Unix()) - vars.StartTime}
	//Checks for status of the given website

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

	request, err2 := http.NewRequest("POST", vars.QUE_GRAPHQL, bytes.NewBuffer(json_data))
	if err2 != nil {
		return vars.Status{}, err2
	}

	client := &http.Client{Timeout: time.Second * 20}
	response, err3 := client.Do(request)
	if err3 != nil {
		return vars.Status{}, err3
	}

	statusRet.CasesApi = response.StatusCode
	defer response.Body.Close()
	//Checks for status of the given website
	resp, err2 := http.Get(vars.QUE_OXFORD_STATUS)
	if err2 != nil {
		return vars.Status{}, err2
	}
	statusRet.PolicyApi = resp.StatusCode

	return statusRet, nil
}
