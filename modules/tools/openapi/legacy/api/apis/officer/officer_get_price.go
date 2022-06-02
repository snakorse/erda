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

package officer

import "github.com/erda-project/erda/modules/tools/openapi/legacy/api/apis"

var OFFICER_GET_PRICE = apis.ApiSpec{
	Path:        "/api/resource/aliyun/price",
	BackendPath: "/api/resource/aliyun/price",
	Method:      "POST",
	Host:        "officer.marathon.l4lb.thisdcos.directory:9029",
	Scheme:      "http",
	CheckLogin:  true,
	Doc:         "summary: 获取阿里云指定资源价格",
}