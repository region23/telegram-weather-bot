package services

import (
	"bytes"
	"io"
	"net/http"
	"testing"
)

type MockHTTPClient struct {
	DoFunc func(req *http.Request) (*http.Response, error)
}

func (m *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	return m.DoFunc(req)
}

func TestGetWeatherByCityName(t *testing.T) {
	tests := []struct {
		name            string
		mockAPIResponse string
		expectedResult  string
		expectError     bool
	}{
		{
			name:            "Successful API Response",
			mockAPIResponse: `{"name":"Moscow","main":{"temp":25.6},"weather":[{"description":"clear sky"}]}`,
			expectedResult:  "Город: Moscow\nОписание погоды: clear sky\nТемпература: 25.60°C",
			expectError:     false,
		},
		{
			name:            "Empty Weather Array",
			mockAPIResponse: `{"name":"Moscow","main":{"temp":25.6},"weather":[]}`,
			expectedResult:  "",
			expectError:     true,
		},
		// Добавьте другие тестовые сценарии по мере необходимости
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Создаем mock сервер
			mockClient := &MockHTTPClient{
				DoFunc: func(req *http.Request) (*http.Response, error) {
					return &http.Response{
						StatusCode: 200,
						Body:       io.NopCloser(bytes.NewBufferString(tt.mockAPIResponse)),
					}, nil
				},
			}

			result, err := GetWeatherByCityName(mockClient, "Moscow", "fakeAPIKey")

			if tt.expectError && err == nil {
				t.Errorf("expected an error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("did not expect an error but got: %v", err)
			}
			if result != tt.expectedResult {
				t.Errorf("expected result %v but got %v", tt.expectedResult, result)
			}
		})
	}
}
