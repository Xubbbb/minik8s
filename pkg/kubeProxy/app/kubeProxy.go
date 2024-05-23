package app

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/MiniK8s-SE3356/minik8s/pkg/kubeProxy/hostsController"
	"github.com/MiniK8s-SE3356/minik8s/pkg/kubeProxy/iptablesController"
	kptype "github.com/MiniK8s-SE3356/minik8s/pkg/kubeProxy/types"
	httpobject "github.com/MiniK8s-SE3356/minik8s/pkg/types/httpObject"
	"github.com/MiniK8s-SE3356/minik8s/pkg/utils/httpRequest"
	"github.com/MiniK8s-SE3356/minik8s/pkg/utils/poller"
)

type KubeProxy struct {
	iptables_controller *(iptablesController.IptablesController)
	hosts_controller    *(hostsController.HostsController)
	services_status     kptype.KpServicesStatus
	dns_status          kptype.KpDnsStatus
	mutex               sync.Mutex
}

var service_need_send bool = false
var dns_need_send bool = false

func NewKubeProxy() *KubeProxy {
	fmt.Printf("Create KubeProxy ...\n")
	iptables_controller := iptablesController.NewIptablesController()
	hosts_controller := hostsController.NewHostsController()
	kube_proxy := &KubeProxy{
		iptables_controller: iptables_controller,
		hosts_controller:    hosts_controller,
		services_status:     kptype.KpServicesStatus{},
		dns_status:          kptype.KpDnsStatus{},
		mutex:               sync.Mutex{},
	}
	return kube_proxy
}

func (kp *KubeProxy) Init() {
	fmt.Printf("Init KubeProxy ...\n")
	kp.iptables_controller.Init()
	kp.hosts_controller.Init()
}

func (kp *KubeProxy) Run() {
	fmt.Printf("Run KubeProxy ...\n")

	go poller.PollerStaticPeriod(10*time.Second, kp.syncServices, true)
	go poller.PollerStaticPeriod(5*time.Second, kp.syncServicesAndDnsToKubelet, true)

	poller.PollerStaticPeriod(10*time.Second, kp.syncDns, true)
}

func (kp *KubeProxy) syncDns() {
	// 获得所有dns
	var dns_list httpobject.HTTPResponse_GetAllDns = httpobject.HTTPResponse_GetAllDns{}
	status, err := httpRequest.GetRequestByObject("http://192.168.1.6:8080/api/v1/GetAllDNS", nil, &dns_list)
	if status != http.StatusOK || err != nil {
		fmt.Printf("EndpointsController routine error get, status %d, return\n", status)
		return
	}

	// 更新本地/etc/hosts
	new_dns_status, err := kp.hosts_controller.SyncEtcHosts(&dns_list)
	if err != nil {
		return
	}

	// 如果正常更新，则可以更新本地dns状态，并准备向kubelet同步
	kp.mutex.Lock()
	kp.dns_status = new_dns_status
	dns_need_send = true
	kp.mutex.Unlock()
}

func (kp *KubeProxy) syncServices() {
	fmt.Printf("KubeProxy routine...\n")

	var service_list httpobject.HTTPResponse_GetAllServices
	status, err := httpRequest.GetRequestByObject("http://192.168.1.6:8080/api/v1/GetAllService", nil, &service_list)

	if status != http.StatusOK || err != nil {
		fmt.Printf("routine error get, status %d, return\n", status)
		return
	}

	var endpoint_list httpobject.HTTPResponse_GetAllEndpoint
	status, err = httpRequest.GetRequestByObject("http://192.168.1.6:8080/api/v1/GetAllEndpoint", nil, &endpoint_list)

	if status != http.StatusOK || err != nil {
		fmt.Printf("routine error get, status %d, return\n", status)
		return
	}

	new_service, err := kp.iptables_controller.SyncConfig(&service_list, &endpoint_list)
	if err != nil {
		return
	}

	err = kp.iptables_controller.SyncIptables()
	if err != nil {
		return
	} else {
		kp.mutex.Lock()
		kp.services_status = new_service
		service_need_send = true
		kp.mutex.Unlock()
	}

}

func (kp *KubeProxy) syncServicesAndDnsToKubelet() {
	// 向kubelet更新本机上的service信息
	if service_need_send || dns_need_send {
		kp.mutex.Lock()
		if dns_need_send {
			// TODO： dns状态报告给kubelet
		}

		if service_need_send {
			// TODO： service状态报告给kubelet
		}
		kp.mutex.Unlock()
	}
}
