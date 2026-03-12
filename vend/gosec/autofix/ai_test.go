package autofix

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/securego/gosec/v2/issue"
)

// MockGenAIClient is a mock of the GenAIClient interface
type MockGenAIClient struct {
	mock.Mock
}

func (m *MockGenAIClient) GenerateSolution(ctx context.Context, prompt string) (string, error) {
	args := m.Called(ctx, prompt)
	return args.String(0), args.Error(1)
}

func TestGenerateSolutionByGemini_Success(t *testing.T) {
	// Arrange
	issues := []*issue.Issue{
		{What: "Example issue 1"},
	}

	mockClient := new(MockGenAIClient)
	mockClient.On("GenerateSolution", mock.Anything, mock.Anything).Return("Autofix for issue 1", nil).Once()

	// Act
	err := generateSolution(mockClient, issues)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, []*issue.Issue{{What: "Example issue 1", Autofix: "Autofix for issue 1"}}, issues)
	mock.AssertExpectationsForObjects(t, mockClient)
}

func TestGenerateSolutionByGemini_NoCandidates(t *testing.T) {
	// Arrange
	issues := []*issue.Issue{
		{What: "Example issue 2"},
	}

	mockClient := new(MockGenAIClient)
	mockClient.On("GenerateSolution", mock.Anything, mock.Anything).Return("", nil).Once()

	// Act
	err := generateSolution(mockClient, issues)

	// Assert
	require.EqualError(t, err, "no autofix returned by gemini")
	mock.AssertExpectationsForObjects(t, mockClient)
}

func TestGenerateSolutionByGemini_APIError(t *testing.T) {
	// Arrange
	issues := []*issue.Issue{
		{What: "Example issue 3"},
	}

	mockClient := new(MockGenAIClient)
	mockClient.On("GenerateSolution", mock.Anything, mock.Anything).Return("", errors.New("API error")).Once()

	// Act
	err := generateSolution(mockClient, issues)

	// Assert
	require.EqualError(t, err, "generating autofix with gemini: API error")
	mock.AssertExpectationsForObjects(t, mockClient)
}

func TestGenerateSolution_UnsupportedProvider(t *testing.T) {
	// Arrange
	issues := []*issue.Issue{
		{What: "Example issue 4"},
	}

	// Act
	// Note: With default OpenAI-compatible fallback, this will attempt to create an OpenAI client
	// The test will fail during client initialization due to missing/invalid API key or base URL
	err := GenerateSolution("custom-model", "", "", false, issues)

	// Assert
	// Expect an error during client initialization or API call
	require.Error(t, err)
}
