package cmd

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/iamgp/hvr/internal/models"
)

func resolveDependencies(name, version string) ([]models.Library, error) {
	url := fmt.Sprintf("http://localhost:8080/resolve?name=%s&version=%s", name, version)
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve dependencies: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to resolve dependencies: %s", resp.Status)
	}

	var dependencies []models.Library
	err = json.NewDecoder(resp.Body).Decode(&dependencies)
	if err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return dependencies, nil
}
