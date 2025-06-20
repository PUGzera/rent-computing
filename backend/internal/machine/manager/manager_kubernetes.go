package machine_manager

import (
	"context"
	"fmt"
	machine "rent-computing/internal/machine/data"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	networkv1 "k8s.io/api/networking/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type MachineManagerK8s struct {
	clientSet *kubernetes.Clientset
	namespace string
	route     string
}

type OptionsK8s struct {
	Config    *rest.Config
	Namespace string
	Route     string
}

func NewK8s(options OptionsK8s) (*MachineManagerK8s, error) {
	clientSet, err := kubernetes.NewForConfig(options.Config)
	if err != nil {
		return nil, err
	}
	return &MachineManagerK8s{
		route:     fmt.Sprintf("%s/connect", options.Route),
		clientSet: clientSet,
		namespace: options.Namespace,
	}, nil
}

func (m *MachineManagerK8s) StartShell(ctx context.Context, machine machine.Machine) (string, error) {
	return "", nil
}

func (m *MachineManagerK8s) StartVNC(ctx context.Context, machine machine.Machine) (string, error) {
	err := m.createVNC(ctx, machine)
	if err != nil {
		return "", err
	}

	err = m.waitForVNC(ctx, machine, time.Minute*5)
	if err != nil {
		return "", err
	}

	return m.getAddress(machine), nil
}

func (m *MachineManagerK8s) Stop(ctx context.Context, machine machine.Machine) error {
	err := m.clientSet.
		AppsV1().
		StatefulSets(m.namespace).
		Delete(ctx, machine.Id, v1.DeleteOptions{})
	if err != nil {
		return err
	}

	err = m.clientSet.
		CoreV1().
		Services(m.namespace).
		Delete(ctx, fmt.Sprintf("svc-%s", machine.Id), v1.DeleteOptions{})
	if err != nil {
		return err
	}

	err = m.clientSet.
		NetworkingV1().
		Ingresses(m.namespace).
		Delete(ctx, machine.Id, v1.DeleteOptions{})
	if err != nil {
		return err
	}

	err = m.clientSet.
		NetworkingV1().
		Ingresses(m.namespace).
		Delete(ctx, fmt.Sprintf("%s-static", machine.Id), v1.DeleteOptions{})

	return err
}

func (m *MachineManagerK8s) getAddress(machine machine.Machine) string {
	return fmt.Sprintf("%s/%s", m.route, machine.Id)
}

func (m *MachineManagerK8s) waitForVNC(ctx context.Context, machine machine.Machine, timeout time.Duration) error {
	return wait.PollUntilContextTimeout(ctx, time.Second*5, timeout, true, func(ctx context.Context) (done bool, err error) {
		set, err := m.clientSet.
			AppsV1().
			StatefulSets(m.namespace).
			Get(ctx, machine.Id, v1.GetOptions{})
		if err != nil {
			return false, err
		}

		if set.Status.ReadyReplicas == *set.Spec.Replicas {
			return true, nil
		}

		return false, nil
	})
}

func (m *MachineManagerK8s) createVNC(ctx context.Context, machine machine.Machine) error {
	statefulSetTemplate := m.vncStatefulSetTemplate(machine)
	_, err := m.clientSet.
		AppsV1().
		StatefulSets(m.namespace).
		Create(ctx, &statefulSetTemplate, v1.CreateOptions{})
	if err != nil {
		return err
	}

	serviceTemplate := m.vncServiceTemplate(machine)
	_, err = m.clientSet.
		CoreV1().
		Services(m.namespace).
		Create(ctx, &serviceTemplate, v1.CreateOptions{})
	if err != nil {
		return err
	}

	ingressTemplate := m.vncIngressTemplate(machine)
	_, err = m.clientSet.
		NetworkingV1().
		Ingresses(m.namespace).
		Create(ctx, &ingressTemplate, v1.CreateOptions{})
	if err != nil {
		return err
	}

	staticIngressTemplate := m.vncStaticIngressTemplate(machine)
	_, err = m.clientSet.
		NetworkingV1().
		Ingresses(m.namespace).
		Create(ctx, &staticIngressTemplate, v1.CreateOptions{})
	return err
}

func (m *MachineManagerK8s) vncIngressTemplate(machine machine.Machine) networkv1.Ingress {
	ingressClass := "nginx"
	pathType := networkv1.PathTypePrefix
	return networkv1.Ingress{
		ObjectMeta: v1.ObjectMeta{
			Name: machine.Id,
			Annotations: map[string]string{
				"nginx.ingress.kubernetes.io/rewrite-target": "/",
			},
		},
		Spec: networkv1.IngressSpec{
			IngressClassName: &ingressClass,
			Rules: []networkv1.IngressRule{
				{
					IngressRuleValue: networkv1.IngressRuleValue{
						HTTP: &networkv1.HTTPIngressRuleValue{
							Paths: []networkv1.HTTPIngressPath{
								{
									Path:     m.getAddress(machine),
									PathType: &pathType,
									Backend: networkv1.IngressBackend{
										Service: &networkv1.IngressServiceBackend{
											Name: fmt.Sprintf("svc-%s", machine.Id),
											Port: networkv1.ServiceBackendPort{
												Number: 80,
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func (m *MachineManagerK8s) vncStaticIngressTemplate(machine machine.Machine) networkv1.Ingress {
	ingressClass := "nginx"
	pathType := networkv1.PathTypeImplementationSpecific
	return networkv1.Ingress{
		ObjectMeta: v1.ObjectMeta{
			Name: fmt.Sprintf("%s-static", machine.Id),
			Annotations: map[string]string{
				"nginx.ingress.kubernetes.io/rewrite-target": "/$1",
				"nginx.ingress.kubernetes.io/use-regex":      "true",
			},
		},
		Spec: networkv1.IngressSpec{
			IngressClassName: &ingressClass,
			Rules: []networkv1.IngressRule{
				{
					IngressRuleValue: networkv1.IngressRuleValue{
						HTTP: &networkv1.HTTPIngressRuleValue{
							Paths: []networkv1.HTTPIngressPath{
								{
									Path:     fmt.Sprintf("%s/(.*)", m.route),
									PathType: &pathType,
									Backend: networkv1.IngressBackend{
										Service: &networkv1.IngressServiceBackend{
											Name: fmt.Sprintf("svc-%s", machine.Id),
											Port: networkv1.ServiceBackendPort{
												Number: 80,
											},
										},
									},
								},
								{
									Path:     fmt.Sprintf("/machines/connect/%s(.*)", machine.Id),
									PathType: &pathType,
									Backend: networkv1.IngressBackend{
										Service: &networkv1.IngressServiceBackend{
											Name: fmt.Sprintf("svc-%s", machine.Id),
											Port: networkv1.ServiceBackendPort{
												Number: 80,
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func (m *MachineManagerK8s) vncServiceTemplate(machine machine.Machine) corev1.Service {
	return corev1.Service{
		ObjectMeta: v1.ObjectMeta{
			Name: fmt.Sprintf("svc-%s", machine.Id),
		},
		Spec: corev1.ServiceSpec{
			Type: corev1.ServiceTypeClusterIP,
			Selector: map[string]string{
				"id":      machine.Id,
				"ownerId": machine.OwnerId,
			},
			Ports: []corev1.ServicePort{
				{
					Name:     "vnc",
					Port:     80,
					Protocol: corev1.ProtocolTCP,
				},
			},
		},
	}
}

func (m *MachineManagerK8s) vncStatefulSetTemplate(machine machine.Machine) appsv1.StatefulSet {
	replicas := int32(1)
	return appsv1.StatefulSet{
		ObjectMeta: v1.ObjectMeta{
			Name: machine.Id,
		},
		Spec: appsv1.StatefulSetSpec{
			ServiceName: machine.Id,
			Selector: &v1.LabelSelector{
				MatchLabels: map[string]string{
					"id":      machine.Id,
					"ownerId": machine.OwnerId,
				},
			},
			Template: vncPodTemplate(machine),
			Replicas: &replicas,
		},
	}
}

func vncPodTemplate(machine machine.Machine) corev1.PodTemplateSpec {
	return corev1.PodTemplateSpec{
		ObjectMeta: v1.ObjectMeta{
			Labels: map[string]string{
				"id":      machine.Id,
				"ownerId": machine.OwnerId,
			},
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:  "vnc",
					Image: "dorowu/ubuntu-desktop-lxde-vnc:latest",
					Ports: []corev1.ContainerPort{
						{
							Name:          "vnc",
							ContainerPort: 80,
						},
					},
					Env: []corev1.EnvVar{
						{
							Name:  "VNC_PASSWORD",
							Value: machine.Password,
						},
						{
							Name:  "USER",
							Value: machine.Username,
						},
						{
							Name:  "PASSWORD",
							Value: machine.Password,
						},
					},
				},
			},
		},
	}
}
