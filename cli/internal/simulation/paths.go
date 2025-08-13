package simulation

import (
	"os"
	"path/filepath"

	"github.com/nitrictech/nitric/cli/internal/version"
)

var (
	DotSugaDir = "./" + version.ConfigDirName
	// Container of all buckets
	BucketsDir = filepath.Join(DotSugaDir, "buckets")
	// Container of all services
	ServicesDir = filepath.Join(DotSugaDir, "services")
	// Container of all service logs
	ServicesLogsDir = filepath.Join(ServicesDir, "logs")
)

// Get the path to the log file for a specific service
func GetServiceLogPath(appDir string, serviceName string) (string, error) {
	serviceLogPath := filepath.Join(appDir, ServicesLogsDir, serviceName+".log")

	err := os.MkdirAll(filepath.Dir(serviceLogPath), os.ModePerm)
	if err != nil {
		return "", err
	}

	return serviceLogPath, nil
}

// Get the path to the bucket directory for a specific bucket
func GetBucketPath(appDir string, bucketName string) (string, error) {
	bucketPath := filepath.Join(appDir, BucketsDir, bucketName)

	err := os.MkdirAll(bucketPath, os.ModePerm)
	if err != nil {
		return "", err
	}

	return bucketPath, nil
}

// Get the path to the blob file for a specific bucket
func GetBlobPath(appDir string, bucketName string, blobName string) (string, error) {
	bucketPath, err := GetBucketPath(appDir, bucketName)
	if err != nil {
		return "", err
	}

	return filepath.Join(bucketPath, blobName), nil
}
