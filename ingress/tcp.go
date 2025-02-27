package ingress

import (
	"context"
	"encoding/json"
	"sort"

	scalablev1alpha1 "github.com/benm-stm/solace-scalable-k8s-operator/api/v1alpha1"
	libs "github.com/benm-stm/solace-scalable-k8s-operator/common"
	"github.com/rung/go-safecast"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// get already existing default haproxy svc and add ports
func NewTcp(
	s *scalablev1alpha1.SolaceScalable,
	ports []corev1.ServicePort,
	d map[string]string,
) *[]corev1.ServicePort {
	// get default
	svcPorts := *GetDefaultTcp(ports)
	var portExist bool
	var portIndex int

	for k := range d {
		portExist = false
		portIndex = 0
		var svcPort corev1.ServicePort

		port, err := safecast.Atoi32(k)
		if err == nil {
			//check if the svc exist
			for i, p := range ports {
				if p.Port == int32(port) {
					portExist = true
					portIndex = i
				}
			}
			if !portExist {
				//create new serviceport
				svcPort = corev1.ServicePort{
					Name:        "tcp-" + k,
					Protocol:    "TCP",
					Port:        int32(port),
					AppProtocol: nil,
				}
			} else {
				svcPort = ports[portIndex]
			}
			svcPorts = append(svcPorts, svcPort)
		}
	}
	return &svcPorts

}

func GetDefaultTcp(servicePorts []corev1.ServicePort) *[]corev1.ServicePort {
	var svcPorts = []corev1.ServicePort{}
	for _, s := range servicePorts {
		if s.Name == "http" || s.Name == "https" || s.Name == "stat" {
			s.AppProtocol = nil
			svcPorts = append(svcPorts, s)
		}
	}
	return &svcPorts
}

func GetTcp(
	solaceScalable *scalablev1alpha1.SolaceScalable,
	serviceName string,
	k k8sClient,
	ctx context.Context,
) (*corev1.Service, error) {
	log := log.FromContext(ctx)
	FoundHaproxySvc := &corev1.Service{}
	if err := k.Get(
		ctx,
		types.NamespacedName{
			Namespace: solaceScalable.Spec.Haproxy.Namespace,
			Name:      serviceName,
		}, FoundHaproxySvc,
	); err != nil {
		log.Info("HAProxy service is not found",
			FoundHaproxySvc.Namespace,
			FoundHaproxySvc.Name,
		)
		return nil, err
	}
	return FoundHaproxySvc, nil
}

func UpdateTcp(
	hashStore *map[string]string,
	FoundHaproxySvc *corev1.Service,
	k k8sClient,
	ctx context.Context,
) error {
	log := log.FromContext(ctx)
	// sort the data (ports cause marshall to fail)
	sort.Slice(FoundHaproxySvc.Spec.Ports, func(i, j int) bool {
		return FoundHaproxySvc.Spec.Ports[i].Name <
			FoundHaproxySvc.Spec.Ports[j].Name
	})
	portsMarshal, _ := json.Marshal(FoundHaproxySvc.Spec.Ports)

	if (*hashStore)[FoundHaproxySvc.Name] == "" ||
		libs.AsSha256(portsMarshal) != (*hashStore)[FoundHaproxySvc.Name] {
		log.Info("Updating Haproxy Svc",
			FoundHaproxySvc.Namespace,
			FoundHaproxySvc.Name,
		)
		if err := k.Update(ctx, FoundHaproxySvc); err != nil {
			return err
		}
		//update hash to not trig update if conf has not changed
		(*hashStore)[FoundHaproxySvc.Name] = libs.AsSha256(portsMarshal)
	}
	return nil
}
