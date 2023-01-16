// Package httpclient ...
// CODE GENERATED AUTOMATICALLY
// DO NOT EDIT
package httpclient

import (
	"context"

	v1 "github.com/rinnedrag/go-transport-generator/example/api/v1"
	"github.com/stretchr/testify/mock"
)

// MockService ...
type MockService struct {
	mock.Mock
}

// CreateMultipartUpload ...
func (s *MockService) CreateMultipartUpload(ctx context.Context, bucket string, key string) (data v1.CreateMultipartUploadData, errorFlag bool, errorText string, additionalErrors *v1.AdditionalErrors, err error) {
	args := s.Called(context.Background(), bucket, key)
	return args.Get(0).(v1.CreateMultipartUploadData), args.Get(1).(bool), args.Get(2).(string), args.Get(3).(*v1.AdditionalErrors), args.Error(4)
}

// UploadPartDocument ...
func (s *MockService) UploadPartDocument(ctx context.Context, bucket string, key string, uploadID string, partNumber int64, document []byte) (err error) {
	args := s.Called(context.Background(), bucket, key, uploadID, partNumber, document)
	return args.Error(0)
}

// CompleteUpload ...
func (s *MockService) CompleteUpload(ctx context.Context, bucket string, key string, uploadID string) (err error) {
	args := s.Called(context.Background(), bucket, key, uploadID)
	return args.Error(0)
}

// UploadDocument ...
func (s *MockService) UploadDocument(ctx context.Context, bucket string, key string, document []byte) (err error) {
	args := s.Called(context.Background(), bucket, key, document)
	return args.Error(0)
}

// DownloadDocument ...
func (s *MockService) DownloadDocument(ctx context.Context, bucket string, key string) (document []byte, err error) {
	args := s.Called(context.Background(), bucket, key)
	return args.Get(0).([]byte), args.Error(1)
}

// GetToken ...
func (s *MockService) GetToken(ctx context.Context, authToken *string, scope string, grantType string) (token string, expiresIn int, err error) {
	args := s.Called(context.Background(), authToken, scope, grantType)
	return args.Get(0).(string), args.Get(1).(int), args.Error(2)
}

// GetBranches ...
func (s *MockService) GetBranches(ctx context.Context, authToken *string, supplierID *string) (branches []int, err error) {
	args := s.Called(context.Background(), authToken, supplierID)
	return args.Get(0).([]int), args.Error(1)
}
