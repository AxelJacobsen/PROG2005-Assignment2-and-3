package sharedfuncs_test

import (
	vars "cloudAss2/code/constants"
	shared "cloudAss2/code/endpoints/sharedFuncs/inner"
	"testing"
)

/*
	Important note about this test is that the function is designed to
	match a given url1 up against a preset url2, so url1 has to be longer
	or same length as url2
*/
func TestCompURL(t *testing.T) {
	//List of urls to be compared
	url1s := []string{
		"thisUrlIs/TooShort",
		"corona/v1/cases/thisisJustRight",
		"here/is/an/example/that/it/doesntCareAbout/length",
		"ThisUrlisnttooshort/itsjustallwrong",
	}

	//List of Urls to be compared to
	url2s := []string{
		"thisUrlIs/TooShort/compared/to/thisOne",
		"corona/v1/cases/",
		"here/is/an/example/",
		"twitch.tv/Cowspace",
	}

	//Excpected results
	results := []bool{
		false,
		true,
		true,
		false,
	}

	//Actual test
	for i := range url1s {
		if shared.CompUrlString(url1s[i], url2s[i]) != results[i] {
			t.Errorf("\nExpected '%s' but got '%s'", url2s[i], url1s[i])
		}
	}
}

/*
	Runs test on alpha3 converter
*/
func TestAlphaCoverter(t *testing.T) {
	shared.MockPath("../../../constants/alpha3.json")
	//List of urls to be compared
	countryNames := []string{
		"NOR",
		"France",
		"NotLegalCountry",
	}

	//List of Urls to be compared to
	result := []string{
		"Norway",
		"FRA",
		"notlegalcountry",
	}

	//Actual test
	for i := range countryNames {
		stow := shared.GetCountNameAlpha(countryNames[i])
		if stow != result[i] {
			t.Errorf("\nExpected '%s' but got '%s'", result[i], stow)
		}
	}
	shared.MockPath(vars.ALPHA3_FILE_PATH)
}

/*
	Runs test on CheckforLegalURL function,
	this is meant to take everything after the "URL header" (localhost/deployment url)
*/
func TestLegalURL(t *testing.T) {
	//List of urls to be compared
	url1s := []string{
		"thisUrlIs/TooShort/compared/to/thisOne",
		"/corona/v1/cases/",
		"/corona/v1/policy/FRA",
		"twitch.tv/Cowspace",
	}

	//List of Urls to be compared to
	url2s := []string{
		vars.COV_CASE,
		vars.COV_CASE,
		vars.COV_POLICY,
		vars.COV_POLICY,
	}

	//Excpected results
	results := []string{
		"Malformed URL",
		"Missing countryname in URL query.",
		"",
		"Malformed URL",
	}

	//Actual test
	for i := range url1s {
		stow := shared.CheckForLegalURL(url1s[i], url2s[i])
		if stow != results[i] {
			t.Errorf("\nExpected '%s' but got '%s'", results[i], stow)
		}
	}
}
