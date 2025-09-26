package main

import (
    "context"
    "flag"
    "fmt"
    "os"

    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    "k8s.io/client-go/kubernetes"
    "k8s.io/client-go/rest"
    "k8s.io/client-go/tools/clientcmd"
)

func getClient() (*kubernetes.Clientset, error) {
    kubeconfig := flag.String("kubeconfig", "", "path to kubeconfig (optional)")
    flag.Parse()
    if *kubeconfig != "" {
        cfg, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
        if err != nil {
            return nil, err
        }
        return kubernetes.NewForConfig(cfg)
    }
    // in-cluster config
    cfg, err := rest.InClusterConfig()
    if err != nil {
        return nil, err
    }
    return kubernetes.NewForConfig(cfg)
}

func main() {
    client, err := getClient()
    if err != nil {
        fmt.Fprintf(os.Stderr, "❌ Failed to connect: %v\n", err)
        os.Exit(1)
    }

    // List nodes
    nodes, err := client.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
    if err != nil {
        panic(err)
    }
    fmt.Printf("Nodes in cluster: %d\n", len(nodes.Items))
    for _, n := range nodes.Items {
        fmt.Println(" -", n.Name)
    }

    // List pods in all namespaces
    pods, err := client.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
    if err != nil {
        panic(err)
    }
    fmt.Printf("Pods in cluster: %d\n", len(pods.Items))
}
