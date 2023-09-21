package inner

import (
	vars "cloudAss2/variables"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

/*
	Takes two strings, splits them and checks if the first components are the same

*/
func CompUrlString(str1 string, str2 string) bool {
	parts1 := strings.Split(str1, "/")
	parts2 := strings.Split(str2, "/")
	for i := 0; i < len(parts2)-1; i++ {
		if parts1[i] != parts2[i] {
			fmt.Println("Discrepancy detected: " + parts1[i] + " and: " + parts2[i] + ", are not the same")
			return false
		}
	}
	return true
}

func GetCountNameAlpha(text string, nameOrAlpha bool) string {
	jsonFile, err0 := os.Open("alpha3.json")
	if err0 != nil {
		log.Fatal(err0)
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	var aCounts vars.Alpha3s

	json.Unmarshal(byteValue, &aCounts)

	if nameOrAlpha {
		for i := 0; i < len(aCounts.Alpha3s); i++ {
			if aCounts.Alpha3s[i].Code == text {
				return aCounts.Alpha3s[i].Name
			}
		}
	} else {
		for i := 0; i < len(aCounts.Alpha3s); i++ {
			if aCounts.Alpha3s[i].Name == text {
				return aCounts.Alpha3s[i].Code
			}
		}
	}

	return text
}
