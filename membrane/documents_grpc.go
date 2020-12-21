package membrane

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	pb "github.com/nitric-dev/membrane/interfaces/nitric/v1/documents"
	"github.com/nitric-dev/membrane/plugins/sdk"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/structpb"
)

// GRPC Interface for registered Nitric Documents Plugins
type DocumentsServer struct {
	pb.UnimplementedDocumentsServer
	// TODO: Support multiple plugin registerations
	// Just need to settle on a way of addressing them on calls
	documentsPlugin sdk.DocumentsPlugin
}

func (s *DocumentsServer) checkPluginRegistered() (bool, error) {
	if s.documentsPlugin == nil {
		return false, status.Errorf(codes.Unimplemented, "Documents plugin not registered")
	}

	return true, nil
}

func (s *DocumentsServer) CreateDocument(ctx context.Context, req *pb.CreateDocumentRequest) (*empty.Empty, error) {
	if ok, err := s.checkPluginRegistered(); ok {
		if err := s.documentsPlugin.CreateDocument(req.GetCollection(), req.GetKey(), req.GetDocument().AsMap()); err == nil {
			return &empty.Empty{}, nil
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

func (s *DocumentsServer) GetDocument(ctx context.Context, req *pb.GetDocumentRequest) (*pb.GetDocumentReply, error) {
	if ok, err := s.checkPluginRegistered(); ok {
		if document, err := s.documentsPlugin.GetDocument(req.GetCollection(), req.GetKey()); err == nil {
			if doc, err := structpb.NewStruct(document); err == nil {
				return &pb.GetDocumentReply{
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

func (s *DocumentsServer) UpdateDocument(ctx context.Context, req *pb.UpdateDocumentRequest) (*empty.Empty, error) {
	if ok, err := s.checkPluginRegistered(); ok {
		if err := s.documentsPlugin.CreateDocument(req.GetCollection(), req.GetKey(), req.GetDocument().AsMap()); err == nil {
			return &empty.Empty{}, nil
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

func (s *DocumentsServer) DeleteDocument(ctx context.Context, req *pb.DeleteDocumentRequest) (*empty.Empty, error) {
	if ok, err := s.checkPluginRegistered(); ok {
		if err := s.documentsPlugin.DeleteDocument(req.GetCollection(), req.GetKey()); err == nil {
			return &empty.Empty{}, nil
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

func NewGrpcDocumentsServer(documentsPlugin sdk.DocumentsPlugin) pb.DocumentsServer {
	return &DocumentsServer{
		documentsPlugin: documentsPlugin,
	}
}
