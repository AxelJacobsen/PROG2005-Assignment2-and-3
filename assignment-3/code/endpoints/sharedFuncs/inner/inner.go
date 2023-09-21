package sharedfuncs

import (
	vars "cloudAss2/code/constants"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

/*
	Takes two strings, splits them and checks if the first components are the same
*/
func CompUrlString(str1 string, str2 string) bool {
	if len(str1) < len(str2) {
		return false
	}

	parts1 := strings.Split(strings.ToLower(str1), "/")
	parts2 := strings.Split(strings.ToLower(str2), "/")

	for i := 0; i < len(parts2)-1; i++ {
		if parts1[i] != parts2[i] {
			return false
		}
	}
	return true
}

var mockAlpha3Path = vars.ALPHA3_FILE_PATH

func MockPath(newPath string) {
	mockAlpha3Path = newPath
}

/*
	Converts a given Alpha3 code to Country name, or visa versa
*/
func GetCountNameAlpha(text string) string {
	//If the cache hasnt been filled, enter
	if len(vars.AllAlpha3s.Alpha3s) == 0 {
		jsonFile, err0 := os.Open(mockAlpha3Path)
		if err0 != nil {
			log.Fatal(err0)
		}
		defer jsonFile.Close()
		//Reads all codes from file.
		byteValue, _ := ioutil.ReadAll(jsonFile)

		json.Unmarshal(byteValue, &vars.AllAlpha3s)
	}

	//Checks if the request is for Alpha3 -> Country Name or opposite
	if len(text) == 3 {
		//Ensures non-cases sensitivity
		text = strings.ToUpper(text)
		for i := 0; i < len(vars.AllAlpha3s.Alpha3s); i++ {
			if vars.AllAlpha3s.Alpha3s[i].Code == text {
				return vars.AllAlpha3s.Alpha3s[i].Name
			}
		}
	} else {
		//Ensures that it isnt case sensitive
		text = strings.ToLower(text)
		for i := 0; i < len(vars.AllAlpha3s.Alpha3s); i++ {
			if strings.ToLower(vars.AllAlpha3s.Alpha3s[i].Name) == text {
				return vars.AllAlpha3s.Alpha3s[i].Code
			}
		}
	}

	return text
}

/*
	Added this to remove some repetetive error handling
*/
func CheckForLegalURL(URL string, controlURL string) string {
	parts := strings.Split(URL, "/")
	if len(parts) < 4 || !CompUrlString(URL, vars.COV_URL+controlURL) {
		return "Malformed URL"
	} else {
		//Exception for Notification EP since it doesnt need a final parameter
		if controlURL != vars.COV_NOTIFY && controlURL != vars.COV_STATUS {
			if 4 < len(parts) && len(parts[4]) != 0 {
				return ""
			} else {
				return "Missing countryname in URL query."
			}
		} else {
			return ""
		}
	}
}
