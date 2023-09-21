package sharedfuncs

import (
	"bytes"
	vars "cloudAss2/code/constants"
	"encoding/json"
	"net/http"
	"time"

	"google.golang.org/api/iterator"
)

/*
	Reads all webhooks from firestore at the start of the program
*/
func ReadHooks() error {
	webHookData := vars.Client.Collection(vars.FS_WEBHOOK_PATH).Documents(vars.Ctx)
	defer webHookData.Stop()
	if webHookData != nil {
		for {
			//Iterates through all files in the firestore
			doc, err := webHookData.Next()
			if err == iterator.Done {
				break
			}
			//Casts webhook into a new webhook struct
			var hook vars.Webhook
			err1 := doc.DataTo(&hook)
			filename := doc.Ref.ID
			if err1 != nil {
				return err1
			} else {
				//Adds webhooks to cache
				vars.Gwebhooks = append(vars.Gwebhooks, hook)
				vars.WebhookRefs = append(vars.WebhookRefs, vars.WebhookRef{Id: hook.Id, Name: filename, Called: 0})
			}
		}
	}
	return nil
}

/*
	Handles updating all webhooks whos country have been called
*/
func IncrementWebhookEntry(countryName string) error {
	var err error
	if countryName != "" {
		for _, it := range vars.Gwebhooks {
			if it.Country == countryName {
				//Calls incrementing function for a given webhook
				triggCheck := IncrementWebhookInner(it.Id)
				//Checks if the webhook now is above the trigger threshhold
				err = CheckForTriggeredWebhook(it, triggCheck)
			}
		}
	}
	return err
}

/*
	Increments a given webhooks "called" value
*/
func IncrementWebhookInner(id string) int {
	//Itterates all webhooks looking for one with the corresponding ID
	for i, refIt := range vars.WebhookRefs {
		if refIt.Id == id {
			//Increments
			vars.WebhookRefs[i].Called++
			return i
		}
	}
	return 0
}

/*
	Checks if a given webhook has reached its trigger threshold
*/
func CheckForTriggeredWebhook(max vars.Webhook, count int) error {
	//"count" variable is used to avoid iterating the whole list again
	// "max" contains the webhook "parent" which holds the trigger limit for the "child" webhook
	// in reality these are the same webhook, they just contain some different data
	if max.Calls <= vars.WebhookRefs[count].Called {
		//Initializes a dummy webhook to print to the screen, filled with "parent" data
		postData := vars.WebhookInvo{Id: max.Id, Country: max.Country, Calls: max.Calls}

		json_data, err := json.Marshal(postData)
		if err != nil {
			return err
		}
		// Honestly this is just how i interpreted it, i set the current calls back to 0
		// Otherwise the webhooks would be spaming their trigger since they will all be increased
		// During testing and such.
		vars.WebhookRefs[count].Called -= max.Calls
		print(max.Url + "\n")
		//Sends the POST request to the given URL in the webhook
		request, err2 := http.NewRequest("POST", max.Url, bytes.NewBuffer(json_data))
		client := &http.Client{Timeout: time.Second * 20}
		response, err3 := client.Do(request)
		if err3 != nil {
			return err3
		}
		defer response.Body.Close()
		return err2
	}
	return nil
}
