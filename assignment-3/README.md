# Assignment 2

## How to use

For the deployed version it functions like the PROG2005 wiki describes,<br />
you can find some example urls further down in the README.<br />

### Notification
For the notifications endpoint i suggest using this site:<br />
https://webhook.site/#!/6a9d60f1-3af2-458b-8439-1a97cfa4adbf/72f4669e-4746-4418-88f4-e8b88a4c786e/1<br />
NOTE: it might give you a different url due to how the site works, but if you use the URL marked:<br />
"Your unique URL (Please copy it from here, not from the address bar!)"<br />
As the post body then you will get a notification on the original site when the webhook triggers<br />
Here is an example body that might work for you if the URL doesnt change as explained above.<br />
```
{
    "url": "https://webhook.site/6a9d60f1-3af2-458b-8439-1a97cfa4adbf",
    "country": "Jamaica",
    "calls": 2
}
```

### Deployed URL Example:
http://34.65.251.24:8080/corona/v1/status/<br /> 
http://34.65.251.24:8080/corona/v1/cases/Norway<br /> 
http://34.65.251.24:8080/corona/v1/policy/France?scope=2022-01-01<br /> 
http://34.65.251.24:8080/corona/v1/notifications/<br /> 

### Local URL Example:
http://localhost:8080/corona/v1/status/<br /> 
http://localhost:8080/corona/v1/cases/Norway<br /> 
http://localhost:8080/corona/v1/policy/France?scope=2022-01-01<br /> 
http://localhost:8080/corona/v1/notifications/<br /> 

### Testing:
To Run testing, run this command in the root directory: <br />
go test ./code/endpoints/... -cover <br />
This will run all test files in the subfolders<br />
<br />
Testing has not been implemented for notifications and policy as of delivery.<br />

## Structure
```
├──code
│   ├── constants
│   │    ├── constants.go
│   │    ├── alpha3.json
│   │    ├── globals.json
│   │    └── structs.go
│   └── endpoints
│        ├── cases
│        │  ├── cases_test.go
│        │  └── cases.go
│        ├── notification
│        │  ├── notification_test.go
│        │  └── notification.go<
│        ├── policy
│        │  ├── policy_test.go
│        │  └── policy.go
│        ├── public
│        │  ├── public_test.go
│        │  └── public.go
│        └── sharedFuncs
│             ├── inner.go
│             │   ├── inner_test.go
│             │   └── inner.go
│             └── webhookFuncs
│                 └── webhookFuncs.go
├── .gitignore
├── ass2-service-account.json
├── Dockerfile
├── go.mod
├── go.sum
├── main.go 
└── README.md
```

### Used APIs:
* Covid 19 Cases API: https://covid19-graphql.vercel.app/
* Corona Policy Stringency API: https://covidtracker.bsg.ox.ac.uk/about-api


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

### Feedback Notes / Updates from Assignment 2 delivery
Structure: I got some feedback about the structuring being a bit flat, and due to testing requiring thing to be in the same folder.<br />
I have changed the structure significantly to have sepperate endpoints in each their own folder paired with their test file.<br />

Commenting: I got feedback that there was inconsistent / not enough commenting so i have added significantly more comments where <br />
i felt that the code might be a little confusing.<br />

Deleting endpoints: Someone reported that deleting endpoints didnt work propperly, i realized that i had made a mistake where only <br />
webhooks from the firestore database could be deleted, and not newly added ones. This was only the case for the deployed version, and<br />
honestly i have no idea how i didnt see it originally. This has been fixed and it should now properly delete any webhook.<br />

Case sensitivity: I used to have case sensitivity in the URL queries, this should now be fixed and the program should be completely,<br />
case insensitive. This includes country names used in the post body of policy queires.<br />

### Known Issues
Though most of them should be fixed now, there might still be some inconsistencies in the alpha3 file,<br />
as of now the only one i cant fix is Kosovo due to it not having a recognized Alpha3 code in the Policy API<br />

Due to adding case insensitivity, there might be some nations where using the country name could result in incorrect queries,<br />
though these countries have been problematic due to naming convention used in the supplied API.<br />
