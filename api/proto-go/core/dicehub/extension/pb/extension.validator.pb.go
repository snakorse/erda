// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: extension.proto

package pb

import (
	fmt "fmt"
	math "math"

	_ "github.com/erda-project/erda-proto-go/common/pb"
	proto "github.com/golang/protobuf/proto"
	github_com_mwitkow_go_proto_validators "github.com/mwitkow/go-proto-validators"
	_ "google.golang.org/genproto/googleapis/api/annotations"
	_ "google.golang.org/protobuf/types/known/structpb"
	_ "google.golang.org/protobuf/types/known/timestamppb"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

func (this *ExtensionSearchRequest) Validate() error {
	return nil
}
func (this *ExtensionSearchResponse) Validate() error {
	// Validation of proto3 map<> fields is unsupported.
	return nil
}
func (this *ExtensionCreateRequest) Validate() error {
	return nil
}
func (this *ExtensionCreateResponse) Validate() error {
	if this.Data != nil {
		if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(this.Data); err != nil {
			return github_com_mwitkow_go_proto_validators.FieldError("Data", err)
		}
	}
	return nil
}
func (this *QueryExtensionsRequest) Validate() error {
	return nil
}
func (this *QueryExtensionsResponse) Validate() error {
	if this.Data != nil {
		if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(this.Data); err != nil {
			return github_com_mwitkow_go_proto_validators.FieldError("Data", err)
		}
	}
	return nil
}
func (this *QueryExtensionsMenuRequest) Validate() error {
	return nil
}
func (this *QueryExtensionsMenuResponse) Validate() error {
	// Validation of proto3 map<> fields is unsupported.
	return nil
}
func (this *Spec) Validate() error {
	// Validation of proto3 map<> fields is unsupported.
	return nil
}
func (this *ExtensionVersionCreateRequest) Validate() error {
	return nil
}
func (this *ExtensionVersionCreateResponse) Validate() error {
	if this.Data != nil {
		if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(this.Data); err != nil {
			return github_com_mwitkow_go_proto_validators.FieldError("Data", err)
		}
	}
	return nil
}
func (this *GetExtensionVersionRequest) Validate() error {
	return nil
}
func (this *GetExtensionVersionResponse) Validate() error {
	if this.Data != nil {
		if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(this.Data); err != nil {
			return github_com_mwitkow_go_proto_validators.FieldError("Data", err)
		}
	}
	return nil
}
func (this *ExtensionQueryRequest) Validate() error {
	return nil
}
func (this *ExtensionQueryResponse) Validate() error {
	for _, item := range this.Data {
		if item != nil {
			if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(item); err != nil {
				return github_com_mwitkow_go_proto_validators.FieldError("Data", err)
			}
		}
	}
	return nil
}
func (this *ExtensionVersionGetRequest) Validate() error {
	return nil
}
func (this *ExtensionVersionQueryRequest) Validate() error {
	return nil
}
func (this *ExtensionVersionGetResponse) Validate() error {
	if this.Data != nil {
		if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(this.Data); err != nil {
			return github_com_mwitkow_go_proto_validators.FieldError("Data", err)
		}
	}
	return nil
}
func (this *ExtensionVersionQueryResponse) Validate() error {
	for _, item := range this.Data {
		if item != nil {
			if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(item); err != nil {
				return github_com_mwitkow_go_proto_validators.FieldError("Data", err)
			}
		}
	}
	return nil
}
func (this *ExtensionMenu) Validate() error {
	for _, item := range this.Items {
		if item != nil {
			if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(item); err != nil {
				return github_com_mwitkow_go_proto_validators.FieldError("Items", err)
			}
		}
	}
	return nil
}
func (this *Extension) Validate() error {
	if this.CreatedAt != nil {
		if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(this.CreatedAt); err != nil {
			return github_com_mwitkow_go_proto_validators.FieldError("CreatedAt", err)
		}
	}
	if this.UpdatedAt != nil {
		if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(this.UpdatedAt); err != nil {
			return github_com_mwitkow_go_proto_validators.FieldError("UpdatedAt", err)
		}
	}
	return nil
}
func (this *ExtensionVersion) Validate() error {
	if this.Spec != nil {
		if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(this.Spec); err != nil {
			return github_com_mwitkow_go_proto_validators.FieldError("Spec", err)
		}
	}
	if this.Dice != nil {
		if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(this.Dice); err != nil {
			return github_com_mwitkow_go_proto_validators.FieldError("Dice", err)
		}
	}
	if this.Swagger != nil {
		if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(this.Swagger); err != nil {
			return github_com_mwitkow_go_proto_validators.FieldError("Swagger", err)
		}
	}
	if this.CreatedAt != nil {
		if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(this.CreatedAt); err != nil {
			return github_com_mwitkow_go_proto_validators.FieldError("CreatedAt", err)
		}
	}
	if this.UpdatedAt != nil {
		if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(this.UpdatedAt); err != nil {
			return github_com_mwitkow_go_proto_validators.FieldError("UpdatedAt", err)
		}
	}
	return nil
}
