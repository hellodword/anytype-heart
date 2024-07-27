package core

import (
	"context"

	"github.com/anyproto/anytype-heart/pb"
)

func (mw *Middleware) ExtensionSetEndpoint(cctx context.Context, request *pb.RpcExtensionSetEndpointRequest) *pb.RpcExtensionSetEndpointResponse {
	return &pb.RpcExtensionSetEndpointResponse{
		Error: &pb.RpcExtensionSetEndpointResponseError{
			Code: pb.RpcExtensionSetEndpointResponseError_NULL,
		},
	}
}

func (mw *Middleware) ExtensionSetMode(cctx context.Context, request *pb.RpcExtensionSetModeRequest) *pb.RpcExtensionSetModeResponse {
	return &pb.RpcExtensionSetModeResponse{
		Error: &pb.RpcExtensionSetModeResponseError{
			Code: pb.RpcExtensionSetModeResponseError_NULL,
		},
	}
}

func (mw *Middleware) ExtensionInstallByID(cctx context.Context, request *pb.RpcExtensionInstallByIDRequest) *pb.RpcExtensionInstallByIDResponse {
	return &pb.RpcExtensionInstallByIDResponse{
		Error: &pb.RpcExtensionInstallByIDResponseError{
			Code: pb.RpcExtensionInstallByIDResponseError_NULL,
		},
	}
}
