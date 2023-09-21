package variables

import (
	"context"

	"cloud.google.com/go/firestore"
)

const CURRENT_VER = "v1"
const COV = "/corona/"
const COV_URL = COV + CURRENT_VER

const COV_CASE = "/cases/"
const COV_POLICY = "/policy/"
const COV_STATUS = "/status/"
const COV_NOTIFY = "/notifications/"

const QUE_GRAPHQL = "https://covid19-graphql.vercel.app/"
const QUE_OXFORD_STATUS = "https://covidtrackerapi.bsg.ox.ac.uk/api/"
const QUE_OXFORD_END = "v2/stringency/actions/"
const QUE_OXFORD = QUE_OXFORD_STATUS + QUE_OXFORD_END

const LOCAL_PORT = "8080"
const LOCAL_HOST_PRE = "http://localhost:" + LOCAL_PORT
const LOCAL_HOST_URL_TOT = LOCAL_HOST_PRE + COV_URL

const ERR_ILLEGAL_METHOD = "The utilized method is not valid for this endpoint."

const FIRESTORE_FILENAME = "ass2-service-account.json"

const FS_WEBHOOK_PATH = "webhooks"

var Client *firestore.Client
var Ctx context.Context
var Gwebhooks []Webhook
var WebhookRefs []WebhookRef

var CasesCache []CasesQuery
var PolicyCache []PolicyQuery

var StartTime int

//How many days the cache is held, only deletes on query, so it wont delete if the server is idle
const MaxCacheTime = (60 * 60 * 24) * 2 //Update last number for amount of days

const MaxCacheSize = 15 //Max limit on how many queries can be in the cache
