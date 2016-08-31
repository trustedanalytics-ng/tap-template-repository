package model

import (
	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/apis/extensions"
)

type HookType string

const (
	HookTypeDeployment  HookType = "deployment"
	HookTypeProvision   HookType = "provision"
	HookTypeDeprovision HookType = "deprovision"
	HookTypeBind        HookType = "bind"
	HookTypeUnbind      HookType = "unbind"
	HookTypeRemoval     HookType = "removal"
)

type Template struct {
	Id    string                `json:"id"`
	Body  KubernetesComponent   `json:"body"`
	Hooks map[HookType]*api.Pod `json:"hooks"`
}

type TemplateMetadata struct {
	Id                  string `json:"id"`
	TemplateDirName     string `json:"templateDirName"`
	TemplatePlanDirName string `json:"templatePlanDirName"`
}

type KubernetesBlueprint struct {
	Id                    int
	SecretsJson           []string
	DeploymentJson        []string
	IngressJson           []string
	ServiceJson           []string
	ServiceAcccountJson   []string
	PersistentVolumeClaim []string
	CredentialsMapping    string
	ReplicaTemplate       string
	UriTemplate           string
	Component             string
	Hooks                 string
}

type ComponentType string

const (
	ComponentTypeBroker   ComponentType = "broker"
	ComponentTypeInstance ComponentType = "instance"
	ComponentTypeBoth     ComponentType = "both"
)

type KubernetesComponent struct {
	Type                   ComponentType                `json:"componentType"`
	PersistentVolumeClaims []*api.PersistentVolumeClaim `json:"persistentVolumeClaims"`
	Deployments            []*extensions.Deployment     `json:"deployments"`
	Ingresses              []*extensions.Ingress        `json:"ingresses"`
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
