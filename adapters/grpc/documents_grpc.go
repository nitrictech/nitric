package grpc

import (
	"context"

	pb "github.com/nitric-dev/membrane/interfaces/nitric/v1"
	"github.com/nitric-dev/membrane/sdk"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/structpb"
)

// GRPC Interface for registered Nitric Documents Plugins
type DocumentsServer struct {
	pb.UnimplementedDocumentServer
	// TODO: Support multiple plugin registerations
	// Just need to settle on a way of addressing them on calls
	documentsPlugin sdk.DocumentService
}

func (s *DocumentsServer) checkPluginRegistered() (bool, error) {
	if s.documentsPlugin == nil {
		return false, status.Errorf(codes.Unimplemented, "Documents plugin not registered")
	}

	return true, nil
}

func (s *DocumentsServer) Create(ctx context.Context, req *pb.DocumentCreateRequest) (*pb.DocumentCreateResponse, error) {
	if ok, err := s.checkPluginRegistered(); ok {
		if err := s.documentsPlugin.Create(req.GetCollection(), req.GetKey(), req.GetDocument().AsMap()); err == nil {
			return &pb.DocumentCreateResponse{}, nil
		} else {
			// Case: Failed to create the document
			// TODO: Translate from internal Documents Plugin Error
			return nil, err
		}
	} else {
		// Case: Plugin was not registered
		return nil, err
	}
}

func (s *DocumentsServer) Get(ctx context.Context, req *pb.DocumentGetRequest) (*pb.DocumentGetResponse, error) {
	if ok, err := s.checkPluginRegistered(); ok {
		if document, err := s.documentsPlugin.Get(req.GetCollection(), req.GetKey()); err == nil {
			if doc, err := structpb.NewStruct(document); err == nil {
				return &pb.DocumentGetResponse{
					Document: doc,
				}, nil
			} else {
				// Case: Failed to create PB struct from stored document
				// TODO: Translate from internal Documents Plugin Error
				return nil, err
			}
		} else {
			// Case: There was an error retrieving the document
			// TODO: Translate from internal Documents Plugin Error
			return nil, err
		}
	} else {
		// Case: The documents plugin was not registered
		// TODO: Translate from internal Documents Plugin Error
		return nil, err
	}
}

func (s *DocumentsServer) Update(ctx context.Context, req *pb.DocumentUpdateRequest) (*pb.DocumentUpdateResponse, error) {
	if ok, err := s.checkPluginRegistered(); ok {
		if err := s.documentsPlugin.Create(req.GetCollection(), req.GetKey(), req.GetDocument().AsMap()); err == nil {
			return &pb.DocumentUpdateResponse{}, nil
		} else {
			// Case: Failed to create the document
			// TODO: Translate from internal Documents Plugin Error
			return nil, err
		}
	} else {
		// Case: Plugin was not registered
		return nil, err
	}
}

func (s *DocumentsServer) Delete(ctx context.Context, req *pb.DocumentDeleteRequest) (*pb.DocumentDeleteResponse, error) {
	if ok, err := s.checkPluginRegistered(); ok {
		if err := s.documentsPlugin.Delete(req.GetCollection(), req.GetKey()); err == nil {
			return &pb.DocumentDeleteResponse{}, nil
		} else {
			// Case: Failed to create the document
			// TODO: Translate from internal Documents Plugin Error
			return nil, err
		}
	} else {
		// Case: Plugin was not registered
		return nil, err
	}
}

func NewDocumentsServer(documentsPlugin sdk.DocumentService) pb.DocumentServer {
	return &DocumentsServer{
		documentsPlugin: documentsPlugin,
	}
}
