package core

import (
	"context"

	"github.com/anyproto/anytype-heart/core/extensions"
	"github.com/anyproto/anytype-heart/pb"
	"github.com/anyproto/anytype-heart/pkg/lib/pb/model"
)

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
	response := func(code pb.RpcExtensionSetDeveloperModeResponseErrorCode, err error) *pb.RpcExtensionSetDeveloperModeResponse {
		m := &pb.RpcExtensionSetDeveloperModeResponse{Error: &pb.RpcExtensionSetDeveloperModeResponseError{Code: code}}
		if err != nil {
			m.Error.Description = err.Error()
		}
		return m
	}

	err := getService[extensions.Service](mw).SetDeveloperMode(cctx, request.GetDeveloperMode())
	if err != nil {
		return response(pb.RpcExtensionSetDeveloperModeResponseError_UNKNOWN_ERROR, err)
	}
	return response(pb.RpcExtensionSetDeveloperModeResponseError_NULL, nil)
}

func (mw *Middleware) ExtensionListInstalled(cctx context.Context, request *pb.RpcExtensionListInstalledRequest) *pb.RpcExtensionListInstalledResponse {
	response := func(code pb.RpcExtensionListInstalledResponseErrorCode, err error, extensions []*model.ExtensionInfo) *pb.RpcExtensionListInstalledResponse {
		m := &pb.RpcExtensionListInstalledResponse{Error: &pb.RpcExtensionListInstalledResponseError{Code: code}}
		if err != nil {
			m.Error.Description = err.Error()
		}
		m.Extensions = extensions
		return m
	}

	extensions, err := getService[extensions.Service](mw).ListInstalled(cctx)
	if err != nil {
		return response(pb.RpcExtensionListInstalledResponseError_UNKNOWN_ERROR, err, nil)
	}
	return response(pb.RpcExtensionListInstalledResponseError_NULL, nil, extensions)
}

func (mw *Middleware) ExtensionInstallFromURL(cctx context.Context, request *pb.RpcExtensionInstallFromURLRequest) *pb.RpcExtensionInstallFromURLResponse {
	response := func(code pb.RpcExtensionInstallFromURLResponseErrorCode, err error, info *model.ExtensionInfo) *pb.RpcExtensionInstallFromURLResponse {
		m := &pb.RpcExtensionInstallFromURLResponse{Error: &pb.RpcExtensionInstallFromURLResponseError{Code: code}}
		if err != nil {
			m.Error.Description = err.Error()
		}
		m.Info = info
		return m
	}

	info, err := getService[extensions.Service](mw).InstallFromURL(cctx, request.GetUrl())
	if err != nil {
		return response(pb.RpcExtensionInstallFromURLResponseError_UNKNOWN_ERROR, err, nil)
	}
	return response(pb.RpcExtensionInstallFromURLResponseError_NULL, nil, info)
}

func (mw *Middleware) ExtensionEnable(cctx context.Context, request *pb.RpcExtensionEnableRequest) *pb.RpcExtensionEnableResponse {
	response := func(code pb.RpcExtensionEnableResponseErrorCode, err error) *pb.RpcExtensionEnableResponse {
		m := &pb.RpcExtensionEnableResponse{Error: &pb.RpcExtensionEnableResponseError{Code: code}}
		if err != nil {
			m.Error.Description = err.Error()
		}
		return m
	}

	err := getService[extensions.Service](mw).Enable(cctx, request.GetExtensionId())
	if err != nil {
		return response(pb.RpcExtensionEnableResponseError_UNKNOWN_ERROR, err)
	}
	return response(pb.RpcExtensionEnableResponseError_NULL, nil)
}

func (mw *Middleware) ExtensionDisable(cctx context.Context, request *pb.RpcExtensionDisableRequest) *pb.RpcExtensionDisableResponse {
	response := func(code pb.RpcExtensionDisableResponseErrorCode, err error) *pb.RpcExtensionDisableResponse {
		m := &pb.RpcExtensionDisableResponse{Error: &pb.RpcExtensionDisableResponseError{Code: code}}
		if err != nil {
			m.Error.Description = err.Error()
		}
		return m
	}

	err := getService[extensions.Service](mw).Disable(cctx, request.GetExtensionId())
	if err != nil {
		return response(pb.RpcExtensionDisableResponseError_UNKNOWN_ERROR, err)
	}
	return response(pb.RpcExtensionDisableResponseError_NULL, nil)
}

func (mw *Middleware) ExtensionCall(cctx context.Context, request *pb.RpcExtensionCallRequest) *pb.RpcExtensionCallResponse {
	response := func(code pb.RpcExtensionCallResponseErrorCode, err error) *pb.RpcExtensionCallResponse {
		m := &pb.RpcExtensionCallResponse{Error: &pb.RpcExtensionCallResponseError{Code: code}}
		if err != nil {
			m.Error.Description = err.Error()
		}
		return m
	}

	err := getService[extensions.Service](mw).Call(cctx, request.GetExtensionId(), request.GetFunctionName(), request.GetBlockId())
	if err != nil {
		return response(pb.RpcExtensionCallResponseError_UNKNOWN_ERROR, err)
	}
	return response(pb.RpcExtensionCallResponseError_NULL, nil)
}
