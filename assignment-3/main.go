package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	vars "cloudAss2/code/constants"
	cases "cloudAss2/code/endpoints/cases"
	notify "cloudAss2/code/endpoints/notification"
	policy "cloudAss2/code/endpoints/policy"
	public "cloudAss2/code/endpoints/public"
	webhooks "cloudAss2/code/endpoints/sharedFuncs/webhookFuncs"
	status "cloudAss2/code/endpoints/status"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

func main() {
	vars.Client, vars.Ctx = getFirestoreClient()
	//Gets startup time for status endpoint
	vars.StartTime = int(time.Now().Unix())

	err := webhooks.ReadHooks()
	if err != nil {
		log.Fatal(err.Error())
	}

	//Public endpoint
	http.HandleFunc(vars.COV_URL, public.IncorrectURL)
	http.HandleFunc("/", public.PublicEntry)

	//Cases endpoint
	http.HandleFunc(vars.COV_URL+vars.COV_CASE, cases.CasesEntry)
	http.HandleFunc(vars.COV_URL+vars.COV_CASE_ALT, cases.CasesEntry)

	//Policy endpoint
	http.HandleFunc(vars.COV_URL+vars.COV_POLICY, policy.PolicyEntry)
	http.HandleFunc(vars.COV_URL+vars.COV_POLICY_ALT, policy.PolicyEntry)

	//Status endpoint
	http.HandleFunc(vars.COV_URL+vars.COV_STATUS, status.StatusEntry)
	http.HandleFunc(vars.COV_URL+vars.COV_STATUS_ALT, status.StatusEntry)

	//Notification endpoint
	http.HandleFunc(vars.COV_URL+vars.COV_NOTIFY, notify.NotifyEntry)
	http.HandleFunc(vars.COV_URL+vars.COV_NOTIFY_ALT, notify.NotifyEntry)

	listenFromServer()
	defer vars.Client.Close()
}

/*
	Initializes server with correct port
*/
func listenFromServer() {
	port := os.Getenv("PORT")
	if port == "" {
		port = vars.LOCAL_PORT
	}

	addr := ":" + port

	log.Printf("Listening on %s ...\n", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		panic(err)
	}
}

func getFirestoreClient() (*firestore.Client, context.Context) {
	// Firebase initialisation
	ctx := context.Background()

	// We use a service account, load credentials file that you downloaded from your project's settings menu.
	// It should reside in your project directory.
	// Make sure this file is git-ignored, since it is the access token to the database.
	sa := option.WithCredentialsFile(vars.FIRESTORE_FILENAME)
	app, err := firebase.NewApp(ctx, nil, sa)
	if err != nil {
		log.Fatalln(err)
	}

	// Instantiate client
	client, err := app.Firestore(ctx)

	if err != nil {
		log.Fatalln(err)
	}
	return client, ctx
}
