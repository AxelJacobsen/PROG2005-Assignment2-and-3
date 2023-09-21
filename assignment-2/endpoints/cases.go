package endpoints

import (
	"bytes"
	inner "cloudAss2/endpoints/inner"
	vars "cloudAss2/variables"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

func CasesEntry(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		parts := strings.Split(r.URL.Path, "/")
		if len(parts) < 4 || !inner.CompUrlString(r.URL.Path, vars.COV_URL+vars.COV_CASE) {
			http.Error(w, "Malformed URL", http.StatusBadRequest)
			return
		} else {
			getCases(w, parts[4])
		}
	default:
		http.Error(w, vars.ERR_ILLEGAL_METHOD, http.StatusMethodNotAllowed)
	}
}

func getCases(w http.ResponseWriter, countName string) {
	http.Header.Add(w.Header(), "content-type", "application/json")
	postString := map[string]string{"placeholder": `placeholder`}
	if len(countName) == 3 {
		// Implies Using country codes so we need to get them from file
		temp := inner.GetCountNameAlpha(countName, true)
		if len(temp) != 0 {
			countName = temp
		} else {
			http.Error(w, "Alpha3 code not recognized", http.StatusBadRequest)
		}
	}

	//CHECK CACHE
	//IF NOT IN CACHE DO NORMAL, ELSE SKIP BIGTIME
	queryData := checkCasesCache(countName)
	if (queryData == vars.CasesQuery{}) {

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

		json_data, err1 := json.Marshal(postString)
		if err1 != nil {
			http.Error(w, err1.Error(), http.StatusInternalServerError)
		}

		request, err2 := http.NewRequest("POST", vars.QUE_GRAPHQL, bytes.NewBuffer(json_data))
		if err2 != nil {
			fmt.Println("ERROR IN: getCases", postString)
			http.Error(w, err2.Error(), http.StatusInternalServerError)
		}

		client := &http.Client{Timeout: time.Second * 20}
		response, err3 := client.Do(request)
		if err3 != nil {
			http.Error(w, err3.Error(), http.StatusInternalServerError)
		}
		defer response.Body.Close()

		parsedData, err4 := ioutil.ReadAll(response.Body)
		if err4 != nil {
			http.Error(w, err4.Error(), http.StatusInternalServerError)
		}

		var queryDataTemp vars.CasesQuery
		err5 := json.Unmarshal(parsedData, &queryDataTemp)
		if err5 != nil {
			http.Error(w, err5.Error(), http.StatusInternalServerError)
		}
		queryDataTemp.Datas.Countries.Recent.Name = queryDataTemp.Datas.Countries.Name
		queryDataTemp.Qdate = int(time.Now().Unix())
		//Prepares data for writing to webpage
		vars.CasesCache = append(vars.CasesCache, queryDataTemp)
		queryData = queryDataTemp
	}
	screenDat, err6 := json.Marshal(&queryData.Datas.Countries.Recent)
	if err6 != nil {
		http.Error(w, err6.Error(), http.StatusInternalServerError)
	}
	inner.IncrementWebhookEntry(countName)
	w.Write(screenDat)
}

func checkCasesCache(countName string) vars.CasesQuery {
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
