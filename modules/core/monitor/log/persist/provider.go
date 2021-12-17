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

package persist

import (
	"context"
	"fmt"
	"time"

	"github.com/erda-project/erda-infra/base/logs"
	"github.com/erda-project/erda-infra/base/servicehub"
	"github.com/erda-project/erda-infra/providers/kafka"
	"github.com/erda-project/erda/modules/core/monitor/log/storage"
	"github.com/erda-project/erda/modules/core/monitor/storekit"
)

type (
	config struct {
		Input           kafka.BatchReaderConfig `file:"input"`
		Parallelism     int                     `file:"parallelism" default:"1"`
		BufferSize      int                     `file:"buffer_size" default:"1024"`
		ReadTimeout     time.Duration           `file:"read_timeout" default:"5s"`
		IDKeys          []string                `file:"id_keys"`
		PrintInvalidLog bool                    `file:"print_invalid_log" default:"false"`
		BenchMarking    bool                    `file:"benchmarking" env:"BENCH_MARKING_TEST" default:"true"`
	}
	provider struct {
		Cfg           *config
		Log           logs.Logger
		Kafka         kafka.Interface `autowired:"kafka"`
		StorageWriter storage.Storage `autowired:"log-storage-writer"`

		storage   storage.Storage
		stats     Statistics
		validator Validator
		metadata  MetadataProcessor
	}
)

func (p *provider) Init(ctx servicehub.Context) (err error) {

	p.validator = newValidator(p.Cfg)
	if runner, ok := p.validator.(servicehub.ProviderRunnerWithContext); ok {
		ctx.AddTask(runner.Run, servicehub.WithTaskName("log validator"))
	}

	p.metadata = newMetadataProcessor(p.Cfg)
	if runner, ok := p.metadata.(servicehub.ProviderRunnerWithContext); ok {
		ctx.AddTask(runner.Run, servicehub.WithTaskName("log metadata processor"))
	}

	if p.Cfg.BenchMarking {
		stats := NewInMemoryStatistics()
		stats.Start()
		p.stats = stats
	} else {
		p.stats = sharedStatistics
	}

	// add consumer task
	for i := 0; i < p.Cfg.Parallelism; i++ {
		ctx.AddTask(func(ctx context.Context) error {
			var r storekit.BatchReader
			if p.Cfg.BenchMarking {
				r = storekit.NewMockBatchReader(ctx, 0, p.Cfg.BufferSize, func() interface{} {
					result, _ := p.decodeLog([]byte(""), []byte(`{"source":"container","id":"c7a8ef453d7d264a7e07835e0e9c6ee36cb12429806fefbb84b8a62062fd091d","stream":"stdout","content":"2021-12-15 19:49:46.135 ERROR [,,] - [Druid-ConnectionPool-Create-1700023750] com.alibaba.druid.pool.DruidDataSource  : create connection RuntimeException\njava.lang.RuntimeException: Error instantiating JsonCustomSchema(name=csauto)\n\tat org.apache.calcite.model.ModelHandler.visit(ModelHandler.java:289)\n\tat org.apache.calcite.model.JsonCustomSchema.accept(JsonCustomSchema.java:45)\n\tat org.apache.calcite.model.ModelHandler.visit(ModelHandler.java:210)\n\tat org.apache.calcite.model.ModelHandler.\u003cinit\u003e(ModelHandler.java:101)\n\tat org.apache.calcite.jdbc.Driver$1.onConnectionInit(Driver.java:98)\n\tat org.apache.calcite.avatica.UnregisteredDriver.connect(UnregisteredDriver.java:139)\n\tat com.alibaba.druid.pool.DruidAbstractDataSource.createPhysicalConnection(DruidAbstractDataSource.java:1643)\n\tat com.alibaba.druid.pool.DruidAbstractDataSource.createPhysicalConnection(DruidAbstractDataSource.java:1709)\n\tat com.alibaba.druid.pool.DruidDataSource$CreateConnectionThread.run(DruidDataSource.java:2715)\nCaused by: java.lang.RuntimeException: Property 'org.apache.calcite.adapter.cassandra.CassandraSchemaFactory' not valid for plugin type org.apache.calcite.schema.SchemaFactory\n\tat org.apache.calcite.avatica.AvaticaUtils.instantiatePlugin(AvaticaUtils.java:241)\n\tat org.apache.calcite.model.ModelHandler.visit(ModelHandler.java:281)\n\t... 8 common frames omitted\nCaused by: java.lang.ClassNotFoundException: org.apache.calcite.adapter.cassandra.CassandraSchemaFactory\n\tat java.net.URLClassLoader.findClass(URLClassLoader.java:382)\n\tat java.lang.ClassLoader.loadClass(ClassLoader.java:419)\n\tat org.springframework.boot.loader.LaunchedURLClassLoader.loadClass(LaunchedURLClassLoader.java:94)\n\tat java.lang.ClassLoader.loadClass(ClassLoader.java:352)\n\tat java.lang.Class.forName0(Native Method)\n\tat java.lang.Class.forName(Class.java:264)\n\tat org.apache.calcite.avatica.AvaticaUtils.instantiatePlugin(AvaticaUtils.java:228)\n\t... 9 common frames omitted","offset":0,"timestamp":1639568986135580707,"tags":{"cluster_name":"terminus-test","component":"fdp-agent","container_id":"c7a8ef453d7d264a7e07835e0e9c6ee36cb12429806fefbb84b8a62062fd091d","container_name":"fdp-agent","dice_cluster_name":"terminus-test","dice_component":"fdp-agent","pod_id":"4aea3331-9e39-446b-8c48-3eac9a27852c","pod_ip":"10.125.36.172","pod_name":"fdp-fdp-agent-64f5db9b6d-s4sxk","pod_namespace":"default"},"labels":{}}`), nil, time.Now())
					return result
				})
			} else {
				r, err = p.Kafka.NewBatchReader(&p.Cfg.Input, kafka.WithReaderDecoder(p.decodeLog))
				if err != nil {
					return err
				}
				defer r.Close()
			}

			w, err := p.StorageWriter.NewWriter(ctx)
			if err != nil {
				return err
			}
			defer w.Close()
			return storekit.BatchConsume(ctx, r, w, &storekit.BatchConsumeOptions{
				BufferSize:          p.Cfg.BufferSize,
				ReadTimeout:         p.Cfg.ReadTimeout,
				ReadErrorHandler:    p.handleReadError,
				WriteErrorHandler:   p.handleWriteError,
				ConfirmErrorHandler: p.confirmErrorHandler,
				Statistics:          p.stats,
			})
		}, servicehub.WithTaskName(fmt.Sprintf("consumer(%d)", i)))
	}
	return nil
}

func init() {
	servicehub.Register("log-persist", &servicehub.Spec{
		ConfigFunc:   func() interface{} { return &config{} },
		Dependencies: []string{"kafka.topic.initializer"},
		Creator: func() servicehub.Provider {
			return &provider{}
		},
	})
}
