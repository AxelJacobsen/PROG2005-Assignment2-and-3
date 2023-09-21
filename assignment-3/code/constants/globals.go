package variables

//Contains all original webhook data in cache:
var Gwebhooks []Webhook

//Contans a refrence to the original webhook in cache:
var WebhookRefs []WebhookRef

//Contains recent seareches for cases and Policy in cache:
var CasesCache []CasesQuery
var PolicyCache []PolicyQuery

//Holds all alpha3 codes
var AllAlpha3s Alpha3s
