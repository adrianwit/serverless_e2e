## Cloud function with e2e testing

**Prerequisites:**

 - install gcloud with alpha components (gcloud components install alpha)
 - go1.11
 - service account for e2e testing dedicated project with enabled API
    * regular secrets file -> ~/.secret/e2e.json from http://console.cloud.google.com/
    * firebase secrets file -> ~/.secret/fbc.json from http://http://console.firebase.google.com/
    * [Setup endly google secrets](https://github.com/viant/endly/tree/master/doc/secrets#gc)
    * [endly e2e runner](https://github.com/viant/endly/releases) (0.29+)



### Running e2e tests with endly docker container

```bash
mkdir -p ~/e2e
mkdir -p ~/.secret
docker run --name endly -v /var/run/docker.sock:/var/run/docker.sock -v ~/e2e:/e2e -v ~/.secret/:/root/.secret/ -p 7722:22  -d endly/endly:latest-ubuntu16.04  

ssh root@127.0.0.1 -p 7722 ## password is dev

#### all operation now taking place in endly docker container

endly -v #to check version

endly -c=localhost  ## create localhost ssh credentials with user root/dev
ls -al /root/.secret/localhost.json ## check encrypted credentials created

## generate google platform project secret -> /root/.secret/e2e.json  
## generate firestore/firebase project secret -> /root/.secret/fbc.json

cd /e2e
git clone https://github.com/adrianwit/serverless_e2e
cd  serverless_e2e/cloud_function/e2e
endly -r=run
```


 
 
#### Introduction


Cloud Functions can be fired by a specific trigger. 
The trigger type determines how and when your function executes. 

Cloud function handler supports the following function signatures: 
 
- HTTP Endpoint
    ```func(http.ResponseWriter, *http.Request)```

- Other Events:
    ```func(context.Context, Event) (error)```



This project provides example for the following native trigger mechanisms:




#### HTTP Endpoint

1. Hello World
- [HelloFn](hello/hello.go)
- [E2E Use Case](e2e/regression/cases/001_hello_world)
- _Reference_: 
    * [Calling HTTP](https://cloud.google.com/functions/docs/calling/http)


2. BigQuery reader
- [QueryFn](query/query_reader.go)
- [E2E Use Case](e2e/regression/cases/002_query_data)
- _Reference_: 
    * [Calling HTTP](https://cloud.google.com/functions/docs/calling/http)


#### Google Storage

- [SplitWorkloadFn](split/split_workload.go)
- [E2E Use Case](e2e/regression/cases/003_split_workload)
- Google Storage Event:
```go
type GSEvent struct {
        Bucket                  string `json:"bucket"`
        Name                    string `json:"name"`
        ContentType             string `json:"contentType"`
        CRC32C                  string `json:"crc32c"`
        Etag                    string `json:"etag"`
        Generation              string `json:"generation"`
        ID                      string `json:"id"`
        Kind                    string `json:"kind"`
        Md5Hash                 string `json:"md5Hash"`
        MediaLink               string `json:"mediaLink"`
        Metageneration          string `json:"metageneration"`
        SelfLink                string `json:"selfLink"`
        Size                    string `json:"size"`
        StorageClass            string `json:"storageClass"`
        TimeCreated             string `json:"timeCreated"`
        TimeStorageClassUpdated string `json:"timeStorageClassUpdated"`
        Updated                 string `json:"updated"`
}
```
- Supported trigger types:
 * google.storage.object.finalize
 * google.storage.object.delete
 * google.storage.object.archive
 * google.storage.object.metadataUpdate
- _References_
  * [Calling Storage](https://cloud.google.com/functions/docs/calling/storage)



#### Pubsub

- [PubSubProxyFn](msg/proxy_msg.go)
- [E2E Use Case](e2e/regression/cases/004_proxy_message)
- _Reference_: 
    * [Calling PubSub](https://cloud.google.com/functions/docs/calling/pubsub)
- Pubsub Event:

```go
type PubSubMessage struct {
    Data []byte json:"data"
}
```



#### Cloud Firestore

- [ModeratePostsFn](moderate/moderate.go)
- [E2E Use Case](e2e/regression/cases/005_moderate_posts)
- Firestore Event

```go
type FirestoreEvent struct {
	Value      *EventValue `json:"value"`
	OldValue   *EventValue `json:"oldValue"`
	UpdateMask *UpdateMask `json:"updateMask"`
}

type EventValue struct {
	CreateTime time.Time              `json:"createTime"`
	UpdateTime time.Time              `json:"updateTime"`
	Fields     map[string]interface{} `json:"fields"`
	Name       string                 `json:"name"`
}

type UpdateMask struct {
	FieldPaths []string `json:"fieldPaths"`
}

```
- _Reference_: 
    * [Calling Firestore](https://cloud.google.com/functions/docs/calling/cloud-firestore)


Event example:
```json
{
  "value": {},
  "oldValue": {
    "createTime": "2019-01-01T22:46:23.742719Z",
    "fields": {
      "id": {
        "integerValue": "2"
      },
      "text": {
        "stringValue": "some  text"
      }
    },
    "name": "projects/abstractdb-154a9/databases/(default)/documents/posts/2",
    "updateTime": "2019-01-01T23:27:19.456879Z"
  },
  "updateMask": {},
  "data": null
}
```



#### Realtime Database

- [CountLikesFn](likes/count_likes.go)
- [E2E Use Case](e2e/regression/cases/006_count_likes)
- Firebase Event
```go
type FirebaseEvent struct{
	Delta []byte json:"delta"
}
```
- _Reference_: 
    * [Calling Realtime Database](https://cloud.google.com/functions/docs/calling/realtime-database)



#### Cloud function meta data.

All events supporting context provide a way to retrieve cloud function meta data. 

```go
meta, err := metadata.FromContext(ctx)
```


```go
// Metadata holds Google Cloud Functions metadata.
type Metadata struct {
	// EventID is a unique ID for the event. For example: "70172329041928".
	EventID string `json:"eventId"`
	// Timestamp is the date/time this event was created.
	Timestamp time.Time `json:"timestamp"`
	// EventType is the type of the event. For example: "google.pubsub.topic.publish".
	EventType string `json:"eventType"`
	// Resource is the resource that triggered the event.
	Resource Resource `json:"resource"`
}

// Resource holds Google Cloud Functions resource metadata.
// Resource values are dependent on the event type they're from.
type Resource struct {
	// Service is the service that triggered the event.
	Service string `json:"service"`
	// Name is the name associated with the event.
	Name string `json:"name"`
	// Type is the type of event.
	Type string `json:"type"`
}

```

