package kubernetes

import "github.com/chewedfeed/automated/internal/config"

type Kubernetes struct {
	Token string
}

func NewKubernetes() (*Kubernetes, error) {
	creds, err := config.GetCredentials()
	if err != nil {
		return nil, err
	}

	return &Kubernetes{
		Token: creds.Kubernetes.Credentials.Token,
	}, nil
}
