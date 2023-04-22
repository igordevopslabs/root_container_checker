package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

var (
	rootContainersGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "root_container_checker_total",
			Help: "Current number of containers running as root.",
		},
		[]string{"pod", "namespace"},
	)
	podHealthGauge = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "root_container_checker_pod_health",
			Help: "Health status of the root container checker Pod (1: healthy, 0: unhealthy).",
		},
	)
)

func init() {
	prometheus.MustRegister(rootContainersGauge)
	prometheus.MustRegister(podHealthGauge)
}

func checkRootContainers(clientset *kubernetes.Clientset) {
	// Set Pod health gauge to 1 (healthy) initially
	podHealthGauge.Set(1)

	pods, err := clientset.CoreV1().Pods("").List(context.Background(), metav1.ListOptions{})
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		// Set Pod health gauge to 0 (unhealthy) if there is an error
		podHealthGauge.Set(0)
		return
	}

	for _, pod := range pods.Items {
		rootContainers := 0
		for _, container := range pod.Spec.Containers {
			if container.SecurityContext == nil || container.SecurityContext.RunAsUser == nil || *container.SecurityContext.RunAsUser == 0 {
				rootContainers++
			}
		}

		if rootContainers > 0 {
			rootContainersGauge.WithLabelValues(pod.Name, pod.Namespace).Set(float64(rootContainers))
		}
	}
}

func main() {
	http.HandleFunc("/_ready", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "OK")
	})

	http.HandleFunc("/_healthy", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "OK")
	})

	http.Handle("/metrics", promhttp.Handler())

	go func() {
		http.ListenAndServe(":8080", nil)
	}()

	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	// Configura o ticker para executar a verificação a cada 1h
	ticker := time.NewTicker(60 * time.Second)
	defer ticker.Stop()

	for {
		checkRootContainers(clientset)
		<-ticker.C
	}
}
