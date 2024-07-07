package core

import (
	"context"

	"github.com/anyproto/anytype-heart/pb"
	extism "github.com/extism/go-sdk"
)

func (mw *Middleware) PluginExample(cctx context.Context, request *pb.RpcPluginExampleRequest) *pb.RpcPluginExampleResponse {
	manifest := extism.Manifest{
		Wasm: []extism.Wasm{
			extism.WasmFile{
				Path: request.Path,
			},
		},
	}
	config := extism.PluginConfig{}
	plugin, err := extism.NewPlugin(cctx, manifest, config, []extism.HostFunction{})
	if err != nil {
		return &pb.RpcPluginExampleResponse{
			Error: &pb.RpcPluginExampleResponseError{
				Code:        pb.RpcPluginExampleResponseError_BAD_INPUT,
				Description: err.Error(),
			},
		}
	}

	data := []byte("Hello, World!")
	_, out, err := plugin.Call("count_vowels", data)
	if err != nil {
		return &pb.RpcPluginExampleResponse{
			Error: &pb.RpcPluginExampleResponseError{
				Code:        pb.RpcPluginExampleResponseError_INTERNAL_ERROR,
				Description: err.Error(),
			},
		}
	}

	return &pb.RpcPluginExampleResponse{
		Error: &pb.RpcPluginExampleResponseError{
			Code: pb.RpcPluginExampleResponseError_NULL,
		},
		Result: string(out),
	}
}
