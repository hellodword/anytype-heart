package extensions

import (
	"context"
	"errors"
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
	ListBuckets(ctx context.Context) ([]*model.ExtensionBucketInfo, error)
	AddBucket(ctx context.Context, bucket *model.ExtensionBucket) (*model.ExtensionBucketInfo, error)
	RemoveBucket(ctx context.Context, bucketId string) error
	GetStatus(ctx context.Context) (bool, error)
	SetDeveloperMode(ctx context.Context, developerMode bool) (bool, error)
	GetByID(ctx context.Context, bucketId, extensionId string) *pb.RpcExtensionGetByIDResponse
	InstallByURL(ctx context.Context, u string) *pb.RpcExtensionInstallByURLResponse
}

func New() Service {
	return &service{}
}

type service struct {
	lock sync.RWMutex

	buckets       map[string]*model.ExtensionBucketInfo
	developerMode bool
}

func (s *service) Name() string {
	return CName
}

func (s *service) Init(a *app.App) error {
	// avoid null while json marsharl
	s.buckets = make(map[string]*model.ExtensionBucketInfo)
	return nil
}

func (s *service) ListBuckets(ctx context.Context) ([]*model.ExtensionBucketInfo, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()
	var buckets []*model.ExtensionBucketInfo
	for _, v := range s.buckets {
		buckets = append(buckets, v)
	}
	return buckets, nil
}

func (s *service) AddBucket(ctx context.Context, bucket *model.ExtensionBucket) (*model.ExtensionBucketInfo, error) {
	if bucket.Name == "" || bucket.Endpoint == "" {
		return nil, errors.New("invalid input")
	}
	var info = &model.ExtensionBucketInfo{
		Name:     bucket.GetName(),
		Endpoint: bucket.GetEndpoint(),
	}

	// TODO fetch bucket info
	// info.Id = ""

	if info.GetId() == "" {
		return nil, errors.New("invalid id")
	}

	s.lock.Lock()
	defer s.lock.Unlock()
	s.buckets[info.GetId()] = info

	return info, nil
}

func (s *service) RemoveBucket(ctx context.Context, bucketId string) error {
	s.lock.Lock()
	defer s.lock.Unlock()
	delete(s.buckets, bucketId)
	return nil
}

func (s *service) GetStatus(ctx context.Context) (bool, error) {
	return s.developerMode, nil
}

func (s *service) SetDeveloperMode(ctx context.Context, developerMode bool) (bool, error) {
	log.Info("SetDeveloperMode", developerMode)
	s.lock.Lock()
	defer s.lock.Unlock()
	s.developerMode = developerMode
	return s.developerMode, nil
}

func (s *service) GetByID(ctx context.Context, bucketId, extensionId string) *pb.RpcExtensionGetByIDResponse {
	return &pb.RpcExtensionGetByIDResponse{}
}

func (s *service) InstallByURL(ctx context.Context, u string) *pb.RpcExtensionInstallByURLResponse {
	log.Info("InstallByURL", zap.String("url", u))

	if !s.developerMode {
		return &pb.RpcExtensionInstallByURLResponse{
			Error: &pb.RpcExtensionInstallByURLResponseError{Code: pb.RpcExtensionInstallByURLResponseError_INTERNAL_ERROR},
		}
	}

	return &pb.RpcExtensionInstallByURLResponse{
		Error: &pb.RpcExtensionInstallByURLResponseError{
			Code: pb.RpcExtensionInstallByURLResponseError_NULL,
		},
	}
}
