package endpoints

import (
	inner "cloudAss2/endpoints/inner"
	vars "cloudAss2/variables"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func NotifyEntry(w http.ResponseWriter, r *http.Request) {
	http.Header.Add(w.Header(), "content-type", "application/json")
	switch r.Method {
	case http.MethodGet:
		parts := strings.Split(r.URL.Path, "/")
		if len(parts) < 5 || !inner.CompUrlString(r.URL.Path, vars.COV_URL+vars.COV_NOTIFY) {
			http.Error(w, "Malformed URL", http.StatusBadRequest)
		} else {
			getEntry, err := getWebhooks(parts[4])
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			} else {
				json.NewEncoder(w).Encode(getEntry)
			}
		}
	case http.MethodPost:
		if r.Body != nil {
			postData, postErr := postHooks(r)
			if postErr != nil {
				http.Error(w, postErr.Error(), http.StatusInternalServerError)
			} else {
				json.NewEncoder(w).Encode(postData)
			}
		} else {
			http.Error(w, "No body in post request", http.StatusBadRequest)
		}

	case http.MethodDelete:
		parts := strings.Split(r.URL.Path, "/")
		if len(parts) < 5 || !inner.CompUrlString(r.URL.Path, vars.COV_URL+vars.COV_NOTIFY) {
			http.Error(w, "Malformed URL", http.StatusBadRequest)
		} else {
			getEntry, err := deleteWebhook(parts[4])
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			} else {
				json.NewEncoder(w).Encode(getEntry)
			}
		}
	default:
		http.Error(w, vars.ERR_ILLEGAL_METHOD, http.StatusMethodNotAllowed)
	}
}

func getWebhooks(id string) ([]vars.Webhook, error) {
	if len(id) != 0 {
		retHook := []vars.Webhook{}
		for _, hookIt := range vars.Gwebhooks {
			if hookIt.Id == id {
				retHook = append(retHook, hookIt)
				return retHook, nil
			}
		}
	} else {
		return vars.Gwebhooks, nil
	}
	return []vars.Webhook{}, nil
}

func postHooks(r *http.Request) (vars.RetId, error) {
	// Expects incoming body in terms of WebhookRegistration struct
	newHook := vars.Webhook{}
	err := json.NewDecoder(r.Body).Decode(&newHook)
	if err != nil {
		return vars.RetId{}, err
	}
	newHook.Id = genWebHookId()
	err1 := writeWebhookToFireStore(newHook)
	if err1 != nil {
		return vars.RetId{}, err1
	}
	vars.Gwebhooks = append(vars.Gwebhooks, newHook)
	id := vars.RetId{Id: newHook.Id}
	return id, nil
}

func deleteWebhook(id string) (string, error) {
	if len(id) != 0 {
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
	} else {
		return "Missing id in DELETE request.", nil
	}
	return "Couldnt find id in registered webhooks", nil
}

func writeWebhookToFireStore(content vars.Webhook) error {
	_, _, err := vars.Client.Collection(vars.FS_WEBHOOK_PATH).Add(vars.Ctx, content)
	return err
}

func deleteWebhookFromFireStore(hookName string) error {
	_, err := vars.Client.Collection(vars.FS_WEBHOOK_PATH).Doc(hookName).Delete(vars.Ctx)
	return err
}

func genWebHookId() string {
	return strconv.Itoa(int(time.Now().UnixNano()))
}
