package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	funcs "cloudAss2/endpoints"
	inner "cloudAss2/endpoints/inner"
	vars "cloudAss2/variables"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

func main() {
	vars.Client, vars.Ctx = getFirestoreClient()

	vars.StartTime = int(time.Now().Unix())

	err := inner.ReadHooks()
	if err != nil {
		log.Fatal(err.Error())
	}

	//Public endpoint
	http.HandleFunc(vars.COV_URL, funcs.PublicEntry)
	http.HandleFunc(vars.COV, funcs.PublicEntry)
	http.HandleFunc("/", funcs.PublicEntry)

	//Cases endpoint
	http.HandleFunc(vars.COV_URL+vars.COV_CASE, funcs.CasesEntry)

	//Policy endpoint
	http.HandleFunc(vars.COV_URL+vars.COV_POLICY, funcs.PolicyEntry)

	//Status endpoint
	http.HandleFunc(vars.COV_URL+vars.COV_STATUS, funcs.StatusEntry)

	//Notification endpoint
	http.HandleFunc(vars.COV_URL+vars.COV_NOTIFY, funcs.NotifyEntry)

	listenFromServer()
	defer vars.Client.Close()
}

func listenFromServer() {
	// Make it Heroku-compatible
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

	// Alternative setup, directly through Firestore (without initial reference to Firebase); but requires Project ID
	// client, err := firestore.NewClient(ctx, projectID)

	if err != nil {
		log.Fatalln(err)
	}

	// Close down client
	/* 	defer func() {
		err := client.Close()
		if err != nil {
			log.Fatal("Closing of the firebase client failed. Error:", err)
		}
	}() */
	return client, ctx
}
