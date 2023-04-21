package main

import (
	"context"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func main() {
	// rest.InClusterConfig() obtém a configuração do cluster a partir do ambiente em que o programa está sendo executado (dentro do cluster).
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err)
	}

	// kubernetes.NewForConfig() cria um cliente para se comunicar com a API do Kubernetes usando a configuração obtida anteriormente.
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	// clientset.CoreV1().Pods("").List() lista todos os pods em todos os namespaces no cluster.
	pods, err := clientset.CoreV1().Pods("").List(context.Background(), metav1.ListOptions{})
	if err != nil {
		panic(err)
	}

	// Itera sobre todos os pods retornados.
	for _, pod := range pods.Items {
		rootContainers := 0

		// Itera sobre todos os containers dentro de cada pod.
		for _, container := range pod.Spec.Containers {
			// Verifica se o container está executando como root.
			if container.SecurityContext == nil || container.SecurityContext.RunAsUser == nil || *container.SecurityContext.RunAsUser == 0 {
				rootContainers++
			}
		}

		// Imprime o nome do pod, o namespace e a quantidade de containers executando como root.
		fmt.Printf("Pod: %s, Namespace: %s, Root Containers: %d\n", pod.Name, pod.Namespace, rootContainers)
	}
}
