package core

import (
	"context"
	"net/url"

	"github.com/anyproto/anytype-heart/pb"
)

func (mw *Middleware) ExtensionListBuckets(cctx context.Context, request *pb.RpcExtensionListBucketsRequest) *pb.RpcExtensionListBucketsResponse {
	return &pb.RpcExtensionListBucketsResponse{
		Error: &pb.RpcExtensionListBucketsResponseError{
			Code: pb.RpcExtensionListBucketsResponseError_NULL,
		},
	}
}

func (mw *Middleware) ExtensionAddBucket(cctx context.Context, request *pb.RpcExtensionAddBucketRequest) *pb.RpcExtensionAddBucketResponse {
	var response = &pb.RpcExtensionAddBucketResponse{
		Error: &pb.RpcExtensionAddBucketResponseError{
			Code: pb.RpcExtensionAddBucketResponseError_NULL,
		},
	}

	// TODO validate inputs with https://github.com/go-playground/validator
	if request.Bucket.Name == "" || request.Bucket.Endpoint == "" {
		response.Error.Code = pb.RpcExtensionAddBucketResponseError_BAD_INPUT
		return response
	}

	if _, err := url.Parse(request.Bucket.Endpoint); err != nil {
		response.Error.Code = pb.RpcExtensionAddBucketResponseError_BAD_INPUT
		return response
	}

	// TODO fetch bucket info here

	return &pb.RpcExtensionAddBucketResponse{
		Error: &pb.RpcExtensionAddBucketResponseError{
			Code: pb.RpcExtensionAddBucketResponseError_NULL,
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
