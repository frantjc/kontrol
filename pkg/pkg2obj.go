package pkg

import (
	"github.com/frantjc/go-fn"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func PackagingToObjects(ref, name, namespace string, packaging *Packaging) []client.Object {
	var (
		objs []client.Object = fn.Map(packaging.CustomResourceDefinitions, func(crd apiextensionsv1.CustomResourceDefinition, _ int) client.Object {
			return &crd
		})
		metadata = metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		}
		selectorLabels = map[string]string{
			"app.kubernetes.io/name": name,
		}
		policyRules = append(
			packaging.PolicyRules,
			fn.Map(packaging.CustomResourceDefinitions, func(crd apiextensionsv1.CustomResourceDefinition, _ int) rbacv1.PolicyRule {
				return rbacv1.PolicyRule{
					Verbs:     []string{"*"},
					APIGroups: []string{crd.Spec.Group},
					Resources: []string{crd.Spec.Names.Kind},
				}
			})...,
		)
		role client.Object = &rbacv1.Role{
			TypeMeta: metav1.TypeMeta{
				APIVersion: rbacv1.SchemeGroupVersion.String(),
				Kind:       "Role",
			},
			ObjectMeta: metadata,
			Rules:      policyRules,
		}
	)

	if fn.Some(packaging.CustomResourceDefinitions, func(crd apiextensionsv1.CustomResourceDefinition, _ int) bool {
		return crd.Spec.Scope != "Namespaced"
	}) {
		role = &rbacv1.ClusterRole{
			TypeMeta: metav1.TypeMeta{
				APIVersion: rbacv1.SchemeGroupVersion.String(),
				Kind:       "ClusterRole",
			},
			ObjectMeta: metadata,
			Rules:      policyRules,
		}
	}

	objs = append(
		objs,
		role,
		&corev1.ServiceAccount{
			TypeMeta: metav1.TypeMeta{
				APIVersion: "v1",
				Kind:       "ServiceAccount",
			},
			ObjectMeta: metadata,
		},
		&rbacv1.ClusterRoleBinding{
			TypeMeta: metav1.TypeMeta{
				APIVersion: rbacv1.SchemeGroupVersion.String(),
				Kind:       "ClusterRoleBinding",
			},
			ObjectMeta: metadata,
			RoleRef: rbacv1.RoleRef{
				Kind: "ClusterRole",
				Name: name,
			},
			Subjects: []rbacv1.Subject{
				{
					Kind:      rbacv1.ServiceAccountKind,
					Name:      name,
					Namespace: namespace,
				},
			},
		},
		&appsv1.Deployment{
			TypeMeta: metav1.TypeMeta{
				APIVersion: appsv1.SchemeGroupVersion.String(),
				Kind:       "Deployment",
			},
			ObjectMeta: metadata,
			Spec: appsv1.DeploymentSpec{
				Replicas: fn.Ptr[int32](1),
				Selector: &metav1.LabelSelector{
					MatchLabels: selectorLabels,
				},
				Template: corev1.PodTemplateSpec{
					ObjectMeta: metav1.ObjectMeta{
						Labels: selectorLabels,
					},
					Spec: corev1.PodSpec{
						Containers: []corev1.Container{
							{
								Name:  name,
								Image: ref,
								Ports: fn.Map(packaging.Overrides.ContainerPorts, func(port int32, _ int) corev1.ContainerPort {
									return corev1.ContainerPort{
										ContainerPort: port,
									}
								}),
							},
						},
					},
				},
			},
		},
		&corev1.Service{
			TypeMeta: metav1.TypeMeta{
				APIVersion: corev1.SchemeGroupVersion.String(),
				Kind:       "Service",
			},
			ObjectMeta: metadata,
			Spec: corev1.ServiceSpec{
				Ports: fn.Map(packaging.Overrides.ContainerPorts, func(port int32, _ int) corev1.ServicePort {
					return corev1.ServicePort{
						Port:       port,
						TargetPort: intstr.FromInt(int(port)),
					}
				}),
				Selector: selectorLabels,
			},
		},
	)

	return objs
}
