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
	return ext.AddBucket(cctx, *request.Bucket)
}

func (mw *Middleware) ExtensionSetMode(cctx context.Context, request *pb.RpcExtensionSetModeRequest) *pb.RpcExtensionSetModeResponse {
	ext := getService[extensions.Service](mw)
	return ext.SetMode(cctx, request.Mode)
}

func (mw *Middleware) ExtensionInstallByID(cctx context.Context, request *pb.RpcExtensionInstallByIDRequest) *pb.RpcExtensionInstallByIDResponse {
	ext := getService[extensions.Service](mw)
	return ext.InstallByID(cctx, request.Id)
}
