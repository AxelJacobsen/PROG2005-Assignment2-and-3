package endpoints

import (
	vars "cloudAss2/code/constants"
	inner "cloudAss2/code/endpoints/sharedFuncs/inner"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"
)

/*
	Error handles for the Notification endpoint and ensures legal URL
*/
func NotifyEntry(w http.ResponseWriter, r *http.Request) {
	http.Header.Add(w.Header(), "content-type", "application/json")
	switch r.Method {
	// If the endpoitn method is Get:
	case http.MethodGet:
		// Splits and checks for legal URL
		parts := strings.Split(r.URL.Path, "/")
		err0 := inner.CheckForLegalURL(r.URL.Path, vars.COV_NOTIFY)
		if len(err0) != 0 {
			http.Error(w, err0, http.StatusBadRequest)
		} else {
			//If there arent any URL problems call get function
			notifyID := ""
			if 4 < len(parts) {
				notifyID = parts[4]
			}
			getEntry, err := getWebhooks(notifyID)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			} else {
				json.NewEncoder(w).Encode(getEntry)
			}
		}
	// If the endpoint method is Post:
	case http.MethodPost:
		if r.Body != nil {
			// If there is a body in the post request, call
			postData, postErr := postHooks(r)
			if postErr != nil {
				http.Error(w, postErr.Error(), http.StatusBadRequest)
			} else {
				//Write to screen if there were no errors
				json.NewEncoder(w).Encode(postData)
			}
		} else {
			http.Error(w, "No body in post request", http.StatusBadRequest)
		}
	// If the endpoint method is Delete:
	case http.MethodDelete:
		//Ensure legal URL
		parts := strings.Split(r.URL.Path, "/")
		err0 := inner.CheckForLegalURL(r.URL.Path, vars.COV_NOTIFY)
		if len(err0) != 0 {
			http.Error(w, err0, http.StatusBadRequest)
		} else {
			//If no errors calls delete function
			getEntry, err := deleteWebhook(parts[4])
			if err != nil {
				http.Error(w, "There was an error int deleting webhook\n."+err.Error(), http.StatusInternalServerError)
			} else {
				//Prints success message to screen
				err1 := json.NewEncoder(w).Encode(getEntry)
				if err1 != nil {
					http.Error(w, "There was an error in printing delete to screen.\n"+err.Error(), http.StatusInternalServerError)
				}
			}
		}
	default:
		http.Error(w, vars.ERR_ILLEGAL_METHOD, http.StatusMethodNotAllowed)
	}
}

/*
	Get endpoint handler, returns errors and content for Entry function to handle
*/
func getWebhooks(id string) ([]vars.Webhook, error) {
	if len(id) != 0 {
		retHook := []vars.Webhook{}
		for _, hookIt := range vars.Gwebhooks {
			if hookIt.Id == id {
				retHook = append(retHook, hookIt)
				return retHook, nil
			}
		}
		err := errors.New("couldnt find a webhook with the supplied ID")
		return retHook, err
	} else {
		if len(vars.Gwebhooks) != 0 {
			return vars.Gwebhooks, nil
		}
		err1 := errors.New("no webhooks registered")
		return []vars.Webhook{}, err1
	}
}

/*
	Post endpoint handler, returns errors and content for Entry function to handle
*/
func postHooks(r *http.Request) (vars.RetId, error) {
	// Expects incoming body in terms of WebhookRegistration struct
	newHook := vars.Webhook{}
	//Decodes data into struct
	err0 := json.NewDecoder(r.Body).Decode(&newHook)
	if err0 != nil {
		return vars.RetId{}, err0
	}
	//Ensures that a legal name has been supplied
	if len(newHook.Country) == 3 {
		newHook.Country = inner.GetCountNameAlpha(newHook.Country)
		if len(newHook.Country) == 3 {
			err := errors.New("unrecognized Alpha3 supplied")
			return vars.RetId{}, err
		}
		//Checks country name up against file to ensure its legal
	} else if len(newHook.Country) != 0 {
		tempNewName := inner.GetCountNameAlpha(newHook.Country)
		if len(tempNewName) != 3 {
			err1 := errors.New("illegal country name supplied")
			return vars.RetId{}, err1
		}
	} else {
		err2 := errors.New("empty country name in post request")
		return vars.RetId{}, err2
	}
	//Ensures capitalization, purely visual due to multiple to lowers where it could matter
	newHook.Country = strings.Title(strings.ToLower(newHook.Country))

	//Ensures a URL is supplied, not picky
	if len(newHook.Url) == 0 {
		err3 := errors.New("no URL supplied in post request")
		return vars.RetId{}, err3
	}
	//Ensures legal calls
	if newHook.Calls <= 0 {
		err4 := errors.New("trigger limit(call) too small, minimum 1")
		return vars.RetId{}, err4
	}

	//Gets a new id based on the current time
	newHook.Id = genWebHookId()
	docName, err1 := writeWebhookToFireStore(newHook)
	if err1 != nil {
		return vars.RetId{}, err1
	}
	tempWebRef := vars.WebhookRef{Id: newHook.Id, Name: docName, Called: 0}
	print("NEW WEBHOOK FILENAME: " + docName)
	//Adds webhook to cache
	vars.Gwebhooks = append(vars.Gwebhooks, newHook)
	vars.WebhookRefs = append(vars.WebhookRefs, tempWebRef)
	id := vars.RetId{Id: newHook.Id}
	return id, nil
}

func deleteWebhook(id string) (string, error) {
	if len(id) == 0 {
		return "Missing id in DELETE request.", nil
	}

	for i, hookIt := range vars.WebhookRefs {
		if hookIt.Id == id {
			err := deleteWebhookFromFireStore(hookIt.Name)
			if i != len(vars.WebhookRefs) {
				vars.WebhookRefs = append(vars.WebhookRefs[:i], vars.WebhookRefs[i+1:]...)
				vars.Gwebhooks = append(vars.Gwebhooks[:i], vars.Gwebhooks[i+1:]...)
			} else {
				var temp1 []vars.WebhookRef
				var temp2 []vars.Webhook
				vars.WebhookRefs = append(temp1, vars.WebhookRefs[:i]...)
				vars.Gwebhooks = append(temp2, vars.Gwebhooks[:i]...)
			}
			return "Webhook with id: " + id + ", deleted.", err
		}
	}
	return "Couldnt find id in registered webhooks", nil
}

/*
	Adds a webhook to the Firestore databse
*/
func writeWebhookToFireStore(content vars.Webhook) (string, error) {
	docRef, _, err := vars.Client.Collection(vars.FS_WEBHOOK_PATH).Add(vars.Ctx, content)
	return docRef.ID, err
}

/*
	Removes a webhook to the Firestore databse
*/
func deleteWebhookFromFireStore(hookName string) error {
	_, err := vars.Client.Collection(vars.FS_WEBHOOK_PATH).Doc(hookName).Delete(vars.Ctx)
	return err
}

/*
	generates a new id for webhooks based on current time
*/
func genWebHookId() string {
	return strconv.Itoa(int(time.Now().UnixNano()))
}
