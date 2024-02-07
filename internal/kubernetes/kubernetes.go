package kubernetes

import (
  "context"
  "encoding/json"
  "fmt"

  v1 "k8s.io/api/core/v1"
  "k8s.io/client-go/kubernetes"
  "k8s.io/client-go/tools/clientcmd"
  "k8s.io/client-go/util/homedir"

  metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Kubernetes struct {
  clientset *kubernetes.Clientset
}

func NewKubernetes() (*Kubernetes, error) {
  cfg, err := clientcmd.BuildConfigFromFlags("", fmt.Sprintf("%s/.kube/config", homedir.HomeDir()))
  if err != nil {
    return nil, err
  }

  clientset, err := kubernetes.NewForConfig(cfg)
  if err != nil {
    return nil, err
  }

  return &Kubernetes{
    clientset: clientset,
  }, nil
}

func (k *Kubernetes) GetNamespaces() ([]string, error) {
  ns, err := k.clientset.CoreV1().Namespaces().List(context.Background(), metav1.ListOptions{})
  if err != nil {
    return nil, err
  }
  names := make([]string, 0)
  for _, n := range ns.Items {
    names = append(names, n.Name)
  }

  return names, nil
}

func (k *Kubernetes) CreateSecret(name, namespace string, data map[string]string) error {
  secretData := make(map[string][]byte)
  for k, v := range data {
    secretData[k] = []byte(v)
  }

  _, err := k.clientset.CoreV1().Secrets(namespace).Create(context.Background(), &v1.Secret{
    ObjectMeta: metav1.ObjectMeta{
      Name: name,
    },
    Type: v1.SecretTypeOpaque,
    Data: secretData,
  }, metav1.CreateOptions{})
  return err
}

func (k *Kubernetes) CreateDockerSecret(server, username, password, email, namespace string) error {
  type dockerAuth struct {
    Username string `json:"username"`
    Password string `json:"password"`
    Email    string `json:"email"`
  }

  type dockerAuths struct {
    Auths map[string]dockerAuth `json:"auths"`
  }

  auth := dockerAuths{
    Auths: map[string]dockerAuth{
      server: {
        Username: username,
        Password: password,
        Email:    email,
      },
    },
  }

  secretData, err := json.Marshal(auth)
  if err != nil {
    return err
  }

  secret := make(map[string]string)
  secret[".dockerconfigjson"] = fmt.Sprintf("%s", secretData)

  return k.CreateSecret("docker-registry-secret", namespace, secret)
}
