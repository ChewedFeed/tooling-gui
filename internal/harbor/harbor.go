package harbor

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/chewedfeed/automated/internal/config"
	"io"
	"net/http"
	"time"
)

type Harbor struct {
	Username string
	Password string
	Server   string
}

type RobotDetails struct {
	Name   string
	Secret string
}

func NewHarbor() (*Harbor, error) {
	creds, err := config.GetCredentials()
	if err != nil {
		return nil, err
	}

	return &Harbor{
		Username: creds.Harbor.Credentials.Username,
		Password: creds.Harbor.Credentials.Password,
		Server:   creds.Harbor.ServiceAddress,
	}, nil
}

func (h *Harbor) GetProjects() ([]string, error) {
	type metadataJSON struct {
		Public string `json:"public"`
	}
	type CVEAllowlistJSON struct {
		CreationTime string   `json:"creation_time"`
		ID           int      `json:"id"`
		ProjectID    int      `json:"project_id"`
		Items        []string `json:"items"`
		UpdateTime   string   `json:"update_time"`
	}
	type projectJSON struct {
		ChartCount         int              `json:"chart_count"`
		CreationTime       string           `json:"creation_time"`
		CurrentUserRoleID  int              `json:"current_user_role_id"`
		CurrentUserRoleIDs []int            `json:"current_user_role_ids"`
		CVEAllowlistJSON   CVEAllowlistJSON `json:"cve_allowlist"`
		MetaData           metadataJSON     `json:"metadata"`
		Name               string           `json:"name"`
		OwnerID            int              `json:"owner_id"`
		OwnerName          string           `json:"owner_name"`
		ProjectID          int              `json:"project_id"`
		RepoCount          int              `json:"repo_count"`
		UpdateTime         string           `json:"update_time"`
	}

	client := &http.Client{}
	auth := base64.StdEncoding.EncodeToString([]byte(h.Username + ":" + h.Password))
	req, _ := http.NewRequest("GET", fmt.Sprintf("%s/api/v2.0/projects?page=1&page_size=30", h.Server), nil)
	req.Header.Set("Authorization", "Basic "+auth)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, err
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	p := []projectJSON{}
	if err = json.Unmarshal(body, &p); err != nil {
		return nil, err
	}
	ret := []string{}
	for _, v := range p {
		ret = append(ret, v.Name)
	}

	return ret, nil
}

func (h *Harbor) GetProjectsErm() ([]string, error) {
	type accessJSON struct {
		Action   string `json:"action"`
		Resource string `json:"resource"`
	}
	type permissionsJSON struct {
		Access    []accessJSON `json:"access"`
		Kind      string       `json:"kind"`
		Namespace string       `json:"namespace"`
	}
	type projectJSON struct {
		CreateTime  string            `json:"create_time"`
		Disable     bool              `json:"disable"`
		Duraction   int               `json:"duraction"`
		ID          int               `json:"id"`
		Name        string            `json:"name"`
		Editable    bool              `json:"editable"`
		ExpiresAt   int               `json:"expires_at"`
		Permissions []permissionsJSON `json:"permissions"`
		UpdateTime  string            `json:"update_time"`
	}

	client := &http.Client{}
	auth := base64.StdEncoding.EncodeToString([]byte(h.Username + ":" + h.Password))
	req, _ := http.NewRequest("GET", fmt.Sprintf("%s/api/v2.0/robots", h.Server), nil)
	req.Header.Set("Authorization", "Basic "+auth)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, err
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	p := []projectJSON{}
	if err = json.Unmarshal(body, &p); err != nil {
		return nil, err
	}

	ret := []string{}
	for _, project := range p {
		for _, permission := range project.Permissions {
			ret = append(ret, permission.Namespace)
		}
	}

	return ret, nil
}

func (h *Harbor) CreateRobot(name, project string, system bool, expireDate time.Time) (*RobotDetails, error) {
	type accessJSON struct {
		Action   string `json:"action"`
		Resource string `json:"resource"`
	}
	type permissionsJSON struct {
		Access    []accessJSON `json:"access"`
		Kind      string       `json:"kind"`
		Namespace string       `json:"namespace"`
	}
	type robotJSON struct {
		Level       string            `json:"level"`
		Name        string            `json:"name"`
		Duration    int64             `json:"duration"`
		Permissions []permissionsJSON `json:"permissions"`
	}

	level := "project"
	if system {
		level = "system"
	}

	var expireSend int64 = 0
	if !expireDate.IsZero() {
		expireSend = expireDate.Unix()
	}

	jd, err := json.Marshal(robotJSON{
		Level:    level,
		Name:     name,
		Duration: expireSend,
		Permissions: []permissionsJSON{
			{
				Access: []accessJSON{
					{
						Action:   "list",
						Resource: "repository",
					},
					{
						Action:   "push",
						Resource: "repository",
					},
					{
						Action:   "delete",
						Resource: "repository",
					},
					{
						Action:   "read",
						Resource: "artifact",
					},
					{
						Action:   "list",
						Resource: "artifact",
					},
					{
						Action:   "delete",
						Resource: "artifact",
					},
					{
						Action:   "create",
						Resource: "artifact-label",
					},
					{
						Action:   "delete",
						Resource: "artifact-label",
					},
					{
						Action:   "create",
						Resource: "tag",
					},
					{
						Action:   "delete",
						Resource: "tag",
					},
					{
						Action:   "list",
						Resource: "tag",
					},
					{
						Action:   "create",
						Resource: "scan",
					},
					{
						Action:   "stop",
						Resource: "scan",
					},
					{
						Action:   "read",
						Resource: "helm-chart",
					},
					{
						Action:   "create",
						Resource: "helm-chart-version",
					},
					{
						Action:   "delete",
						Resource: "helm-chart-version",
					},
					{
						Action:   "create",
						Resource: "helm-chart-version-label",
					},
					{
						Action:   "pull",
						Resource: "repository",
					},
					{
						Action:   "delete",
						Resource: "helm-chart-version-label",
					},
				},
				Kind:      "project",
				Namespace: project,
			},
		},
	})
	if err != nil {
		return nil, err
	}
	auth := base64.StdEncoding.EncodeToString([]byte(h.Username + ":" + h.Password))
	client := &http.Client{}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v2.0/robots", h.Server), bytes.NewBuffer(jd))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Basic "+auth)
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 201 {
		return nil, err
	}

	type robotDetailsJSON struct {
		CreateTime string `json:"create_time"`
		ExpiresAt  int    `json:"expires_at"`
		ID         int    `json:"id"`
		Name       string `json:"name"`
		Secret     string `json:"secret"`
	}
	rdj := robotDetailsJSON{}
	if err = json.NewDecoder(resp.Body).Decode(&rdj); err != nil {
		return nil, err
	}

	return &RobotDetails{
		Name:   rdj.Name,
		Secret: rdj.Secret,
	}, nil
}

func (h *Harbor) ValidateRobot(name, secret string) (bool, error) {
	username := fmt.Sprintf("robot$%s", name)
	auth := base64.StdEncoding.EncodeToString([]byte(username + ":" + secret))
	client := &http.Client{}
	req, _ := http.NewRequest("GET", fmt.Sprintf("%s/api/v2.0/projects?page=1&page_size=30", h.Server), nil)
	req.Header.Set("Authorization", "Basic "+auth)
	resp, err := client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return false, nil
	}

	return true, nil
}
