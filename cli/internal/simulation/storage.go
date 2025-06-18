package simulation

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/nitrictech/nitric/cli/internal/netx"
	"github.com/nitrictech/nitric/cli/internal/simulation/middleware"
	storagepb "github.com/nitrictech/nitric/proto/storage/v2"
	"github.com/spf13/afero"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	signingSecret     *string = nil
	signingSecretLock sync.Mutex
)

func getSigningSecret() ([]byte, error) {
	signingSecretLock.Lock()
	defer signingSecretLock.Unlock()

	if signingSecret == nil {
		key := make([]byte, 32) // Generate a 256-bit key

		_, err := rand.Read(key)
		if err != nil {
			return nil, err
		}

		secret := base64.StdEncoding.EncodeToString(key)
		signingSecret = &secret
	}

	return []byte(*signingSecret), nil
}

func tokenFromRequest(req *storagepb.StoragePreSignUrlRequest) *jwt.Token {
	return jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp": time.Now().Add(req.Expiry.AsDuration()).Unix(),
		"request": map[string]string{
			"bucket": req.BucketName,
			"key":    req.Key,
			"op":     req.Operation.String(),
		},
	})
}

func requestFromToken(token string) (*storagepb.StoragePreSignUrlRequest, error) {
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return getSigningSecret()
	}, jwt.WithExpirationRequired())
	if err != nil {
		return nil, err
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("could not convert claims to map")
	}

	requestMap, ok := claims["request"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("could not convert request to map")
	}

	return &storagepb.StoragePreSignUrlRequest{
		BucketName: requestMap["bucket"].(string),
		Key:        requestMap["key"].(string),
		Operation:  storagepb.StoragePreSignUrlRequest_Operation(storagepb.StoragePreSignUrlRequest_Operation_value[requestMap["op"].(string)]),
	}, nil
}

func (s *SimulationServer) Delete(ctx context.Context, req *storagepb.StorageDeleteRequest) (*storagepb.StorageDeleteResponse, error) {
	path, err := GetBlobPath(s.appDir, req.BucketName, req.Key)
	if err != nil {
		return nil, err
	}

	err = s.fs.Remove(path)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "%v", err)
	}

	return &storagepb.StorageDeleteResponse{}, nil
}

func (s *SimulationServer) Exists(ctx context.Context, req *storagepb.StorageExistsRequest) (*storagepb.StorageExistsResponse, error) {
	path, err := GetBlobPath(s.appDir, req.BucketName, req.Key)
	if err != nil {
		return nil, err
	}

	exists, err := afero.Exists(s.fs, path)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "%v", err)
	}

	return &storagepb.StorageExistsResponse{
		Exists: exists,
	}, nil
}

func (s *SimulationServer) ListBlobs(ctx context.Context, req *storagepb.StorageListBlobsRequest) (*storagepb.StorageListBlobsResponse, error) {
	path, err := GetBucketPath(s.appDir, req.BucketName)
	if err != nil {
		return nil, err
	}

	files, err := afero.ReadDir(s.fs, path)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "%v", err)
	}

	blobs := make([]*storagepb.Blob, 0, len(files))
	for _, file := range files {
		if file.IsDir() {
			continue // Skip directories
		}

		blobs = append(blobs, &storagepb.Blob{
			Key: file.Name(),
		})
	}

	return &storagepb.StorageListBlobsResponse{
		Blobs: blobs,
	}, nil
}

func (s *SimulationServer) PreSignUrl(ctx context.Context, req *storagepb.StoragePreSignUrlRequest) (*storagepb.StoragePreSignUrlResponse, error) {
	// Call to ensure the bucket exists
	_, err := GetBlobPath(s.appDir, req.BucketName, req.Key)
	if err != nil {
		return nil, err
	}

	var address string = ""

	token := tokenFromRequest(req)

	secret, err := getSigningSecret()
	if err != nil {
		return nil, status.Error(codes.Internal, "error generating presigned url, could not get signing secret")
	}

	tokenString, err := token.SignedString(secret)
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("error generating presigned url, could not sign token: %v", err))
	}

	// XXX: Do not URL encode keys (path needs to be preserved)
	// TODO: May need to re-write slashes to a non-escapable character format
	switch req.Operation {
	case storagepb.StoragePreSignUrlRequest_WRITE:
		address = fmt.Sprintf("http://localhost:%d/write/%s", s.fileServerPort, tokenString)
	case storagepb.StoragePreSignUrlRequest_READ:
		address = fmt.Sprintf("http://localhost:%d/read/%s", s.fileServerPort, tokenString)
	}

	if address == "" {
		return nil, status.Error(codes.Internal, "error generating presigned url, unknown operation")
	}

	return &storagepb.StoragePreSignUrlResponse{
		Url: address,
	}, nil
}

func (s *SimulationServer) Read(ctx context.Context, req *storagepb.StorageReadRequest) (*storagepb.StorageReadResponse, error) {
	path, err := GetBlobPath(s.appDir, req.BucketName, req.Key)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "%v", err)
	}

	data, err := afero.ReadFile(s.fs, path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, status.Errorf(codes.NotFound, "Blob %s not found in bucket %s", req.Key, req.BucketName)
		}
		return nil, status.Errorf(codes.Internal, "%v", err)
	}
	return &storagepb.StorageReadResponse{
		Body: data,
	}, nil
}

func (s *SimulationServer) Write(ctx context.Context, req *storagepb.StorageWriteRequest) (*storagepb.StorageWriteResponse, error) {
	path, err := GetBlobPath(s.appDir, req.BucketName, req.Key)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "%v", err)
	}

	err = afero.WriteFile(s.fs, path, req.Body, os.ModePerm)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "%v", err)
	}

	return &storagepb.StorageWriteResponse{}, nil
}

const (
	BUCKET_MIN_PORT = 5000
	BUCKET_MAX_PORT = 5999
)

func DetectContentType(filename string, content []byte) string {
	contentType := mime.TypeByExtension(filepath.Ext(filename))
	if contentType != "" {
		return contentType
	}

	return http.DetectContentType(content)
}

func (s *SimulationServer) startBuckets() error {
	bucketIntents := s.appSpec.GetBucketIntents()

	for bucketName, bucketIntent := range bucketIntents {
		if bucketIntent.ContentPath != "" {
			bucketDir, err := s.ensureBucketDir(bucketName)
			if err != nil {
				return err
			}

			err = s.CopyDir(bucketDir, filepath.Join(s.appDir, bucketIntent.ContentPath))
			if err != nil {
				return err
			}
		}
	}

	router := http.NewServeMux()

	// Serve files (for presigned URLS)
	// TODO: Add origin check for an entrypoint
	httpFs := afero.NewHttpFs(s.fs)
	router.Handle("GET /", http.FileServer(httpFs.Dir(BucketsDir)))

	router.HandleFunc("GET /read/{token}", func(w http.ResponseWriter, r *http.Request) {
		token := r.PathValue("token")
		if token == "" {
			http.Error(w, "Token is required", http.StatusBadRequest)
			return
		}

		req, err := requestFromToken(token)
		if err != nil {
			http.Error(w, "Invalid token", http.StatusBadRequest)
			return
		}

		if req.Operation != storagepb.StoragePreSignUrlRequest_READ {
			http.Error(w, "Invalid operation", http.StatusBadRequest)
			return
		}

		resp, err := s.Read(r.Context(), &storagepb.StorageReadRequest{
			BucketName: req.BucketName,
			Key:        req.Key,
		})
		if err != nil {
			http.Error(w, fmt.Sprintf("Error reading file: %v", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", DetectContentType(req.Key, resp.Body))

		w.Write(resp.Body)
	})

	router.HandleFunc("PUT /write/{token}", func(w http.ResponseWriter, r *http.Request) {
		token := r.PathValue("token")
		if token == "" {
			http.Error(w, "Token is required", http.StatusBadRequest)
			return
		}

		req, err := requestFromToken(token)
		if err != nil {
			http.Error(w, "Invalid token", http.StatusBadRequest)
			return
		}

		if req.Operation != storagepb.StoragePreSignUrlRequest_WRITE {
			http.Error(w, "Invalid operation", http.StatusBadRequest)
			return
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error reading body", http.StatusInternalServerError)
			return
		}

		_, err = s.Write(r.Context(), &storagepb.StorageWriteRequest{
			BucketName: req.BucketName,
			Key:        req.Key,
			Body:       body,
		})
		if err != nil {
			http.Error(w, "Error writing file", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
	})

	// Wrap the router with CORS middleware
	handler := middleware.CORS(middleware.DefaultCORSOptions)(router)

	reservedPort, err := netx.GetNextPort(netx.MinPort(BUCKET_MIN_PORT), netx.MaxPort(BUCKET_MAX_PORT))
	if err != nil {
		return err
	}

	go http.ListenAndServe(fmt.Sprintf(":%d", reservedPort), handler)

	s.fileServerPort = int(reservedPort)

	return nil
}
