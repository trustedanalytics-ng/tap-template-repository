package brokerapi

import (
	"github.com/gocraft/web"
	"github.com/tapng/broker/catalog"
	webutils "github.com/tapng/broker/webutils"
	"net/http"
)

type Context struct {
	ServiceTemplates []catalog.ServiceTemplate
	SecretsTemplates []string
}

func SetupMainRoutes(router web.Router) error {
	router.Get("/", (*Context).Index)
	return nil
}

func (c *Context) Index(rw web.ResponseWriter, req *web.Request) {
	webutils.WriteJson(rw, "OK", http.StatusOK)
}

