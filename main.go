/**
 * Copyright (c) 2016 Intel Corporation
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
	"math/rand"

	"os"
	"time"

	"github.com/gocraft/web"

	httpGoCommon "github.com/trustedanalytics/tapng-go-common/http"
	"github.com/trustedanalytics/tapng-go-common/logger"
	"github.com/trustedanalytics/tapng-template-repository/api"
	"github.com/trustedanalytics/tapng-template-repository/catalog"
)

type appHandler func(web.ResponseWriter, *web.Request) error

var logger = logger_wrapper.InitLogger("main")

func main() {
	rand.Seed(time.Now().UnixNano())
	catalog.LoadAvailableTemplates()

	r := web.New(api.Context{})
	r.Middleware(web.LoggerMiddleware)

	basicAuthRouter := r.Subrouter(api.Context{}, "/api/v1")
	route(basicAuthRouter)
	v1AliasRouter := r.Subrouter(api.Context{}, "/api/v1.0")
	route(v1AliasRouter)

	if os.Getenv("CONSOLE_SERVICE_SSL_CERT_FILE_LOCATION") != "" {
		httpGoCommon.StartServerTLS(os.Getenv("TEMPLATE_REPOSITORY_SSL_CERT_FILE_LOCATION"),
			os.Getenv("TEMPLATE_REPOSITORY_SSL_KEY_FILE_LOCATION"), r)
	} else {
		httpGoCommon.StartServer(r)
	}

}

func route(router *web.Router) {
	router.Middleware((*api.Context).BasicAuthorizeMiddleware)
	router.Get("/templates", (*api.Context).Templates)
	router.Get("/templates/:templateId", (*api.Context).GetCustomTemplate)
	router.Get("/parsed_template/:templateId/", (*api.Context).GenerateParsedTemplate)

	//TODO: change to jwtRouter after UAA integration
	router.Post("/templates", (*api.Context).CreateCustomTemplate)
	router.Delete("/templates/:templateId", (*api.Context).DeleteCustomTemplate)
}
