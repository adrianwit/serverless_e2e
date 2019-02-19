## Cloud functions with e2e testing

**Prerequisites:**

 - install gcloud
 - go1.11
 - service account for e2e testing dedicated project with enabled API
    * regular secrets file -> ~/.secret/e2e.json from http://console.cloud.google.com/
    * firebase secrets file -> ~/.secret/fbc.json from http://http://console.firebase.google.com/
    * [Setup endly google secrets](https://github.com/viant/endly/tree/master/doc/secrets#gc)
    * [endly e2e runner](https://github.com/viant/endly/releases) (0.30+)

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


Cloud Function is a serverless compute service known also as Function-as-a-Service (FaaS) that runs a provided code in response to events 
and automatically manages the underlying compute resources.

The event can be fired by a specific trigger, which determines how and when your function executes. 
Google cloud platform supports two following function types:
 
- HTTP Functions
    ```func(http.ResponseWriter, *http.Request)```

- Background functions:
    ```func(context.Context, Event) (error)```


**Advantages:**
- infrastructure management is handled by cloud provider
- lightweight and modular
- event driven
- scales up automatically
- fault tolerant

**Challenges:**
- debugging
- transparency
- cold start/caching
- cost control
- end to end testing


### E2E automation workflow:

This project uses [endly](http://github.com/viant/endly/) as e2e test runner.

* Basic workflow demystified

[@deploy.yaml](e2e/regression/cases/001_hello_world/deploy.yaml)
```yaml
init:
   srcCodeLocation: ${appPath}/hello/
   target:
      URL: ssh://127.0.0.1
   credentials: dev   
defaults:
  credentials: gcp
pipeline:
  
  package:
    action: exec:run
    comments: vendor build for deployment speedup using SSH service
    target: $target
    commands:
      - cd ${appPath}/split
      - go mod vendor

  deploy:
    action: gcp/cloudfunctions:deploy
    '@name': HelloWorld
    entryPoint: HelloWorldFn
    runtime: go111
    source:
      URL: $srcCodeLocation
```

1. 'init' node: variable initialization
2. 'pipeline' node: tasks definition
3. package/deploy is arbitrary task name.
4. defaults node appends specified attributes to all action nodes (no override).
5. node with 'action' attribute triggers API call to specified service:method
   * package
   ```javascript
   sshExecutorService.run({
       target: {
           URL:"ssh://127.0.0.1",
           credentials: "dev"
       },
       commands: [
           'cd ${appPath}/split',
           'go mod vendor' 
       ]
   })
   ```
   * deploy
   ```javascript
   gcpCloudFunctionsService.deploy({
        name: "HelloWorld",
        entryPoint: "HelloWorldFn",
        runtime: "go111",
        source: {
           URL: '/somepath//hello/'
        }
   })
    ```
6. API contracts details
```bash
    endly -s='gcp/cloudfunctions' -a='deploy'
    endly -s='gcp/cloudfunctions' 
    endly -s='*' 
```   
7. [Find out more](https://github.com/adrianwit/endly-introduction) about E2E endly workflows.




### HTTP Functions

1. Deployment: 
    ```bash
    gcloud functions deploy MyFunction --entry-point MyFunctionFn \
           --runtime go111 --trigger-http
    ```
2. _Reference_: 
    * [Calling HTTP](https://cloud.google.com/functions/docs/calling/http)
 

#### HTTP Examples:
1. Hello World
    - [Source code](hello/hello.go)
    - [E2E Use Case](e2e/regression/cases/001_hello_world)
      * ```endly -i=Test_cases001_hello_world```
2. BigQuery reader
    - [Source code](query/query_reader.go)
    - [E2E Use Case](e2e/regression/cases/002_query_data)
      * ```endly -i=Test_cases002_query_data```



### Background Functions
1. Deployment: 
```bash
    gcloud functions deploy MyFunction --entry-point MyFunctionFn \
            --trigger-resource $Resource \
            --trigger-event $EventType  \
            --runtime go111

```
2. Resource meta:
    * Retrieve cloud function event resource info. 
    ```go
        meta, err := metadata.FromContext(ctx)
    ```
    * Metadata contract
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

### Cloud function runtime environment

The following variables are automatically set by the Cloud Functions runtime.
- GCLOUD_PROJECT 
- FUNCTION_NAME
- FUNCTION_REGION

```go
    projectID    = os.Getenv("GCLOUD_PROJECT")
    functionName = os.Getenv("FUNCTION_NAME")
    region       = os.Getenv("FUNCTION_REGION")
```
        


#### BigQuery 
1. Deployment: 
    ```bash
    gcloud functions deploy MyFunction --entry-point MyFunctionFn \
           --trigger-resource projects/_/jobs/{jobId}  \
           --trigger-event google.cloud.bigquery.job.complete \
           --runtime go111
    ```
2. Supported trigger types:
    * google.cloud.bigquery.job.create
    * google.cloud.bigquery.job.complete
    
3. Resource: projects/_/jobs/{jobId} 

##### BigQueryExample:
    
1. BigQuery Event dispatcher
    - [Source code](dispatcher/bq.go)
    - [E2E Use Case](e2e/regression/cases/003_bqdispatcher)
      * ```endly -i=Test_cases003_bqdispatcher```

#### Google Storage
1. Deployment: 
    ```bash
    gcloud functions deploy SplitWorkload --entry-point SplitWorkloadFn \
       --trigger-resource $myBucket \
       --trigger-event google.storage.object.finalize \
       --runtime go111
    ```
2. Supported trigger types:
     * google.storage.object.finalize
     * google.storage.object.delete
     * google.storage.object.archive
     * google.storage.object.metadataUpdate
3. Resource:  projects/_/buckets/BUCKET
4. _Google Storage_ Event:
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
5. _References_
  * [Calling Storage](https://cloud.google.com/functions/docs/calling/storage)


##### Storage Examples:

1. Storage Event dispatcher
    - [Source code](dispatcher/gs.go)
    - [E2E Use Case](e2e/regression/cases/004_gsdispatcher)
      * ```endly -i=Test_cases004_gsdispatcher```
2. SplitWorkload
    - [Source code](split/split_workload.go)
    - [E2E Use Case](e2e/regression/cases/005_split_workload)
      * ```endly -i=Test_cases005_split_workload```

####  Pubsub
1. Deployment: 
    ```bash
    gcloud functions deploy MyFunction --entry-point MyFunctionFn \
       --trigger-resource myTopic \
       --trigger-event google.pubsub.topic.publish \
       --runtime go111
    ```
2. Supported trigger types:
   * google.pubsub.topic.publish
3. Resource: projects/_/topics/TOPIC
4. Pubsub Event:
    ```go
    type PubSubMessage struct {
        Data []byte json:"data"
    }
    ```
5. _Reference_: 
    * [Calling PubSub](https://cloud.google.com/functions/docs/calling/pubsub)


##### Pubsub Examples:
1. Pubsub proxy
    - [Source code](msg/proxy_msg.go)
    - [E2E Use Case](e2e/regression/cases/006_proxy_message)
      * ```endly -i=Test_cases006_proxy_message```


#### Realtime Database
1. Deployment: 
    ```bash
    gcloud functions deploy MyFunction --entry-point MyFunctionFn  \
       --trigger-event providers/cloud.firestore/eventTypes/document.write  \
       --trigger-resource  'projects/${gcp.projectID}/databases/(default)/documents/posts/{doc}' \
       --runtime go111
    ```
2. Supported trigger types:
   * providers/google.firebase.database/eventTypes/ref.create
   * providers/google.firebase.database/eventTypes/ref.update  
   * providers/google.firebase.database/eventTypes/ref.delete
   * providers/google.firebase.database/eventTypes/ref.write
3. Resource: projects/_/topics/TOPIC
4. Firebase Event
    ```go
    type FirebaseEvent struct{
        Delta []byte json:"delta"
    }
    ```
5. _Reference_: 
    * [Calling Realtime Database](https://cloud.google.com/functions/docs/calling/realtime-database)



##### Realtime Database Examples

1. Counting likes
    - [Source code](likes/count_likes.go)
    - [E2E Use Case](e2e/regression/cases/007_count_likes)
      * ```endly -i=Test_cases007_count_likes```

####  Cloud Firestore
1. Deployment: 
    ```bash
    gcloud functions deploy MyFunction --entry-point MyFunctionFn \
       --trigger-event providers/cloud.firestore/eventTypes/document.write  \
       --trigger-resource  'projects/${gcp.projectID}/databases/(default)/documents/posts/{doc}' \
       --runtime go111
    ```
2. Supported trigger types:
   * providers/cloud.firestore/eventTypes/document.create
   * providers/cloud.firestore/eventTypes/document.delete
   * providers/cloud.firestore/eventTypes/document.update  
   * providers/cloud.firestore/eventTypes/document.write
3. Resource: projects/_/topics/TOPIC
4. Firestore Event
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
    * Example Event
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
        "name": "projects/_/databases/(default)/documents/posts/2",
        "updateTime": "2019-01-01T23:27:19.456879Z"
      },
      "updateMask": {},
      "data": null
    }
    ```
- _Reference_: 
    * [Calling Firestore](https://cloud.google.com/functions/docs/calling/cloud-firestore)


##### Cloud Firestore Examples:
1. Moderate posts:
    - [Source code](moderate/moderate.go)
    - [E2E Use Case](e2e/regression/cases/008_moderate_posts)
      * ```endly -i=Test_cases008_moderate_posts```


####  Google Compute Engine
1. Deployment: 
    ```bash
    gcloud functions deploy MyFunction --entry-point MyFunctionFn \
       --trigger-resource projects/PROJECT/zones/ZONE/instances/{id}  \
       --trigger-event google.compute.instance.stop \
       --runtime go111
    ```
2. Instances 
    -  Supported trigger types:
        * google.compute.instance.create
        * google.compute.instance.delete
        * google.compute.instance.start
        * google.compute.instance.stop
        * google.compute.instance.reset
    - Resource: 
        * projects/{project}/zones/{zone}/instances/{id}
3. instanceGroupManagers
    -  Supported trigger types:
        * google.compute.instancegroupmanager.create
        * google.compute.instancegroupmanager.update
        * google.compute.instancegroupmanager.resize
        * google.compute.instancegroupmanager.delete
    - Resource
        * projects/{project}/{locationType}/{location}/instanceGroupManagers/{manager} 
4. Firewall     
    -  Supported trigger types:
        * google.compute.firewall.create
        * google.compute.firewall.update
        * google.compute.firewall.delete
5. Autoscaler    
    -  Supported trigger types:
        * google.compute.autoscaler.create
        * google.compute.autoscaler.update
        * google.compute.autoscaler.delete
    
### External projects with e2e testing using cloud functions
- [BigQuery Windowed Tables](http://github.com/viant/bqwt)




