package brokerapi

import (
	"github.com/gocraft/web"
	"net/http"
	webutils "github.com/tapng/broker/webutils"
)

type BasicAuthContext struct {
	*Context
}



func SetupApiV1Routes(router web.Router) error {
	router.Get("/catalog", (*BasicAuthContext).Catalog)

	router.Get("/service_definition/:service_id", (*BasicAuthContext).GetServiceDefinitionById)
	router.Get("/service_instances/:service_id", (*BasicAuthContext).GetServiceInstancesById)
	router.Get("/secrets", (*BasicAuthContext).GetSecretByType)
	router.Get("/secrets/:secret_id", (*BasicAuthContext).GetSecretById)


	router.Put("/service_instances", (*BasicAuthContext).CreateService)
	router.Put("/secrets", (*BasicAuthContext).CreateSecret)
	return nil

}

func (c *BasicAuthContext) BasicAuthRequired(rw web.ResponseWriter, r *web.Request, next web.NextMiddlewareFunc) {
	if true {
		//c.User = user
		next(rw, r)
	} else {
		rw.Header().Set("Location", "/")
		rw.WriteHeader(http.StatusMovedPermanently)
	}
}

func (c *BasicAuthContext) Catalog(rw web.ResponseWriter, req *web.Request) {

	webutils.WriteJson(rw, c.ServiceTemplates , http.StatusOK)
}
func (c *BasicAuthContext) GetServiceDefinitionById(rw web.ResponseWriter, req *web.Request) {
	webutils.WriteJson(rw, c.ServiceTemplates , http.StatusOK)
}
func (c *BasicAuthContext) GetServiceInstancesById(rw web.ResponseWriter, req *web.Request) {
	webutils.WriteJson(rw, c.ServiceTemplates , http.StatusOK)
}
func (c *BasicAuthContext) GetSecretByType(rw web.ResponseWriter, req *web.Request) {
	webutils.WriteJson(rw, c.ServiceTemplates , http.StatusOK)
}
func (c *BasicAuthContext) CreateService(rw web.ResponseWriter, req *web.Request) {
	webutils.WriteJson(rw, c.ServiceTemplates , http.StatusOK)
}
func (c *BasicAuthContext) GetSecretById(rw web.ResponseWriter, req *web.Request) {
	webutils.WriteJson(rw, c.ServiceTemplates , http.StatusOK)
}
func (c *BasicAuthContext) CreateSecret(rw web.ResponseWriter, req *web.Request) {
	webutils.WriteJson(rw, c.ServiceTemplates , http.StatusOK)
}
