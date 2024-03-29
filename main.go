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
	"time"

	"github.com/gocraft/web"

	httpGoCommon "github.com/trustedanalytics-ng/tap-go-common/http"
	"github.com/trustedanalytics-ng/tap-template-repository/api"
	"github.com/trustedanalytics-ng/tap-template-repository/catalog"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	context := api.Context{
		TemplateApi: &catalog.TemplateApiConnector{},
	}
	context.TemplateApi.GetTemplatesPaths()

	r := web.New(context)
	r.Middleware(web.LoggerMiddleware)

	r.Get("/healthz", context.GetTemplateRepositoryHealth)

	basicAuthRouter := r.Subrouter(context, "/api/v1")
	route(basicAuthRouter, &context)
	v1AliasRouter := r.Subrouter(context, "/api/v1.0")
	route(v1AliasRouter, &context)

	httpGoCommon.StartServer(r)
}

func route(router *web.Router, context *api.Context) {
	router.Middleware((*context).BasicAuthorizeMiddleware)
	router.Get("/templates", (*context).Templates)
	router.Get("/templates/:templateId", (*context).GetRawTemplate)
	router.Get("/parsed_template/:templateId", (*context).GenerateParsedTemplate)

	router.Post("/templates", (*context).CreateCustomTemplate)
	router.Delete("/templates/:templateId", (*context).DeleteCustomTemplate)
}
