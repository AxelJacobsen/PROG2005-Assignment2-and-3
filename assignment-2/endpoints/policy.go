package endpoints

import (
	inner "cloudAss2/endpoints/inner"
	vars "cloudAss2/variables"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

func PolicyEntry(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		parts := strings.Split(r.URL.Path, "/")
		if len(parts) < 4 || !inner.CompUrlString(r.URL.Path, vars.COV_URL+vars.COV_POLICY) {
			http.Error(w, "Malformed URL", http.StatusBadRequest)
			return
		} else {
			getPolicies(w, r, parts[4])
		}
	default:
		http.Error(w, vars.ERR_ILLEGAL_METHOD, http.StatusMethodNotAllowed)
	}
}

func getPolicies(w http.ResponseWriter, r *http.Request, countName string) {
	http.Header.Add(w.Header(), "content-type", "application/json")
	if len(countName) != 3 {
		// Implies Using country codes so we need to get them from file
		temp := inner.GetCountNameAlpha(countName, false)
		if len(temp) != 0 {
			countName = temp
		} else {
			http.Error(w, "Country name not recognized", http.StatusBadRequest)
		}
	}

	potentialPolicyQuery := checkPolicyCache(countName)

	if (potentialPolicyQuery.Stringencies == vars.PolicyQuery{}.Stringencies) {
		querystring := vars.QUE_OXFORD + countName + "/" + time.Now().Format("2006-01-02")

		dateScope, exists := r.URL.Query()["scope"]
		if exists && len(dateScope[0]) != 0 {
			_, err := time.Parse("2006-02-01", dateScope[0])
			if err != nil {
				log.Print("Error, incorrect date format in 'scope' variable.")
			} else {
				querystring = vars.QUE_OXFORD + countName + "/" + dateScope[0]
			}
		}

		data, err1 := http.Get(querystring)
		if err1 != nil {
			fmt.Println("ERROR IN: getPolicies", data)
			http.Error(w, err1.Error(), http.StatusInternalServerError)
		}
		defer data.Body.Close()
		//Prepears data for writing to screen
		parsedData, err2 := ioutil.ReadAll(data.Body)
		if err2 != nil {
			http.Error(w, err2.Error(), http.StatusInternalServerError)
		}

		//Define empty University list
		var queryData vars.PolicyQuery
		err5 := json.Unmarshal(parsedData, &queryData)
		if err5 != nil {
			http.Error(w, err5.Error(), http.StatusInternalServerError)
		}
		potentialPolicyQuery = queryData
		potentialPolicyQuery.DateQueried = int(time.Now().Unix())
		vars.PolicyCache = append(vars.PolicyCache, potentialPolicyQuery)
	}

	var printData vars.PolicyResponse
	printData = setPrintData(printData, potentialPolicyQuery)

	//Prepares data for writing to webpage
	screenDat, err6 := json.Marshal(&printData)
	if err6 != nil {
		http.Error(w, err6.Error(), http.StatusInternalServerError)
	}

	inner.IncrementWebhookEntry(inner.GetCountNameAlpha(countName, true))
	w.Write(screenDat)
}

func setPrintData(retDat vars.PolicyResponse, contDat vars.PolicyQuery) vars.PolicyResponse {
	retDat.Alpha = contDat.Stringencies.Alpha
	retDat.Date = contDat.Stringencies.Date

	if contDat.Policies[0].TypeCode != "NONE" {
		retDat.Policies = len(contDat.Policies)
	} else {
		retDat.Policies = 0
	}
	if contDat.Stringencies.StringA != 0 {
		retDat.Stringency = contDat.Stringencies.StringA
	} else if contDat.Stringencies.Stringency != 0 {
		retDat.Stringency = contDat.Stringencies.Stringency
	} else {
		retDat.Stringency = -1
	}
	return retDat
}

func checkPolicyCache(countName string) vars.PolicyQuery {
	retCase := vars.PolicyQuery{}
	var newCache []vars.PolicyQuery
	notFound := true
	for _, entry := range vars.PolicyCache {
		if notFound && entry.Stringencies.Alpha == countName {
			retCase = entry
			notFound = false
		}
		if entry.DateQueried-vars.StartTime < vars.MaxCacheTime {
			newCache = append(newCache, entry)
		}
		if vars.MaxCacheSize <= len(newCache) {
			break
		}
	}
	vars.PolicyCache = newCache
	return retCase
}
