## Cloud function with e2e testing

**Prerequisites:**

 - install gcloud with alpha components (gcloud components install)
 - go1.11
 - service account for e2e testing dedicated project with enabled API
    * regular secrets file -> ~/.secret/e2e.json from http://console.cloud.google.com/
    * firebase secrets file -> ~/.secret/fbc.json from http://http://console.firebase.google.com/
    * [Setup endly google secrets](https://github.com/viant/endly/tree/master/doc/secrets#gc)
 
 
#### Introduction


Cloud Functions is a serverless compute service that runs a specified  code in response to events and automatically manages the underlying compute resources.
The event can be fired by a specific trigger, which determines how and when your function executes.

Cloud function handler supports the following function signatures: 
 
- HTTP Endpoint
    func(http.ResponseWriter, *http.Request)

- Other Events:
    func(context.Context, Event) (error)


TODO: update this project with defined examples once cloud function goes to Beta
