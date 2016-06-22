package model

import (
	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/apis/extensions"
)

type JobType string

const (
	JobTypeOnCreateInstance JobType = "onCreateInstance"
	JobTypeOnDeleteInstance JobType = "onDeleteInstance"
	JobTypeOnBindInstance   JobType = "onBindInstance"
	JobTypeOnUnbindInstance JobType = "onUnbindInstance"
)

type Template struct {
	Id    string              `json:"id"`
	Body  KubernetesComponent `json:"body"`
	Hooks []*JobHook          `json:"hooks"`
}

type TemplateMetadata struct {
	Id                  string `json:"id"`
	TemplateDirName     string `json:"templateDirName"`
	TemplatePlanDirName string `json:"templatePlanDirName"`
}

type JobHook struct {
	Type JobType        `json:"type"`
	Job  extensions.Job `json:"job"`
}

type KubernetesBlueprint struct {
	Id                    int
	SecretsJson           []string
	DeploymentJson        []string
	ServiceJson           []string
	ServiceAcccountJson   []string
	PersistentVolumeClaim []string
	CredentialsMapping    string
	ReplicaTemplate       string
	UriTemplate           string
}

type KubernetesComponent struct {
	PersistentVolumeClaims []*api.PersistentVolumeClaim `json:"persistentVolumeClaims"`
	Deployments            []*extensions.Deployment     `json:"deployments"`
	Services               []*api.Service               `json:"services"`
	ServiceAccounts        []*api.ServiceAccount        `json:"serviceAccounts"`
	Secrets                []*api.Secret                `json:"secrets"`
}

type ServicesMetadata struct {
	Services []ServiceMetadata `json:"services"`
}

type ServiceMetadata struct {
	Id          string         `json:"id"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Bindable    bool           `json:"bindable"`
	Tags        []string       `json:"tags"`
	Plans       []PlanMetadata `json:"plans"`
	InternalId  string         `json:"-"`
}

type PlanMetadata struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Free        bool   `json:"free"`
	InternalId  string `json:"-"`
}
