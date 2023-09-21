package inner

import (
	"bytes"
	vars "cloudAss2/variables"
	"encoding/json"
	"net/http"
	"time"

	"google.golang.org/api/iterator"
)

func ReadHooks() error {
	webHookData := vars.Client.Collection(vars.FS_WEBHOOK_PATH).Documents(vars.Ctx)
	defer webHookData.Stop()
	if webHookData != nil {
		for {
			doc, err := webHookData.Next()
			if err == iterator.Done {
				break
			}
			var hook vars.Webhook
			err1 := doc.DataTo(&hook)
			filename := doc.Ref.ID
			if err1 != nil {
				return err1
			} else {
				vars.Gwebhooks = append(vars.Gwebhooks, hook)
				vars.WebhookRefs = append(vars.WebhookRefs, vars.WebhookRef{Id: hook.Id, Name: filename, Called: 0})
			}
		}
	}
	return nil
}

func IncrementWebhookEntry(countryName string) error {
	var err error
	if countryName != "" {
		for _, it := range vars.Gwebhooks {
			if it.Country == countryName {
				triggCheck := IncrementWebhookInner(it.Id)
				err = CheckForTriggeredWebhook(it, triggCheck)
			}
		}
	}
	return err
}

func IncrementWebhookInner(id string) int {
	for i, refIt := range vars.WebhookRefs {
		if refIt.Id == id {
			vars.WebhookRefs[i].Called++
			return i
		}
	}
	return 0
}

func CheckForTriggeredWebhook(max vars.Webhook, count int) error {
	if max.Calls <= vars.WebhookRefs[count].Called {

		postData := vars.WebhookInvo{Id: max.Id, Country: max.Country, Calls: max.Calls}

		json_data, err := json.Marshal(postData)
		if err != nil {
			return err
		}
		vars.WebhookRefs[count].Called = 0
		print(max.Url + "\n")
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
