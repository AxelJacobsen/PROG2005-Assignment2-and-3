package endpoints

import (
	vars "cloudAss2/code/constants"
	inner "cloudAss2/code/endpoints/sharedFuncs/inner"
	webhookFuncs "cloudAss2/code/endpoints/sharedFuncs/webhookFuncs"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

/*
	Handles ensuring the URL follows requirements
*/
func PolicyEntry(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		parts := strings.Split(r.URL.Path, "/")
		err1 := inner.CheckForLegalURL(r.URL.Path, vars.COV_POLICY)
		if len(err1) != 0 {
			http.Error(w, err1, http.StatusBadRequest)
			return
		} else {
			getPolicies(w, r, parts[4])
		}
	default:
		http.Error(w, vars.ERR_ILLEGAL_METHOD, http.StatusMethodNotAllowed)
	}
}

/*
	GET function for policies EP
*/
func getPolicies(w http.ResponseWriter, r *http.Request, countName string) {
	http.Header.Add(w.Header(), "content-type", "application/json")
	if len(countName) != 3 {
		// Implies Using country codes so we need to get them from file
		temp := inner.GetCountNameAlpha(countName)
		if len(temp) != 0 {
			countName = temp
		} else {
			http.Error(w, "Country name not recognized", http.StatusBadRequest)
		}
	}

	var potentialPolicyQuery vars.PolicyQuery

	dateScope, exists := r.URL.Query()["scope"]
	if exists && len(dateScope[0]) != 0 {
		_, err := time.Parse("2006-02-01", dateScope[0])
		if err != nil {
			http.Error(w, "Error, incorrect date format in 'scope' variable.", http.StatusBadRequest)
		}
		potentialPolicyQuery = checkPolicyCache(countName, dateScope[0])
	} else {
		potentialPolicyQuery = checkPolicyCache(countName, time.Now().Format("2006-01-02"))
	}

	//Checks the cache to see if the query has recently been searched

	if (potentialPolicyQuery.Stringencies == vars.PolicyQuery{}.Stringencies) {
		querystring := vars.QUE_OXFORD + countName + "/" + time.Now().Format("2006-01-02")
		//Grabs date from URL query if there is one
		if exists && len(dateScope[0]) != 0 {
			querystring = vars.QUE_OXFORD + countName + "/" + dateScope[0]
		}

		//Perform get request
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

		//Define empty PolicyQueries var
		var queryData vars.PolicyQuery
		err5 := json.Unmarshal(parsedData, &queryData)
		if err5 != nil {
			http.Error(w, err5.Error(), http.StatusInternalServerError)
		}
		potentialPolicyQuery = queryData
		potentialPolicyQuery.DateQueried = int(time.Now().Unix())

		if len(potentialPolicyQuery.Stringencies.Alpha) == 0 {
			potentialPolicyQuery.Stringencies.Alpha = countName
		}

		if len(potentialPolicyQuery.Stringencies.Date) == 0 {
			potentialPolicyQuery.Stringencies.Date = time.Now().Format("2006-01-02")
		}

		vars.PolicyCache = append(vars.PolicyCache, potentialPolicyQuery)
	}

	var printData vars.PolicyResponse
	printData = setPrintData(printData, potentialPolicyQuery)

	//Prepares data for writing to webpage
	screenDat, err6 := json.Marshal(&printData)
	if err6 != nil {
		http.Error(w, err6.Error(), http.StatusInternalServerError)
	}
	// Starts new thread to update webhook since nothing is dependent on this happening
	// Before anything else
	go webhookFuncs.IncrementWebhookEntry(inner.GetCountNameAlpha(countName))
	w.Write(screenDat)
}

/*
	Manually sets data into correct struct format so it can be printed later
*/
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

/*
	Cache function for Policy
*/
func checkPolicyCache(countName string, date string) vars.PolicyQuery {
	retCase := vars.PolicyQuery{}
	currentTime := int(time.Now().Unix())
	var newCache []vars.PolicyQuery
	notFound := true
	countName = strings.ToUpper(countName)
	//Iterates a policy query cache
	for _, entry := range vars.PolicyCache {
		if notFound && entry.Stringencies.Alpha == countName && entry.Stringencies.Date == date {
			retCase = entry
			notFound = false
		}
		//If a query is older than permitted age it will not be carried onwards
		if currentTime-entry.DateQueried < vars.MaxCacheTime {
			newCache = append(newCache, entry)
		}
		if vars.MaxCacheSize <= len(newCache) {
			break
		}
	}
	vars.PolicyCache = newCache
	return retCase
}
