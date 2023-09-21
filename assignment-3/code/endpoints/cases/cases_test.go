package endpoints_test

import (
	vars "cloudAss2/code/constants"
	cases "cloudAss2/code/endpoints/cases"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestCases(t *testing.T) {
	methods := []string{http.MethodGet, http.MethodGet, http.MethodGet, http.MethodPut}

	testNames := []string{"Expected query", "Missing Name query", "Illegal name query", "Illegal method"}

	mock := []string{`{
	"data":{
		"country": {
			"name": "Jamaica",
			"mostRecent": {
				"date": "2022-04-30",
				"confirmed": 129978,
				"deaths": 2962,
				"recovered": 0,
				"growthRate": 0
				}
			}
		}
	}`, `{
		"country": {
			"name": null
		}
	}`, `{
		"country": {
			"name": null
		}
	}`, ""}

	countryName := []string{"Jamaica", "", "nonExistantCountry", ""}

	results := []string{`{"name":"Jamaica","date":"2022-04-30","confirmed":129978,"recovered":0,"deaths":2962,"growthRate":0}`,
		"Missing countryname in URL query.",
		"Couldnt find any data on the given country: nonExistantCountry",
		vars.ERR_ILLEGAL_METHOD,
	}

	// Test subtests
	for i, mockString := range mock {
		t.Run(testNames[i], func(t *testing.T) {
			//Update mock in cases file
			cases.Query = []byte(mockString)

			// Send request
			req := httptest.NewRequest(methods[i], vars.LOCAL_HOST_URL_TOT+vars.COV_CASE+countryName[i], nil)

			w := httptest.NewRecorder()
			cases.CasesEntry(w, req)
			res := w.Result()
			defer res.Body.Close()

			//Read respons data
			data, err := ioutil.ReadAll(res.Body)
			if err != nil {
				t.Error("Error reading body of Cases test result", err)
			}
			strDat := string(data)
			if i != 0 {
				strDat = strDat[:len(strDat)-1]
			}
			if strDat != results[i] {
				t.Errorf("Expected '%s' but got '%v'", results[i], strDat)
			}
		})
	}
	cases.Query = []byte{}
}

/*
	Tests the CheckCasesCache function in the cases file
*/
func TestCasesCache(t *testing.T) {
	//Establish template data
	tempQuery := vars.CasesQuery{
		Datas: vars.Data{
			Countries: vars.Country{
				Name: "Bhutan",
				Recent: vars.MostRecent{
					Name:       "Bhutan",
					Date:       "9-11-2001",
					Confirmed:  25000,
					Recovered:  0,
					Deaths:     2996,
					GrowthRate: 0,
				},
			},
		},
		//ensures that the cache will delete the hook in 20 sek if it somehow gets retained
		Qdate: int(time.Now().Unix()) - (vars.MaxCacheTime - 20),
	}
	//Adds new query to cache
	vars.CasesCache = append(vars.CasesCache, tempQuery)
	//checks for the newly added country
	result := cases.CheckCasesCache("Bhutan")
	if result != tempQuery {
		t.Errorf("Expected '%s' but got '%v'", tempQuery.Datas.Countries.Name, result.Datas.Countries.Name)
	}
	//Runs a check for a non existent hook
	result2 := cases.CheckCasesCache("NonExsistenAndVeryImaginaryCountry")
	empty := vars.CasesQuery{}
	if result2 != empty {
		t.Errorf("Expected 'Empty data' but got '%v'", result2.Datas.Countries.Name)
	}
}
