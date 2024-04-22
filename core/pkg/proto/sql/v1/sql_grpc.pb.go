// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.4
// source: nitric/proto/sql/v1/sql.proto

package sqlpb

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// SqlClient is the client API for Sql service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type SqlClient interface {
	// Retrieve the connection string for a given database
	ConnectionString(ctx context.Context, in *SqlConnectionStringRequest, opts ...grpc.CallOption) (*SqlConnectionStringResponse, error)
}

type sqlClient struct {
	cc grpc.ClientConnInterface
}

func NewSqlClient(cc grpc.ClientConnInterface) SqlClient {
	return &sqlClient{cc}
}

func (c *sqlClient) ConnectionString(ctx context.Context, in *SqlConnectionStringRequest, opts ...grpc.CallOption) (*SqlConnectionStringResponse, error) {
	out := new(SqlConnectionStringResponse)
	err := c.cc.Invoke(ctx, "/nitric.proto.sql.v1.Sql/ConnectionString", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// SqlServer is the server API for Sql service.
// All implementations should embed UnimplementedSqlServer
// for forward compatibility
type SqlServer interface {
	// Retrieve the connection string for a given database
	ConnectionString(context.Context, *SqlConnectionStringRequest) (*SqlConnectionStringResponse, error)
}

// UnimplementedSqlServer should be embedded to have forward compatible implementations.
type UnimplementedSqlServer struct {
}

func (UnimplementedSqlServer) ConnectionString(context.Context, *SqlConnectionStringRequest) (*SqlConnectionStringResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ConnectionString not implemented")
}

// UnsafeSqlServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to SqlServer will
// result in compilation errors.
type UnsafeSqlServer interface {
	mustEmbedUnimplementedSqlServer()
}

func RegisterSqlServer(s grpc.ServiceRegistrar, srv SqlServer) {
	s.RegisterService(&Sql_ServiceDesc, srv)
}

func _Sql_ConnectionString_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SqlConnectionStringRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SqlServer).ConnectionString(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/nitric.proto.sql.v1.Sql/ConnectionString",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SqlServer).ConnectionString(ctx, req.(*SqlConnectionStringRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Sql_ServiceDesc is the grpc.ServiceDesc for Sql service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Sql_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "nitric.proto.sql.v1.Sql",
	HandlerType: (*SqlServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ConnectionString",
			Handler:    _Sql_ConnectionString_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "nitric/proto/sql/v1/sql.proto",
}