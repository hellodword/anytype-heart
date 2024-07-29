package extensions

import (
	"archive/zip"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/anyproto/any-sync/app"
	"github.com/anyproto/anytype-heart/core/wallet"
	"github.com/anyproto/anytype-heart/pkg/lib/logging"
	"github.com/anyproto/anytype-heart/pkg/lib/pb/model"
	extism "github.com/extism/go-sdk"
	"go.uber.org/zap"
)

const CName = "extensions"

var log = logging.Logger(CName)

type Service interface {
	app.Component
	ListBuckets(ctx context.Context) ([]*model.ExtensionBucketInfo, error)
	AddBucket(ctx context.Context, bucket *model.ExtensionBucket) (*model.ExtensionBucketInfo, error)
	RemoveBucket(ctx context.Context, bucketId string) error
	GetDeveloperMode(ctx context.Context) (bool, error)
	SetDeveloperMode(ctx context.Context, developerMode bool) (bool, error)
	GetByID(ctx context.Context, bucketId, extensionId string) (*model.ExtensionBucketExtensionInfo, error)
	InstallByURL(ctx context.Context, u string) (*model.ExtensionManifest, error)
	Call(ctx context.Context, extensionId, functionName, blockId string) error
}

func New() Service {
	return &service{}
}

type service struct {
	lock sync.RWMutex

	buckets       map[string]*model.ExtensionBucketInfo
	developerMode bool

	plugins map[string]*extism.Plugin

	repoPath         string
	rootPath         string
	downloadRootPath string
	extractRootPath  string
}

func (s *service) Name() string {
	return CName
}

func (s *service) Init(a *app.App) error {
	// avoid null while json marsharl
	s.buckets = make(map[string]*model.ExtensionBucketInfo)
	s.plugins = make(map[string]*extism.Plugin)

	s.repoPath = app.MustComponent[wallet.Wallet](a).RepoPath()
	s.rootPath = app.MustComponent[wallet.Wallet](a).RootPath()
	s.downloadRootPath = filepath.Join(s.rootPath, "extensions", "download")
	s.extractRootPath = filepath.Join(s.rootPath, "extensions", "extract")
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

func (s *service) GetDeveloperMode(ctx context.Context) (bool, error) {
	return s.developerMode, nil
}

func (s *service) SetDeveloperMode(ctx context.Context, developerMode bool) (bool, error) {
	log.Info("SetDeveloperMode", developerMode)
	s.lock.Lock()
	defer s.lock.Unlock()
	s.developerMode = developerMode
	return s.developerMode, nil
}

func (s *service) GetByID(ctx context.Context, bucketId, extensionId string) (*model.ExtensionBucketExtensionInfo, error) {
	panic("TODO")
}

func (s *service) InstallByURL(ctx context.Context, u string) (*model.ExtensionManifest, error) {
	log.Info("InstallByURL", zap.String("url", u))

	if !s.developerMode {
		return nil, errors.New("not in developer mode")
	}

	filename, err := downloadFile(ctx, u, s.downloadRootPath)
	if err != nil {
		return nil, err
	}

	extracted := filepath.Join(s.extractRootPath, strings.TrimSuffix(filename, ".ext"))
	err = Unzip(filepath.Join(s.downloadRootPath, filename), extracted)
	if err != nil {
		return nil, err
	}

	// manifest := extism.Manifest{
	// 	Wasm: []extism.Wasm{
	// 		extism.WasmUrl{
	// 			Url: u,
	// 		},
	// 	},
	// }

	// config := extism.PluginConfig{}
	// if runtime.GOOS == "ios" {
	// 	config.RuntimeConfig = wazero.NewRuntimeConfigInterpreter()
	// }

	// plugin, err := extism.NewPlugin(context.Background(), manifest, config, []extism.HostFunction{})
	// if err != nil {
	// 	return err
	// }

	// s.lock.Lock()
	// defer s.lock.Unlock()
	// s.plugins[u] = plugin
	// s.plugins["id1"] = plugin

	return &model.ExtensionManifest{}, nil
}

func (s *service) Call(ctx context.Context, extensionId, functionName, blockId string) error {
	return nil
}

func downloadFile(ctx context.Context, u, dir string) (filename string, err error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return
	}

	filename = filepath.Base(req.URL.Path)
	if filename == "" || !strings.HasSuffix(filename, ".ext") || strings.TrimSuffix(filename, ".ext") == "" {
		err = fmt.Errorf("invalid filename: %s", filename)
		return
	}

	f, err := os.OpenFile(filepath.Join(dir, filename), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return
	}
	defer f.Close()

	req.Header.Set("User-Agent", "Mozilla/5.0 (AnytypeExtensionDownloader/1.0)")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}

	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("http status code: %d", resp.StatusCode)
		return
	}

	defer resp.Body.Close()

	_, err = io.Copy(f, resp.Body)
	if err != nil {
		return
	}

	return
}

// https://gist.github.com/paulerickson/6d8650947ee4e3f3dbcc28fde10eaae7
/**
 * Extract a zip file named source to directory destination.  Handles cases where destination dir…
 *  - does not exist (creates it)
 *  - is empty
 *  - already has source archive extracted into it (files are overwritten)
 *  - has other files in it, not in source archive (not overwritten)
 * But is expected to fail if it…
 *  - is not writable
 *  - contains a non-empty directory with the same path as a file in source archive (that's not a simple overwrite)
 */
func Unzip(source, destination string) error {
	archive, err := zip.OpenReader(source)
	if err != nil {
		return err
	}
	defer archive.Close()
	for _, file := range archive.Reader.File {
		reader, err := file.Open()
		if err != nil {
			return err
		}
		defer reader.Close()
		path := filepath.Join(destination, file.Name)
		// Remove file if it already exists; no problem if it doesn't; other cases can error out below
		_ = os.Remove(path)
		// Create a directory at path, including parents
		err = os.MkdirAll(path, os.ModePerm)
		if err != nil {
			return err
		}
		// If file is _supposed_ to be a directory, we're done
		if file.FileInfo().IsDir() {
			continue
		}
		// otherwise, remove that directory (_not_ including parents)
		err = os.Remove(path)
		if err != nil {
			return err
		}
		// and create the actual file.  This ensures that the parent directories exist!
		// An archive may have a single file with a nested path, rather than a file for each parent dir
		writer, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			return err
		}
		defer writer.Close()
		_, err = io.Copy(writer, reader)
		if err != nil {
			return err
		}
	}
	return nil
}
