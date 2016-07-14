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
	"net/http"
	"os"
	"time"

	"github.com/gocraft/web"

	"github.com/trustedanalytics/tapng-go-common/logger"
	"github.com/trustedanalytics/tapng-template-repository/api"
	"github.com/trustedanalytics/tapng-template-repository/catalog"
)

type appHandler func(web.ResponseWriter, *web.Request) error

var logger = logger_wrapper.InitLogger("main")

func main() {
	rand.Seed(time.Now().UnixNano())
	context := api.Context{
		Template: &catalog.Template{},
	}
	context.Template.LoadAvailableTemplates()

	r := web.New(context)
	r.Middleware(web.LoggerMiddleware)

	basicAuthRouter := r.Subrouter(context, "/api/v1")
	basicAuthRouter.Middleware(context.BasicAuthorizeMiddleware)

	jwtRouter := r.Subrouter(context, "/api/v1")
	jwtRouter.Middleware(context.JWTAuthorizeMiddleware)

	basicAuthRouter.Get("/templates", context.Templates)
	basicAuthRouter.Get("/templates/:templateId", context.GetCustomTemplate)
	basicAuthRouter.Get("/parsed_template/:templateId/", context.GenerateParsedTemplate)

	//TODO: change to jwtRouter after UAA integration
	basicAuthRouter.Post("/templates", (*api.Context).CreateCustomTemplate)
	basicAuthRouter.Delete("/templates/:templateId", context.DeleteCustomTemplate)

	port := os.Getenv("TEMPLATE_REPOSITORY_PORT")
	logger.Info("Will listen on:", port)

	var err error
	if os.Getenv("TEMPLATE_REPOSITORY_SSL_CERT_FILE_LOCATION") != "" {
		err = http.ListenAndServeTLS(":"+port, os.Getenv("TEMPLATE_REPOSITORY_SSL_CERT_FILE_LOCATION"),
			os.Getenv("TEMPLATE_REPOSITORY_SSL_KEY_FILE_LOCATION"), r)
	} else {
		err = http.ListenAndServe(":"+port, r)
	}

	if err != nil {
		logger.Critical("Couldn't serve app on port:", port, " Error:", err)
	}
}
