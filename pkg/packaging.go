package pkg

import (
	apiv1alpha1 "github.com/frantjc/kontrol/api/v1alpha1"
	rbacv1 "k8s.io/api/rbac/v1"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
)

type Packaging struct {
	PolicyRules               []rbacv1.PolicyRule                        `json:"rules,omitempty"`
	CustomResourceDefinitions []apiextensionsv1.CustomResourceDefinition `json:"crds,omitempty"`
	Overrides                 apiv1alpha1.Overrides                      `json:",inline"`
}
