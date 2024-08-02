package extensions

import (
	"archive/zip"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/anyproto/any-sync/app"
	"github.com/anyproto/anytype-heart/core/files"
	"github.com/anyproto/anytype-heart/core/files/fileobject"
	"github.com/anyproto/anytype-heart/core/wallet"
	"github.com/anyproto/anytype-heart/pkg/lib/logging"
	"github.com/anyproto/anytype-heart/pkg/lib/pb/model"
	extism "github.com/extism/go-sdk"
	"github.com/google/uuid"
)

const CName = "extensions"

const (
	manifestFileName   = "manifest.json"
	extensionsFileName = "extensions.json"
)

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
	Enabled  bool                     `json:"enabled"`
	core     *extism.Plugin           `json:"-"`
	manifest *model.ExtensionManifest `json:"-"`
}

type service struct {
	lock sync.RWMutex

	developerMode bool

	extensions map[string]*Extension

	repoPath           string
	rootPath           string
	downloadRootPath   string
	extensionsRootPath string

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
	s.extensionsRootPath = filepath.Join(s.repoPath, "extensions")

	if err := os.MkdirAll(s.downloadRootPath, os.ModePerm); err != nil {
		return err
	}

	if err := os.MkdirAll(s.extensionsRootPath, os.ModePerm); err != nil {
		return err
	}

	if err := s.loadConfig(); err != nil {
		return err
	}

	for k, v := range s.extensions {
		manifest, err := s.parse(filepath.Join(s.extensionsRootPath, k, manifestFileName))
		if err != nil {
			delete(s.extensions, k)
			continue
		}
		v.manifest = manifest

		if v.Enabled {
			if err := s.load(k); err != nil {
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
	conf := filepath.Join(s.repoPath, extensionsFileName)
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
	conf := filepath.Join(s.repoPath, extensionsFileName)
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
			Enabled:  v.Enabled,
			Manifest: v.manifest,
			// Functions: ,
		})
	}
	return result, nil
}

func (s *service) InstallFromURL(ctx context.Context, u string) (*model.ExtensionInfo, error) {
	targetPath := filepath.Join(s.downloadRootPath, uuid.New().String()+".dl")
	defer os.Remove(targetPath)

	err := downloadFile(ctx, u, targetPath)
	if err != nil {
		return nil, err
	}

	return s.installFromZip(ctx, targetPath)
}

func (s *service) InstallFromZip(ctx context.Context, zipPath string) (*model.ExtensionInfo, error) {
	return s.installFromZip(ctx, zipPath)
}

func (s *service) installFromZip(ctx context.Context, zipPath string) (*model.ExtensionInfo, error) {
	var manifest model.ExtensionManifest
	err := readSingleFromZip(zipPath, func(f *zip.File) error {
		if f.Name != manifestFileName {
			return nil
		}

		rc, e := f.Open()
		if e != nil {
			return e
		}
		defer rc.Close()

		return json.NewDecoder(rc).Decode(&manifest)
	})
	if err != nil {
		return nil, err
	}

	if !isValidUUID(manifest.Id) {
		err = fmt.Errorf("invalid id %s", manifest.Id)
		return nil, err
	}

	extracted := filepath.Join(s.extensionsRootPath, manifest.Id)

	err = unzip(zipPath, extracted)
	if err != nil {
		return nil, err
	}

	return s.install(ctx, manifest.Id)
}

func (s *service) install(ctx context.Context, extensionId string) (*model.ExtensionInfo, error) {
	manifest, err := s.parse(filepath.Join(s.extensionsRootPath, extensionId, manifestFileName))
	if err != nil {
		return nil, err
	}

	extension := &Extension{
		Enabled:  true,
		manifest: manifest,
	}

	s.lock.Lock()
	s.extensions[manifest.Id] = extension
	s.lock.Unlock()

	_ = s.saveConfig()

	err = s.load(extensionId)
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

	return &model.ExtensionInfo{
		Enabled:  extension.Enabled,
		Manifest: extension.manifest,
		// Functions: ,
	}, nil
}

func (s *service) Enable(ctx context.Context, extensionId string) error {
	if _, ok := s.extensions[extensionId]; !ok {
		err := errors.New("extension not exist")
		return err
	}

	s.lock.Lock()
	s.extensions[extensionId].Enabled = true
	s.lock.Unlock()

	_ = s.saveConfig()

	return s.load(extensionId)
}

func (s *service) load(extensionId string) error {
	return errors.New("TODO")
}

func (s *service) Disable(ctx context.Context, extensionId string) error {
	if _, ok := s.extensions[extensionId]; !ok {
		err := errors.New("extension not exist")
		return err
	}

	s.lock.Lock()
	s.extensions[extensionId].Enabled = false
	s.lock.Unlock()

	_ = s.saveConfig()

	return s.unload(extensionId)
}

func (s *service) unload(extensionId string) error {
	return errors.New("TODO")
}

func (s *service) parse(manifestFile string) (*model.ExtensionManifest, error) {
	f, err := os.Open(manifestFile)
	if err != nil {
		return nil, err
	}

	var manifest model.ExtensionManifest
	err = json.NewDecoder(f).Decode(&manifest)
	if err != nil {
		return nil, err
	}

	return &manifest, nil
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
