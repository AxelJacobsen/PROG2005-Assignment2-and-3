# Assignment 2

# Notes

### Alpha3 codes
The list of countries and their alpha3 codes was gotten from here:<br />
https://gist.github.com/bensquire/1ba2037079b69e38bb0d6aea4c4a0229<br />
This list used the formal names of countries, something the grphQl site didnt use<br />
due to no full list of this in the documentation on the site, i went through the added file<br />
and moderated it manualy. The file contains more countries than there is statistics for,<br />
but i think i fixed most names such as: Russian Federation -> Russia.<br />
If you find any more country names that are incorrect, please note them in the review, thanks!<br /><br />

### Policy
Since i had created a funciton that converts country alphacode to name, i decided to apply that here as well.<br />
If you decide to query with full name such as:<br />
http://localhost:8080/corona/v1/policy/France?scope=01-01-2022 <br />
Then it will simply convert the name to Alpha before performaing the query.<br />
I dont think this necessarilly goes against the task, and is extremely easy to remove, so i kept it in.<br /><br />
Another thing is that if you choose to have to scope parameter, but dont follow the YYYY-MM-DD format,<br />
or simply supply it with gibberish, it will give you the results of a no-date query.<br /><br />

### When querying, its case sensitive, so using full name uses capital letter, while Alpha3 uses all capital
Example: 
France == FRA

## Structure
>├── endpoints<br />
>│   ├── internal<br />
>│   |  ├── inner.go<br />
>│   │  └── webhookFuncs.go<br />
>│   ├── cases.go<br />
>│   ├── notification.go<br />
>│   ├── policy.go<br />
>│   ├── public.go<br />
>│   └── status.go<br />
>├── variables<br />
>│   ├── constants.go<br />
>│   └── structs.go<br />
>├── .gitignore<br />
>├── alpha3.json<br />
>├── ass2-service-account.json<br />
>├── Dockerfile<br />
>├── go.mod<br />
>├── go.sum<br />
>├── main.go <br />
>└── README.md<br />

## Endpoints
>/corona/v1/country/<br />
>/corona/v1/policy/<br />
>/corona/v1/diag/<br />
>/corona/v1/notifications/
### Deployed URL Example:
http://10.212.138.83:8080/corona/v1/status/
### Local URL Example:
http://localhost:8080/corona/v1/status/