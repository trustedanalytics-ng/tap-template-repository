# TAP Template Repository

Template repository is a microservice developed to be part of TAP platform. 
It is used to store json files used by container-broker and to parse provided templates.

It has responsibilities for:

* storing kubernetes objects as templates in folder that can be later used
* providing information about stored templates
* parsing templates with provided query parameters (template variables are visible as `$foo`)
 

### REQUIREMENTS

### Binary
There is no requirements for binary app.

### Compilation
* git (for pulling repository only)
* go >= 1.6
* clone this repo
* in directory of just cloned repository invoke: `make build_anywhere`
* binaries are available in ./application/

### USAGE

To build and run project:

```
  git clone https://github.com/intel-data/tap-template-repository
  cd tap-template-repository
  make build_anywhere
  TEMPLATE_REPOSITORY_USER=admin TEMPLATE_REPOSITORY_PASS=password PORT=8082 ./application/tap-template-repository
```

Template repository provides few endpoints

#### Creating new template 

```
  curl -v -XPOST -H 'Content-type: application/json' admin:password@localhost:8082/api/v1/templates -d "@examples/add_template_with_body.json"
```

It is expected that your template will be added to catalogData/custom folder
You can also validate through endpoint that template has been added.
To display all templates:

```
  curl -v  admin:password@localhost:8082/api/v1/templates
```

#### Displaying created template

To display just one provide id:

```
  curl -v  admin:password@localhost:8082/api/v1/templates/test1
```

#### Parsing created template

You can also validate that parsing template will work if you provide query parameters.
Each `$foo` will be replaced with `bar` if query param will have format `/parsed_template/:templateId?foo=bar`
Template to be parsed requires instanceId in query param to be valid UUID

```
  curl -v  admin:password@localhost:8082/api/v1/parsed_template/test1?idx_and_short_instance_id=NewValue&instanceId=27523a96-63a1-11e6-bc3a-00155d3d8807
```

As an output you should see something similar to:

```
{"id":"test1","body":{"componentType":"","persistentVolumeClaims":null,"deployments":[{"kind":"Deployment","apiVersion":"extensions/v1beta1","metadata":{"name":"x27523a9663a11","creationTimestamp":null,"labels":{"managed_by":"TAP","service_id":"27523a96-63a1-11e6-bc3a-00155d3d8807"}},"spec":{"replicas":1,"selector":{"matchLabels":{"managed_by":"TAP","service_id":"27523a96-63a1-11e6-bc3a-00155d3d8807"}},"template":{"metadata":{"creationTimestamp":null,"labels":{"managed_by":"TAP","service_id":"27523a96-63a1-11e6-bc3a-00155d3d8807"}},"spec":{"volumes":null,"containers":[{"name":"k-mongodb30","image":"frodenas/mongodb:3.0","command":["/scripts/run.sh","--smallfiles","--httpinterface"],"ports":[{"containerPort":27017,"protocol":"TCP"}],"resources":{},"imagePullPolicy":"IfNotPresent"}],"restartPolicy":"Always","dnsPolicy":"ClusterFirst","serviceAccountName":""}},"strategy":{}},"status":{}}],"ingresses":null,"services":null,"serviceAccounts":null,"secrets":null},"hooks":null}
```

Please bear in mind that in order to have name compliant with kubernetes UUID will be truncated to proper DNS label

#### Removing template

Finally you can remove template with

```
curl -v -XDELETE admin:password@localhost:8082/api/v1/templates/test1
```

It is expected that template will be removed from folder and not available on GET requests.
