package endpoints

import (
	vars "cloudAss2/code/constants"
	"net/http"
)

/*
	Any query before /corona/v1 will display the bellow message

	Written as an error because it was simpler
*/
func PublicEntry(w http.ResponseWriter, r *http.Request) {
	http.Error(w, vars.PUBLIC_SHORT, http.StatusOK)
}

/*
	Any illegal query after /corona/v1/ will display the bellow message

	This is just because if the user did something after the initial
	/corona/v1, then likely they have missedtyped an endpoint or somehting similar

*/
func IncorrectURL(w http.ResponseWriter, r *http.Request) {
	http.Error(w, vars.PUBLIC_ILLEGAL, http.StatusBadRequest)
}
