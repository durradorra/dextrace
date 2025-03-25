package infrastructure

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/brkss/dextrace/internal/domain"
)

type SibionicRepository struct {
	apiURL string
}

func NewSibionicRepository(apiURL string) *SibionicRepository {
	return &SibionicRepository{
		apiURL: apiURL,
	}
}

func (r *SibionicRepository) Login(user domain.User) (string, error) {
	client := &http.Client{}
	loginURL := fmt.Sprintf("%s/auth/app/user/login", r.apiURL)

	jsonData, err := json.Marshal(user)
	if err != nil {
		return "", fmt.Errorf("error marshaling user data: %v", err)
	}

	req, err := http.NewRequest("POST", loginURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	var loginResp domain.LoginResponse
	if err := json.NewDecoder(resp.Body).Decode(&loginResp); err != nil {
		return "", fmt.Errorf("error decoding response: %v", err)
	}

	if !loginResp.Success {
		return "", fmt.Errorf("login failed: %s", loginResp.Msg)
	}
	return loginResp.Data.AccessToken, nil
}

func (r *SibionicRepository) GetData(token string, userID string) (*domain.GlucoseDataResponse, error) {
	client := &http.Client{}
	dataURL := fmt.Sprintf("%s/user/app/follow/deviceGlucose", r.apiURL)

	body, err := json.Marshal(map[string]string{
		"id":    userID,
		"range": "24",
	})

	if err != nil {
		return nil, fmt.Errorf("error marshaling body: %v", err)
	}

	req, err := http.NewRequest("POST", dataURL, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	var dataResp domain.GlucoseDataResponse
	if err := json.NewDecoder(resp.Body).Decode(&dataResp); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	return &dataResp, nil
}