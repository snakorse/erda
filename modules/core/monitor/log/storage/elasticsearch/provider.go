// Copyright (c) 2021 Terminus, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package elasticsearch

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"bou.ke/monkey"

	"github.com/olivere/elastic"

	"github.com/erda-project/erda-infra/base/logs"
	"github.com/erda-project/erda-infra/base/servicehub"
	"github.com/erda-project/erda-infra/providers/elasticsearch"
	"github.com/erda-project/erda/modules/core/monitor/log"
	"github.com/erda-project/erda/modules/core/monitor/log/storage"
	retention "github.com/erda-project/erda/modules/core/monitor/settings/retention-strategy"
	"github.com/erda-project/erda/modules/core/monitor/storekit"
	"github.com/erda-project/erda/modules/core/monitor/storekit/elasticsearch/index/creator"
	"github.com/erda-project/erda/modules/core/monitor/storekit/elasticsearch/index/loader"
)

type (
	config struct {
		QueryTimeout        time.Duration `file:"query_timeout" default:"1m"`
		WriteTimeout        time.Duration `file:"write_timeout" default:"1m"`
		ReadPageSize        int           `file:"read_page_size" default:"1024"`
		IndexType           string        `file:"index_type" default:"log"`
		BenchMarking        bool          `file:"benchmarking" env:"BENCHMARKING_TEST" default:"true"`
		BenchMarkingTimeout time.Duration `file:"benchmarking_timeout" env:"BENCHMARKING_TIMEOUT" default:"200ms"`
		BenchMarkingLogging bool          `file:"benchmarking_logging" env:"BENCHMARKING_LOGGING" default:"false"`
	}
	provider struct {
		Cfg          *config
		Log          logs.Logger
		ES1          elasticsearch.Interface `autowired:"elasticsearch@log" optional:"true"`
		ES2          elasticsearch.Interface `autowired:"elasticsearch" optional:"true"`
		Loader       loader.Interface        `autowired:"elasticsearch.index.loader@log" optional:"true"`
		Creator      creator.Interface       `autowired:"elasticsearch.index.creator@log" optional:"true"`
		Retention    retention.Interface     `autowired:"storage-retention-strategy@log" optional:"true"`
		es           elasticsearch.Interface
		client       *elastic.Client
		queryTimeout string
	}
)

func (p *provider) Init(ctx servicehub.Context) (err error) {
	if p.ES1 != nil {
		p.es = p.ES1
	} else if p.ES2 != nil {
		p.es = p.ES2
	} else {
		return fmt.Errorf("elasticsearch is required")
	}
	p.client = p.es.Client()
	if p.Retention != nil {
		ctx.AddTask(func(c context.Context) error {
			p.Retention.Loading(ctx)
			return nil
		})
	}
	if p.Cfg.BenchMarking {
		monkey.Patch((*elastic.Client).PerformRequest, func(c *elastic.Client, ctx context.Context, opt elastic.PerformRequestOptions) (*elastic.Response, error) {
			if p.Cfg.BenchMarkingLogging {
				fmt.Printf("benchmarking about to write es: %+v", opt.Body)
			}
			if p.Cfg.BenchMarkingTimeout > 0 {
				time.Sleep(p.Cfg.BenchMarkingTimeout)
			}
			return &elastic.Response{
				StatusCode: 200,
				Header:     map[string][]string{"Content-Type": {"application/json"}},
				Body: []byte(`{
													   "took": 30,
													   "errors": false,
													   "items": []
													}`),
			}, nil
		})
	}
	return nil
}

var _ storage.Storage = (*provider)(nil)

func (p *provider) NewWriter(ctx context.Context) (storekit.BatchWriter, error) {
	if !p.Cfg.BenchMarking && (p.Creator == nil || p.Retention == nil) {
		return nil, fmt.Errorf("elasticsearch.index.creator@log and storage-retention-strategy@log is required for Writer")
	}
	w := p.es.NewWriter(&elasticsearch.WriteOptions{
		Timeout: p.Cfg.WriteTimeout,
		Enc: func(val interface{}) (index, id, typ string, body interface{}, err error) {
			data := val.(*log.LabeledLog)
			if p.Cfg.BenchMarking {
				index = "benchmarking_log"
			} else {
				var wait <-chan error
				if p.Retention == nil {
					wait, index = p.Creator.Ensure(data.Tags["dice_org_name"])
				} else {
					key := p.Retention.GetConfigKey(data.Source, data.Tags)
					if len(key) > 0 {
						wait, index = p.Creator.Ensure(data.Tags["dice_org_name"], key)
					} else {
						wait, index = p.Creator.Ensure(data.Tags["dice_org_name"])
					}
				}
				if wait != nil {
					select {
					case <-wait:
					case <-ctx.Done():
						return "", "", "", nil, storekit.ErrExitConsume
					}
				}
			}
			id = data.ID
			if len(id) > 12 {
				id = id[:12]
			}
			return index, id + strconv.FormatInt(data.Timestamp, 36) + "-" + strconv.FormatInt(data.Offset, 36), p.Cfg.IndexType, &Document{
				Log:  &data.Log,
				Date: getUnixMillisecond(data.Timestamp),
			}, nil
		},
	})
	return w, nil
}

// Document .
type Document struct {
	*log.Log
	Date int64 `json:"@timestamp"`
}

const maxUnixMillisecond int64 = 9999999999999

func getUnixMillisecond(ts int64) int64 {
	if ts > maxUnixMillisecond {
		return ts / int64(time.Millisecond)
	}
	return ts
}

func init() {
	servicehub.Register("log-storage-elasticsearch", &servicehub.Spec{
		Services:     []string{"log-storage-elasticsearch-reader", "log-storage-writer"},
		Dependencies: []string{"elasticsearch"},
		ConfigFunc:   func() interface{} { return &config{} },
		Creator:      func() servicehub.Provider { return &provider{} },
	})
}
