package core

import (
	"context"

	"github.com/anyproto/anytype-heart/core/extensions"
	"github.com/anyproto/anytype-heart/pb"
	"github.com/anyproto/anytype-heart/pkg/lib/pb/model"
)

func (mw *Middleware) ExtensionListBuckets(cctx context.Context, request *pb.RpcExtensionListBucketsRequest) *pb.RpcExtensionListBucketsResponse {
	response := func(code pb.RpcExtensionListBucketsResponseErrorCode, err error, buckets []*model.ExtensionBucketInfo) *pb.RpcExtensionListBucketsResponse {
		m := &pb.RpcExtensionListBucketsResponse{Error: &pb.RpcExtensionListBucketsResponseError{Code: code}}
		if err != nil {
			m.Error.Description = err.Error()
		}
		m.Buckets = buckets
		return m
	}

	buckets, err := getService[extensions.Service](mw).ListBuckets(cctx)
	if err != nil {
		return response(pb.RpcExtensionListBucketsResponseError_UNKNOWN_ERROR, err, nil)
	}
	return response(pb.RpcExtensionListBucketsResponseError_NULL, nil, buckets)
}

func (mw *Middleware) ExtensionAddBucket(cctx context.Context, request *pb.RpcExtensionAddBucketRequest) *pb.RpcExtensionAddBucketResponse {
	response := func(code pb.RpcExtensionAddBucketResponseErrorCode, err error, bucket *model.ExtensionBucketInfo) *pb.RpcExtensionAddBucketResponse {
		m := &pb.RpcExtensionAddBucketResponse{Error: &pb.RpcExtensionAddBucketResponseError{Code: code}}
		if err != nil {
			m.Error.Description = err.Error()
		}
		m.Bucket = bucket
		return m
	}

	bucket, err := getService[extensions.Service](mw).AddBucket(cctx, request.GetBucket())
	if err != nil {
		return response(pb.RpcExtensionAddBucketResponseError_UNKNOWN_ERROR, err, nil)
	}
	return response(pb.RpcExtensionAddBucketResponseError_NULL, nil, bucket)
}

func (mw *Middleware) ExtensionRemoveBucket(cctx context.Context, request *pb.RpcExtensionRemoveBucketRequest) *pb.RpcExtensionRemoveBucketResponse {
	// ext := getService[extensions.Service](mw)
	// return ext.RemoveBucket(cctx, request.GetBucketId())

	response := func(code pb.RpcExtensionRemoveBucketResponseErrorCode, err error) *pb.RpcExtensionRemoveBucketResponse {
		m := &pb.RpcExtensionRemoveBucketResponse{Error: &pb.RpcExtensionRemoveBucketResponseError{Code: code}}
		if err != nil {
			m.Error.Description = err.Error()
		}
		return m
	}

	err := getService[extensions.Service](mw).RemoveBucket(cctx, request.GetBucketId())
	if err != nil {
		return response(pb.RpcExtensionRemoveBucketResponseError_UNKNOWN_ERROR, err)
	}
	return response(pb.RpcExtensionRemoveBucketResponseError_NULL, nil)
}

func (mw *Middleware) ExtensionGetDeveloperMode(cctx context.Context, request *pb.RpcExtensionGetDeveloperModeRequest) *pb.RpcExtensionGetDeveloperModeResponse {
	response := func(code pb.RpcExtensionGetDeveloperModeResponseErrorCode, err error, developerMode bool) *pb.RpcExtensionGetDeveloperModeResponse {
		m := &pb.RpcExtensionGetDeveloperModeResponse{Error: &pb.RpcExtensionGetDeveloperModeResponseError{Code: code}}
		if err != nil {
			m.Error.Description = err.Error()
		}
		m.DeveloperMode = developerMode
		return m
	}

	mode, err := getService[extensions.Service](mw).GetDeveloperMode(cctx)
	if err != nil {
		return response(pb.RpcExtensionGetDeveloperModeResponseError_UNKNOWN_ERROR, err, mode)
	}
	return response(pb.RpcExtensionGetDeveloperModeResponseError_NULL, nil, mode)
}

func (mw *Middleware) ExtensionSetDeveloperMode(cctx context.Context, request *pb.RpcExtensionSetDeveloperModeRequest) *pb.RpcExtensionSetDeveloperModeResponse {
	response := func(code pb.RpcExtensionSetDeveloperModeResponseErrorCode, err error, developerMode bool) *pb.RpcExtensionSetDeveloperModeResponse {
		m := &pb.RpcExtensionSetDeveloperModeResponse{Error: &pb.RpcExtensionSetDeveloperModeResponseError{Code: code}}
		if err != nil {
			m.Error.Description = err.Error()
		}
		m.DeveloperMode = developerMode
		return m
	}

	mode, err := getService[extensions.Service](mw).SetDeveloperMode(cctx, request.GetDeveloperMode())
	if err != nil {
		return response(pb.RpcExtensionSetDeveloperModeResponseError_UNKNOWN_ERROR, err, mode)
	}
	return response(pb.RpcExtensionSetDeveloperModeResponseError_NULL, nil, mode)
}

func (mw *Middleware) ExtensionGetByID(cctx context.Context, request *pb.RpcExtensionGetByIDRequest) *pb.RpcExtensionGetByIDResponse {
	ext := getService[extensions.Service](mw)
	return ext.GetByID(cctx, request.GetBucketId(), request.GetExtensionId())
}

func (mw *Middleware) ExtensionInstallByURL(cctx context.Context, request *pb.RpcExtensionInstallByURLRequest) *pb.RpcExtensionInstallByURLResponse {
	ext := getService[extensions.Service](mw)
	return ext.InstallByURL(cctx, request.GetUrl())
}
