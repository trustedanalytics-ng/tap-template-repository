/**
 * Copyright (c) 2015 Intel Corporation
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
	"log"
	"math/rand"
	"net/http"
	"time"
	"strings"
	"strconv"
	"regexp"
	"github.com/tapng/broker/brokerapi"
	webutils "github.com/tapng/broker/webutils"
	"github.com/gocraft/web"
	"k8s.io/kubernetes/pkg/api"
	"github.com/tapng/broker/catalog"
	"github.com/satori/go.uuid"
	"encoding/json"
	k8sClient "k8s.io/kubernetes/pkg/client/unversioned"
	"k8s.io/kubernetes/pkg/client/restclient"
	"k8s.io/kubernetes/pkg/labels"
	"k8s.io/kubernetes/pkg/util/sets"
	b64 "encoding/base64"
	"os"
)



type appHandler func(web.ResponseWriter, *web.Request) error


func _main() {
	tapServiceTemplates,_:=catalog.GetTapServices()
	tapSecretTemplates,_:=catalog.GetTapSecrets()

	apiContext := brokerapi.Context{ServiceTemplates: tapServiceTemplates.Services, SecretsTemplates: tapSecretTemplates}

	routeRouter := web.New(apiContext)
	routeRouter.Middleware(web.LoggerMiddleware)

	basicContext :=brokerapi.BasicAuthContext{&apiContext}

	apiRouter := routeRouter.Subrouter(basicContext, "/api/v1")
	apiRouter.Middleware((* brokerapi.BasicAuthContext).BasicAuthRequired)


	port := os.Getenv("PORT")
	log.Println("Will listen on:", "0.0.0.0:"+port)
	err := http.ListenAndServe("0.0.0.0:"+port, routeRouter)
	if err != nil {
		log.Println("Couldn't serve app on port ", port, " Application will be closed now.")
	}
}
type Context struct {
	Namespace string
}

func main() {


	rand.Seed(time.Now().UnixNano())

	ctx := Context{os.Getenv("KUBE_APPS_NS")}


	r := web.New(ctx)
	r.Middleware(web.LoggerMiddleware)


	r.Error((*Context).Error)

	r.Get("/", (*Context).Index)
	r.Get("/catalog", (*Context).Catalog)
	r.Get("/service/:service_id", (*Context).GetService)
	r.Get("/service/:service_id/instances", (*Context).GetRunningServices)
	r.Put("/service", (*Context).ServiceCreate)

	r.Get("/secrets/:secret_type", (*Context).GetSecretsByType)

	port := os.Getenv("PORT")
	log.Println("Will listen on:", "0.0.0.0:"+port)
	err := http.ListenAndServe("0.0.0.0:"+port, r)
	if err != nil {
		log.Println("Couldn't serve app on port ", port, " Application will be closed now.")
	}

}

type MySecret struct{
	Id string `json:"id"`
	Name string `json:"name"`
}
func (c *Context) GetSecretsByType(rw web.ResponseWriter, req *web.Request) {
	secret_type:= req.PathParams["secret_type"]

	client, _ := getK8SClient()

	list_options := api.ListOptions {}
	list_options.LabelSelector = labels.NewSelector()
	managedByReq, err :=  labels.NewRequirement("secretClass", labels.EqualsOperator,  sets.NewString(secret_type))
	if err != nil {
		log.Println(err)
	}
	list_options.LabelSelector.Add(*managedByReq)

	secrets, _ := client.Secrets(c.Namespace).List(list_options)



	b := []MySecret{}
	for _, secr := range secrets.Items {
		if secr.ObjectMeta.Labels["secretClass"] == secret_type{
			id := strings.Replace(secr.ObjectMeta.Name, secr.ObjectMeta.Labels["secretClass"]+".", "", -1)
			b = append(b, MySecret{Id: id, Name: secr.ObjectMeta.Labels["name"]})
		}
	}

	webutils.WriteJson(rw, b, http.StatusOK)
}

func (c *Context) Error(rw web.ResponseWriter, r *web.Request, err interface{}) {
	log.Println("Respond500: reason: error ", err)
	rw.WriteHeader(http.StatusInternalServerError)
}
func (c *Context) Index(rw web.ResponseWriter, req *web.Request) {
	webutils.WriteJson(rw, "I'm OK", http.StatusOK)
}

func (c *Context) Catalog(rw web.ResponseWriter, req *web.Request) {
	catalogServices, _ := catalog.GetTapServices()
	webutils.WriteJson(rw, catalogServices.Services , http.StatusOK)
}

func (c *Context) GetService(rw web.ResponseWriter, req *web.Request) {
	catalogServices, _ := catalog.GetTapServices()
	service_id := req.PathParams["service_id"]
	service, err := catalog.FindServiceByIdInServices(service_id, catalogServices.Services)
	if err != nil {
		log.Println(err)
	}

	webutils.WriteJson(rw, service, http.StatusOK)
}

func (c *Context) GetRunningServices(rw web.ResponseWriter, req *web.Request) {
	service_id := req.PathParams["service_id"]
	client, _ := getK8SClient()

	log.Println("trying to find:", service_id)

	//TODO not working as expected
	list_options := api.ListOptions {}
	list_options.LabelSelector = labels.NewSelector()
	managedByReq, err :=  labels.NewRequirement("catalog_service_id", labels.EqualsOperator,  sets.NewString(service_id))
	if err != nil {
		log.Println(err)
	}
	list_options.LabelSelector.Add(*managedByReq)

	services, _ := client.Services(c.Namespace).List(api.ListOptions {})


	b := api.ServiceList{}
	for _, service := range services.Items {
		if service.ObjectMeta.Labels["catalog_service_id"] == service_id{
			b.Items = append(b.Items, service)
		}
	}

	webutils.WriteJson(rw, b.Items, http.StatusOK)
}

type ServiceInstancesPutRequest struct{
	Service_id string `json:"service_id"`
	Plan_id string `json:"plan_id"`
	Secret_id []string `json:"secret_id"`
	Cluster_ip string `json:cluster_ip`
	Name string `json:name`
}


func (c *Context) ServiceCreate(rw web.ResponseWriter, req *web.Request) {

	log.Println(req)
	req_json := ServiceInstancesPutRequest{}
	err := webutils.ReadJson(req, &req_json)
	if err != nil {
		log.Println("err when put", err)
	} else {
		log.Println(req_json)
	}


	service_id := req_json.Service_id
	plan_id := req_json.Plan_id
	secret_id := req_json.Secret_id
	secret_provided := true;

	catalogServices, _ := catalog.GetTapServices()
	secrets, _ := catalog.GetTapSecrets()

	log.Println("looking for", service_id)
	service, err := catalog.FindServiceByIdInServices(service_id, catalogServices.Services)
	if err != nil {
		log.Println(err)
	}
	log.Println("looking for plan", plan_id)
	plan, err := catalog.FindPlanByIdInService(plan_id, service)
	if err != nil {
		log.Println(err)
	}
	client, _ := getK8SClient()

	u1 := uuid.NewV4()
	if len(secret_id)==0 || secret_id[0] =="" {

		secret_id = append(secret_id,"x"+strings.Replace(u1.String()[0:15], "-", "", -1))
		secret_provided = false
	}

	req_json.Secret_id = secret_id
	req_json.Plan_id = plan_id
	req_json.Service_id = service_id

	//log.Println("catalogServices", catalogServices)
	log.Println("service", service)
	log.Println("plan", plan)

	k8sService := api.Service{}
	unmarshalAdjustedJson(plan.K8SServicesTemplates[0],req_json, &k8sService, u1)
	serv, err := client.Services(c.Namespace).Create(&k8sService)
	if err != nil {
		log.Println(err)
	} else {
		log.Println("service created")
	}

	req_json.Cluster_ip = serv.Spec.ClusterIP


	k8sSecrets := []api.Secret{}
	for _, strSecret := range secrets {
		apiSecret := api.Secret{}
		unmarshalAdjustedJson(strSecret, req_json, &apiSecret, u1)
		k8sSecrets = append(k8sSecrets, apiSecret)
	}




	k8sReplicationController := api.ReplicationController{}
	log.Println(plan.ReplicationControllersTemplates[0])
	unmarshalAdjustedJson(plan.ReplicationControllersTemplates[0],req_json, &k8sReplicationController, u1)

	//find proper secret

	//requiredSecrets := k8sReplicationController.Spec.Template.ObjectMeta.Labels["secretsRequired"]
	//createsSecrets := k8sReplicationController.Spec.Template.ObjectMeta.Labels["createsSecrets"]
	createsSecrets := "-"
	if len(plan.CreatesSecrets) == 1{
		createsSecrets = plan.CreatesSecrets[0]
	}


	log.Println("Creates secret: ", createsSecrets)

	for _, secret := range k8sSecrets {
		log.Println("comparing", createsSecrets, secret.ObjectMeta.Labels["secretClass"])
		if secret.ObjectMeta.Labels["secretClass"] == createsSecrets {
			if secret_provided == true {
				log.Println("secret provided")

			} else {
				log.Println("creating new secret")

				_, err = client.Secrets(c.Namespace).Create(&secret)
				if err != nil {
					log.Println(err)
				} else {
					log.Println("secret created")
					break
				}
			}
		}
	}


	//jsss, _ := json.Marshal(k8sReplicationController)
	//log.Println(string(jsss))
	_, err = client.ReplicationControllers(c.Namespace).Create(&k8sReplicationController)
	if err != nil {
		log.Println(err)
	}
	webutils.WriteJson(rw, "I'm OK" , http.StatusOK)

}

func unmarshalAdjustedJson(jsonString string, params ServiceInstancesPutRequest, structType interface{}, uid uuid.UUID) {

	jsonparsed := adjust_params(jsonString,params, uid)
	err := json.Unmarshal([]byte(jsonparsed), structType)
	if err != nil {
		log.Println("Couldn't unmarshal json ", err)
	}
}


func getK8SClient() (*k8sClient.Client, error){
	tls := restclient.TLSClientConfig{}
	tls.CertData = []byte(os.Getenv("KUBERNETES_CERT_PEM_STRING"))
	tls.KeyData = []byte(os.Getenv("KUBERNETES_KEY_PEM_STRING"))
	tls.CAData = []byte(os.Getenv("KUBERNETES_CA_PEM_STRING"))

	config := &restclient.Config{
		Host:   os.Getenv("KUBERNETES_URL"),
		TLSClientConfig:tls,
	}
	return k8sClient.New(config)
}


func adjust_params(content string,params ServiceInstancesPutRequest, uid uuid.UUID) string {
	uidstring := "x"+strings.Replace(uid.String()[0:15], "-", "", -1)
	f := content
	f = strings.Replace(f, "$org", "default", -1)
	f = strings.Replace(f, "$space", "default", -1)
	f = strings.Replace(f, "$catalog_service_id", params.Service_id, -1)
	f = strings.Replace(f, "$catalog_plan_id", params.Plan_id, -1)
	f = strings.Replace(f, "$name", params.Name, -1)
	f = strings.Replace(f, "$service_id", uidstring+"-"+params.Name, -1)

	f = strings.Replace(f, "$clusterip", params.Cluster_ip, -1)



	f = strings.Replace(f, "$idx_and_short_serviceid", uidstring, -1)
	f = strings.Replace(f, "$secret_idx_and_short_serviceid", params.Secret_id[0], -1)


	for i := 0; i < 9; i++ {
		rnd:=get_random_string(10)
		f = strings.Replace(f, "$random"+strconv.Itoa(i), rnd, -1)
	}

	rp := regexp.MustCompile(`\$base64\-(.*)\"`)
	fs := rp.FindAllString(f,-1)
	log.Println("znalazlem w ",f,fs) //[$base64-LEP2cvKxmP" $base64-p1a0Z4w0R7" $base64-7BObjkcWZG"]
	for _, sub := range fs{
		sub = strings.Replace(sub, "$base64-","",-1)
		sub = strings.Replace(sub, "\"","",-1)
		f = strings.Replace(f, "$base64-"+sub, b64.StdEncoding.EncodeToString([]byte(sub)), -1)

	}

	return f
}

func get_random_string(n int) string {
	possible_rand_chars := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")
	b := make([]rune, n)
	for i := range b {
		b[i] = possible_rand_chars[rand.Intn(len(possible_rand_chars))]
	}
	return string(b)
}
