package model

import (
	"encoding/json"
	"errors"
)

func (p *ExtensionManifestPlatformType) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	switch s {
	case "ios":
		*p = ExtensionManifestPlatformType_PLATFORM_IOS
	case "android":
		*p = ExtensionManifestPlatformType_PLATFORM_ANDROID
	case "linux":
		*p = ExtensionManifestPlatformType_PLATFORM_LINUX
	case "windows":
		*p = ExtensionManifestPlatformType_PLATFORM_WINDOWS
	case "darwin":
		*p = ExtensionManifestPlatformType_PLATFORM_DARWIN
	default:
		return errors.New("invalid platform type")
	}

	return nil
}

func (p *ExtensionManifestScopeType) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	switch s {
	case "space":
		*p = ExtensionManifestScopeType_SCOPE_SPACE
	case "set":
		*p = ExtensionManifestScopeType_SCOPE_SET
	case "page":
		*p = ExtensionManifestScopeType_SCOPE_PAGE
	case "object":
		*p = ExtensionManifestScopeType_SCOPE_OBJECT
	default:
		return errors.New("invalid scope type")
	}

	return nil
}
