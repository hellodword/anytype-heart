package extensions

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/anyproto/any-sync/app"
	"github.com/anyproto/anytype-heart/core/files"
	"github.com/anyproto/anytype-heart/core/files/fileobject"
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
	GetDeveloperMode(ctx context.Context) (bool, error)
	SetDeveloperMode(ctx context.Context, developerMode bool) error
	ListInstalled(ctx context.Context) ([]*model.ExtensionInfo, error)
	InstallFromURL(ctx context.Context, u string) (*model.ExtensionInfo, error)
	InstallFromZip(ctx context.Context, zipPath string) (*model.ExtensionInfo, error)
	Enable(ctx context.Context, ExtensionId string) error
	Disable(ctx context.Context, extensionId string) error
	Call(ctx context.Context, extensionId, functionName, blockId string) error
}

func New() Service {
	return &service{}
}

type Extension struct {
	ID       string                   `json:"id"`
	Enabled  bool                     `json:"enabled"`
	core     *extism.Plugin           `json:"-"`
	manifest *model.ExtensionManifest `json:"-"`
}

type service struct {
	lock sync.RWMutex

	developerMode bool

	extensions map[string]*Extension

	repoPath          string
	rootPath          string
	downloadRootPath  string
	extractedRootPath string

	fileService       files.Service
	fileObjectService fileobject.Service
}

func (s *service) Name() string {
	return CName
}

func (s *service) Init(a *app.App) error {
	s.extensions = make(map[string]*Extension)

	s.fileService = app.MustComponent[files.Service](a)
	s.fileObjectService = app.MustComponent[fileobject.Service](a)

	s.repoPath = app.MustComponent[wallet.Wallet](a).RepoPath()
	s.rootPath = app.MustComponent[wallet.Wallet](a).RootPath()
	s.downloadRootPath = filepath.Join(s.rootPath, "extensions", "download")
	s.extractedRootPath = filepath.Join(s.rootPath, "extensions", "extracted")

	if err := s.loadConfig(); err != nil {
		return err
	}

	for k, v := range s.extensions {
		manifest, err := s.parse(k)
		if err != nil {
			delete(s.extensions, k)
			continue
		}
		v.manifest = manifest

		if v.Enabled {
			if err := s.enable(v.ID); err != nil {
				delete(s.extensions, k)
				continue
			}
		}
	}

	if err := s.saveConfig(); err != nil {
		return err
	}

	return nil
}

func (s *service) loadConfig() error {
	s.lock.Lock()
	defer s.lock.Unlock()

	conf := filepath.Join(s.repoPath, "extensions.json")
	f, err := os.Open(conf)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return err
	}
	defer f.Close()

	return json.NewDecoder(f).Decode(&s.extensions)
}

func (s *service) saveConfig() error {
	s.lock.Lock()
	defer s.lock.Unlock()

	conf := filepath.Join(s.repoPath, "extensions.json")
	f, err := os.OpenFile(conf, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	return json.NewEncoder(f).Encode(s.extensions)
}

func (s *service) GetDeveloperMode(ctx context.Context) (bool, error) {
	return s.developerMode, nil
}

func (s *service) SetDeveloperMode(ctx context.Context, developerMode bool) error {
	log.Info("SetDeveloperMode", developerMode)
	s.lock.Lock()
	defer s.lock.Unlock()
	s.developerMode = developerMode
	return nil
}

func (s *service) ListInstalled(ctx context.Context) ([]*model.ExtensionInfo, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	var result = []*model.ExtensionInfo{}
	for _, v := range s.extensions {
		result = append(result, &model.ExtensionInfo{
			ExtensionId: v.ID,
			Enabled:     v.Enabled,
			Manifest:    v.manifest,
			// Functions: ,
		})
	}
	return result, nil
}

func (s *service) InstallFromURL(ctx context.Context, u string) (*model.ExtensionInfo, error) {
	log.Info("InstallFromURL", zap.String("url", u))

	if !s.developerMode {
		return nil, errors.New("not in developer mode")
	}

	filename, err := downloadFile(ctx, u, s.downloadRootPath)
	if err != nil {
		return nil, err
	}

	zipPath := filepath.Join(s.downloadRootPath, filename)
	extracted := filepath.Join(s.extractedRootPath, strings.TrimSuffix(filename, ".ext"))
	err = unzip(zipPath, extracted)
	if err != nil {
		return nil, err
	}

	return s.install(ctx, extracted)
}

func (s *service) InstallFromZip(ctx context.Context, zipPath string) (*model.ExtensionInfo, error) {
	filename := filepath.Base(zipPath)
	extracted := filepath.Join(s.extractedRootPath, strings.TrimSuffix(filename, ".ext"))
	err := unzip(zipPath, extracted)
	if err != nil {
		return nil, err
	}

	return s.install(ctx, extracted)
}

func (s *service) install(ctx context.Context, extracted string) (*model.ExtensionInfo, error) {
	// parse
	// move
	// save config
	// enable

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

	return &model.ExtensionInfo{}, nil
}

func (s *service) Enable(ctx context.Context, extensionId string) error {
	return s.enable(extensionId)
}

func (s *service) parse(extensionId string) (*model.ExtensionManifest, error) {
	return nil, errors.New("TODO")
}

func (s *service) enable(extensionId string) error {
	return errors.New("TODO")
}

func (s *service) Disable(ctx context.Context, extensionId string) error {
	return errors.New("TODO")
}

func (s *service) Call(ctx context.Context, extensionId, functionName, blockId string) error {
	// id, err := s.fileObjectService.GetFileIdFromObject(blockId)
	// if err != nil {
	// 	return fmt.Errorf("get file hash from object: %w", err)
	// }
	// image, err := s.fileService.ImageByHash(ctx, id)
	// if err != nil {
	// 	return err
	// }

	// f, err := image.GetOriginalFile()
	// if err != nil {
	// 	return err
	// }

	// r, err := f.Reader(ctx)
	// if err != nil {
	// 	return err
	// }

	// _ = r

	return nil
}
