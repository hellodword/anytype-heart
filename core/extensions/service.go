package extensions

import (
	"context"
	"sync"

	"github.com/anyproto/any-sync/app"
	"github.com/anyproto/anytype-heart/pb"
	"github.com/anyproto/anytype-heart/pkg/lib/logging"
	"github.com/anyproto/anytype-heart/pkg/lib/pb/model"
	"go.uber.org/zap"
)

const CName = "extensions"

var log = logging.Logger(CName)

type Service interface {
	app.Component
	ListBuckets(ctx context.Context) *pb.RpcExtensionListBucketsResponse
	AddBucket(ctx context.Context, bucket *model.ExtensionBucket) *pb.RpcExtensionAddBucketResponse
	RemoveBucket(ctx context.Context, bucketId string) *pb.RpcExtensionRemoveBucketResponse
	SetMode(ctx context.Context, mode model.ExtensionMode) *pb.RpcExtensionSetModeResponse
	GetByID(ctx context.Context, bucketId, extensionId string) *pb.RpcExtensionGetByIDResponse
	InstallByURL(ctx context.Context, u string) *pb.RpcExtensionInstallByURLResponse
}

func New() Service {
	return &service{}
}

type service struct {
	lock sync.RWMutex

	buckets map[string]*model.ExtensionBucketInfo
	mode    model.ExtensionMode
}

func (s *service) Name() string {
	return CName
}

func (s *service) Init(a *app.App) error {
	// avoid null while json marsharl
	s.buckets = make(map[string]*model.ExtensionBucketInfo)
	return nil
}

func (s *service) ListBuckets(ctx context.Context) *pb.RpcExtensionListBucketsResponse {
	s.lock.RLock()
	defer s.lock.RUnlock()
	var buckets []*model.ExtensionBucketInfo
	for _, v := range s.buckets {
		buckets = append(buckets, v)
	}
	return &pb.RpcExtensionListBucketsResponse{Buckets: buckets}
}

func (s *service) AddBucket(ctx context.Context, bucket *model.ExtensionBucket) *pb.RpcExtensionAddBucketResponse {
	if bucket.Name == "" || bucket.Endpoint == "" {
		return &pb.RpcExtensionAddBucketResponse{
			Error: &pb.RpcExtensionAddBucketResponseError{Code: pb.RpcExtensionAddBucketResponseError_BAD_INPUT},
		}
	}
	var info = &model.ExtensionBucketInfo{
		Name:     bucket.GetName(),
		Endpoint: bucket.GetEndpoint(),
	}

	// TODO fetch bucket info
	// info.Id = ""

	if info.GetId() == "" {
		return &pb.RpcExtensionAddBucketResponse{
			Error: &pb.RpcExtensionAddBucketResponseError{Code: pb.RpcExtensionAddBucketResponseError_INTERNAL_ERROR},
		}
	}

	s.lock.Lock()
	defer s.lock.Unlock()
	s.buckets[info.GetId()] = info

	return &pb.RpcExtensionAddBucketResponse{Bucket: info}
}

func (s *service) RemoveBucket(ctx context.Context, bucketId string) *pb.RpcExtensionRemoveBucketResponse {
	s.lock.Lock()
	defer s.lock.Unlock()
	delete(s.buckets, bucketId)
	return &pb.RpcExtensionRemoveBucketResponse{}
}

func (s *service) SetMode(ctx context.Context, mode model.ExtensionMode) *pb.RpcExtensionSetModeResponse {
	log.Info("SetMode", zap.Int32("mode", int32(mode)))
	s.lock.Lock()
	defer s.lock.Unlock()
	s.mode = mode
	return &pb.RpcExtensionSetModeResponse{Mode: s.mode}
}

func (s *service) GetByID(ctx context.Context, bucketId, extensionId string) *pb.RpcExtensionGetByIDResponse {
	return &pb.RpcExtensionGetByIDResponse{}
}

func (s *service) InstallByURL(ctx context.Context, u string) *pb.RpcExtensionInstallByURLResponse {
	switch s.mode {
	case model.ExtensionMode_Developer:
		break
	default:
		return &pb.RpcExtensionInstallByURLResponse{
			Error: &pb.RpcExtensionInstallByURLResponseError{Code: pb.RpcExtensionInstallByURLResponseError_INTERNAL_ERROR},
		}
	}

	return &pb.RpcExtensionInstallByURLResponse{}
}
