// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.19.2
// source: blog/blogpb/blog.proto

package blogpb

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

// BlogServicClient is the client API for BlogServic service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type BlogServicClient interface {
	CreateBlog(ctx context.Context, in *CreateBlogRequest, opts ...grpc.CallOption) (*CreateBlogResponse, error)
	//return NOT_FOUND if not found
	ReadBlog(ctx context.Context, in *ReadBlogRequest, opts ...grpc.CallOption) (*ReadBlogResponse, error)
}

type blogServicClient struct {
	cc grpc.ClientConnInterface
}

func NewBlogServicClient(cc grpc.ClientConnInterface) BlogServicClient {
	return &blogServicClient{cc}
}

func (c *blogServicClient) CreateBlog(ctx context.Context, in *CreateBlogRequest, opts ...grpc.CallOption) (*CreateBlogResponse, error) {
	out := new(CreateBlogResponse)
	err := c.cc.Invoke(ctx, "/blog.BlogServic/CreateBlog", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *blogServicClient) ReadBlog(ctx context.Context, in *ReadBlogRequest, opts ...grpc.CallOption) (*ReadBlogResponse, error) {
	out := new(ReadBlogResponse)
	err := c.cc.Invoke(ctx, "/blog.BlogServic/ReadBlog", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// BlogServicServer is the server API for BlogServic service.
// All implementations must embed UnimplementedBlogServicServer
// for forward compatibility
type BlogServicServer interface {
	CreateBlog(context.Context, *CreateBlogRequest) (*CreateBlogResponse, error)
	//return NOT_FOUND if not found
	ReadBlog(context.Context, *ReadBlogRequest) (*ReadBlogResponse, error)
	mustEmbedUnimplementedBlogServicServer()
}

// UnimplementedBlogServicServer must be embedded to have forward compatible implementations.
type UnimplementedBlogServicServer struct {
}

func (UnimplementedBlogServicServer) CreateBlog(context.Context, *CreateBlogRequest) (*CreateBlogResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateBlog not implemented")
}
func (UnimplementedBlogServicServer) ReadBlog(context.Context, *ReadBlogRequest) (*ReadBlogResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ReadBlog not implemented")
}
func (UnimplementedBlogServicServer) mustEmbedUnimplementedBlogServicServer() {}

// UnsafeBlogServicServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to BlogServicServer will
// result in compilation errors.
type UnsafeBlogServicServer interface {
	mustEmbedUnimplementedBlogServicServer()
}

func RegisterBlogServicServer(s grpc.ServiceRegistrar, srv BlogServicServer) {
	s.RegisterService(&BlogServic_ServiceDesc, srv)
}

func _BlogServic_CreateBlog_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateBlogRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BlogServicServer).CreateBlog(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/blog.BlogServic/CreateBlog",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BlogServicServer).CreateBlog(ctx, req.(*CreateBlogRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _BlogServic_ReadBlog_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ReadBlogRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BlogServicServer).ReadBlog(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/blog.BlogServic/ReadBlog",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BlogServicServer).ReadBlog(ctx, req.(*ReadBlogRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// BlogServic_ServiceDesc is the grpc.ServiceDesc for BlogServic service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var BlogServic_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "blog.BlogServic",
	HandlerType: (*BlogServicServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateBlog",
			Handler:    _BlogServic_CreateBlog_Handler,
		},
		{
			MethodName: "ReadBlog",
			Handler:    _BlogServic_ReadBlog_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "blog/blogpb/blog.proto",
}
