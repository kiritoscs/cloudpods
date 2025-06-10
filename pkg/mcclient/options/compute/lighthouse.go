package compute

import (
	"yunion.io/x/jsonutils"

	"yunion.io/x/onecloud/pkg/mcclient/options"
)

type LighthouseListOptions struct {
	options.BaseListOptions
}

func (opts *LighthouseListOptions) Params() (jsonutils.JSONObject, error) {
	return options.ListStructToParams(opts)
}

type LighthouseIdOption struct {
	ID string `help:"Lighthouse Id"`
}

func (opts *LighthouseIdOption) GetId() string {
	return opts.ID
}

func (opts *LighthouseIdOption) Params() (jsonutils.JSONObject, error) {
	return nil, nil
}
