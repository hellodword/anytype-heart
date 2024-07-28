package core

import (
	"context"

	"github.com/anyproto/anytype-heart/core/extensions"
	"github.com/anyproto/anytype-heart/pb"
)

func (mw *Middleware) ExtensionListBuckets(cctx context.Context, request *pb.RpcExtensionListBucketsRequest) *pb.RpcExtensionListBucketsResponse {
	ext := getService[extensions.Service](mw)
	return ext.ListBuckets(cctx)
}

func (mw *Middleware) ExtensionAddBucket(cctx context.Context, request *pb.RpcExtensionAddBucketRequest) *pb.RpcExtensionAddBucketResponse {
	ext := getService[extensions.Service](mw)
	return ext.AddBucket(cctx, request.GetBucket())
}

func (mw *Middleware) ExtensionRemoveBucket(cctx context.Context, request *pb.RpcExtensionRemoveBucketRequest) *pb.RpcExtensionRemoveBucketResponse {
	ext := getService[extensions.Service](mw)
	return ext.RemoveBucket(cctx, request.GetBucketId())
}

func (mw *Middleware) ExtensionSetMode(cctx context.Context, request *pb.RpcExtensionSetModeRequest) *pb.RpcExtensionSetModeResponse {
	ext := getService[extensions.Service](mw)
	return ext.SetMode(cctx, request.GetMode())
}

func (mw *Middleware) ExtensionGetByID(cctx context.Context, request *pb.RpcExtensionGetByIDRequest) *pb.RpcExtensionGetByIDResponse {
	ext := getService[extensions.Service](mw)
	return ext.GetByID(cctx, request.GetBucketId(), request.GetExtensionId())
}

func (mw *Middleware) ExtensionInstallByURL(cctx context.Context, request *pb.RpcExtensionInstallByURLRequest) *pb.RpcExtensionInstallByURLResponse {
	ext := getService[extensions.Service](mw)
	return ext.InstallByURL(cctx, request.GetUrl())
}
