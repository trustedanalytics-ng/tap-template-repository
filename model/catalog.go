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

const RAW_TEMPLATE_ID_FIELD = "id"

type RawTemplate map[string]interface{}

type Template struct {
	Id    string                `json:"id"`
	Body  KubernetesComponent   `json:"body"`
	Hooks map[HookType]*api.Pod `json:"hooks"`
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
