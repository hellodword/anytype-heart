package extensions

import (
	"context"
	"sync"

	"github.com/anyproto/any-sync/app"
	"github.com/anyproto/anytype-heart/pb"
	"github.com/anyproto/anytype-heart/pkg/lib/pb/model"
)

const CName = "extensions"

type Service interface {
	app.Component
	ListBuckets(ctx context.Context) *pb.RpcExtensionListBucketsResponse
	AddBucket(ctx context.Context, bucket model.ExtensionBucket) *pb.RpcExtensionAddBucketResponse
	SetMode(ctx context.Context, mode model.ExtensionMode) *pb.RpcExtensionSetModeResponse
	InstallByID(ctx context.Context, id string) *pb.RpcExtensionInstallByIDResponse
}

func New() Service {
	return &service{}
}

type service struct {
	lock sync.RWMutex

	buckets []*model.ExtensionBucket
	mode    model.ExtensionMode
}

func (s *service) Name() string {
	return CName
}

func (s *service) Init(a *app.App) error {
	// avoid null while json marsharl
	s.buckets = make([]*model.ExtensionBucket, 0)
	return nil
}

func (s *service) ListBuckets(ctx context.Context) *pb.RpcExtensionListBucketsResponse {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return &pb.RpcExtensionListBucketsResponse{Buckets: s.buckets}
}

func (s *service) AddBucket(ctx context.Context, bucket model.ExtensionBucket) *pb.RpcExtensionAddBucketResponse {
	if bucket.Name == "" || bucket.Endpoint == "" {

	}
	// TODO fetch bucket info
	s.lock.Lock()
	defer s.lock.Unlock()
	s.buckets = append(s.buckets, &bucket)
	return nil
}

func (s *service) SetMode(ctx context.Context, mode model.ExtensionMode) *pb.RpcExtensionSetModeResponse {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.mode = mode
	return &pb.RpcExtensionSetModeResponse{Mode: s.mode}
}

func (s *service) InstallByID(ctx context.Context, id string) *pb.RpcExtensionInstallByIDResponse {
	return &pb.RpcExtensionInstallByIDResponse{}
}
