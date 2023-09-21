package endpoints

import (
	"bytes"
	vars "cloudAss2/code/constants"
	inner "cloudAss2/code/endpoints/sharedFuncs/inner"
	webhooks "cloudAss2/code/endpoints/sharedFuncs/webhookFuncs"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

var Query []byte

/*
	Entry function for cases endpoint, ensures legal URL
*/
func CasesEntry(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		parts := strings.Split(r.URL.Path, "/")
		err1 := inner.CheckForLegalURL(r.URL.Path, vars.COV_CASE)
		if len(err1) != 0 {
			http.Error(w, err1, http.StatusBadRequest)
			return
		} else {
			if 5 <= len(parts) && parts[4] != "" {
				getCases(w, parts[4])
			} else {
				http.Error(w, "Missing countryname in URL query.", http.StatusBadRequest)
			}
		}
	default:
		http.Error(w, vars.ERR_ILLEGAL_METHOD, http.StatusMethodNotAllowed)
	}
}

/*
	Handles cases GET requests
*/
func getCases(w http.ResponseWriter, countName string) {
	http.Header.Add(w.Header(), "content-type", "application/json")
	postString := map[string]string{"placeholder": `placeholder`}

	if len(countName) == 3 {
		// Implies Using country codes so we need to get them from file
		temp := inner.GetCountNameAlpha(countName)
		if temp != countName {
			countName = temp
		} else {
			http.Error(w, "Alpha3 code not recognized", http.StatusBadRequest)
			return
		}
	}
	//Check if this is a test call, if it is we dont want to check cache
	mock := false
	queryData := vars.CasesQuery{}
	if len(Query) != 0 {
		mock = true
	} else {
		//CHECK CACHE
		//IF NOT IN CACHE DO NORMAL, ELSE SKIP BIGTIME
		queryData = CheckCasesCache(countName)
	}

	if (queryData == vars.CasesQuery{}) {
		//Prepares the post string with the relevant country name
		postString = map[string]string{
			"query": `{
			country(name: "` + countName + `") {
			name
			mostRecent {
				date(format: "yyyy-MM-dd")
				confirmed
				deaths
				recovered
				growthRate
			}
			}
		}`,
		}
		if !mock {
			//prepares string for post request
			json_data, err1 := json.Marshal(postString)
			if err1 != nil {
				http.Error(w, err1.Error(), http.StatusInternalServerError)
				return
			}
			//Sends request to graphql api
			request, err2 := http.NewRequest("POST", vars.QUE_GRAPHQL, bytes.NewBuffer(json_data))
			if err2 != nil {
				http.Error(w, err2.Error(), http.StatusInternalServerError)
				return
			}
			client := &http.Client{Timeout: time.Second * 20}
			response, err3 := client.Do(request)
			if err3 != nil {
				http.Error(w, err3.Error(), http.StatusInternalServerError)
				return
			}
			defer response.Body.Close()

			//preps data from api response
			var err4 error
			Query, err4 = ioutil.ReadAll(response.Body)
			if err4 != nil {
				http.Error(w, err4.Error(), http.StatusInternalServerError)
				return
			}
		}
		//If this is a test then everyhitng before this will have been skipped

		//Adds data to struct from api
		var queryDataTemp vars.CasesQuery
		err5 := json.Unmarshal(Query, &queryDataTemp)
		if err5 != nil {
			Query = []byte{}
			http.Error(w, err5.Error(), http.StatusInternalServerError)
			return
		}
		//Empty query
		Query = []byte{}
		//Body is empty AKA couldnt find any data on the given countryname
		if len(queryDataTemp.Datas.Countries.Name) == 0 {
			http.Error(w, "Couldnt find any data on the given country: "+countName, http.StatusNotFound)
			return
		}
		//Shuffles some data and adds timestamp
		queryDataTemp.Datas.Countries.Recent.Name = queryDataTemp.Datas.Countries.Name
		queryDataTemp.Qdate = int(time.Now().Unix())
		//Prepares data for writing to webpage
		if !mock {
			vars.CasesCache = append(vars.CasesCache, queryDataTemp)
		}
		queryData = queryDataTemp
	}
	screenDat, err6 := json.Marshal(&queryData.Datas.Countries.Recent)
	if err6 != nil {
		http.Error(w, err6.Error(), http.StatusInternalServerError)
		return
	}
	if !mock {
		go webhooks.IncrementWebhookEntry(countName)
	}
	w.Write(screenDat)
}

/*
	Cache function for Cases
*/
func CheckCasesCache(countName string) vars.CasesQuery {
	retCase := vars.CasesQuery{}
	currentTime := int(time.Now().Unix())
	var newCache []vars.CasesQuery
	notFound := true
	for _, entry := range vars.CasesCache {
		if notFound && entry.Datas.Countries.Recent.Name == countName {
			retCase = entry
			notFound = false
		}
		if currentTime-entry.Qdate < vars.MaxCacheTime {
			newCache = append(newCache, entry)
		}
		if vars.MaxCacheSize <= len(newCache) {
			break
		}
	}
	vars.CasesCache = newCache
	return retCase
}
