package main

import (
	"context"
	"fmt"

	"github.com/ghodss/yaml"
	v1 "istio.io/api/alibabacloud-servicemesh/v1"
	versionedclient "istio.io/client-go/asm/pkg/clientset"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/clientcmd"
)

const (
	swimlanegroupYaml = `apiVersion: istio.alibabacloud.com/v1
kind: ASMSwimLaneGroup
metadata:
  name: mock
spec:
  ingress:
    gateway:
      name: ingressgateway
      namespace: istio-system
      type: ASM
  services:
  - name: mocka
    namespace: default
  - name: mockb
    namespace: default
  - name: mockc
    namespace: default`
	swimlaneYaml = `apiVersion: istio.alibabacloud.com/v1
kind: ASMSwimLane
metadata:
  labels:
    swimlane-group: mock
  name: v1
spec:
  ingressRules:
  - hosts:
    - '*'
    match:
      uri:
        exact: /mock
    name: r1
    online: true
    route:
      destination:
        host: mocka.default.svc.cluster.local
  labelSelector:
    version: v1`
)

func main() {
	const NAMESPACE = "default"
	restConfig, err := clientcmd.BuildConfigFromFlags("", "{your kubeconfig path here}") // 将{}内容替换成Kube Config实际路径，默认是$HOME/.kube/config。
	if err != nil {
		fmt.Println("获取连接配置失败")
		return
	}
	clientset, err := versionedclient.NewForConfig(restConfig)
	if err != nil {
		fmt.Println("创建客户端失败")
		return
	}
	// 此处通过YAML反序列化构造VirtualService结构体，也可以直接构造。
	asmswimlanegroup := &v1.ASMSwimLaneGroup{}
	err = yaml.Unmarshal([]byte(swimlanegroupYaml), &asmswimlanegroup)
	if err != nil {
		fmt.Println("反序列化失败")
		return
	}
	_, err = clientset.IstioV1().ASMSwimLaneGroups().Create(context.TODO(), asmswimlanegroup, metav1.CreateOptions{})
	if err != nil {
		fmt.Println("创建泳道组失败")
	}

	asmswimlane := &v1.ASMSwimLane{}
	err = yaml.Unmarshal([]byte(swimlaneYaml), &asmswimlane)
	if err != nil {
		fmt.Println("反序列化失败")
		return
	}
	_, err = clientset.IstioV1().ASMSwimLanes().Create(context.TODO(), asmswimlane, metav1.CreateOptions{})
	if err != nil {
		fmt.Println("创建泳道失败")
	}
	printAllASMSwimlaneGroups(clientset, NAMESPACE)
}

// 打印命名空间下全部的泳道组。
func printAllASMSwimlaneGroups(clientset *versionedclient.Clientset, namespace string) {
	slgList, err := clientset.IstioV1().ASMSwimLaneGroups().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		fmt.Println("获取泳道组失败")
	}
	for _, vs := range slgList.Items {
		fmt.Println(vs)
	}
	slList, err := clientset.IstioV1().ASMSwimLanes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		fmt.Println("获取泳道失败")
	}
	for _, vs := range slList.Items {
		fmt.Println(vs)
	}
}
