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

package table

import (
	"reflect"

	"github.com/erda-project/erda-infra/base/logs"
	"github.com/erda-project/erda-infra/base/servicehub"
	"github.com/erda-project/erda-infra/providers/component-protocol/components/table"
	"github.com/erda-project/erda-infra/providers/component-protocol/components/table/impl"
	"github.com/erda-project/erda-infra/providers/component-protocol/cpregister"
	"github.com/erda-project/erda-infra/providers/component-protocol/cptype"
	"github.com/erda-project/erda-infra/providers/component-protocol/protocol"
	"github.com/erda-project/erda-infra/providers/component-protocol/utils/cputil"
	"github.com/erda-project/erda-infra/providers/i18n"
	metricpb "github.com/erda-project/erda-proto-go/core/monitor/metric/pb"
	"github.com/erda-project/erda-proto-go/msp/apm/trace/pb"
	"github.com/erda-project/erda/modules/msp/apm/trace/query"
	"github.com/erda-project/erda/modules/msp/apm/trace/query/commom/custom"
	"github.com/erda-project/erda/modules/msp/apm/trace/query/commom/trace"
)

type provider struct {
	impl.DefaultTable
	custom.TraceInParams
	Log          logs.Logger
	I18n         i18n.Translator              `autowired:"i18n" translator:"msp-i18n"`
	TraceService *query.TraceService          `autowired:"erda.msp.apm.trace.TraceService"`
	Metric       metricpb.MetricServiceServer `autowired:"erda.core.monitor.metric.MetricService"`
}

// RegisterInitializeOp .
func (p *provider) RegisterInitializeOp() (opFunc cptype.OperationFunc) {
	return func(sdk *cptype.SDK) {
		params := p.TraceInParams.InParamsPtr
		pageNo, pageSize := trace.GetPagingFromGlobalState(*sdk.GlobalState)
		//sorts := trace.GetSortsFromGlobalState(*sdk.GlobalState)
		//fmt.Println(sorts)
		_, err := p.TraceService.GetTraces(sdk.Ctx, &pb.GetTracesRequest{
			TenantID:    params.TenantId,
			Status:      params.Status,
			StartTime:   params.StartTime,
			EndTime:     params.EndTime,
			Limit:       params.Limit,
			TraceID:     params.TraceId,
			DurationMin: params.DurationMin,
			DurationMax: params.DurationMax,
			Sort:        "TODO params.",
			ServiceName: params.ServiceName,
			RpcMethod:   params.RpcMethod,
			HttpPath:    params.HttpPath,
			PageNo:      int64(pageNo),
			PageSize:    int64(pageSize),
		})
		if err != nil {
			p.Log.Error(err)
			return
		}
		p.StdDataPtr = &table.Data{
			//Table: *data,
			Operations: map[cptype.OperationKey]cptype.Operation{
				table.OpTableChangePage{}.OpKey(): cputil.NewOpBuilder().WithServerDataPtr(&table.OpTableChangePageServerData{}).Build(),
				table.OpTableChangeSort{}.OpKey(): cputil.NewOpBuilder().Build(),
			}}
	}
}

func (p *provider) RegisterTablePagingOp(opData table.OpTableChangePage) (opFunc cptype.OperationFunc) {
	return func(sdk *cptype.SDK) {
		(*sdk.GlobalState)[trace.StateKeyTracePaging] = opData.ClientData
		p.RegisterInitializeOp()(sdk)
	}
}

func (p *provider) RegisterTableChangePageOp(opData table.OpTableChangePage) (opFunc cptype.OperationFunc) {
	return func(sdk *cptype.SDK) {
		(*sdk.GlobalState)[trace.StateKeyTracePaging] = opData.ClientData
		p.RegisterInitializeOp()(sdk)
	}
}

func (p *provider) RegisterTableSortOp(opData table.OpTableChangeSort) (opFunc cptype.OperationFunc) {
	return func(sdk *cptype.SDK) {
		(*sdk.GlobalState)[trace.StateKeyTraceSort] = opData.ClientData
		p.RegisterInitializeOp()(sdk)
	}
}

func (p *provider) RegisterBatchRowsHandleOp(opData table.OpBatchRowsHandle) (opFunc cptype.OperationFunc) {
	return nil
}

func (p *provider) RegisterRowSelectOp(opData table.OpRowSelect) (opFunc cptype.OperationFunc) {
	return nil
}

func (p *provider) RegisterRowAddOp(opData table.OpRowAdd) (opFunc cptype.OperationFunc) {
	return nil
}

func (p *provider) RegisterRowEditOp(opData table.OpRowEdit) (opFunc cptype.OperationFunc) {
	return nil
}

func (p *provider) RegisterRowDeleteOp(opData table.OpRowDelete) (opFunc cptype.OperationFunc) {
	return nil
}

// RegisterRenderingOp .
func (p *provider) RegisterRenderingOp() (opFunc cptype.OperationFunc) {
	return p.RegisterInitializeOp()
}

// Init .
func (p *provider) Init(ctx servicehub.Context) error {
	p.DefaultTable = impl.DefaultTable{}
	v := reflect.ValueOf(p)
	v.Elem().FieldByName("Impl").Set(v)
	compName := "table"
	if ctx.Label() != "" {
		compName = ctx.Label()
	}
	protocol.MustRegisterComponent(&protocol.CompRenderSpec{
		Scenario: "trace-query",
		CompName: compName,
		Creator:  func() cptype.IComponent { return p },
	})
	return nil
}

// Provide .
func (p *provider) Provide(ctx servicehub.DependencyContext, args ...interface{}) interface{} {
	return p
}

func init() {
	name := "component-protocol.components.trace-query.table"
	cpregister.AllExplicitProviderCreatorMap[name] = nil
	servicehub.Register(name, &servicehub.Spec{
		Creator: func() servicehub.Provider { return &provider{} },
	})
}
