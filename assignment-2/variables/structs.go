package variables

/*
	Cases related structs
*/
type CasesQuery struct {
	Datas Data `json:"data"`
	Qdate int  `json:"date_queried"`
}

type Data struct {
	Countries Country `json:"country"`
}

type Country struct {
	Name   string     `json:"name"`
	Recent MostRecent `json:"mostRecent"`
}

type MostRecent struct {
	Name       string  `json:"name"`
	Date       string  `json:"date"`
	Confirmed  int     `json:"confirmed"`
	Recovered  int     `json:"recovered"`
	Deaths     int     `json:"deaths"`
	GrowthRate float64 `json:"growthRate"`
}

/*
	Policy related structs
*/

type PolicyQuery struct {
	Policies     []PolicyType     `json:"policyActions"`
	Stringencies PolicyDataHolder `json:"stringencyData"`
	DateQueried  int              `json:"date_queried"`
}

type PolicyType struct {
	TypeCode string `json:"policy_type_code"`
}

type PolicyDataHolder struct {
	Alpha      string  `json:"country_code"`
	Date       string  `json:"date_value"`
	Stringency float64 `json:"stringency"`
	StringA    float64 `json:"stringency_actual"`
}

type PolicyResponse struct {
	Alpha      string  `json:"country_code"`
	Date       string  `json:"scope"`
	Stringency float64 `json:"stringency"`
	Policies   int     `json:"polices"`
}

/*
	Webhook related Structs
*/
type Webhook struct {
	Id      string `json:"webhook_id"`
	Url     string `json:"url"`
	Country string `json:"country"`
	Calls   int    `json:"calls"`
}

type WebhookInvo struct {
	Id      string `json:"webhook_id"`
	Country string `json:"country"`
	Calls   int    `json:"calls"`
}

type WebhookRef struct {
	Id     string `json:"webhook_id"`
	Name   string `json:"webhook_name"`
	Called int    `json:"called_times"`
}

type RetId struct {
	Id string `json:"webhook_id"`
}

/*
	Status endpoint struct
*/
type Status struct {
	CasesApi  int    `json:"cases_api"`
	PolicyApi int    `json:"policy_api"`
	WebhookA  int    `json:"webhooks"`
	Version   string `json:"version"`
	Uptime    int    `json:"uptime"`
}

/*
Alpha3 convertion related structs
*/
type Alpha3s struct {
	Alpha3s []Alpha3 `json:"countries"`
}

type Alpha3 struct {
	Name string `json:"name"`
	Code string `json:"code"`
}
