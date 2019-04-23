# Serverless e2e testing examples

This project contains example applications using serverless technology
You can run all e2e test automation runner. 


- [Lambda](lambda)
- [Cloud function](cloud_function)
- [Dataflow](dataflow)
- [Cloudrun](cloudrun)



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
  credentials: gcp-e2e
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
        credentials: "gcp-e2e",
        entryPoint: "HelloWorldFn",
        runtime: "go111",
        source: {
           URL: '/somepath//hello/'
        }
   })
    ```
6. API contracts details
```bash
     endly -s='gcp/cloudfunctions'
     endly -s='gcp/cloudfunctions:deploy'
     
     endly -s='aws/lambda'
     endly -s='aws/lambda:deploy'

     endly -s='*'
      
```   
7. [Find out more](https://github.com/adrianwit/endly-introduction) about E2E endly workflows.


