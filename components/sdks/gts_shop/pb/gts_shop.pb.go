// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        v4.25.0
// source: sdks/gts_shop/pb/gts_shop.proto

package pb

import (
	vo "components/sdks/gts_shop/pb/vo"
	_ "google.golang.org/genproto/googleapis/api/annotations"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

var File_sdks_gts_shop_pb_gts_shop_proto protoreflect.FileDescriptor

var file_sdks_gts_shop_pb_gts_shop_proto_rawDesc = []byte{
	0x0a, 0x1f, 0x73, 0x64, 0x6b, 0x73, 0x2f, 0x67, 0x74, 0x73, 0x5f, 0x73, 0x68, 0x6f, 0x70, 0x2f,
	0x70, 0x62, 0x2f, 0x67, 0x74, 0x73, 0x5f, 0x73, 0x68, 0x6f, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x12, 0x1b, 0x63, 0x6f, 0x6d, 0x70, 0x6f, 0x6e, 0x65, 0x6e, 0x74, 0x73, 0x2e, 0x73, 0x64,
	0x6b, 0x73, 0x2e, 0x67, 0x74, 0x73, 0x5f, 0x73, 0x68, 0x6f, 0x70, 0x2e, 0x70, 0x62, 0x1a, 0x1c,
	0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x22, 0x73, 0x64,
	0x6b, 0x73, 0x2f, 0x67, 0x74, 0x73, 0x5f, 0x73, 0x68, 0x6f, 0x70, 0x2f, 0x70, 0x62, 0x2f, 0x76,
	0x6f, 0x2f, 0x67, 0x74, 0x73, 0x5f, 0x73, 0x68, 0x6f, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x32, 0xc9, 0x0b, 0x0a, 0x07, 0x47, 0x74, 0x73, 0x53, 0x68, 0x6f, 0x70, 0x12, 0xa2, 0x01, 0x0a,
	0x08, 0x41, 0x75, 0x74, 0x68, 0x43, 0x6f, 0x64, 0x65, 0x12, 0x35, 0x2e, 0x63, 0x6f, 0x6d, 0x70,
	0x6f, 0x6e, 0x65, 0x6e, 0x74, 0x73, 0x2e, 0x73, 0x64, 0x6b, 0x73, 0x2e, 0x67, 0x74, 0x73, 0x5f,
	0x73, 0x68, 0x6f, 0x70, 0x2e, 0x70, 0x62, 0x2e, 0x41, 0x75, 0x74, 0x68, 0x6f, 0x72, 0x69, 0x7a,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x43, 0x6f, 0x64, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x1a, 0x2f, 0x2e, 0x63, 0x6f, 0x6d, 0x70, 0x6f, 0x6e, 0x65, 0x6e, 0x74, 0x73, 0x2e, 0x73, 0x64,
	0x6b, 0x73, 0x2e, 0x67, 0x74, 0x73, 0x5f, 0x73, 0x68, 0x6f, 0x70, 0x2e, 0x70, 0x62, 0x2e, 0x47,
	0x65, 0x74, 0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x22, 0x2e, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x28, 0x3a, 0x01, 0x2a, 0x22, 0x23, 0x2f, 0x62,
	0x61, 0x63, 0x6b, 0x65, 0x6e, 0x64, 0x2f, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x5f, 0x75, 0x73, 0x65,
	0x72, 0x2f, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x5f, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x69, 0x6e, 0x66,
	0x6f, 0x12, 0xce, 0x01, 0x0a, 0x14, 0x53, 0x65, 0x6e, 0x64, 0x56, 0x65, 0x72, 0x69, 0x66, 0x69,
	0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x43, 0x6f, 0x64, 0x65, 0x12, 0x38, 0x2e, 0x63, 0x6f, 0x6d,
	0x70, 0x6f, 0x6e, 0x65, 0x6e, 0x74, 0x73, 0x2e, 0x73, 0x64, 0x6b, 0x73, 0x2e, 0x67, 0x74, 0x73,
	0x5f, 0x73, 0x68, 0x6f, 0x70, 0x2e, 0x70, 0x62, 0x2e, 0x53, 0x65, 0x6e, 0x64, 0x56, 0x65, 0x72,
	0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x43, 0x6f, 0x64, 0x65, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x1a, 0x39, 0x2e, 0x63, 0x6f, 0x6d, 0x70, 0x6f, 0x6e, 0x65, 0x6e, 0x74,
	0x73, 0x2e, 0x73, 0x64, 0x6b, 0x73, 0x2e, 0x67, 0x74, 0x73, 0x5f, 0x73, 0x68, 0x6f, 0x70, 0x2e,
	0x70, 0x62, 0x2e, 0x53, 0x65, 0x6e, 0x64, 0x56, 0x65, 0x72, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x43, 0x6f, 0x64, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22,
	0x41, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x3b, 0x3a, 0x01, 0x2a, 0x22, 0x36, 0x2f, 0x62, 0x61, 0x63,
	0x6b, 0x65, 0x6e, 0x64, 0x2f, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x5f, 0x75, 0x73, 0x65, 0x72, 0x2f,
	0x7b, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x7d, 0x2f, 0x73, 0x65, 0x6e, 0x64, 0x5f, 0x76, 0x65,
	0x72, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x6d, 0x65, 0x73, 0x73, 0x61,
	0x67, 0x65, 0x12, 0x9b, 0x01, 0x0a, 0x0a, 0x56, 0x65, 0x72, 0x69, 0x66, 0x79, 0x43, 0x6f, 0x64,
	0x65, 0x12, 0x2e, 0x2e, 0x63, 0x6f, 0x6d, 0x70, 0x6f, 0x6e, 0x65, 0x6e, 0x74, 0x73, 0x2e, 0x73,
	0x64, 0x6b, 0x73, 0x2e, 0x67, 0x74, 0x73, 0x5f, 0x73, 0x68, 0x6f, 0x70, 0x2e, 0x70, 0x62, 0x2e,
	0x56, 0x65, 0x72, 0x69, 0x66, 0x79, 0x43, 0x6f, 0x64, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x1a, 0x2f, 0x2e, 0x63, 0x6f, 0x6d, 0x70, 0x6f, 0x6e, 0x65, 0x6e, 0x74, 0x73, 0x2e, 0x73,
	0x64, 0x6b, 0x73, 0x2e, 0x67, 0x74, 0x73, 0x5f, 0x73, 0x68, 0x6f, 0x70, 0x2e, 0x70, 0x62, 0x2e,
	0x56, 0x65, 0x72, 0x69, 0x66, 0x79, 0x43, 0x6f, 0x64, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x22, 0x2c, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x26, 0x3a, 0x01, 0x2a, 0x22, 0x21, 0x2f,
	0x62, 0x61, 0x63, 0x6b, 0x65, 0x6e, 0x64, 0x2f, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x5f, 0x75, 0x73,
	0x65, 0x72, 0x2f, 0x63, 0x68, 0x65, 0x63, 0x6b, 0x5f, 0x71, 0x72, 0x5f, 0x63, 0x6f, 0x64, 0x65,
	0x12, 0xaa, 0x01, 0x0a, 0x0b, 0x55, 0x73, 0x65, 0x72, 0x50, 0x72, 0x6f, 0x66, 0x69, 0x6c, 0x65,
	0x12, 0x2f, 0x2e, 0x63, 0x6f, 0x6d, 0x70, 0x6f, 0x6e, 0x65, 0x6e, 0x74, 0x73, 0x2e, 0x73, 0x64,
	0x6b, 0x73, 0x2e, 0x67, 0x74, 0x73, 0x5f, 0x73, 0x68, 0x6f, 0x70, 0x2e, 0x70, 0x62, 0x2e, 0x55,
	0x73, 0x65, 0x72, 0x50, 0x72, 0x6f, 0x66, 0x69, 0x6c, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x1a, 0x2f, 0x2e, 0x63, 0x6f, 0x6d, 0x70, 0x6f, 0x6e, 0x65, 0x6e, 0x74, 0x73, 0x2e, 0x73,
	0x64, 0x6b, 0x73, 0x2e, 0x67, 0x74, 0x73, 0x5f, 0x73, 0x68, 0x6f, 0x70, 0x2e, 0x70, 0x62, 0x2e,
	0x47, 0x65, 0x74, 0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x22, 0x39, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x33, 0x12, 0x31, 0x2f, 0x62, 0x61, 0x63,
	0x6b, 0x65, 0x6e, 0x64, 0x2f, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x5f, 0x75, 0x73, 0x65, 0x72, 0x2f,
	0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x5f, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x69, 0x6e, 0x66, 0x6f, 0x2f,
	0x7b, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x55, 0x73, 0x65, 0x72, 0x49, 0x64, 0x7d, 0x12, 0x9e, 0x01,
	0x0a, 0x09, 0x4c, 0x69, 0x73, 0x74, 0x41, 0x73, 0x73, 0x65, 0x74, 0x12, 0x2d, 0x2e, 0x63, 0x6f,
	0x6d, 0x70, 0x6f, 0x6e, 0x65, 0x6e, 0x74, 0x73, 0x2e, 0x73, 0x64, 0x6b, 0x73, 0x2e, 0x67, 0x74,
	0x73, 0x5f, 0x73, 0x68, 0x6f, 0x70, 0x2e, 0x70, 0x62, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x41, 0x73,
	0x73, 0x65, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x2e, 0x2e, 0x63, 0x6f, 0x6d,
	0x70, 0x6f, 0x6e, 0x65, 0x6e, 0x74, 0x73, 0x2e, 0x73, 0x64, 0x6b, 0x73, 0x2e, 0x67, 0x74, 0x73,
	0x5f, 0x73, 0x68, 0x6f, 0x70, 0x2e, 0x70, 0x62, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x41, 0x73, 0x73,
	0x65, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x32, 0x82, 0xd3, 0xe4, 0x93,
	0x02, 0x2c, 0x12, 0x2a, 0x2f, 0x62, 0x61, 0x63, 0x6b, 0x65, 0x6e, 0x64, 0x2f, 0x75, 0x73, 0x65,
	0x72, 0x5f, 0x67, 0x70, 0x2f, 0x77, 0x61, 0x6c, 0x6c, 0x65, 0x74, 0x2f, 0x71, 0x75, 0x65, 0x72,
	0x79, 0x2f, 0x7b, 0x72, 0x65, 0x61, 0x6c, 0x55, 0x73, 0x65, 0x72, 0x49, 0x64, 0x7d, 0x12, 0x8d,
	0x01, 0x0a, 0x08, 0x41, 0x64, 0x64, 0x41, 0x73, 0x73, 0x65, 0x74, 0x12, 0x2c, 0x2e, 0x63, 0x6f,
	0x6d, 0x70, 0x6f, 0x6e, 0x65, 0x6e, 0x74, 0x73, 0x2e, 0x73, 0x64, 0x6b, 0x73, 0x2e, 0x67, 0x74,
	0x73, 0x5f, 0x73, 0x68, 0x6f, 0x70, 0x2e, 0x70, 0x62, 0x2e, 0x41, 0x64, 0x64, 0x41, 0x73, 0x73,
	0x65, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x2d, 0x2e, 0x63, 0x6f, 0x6d, 0x70,
	0x6f, 0x6e, 0x65, 0x6e, 0x74, 0x73, 0x2e, 0x73, 0x64, 0x6b, 0x73, 0x2e, 0x67, 0x74, 0x73, 0x5f,
	0x73, 0x68, 0x6f, 0x70, 0x2e, 0x70, 0x62, 0x2e, 0x41, 0x64, 0x64, 0x41, 0x73, 0x73, 0x65, 0x74,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x24, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x1e,
	0x3a, 0x01, 0x2a, 0x22, 0x19, 0x2f, 0x62, 0x61, 0x63, 0x6b, 0x65, 0x6e, 0x64, 0x2f, 0x75, 0x73,
	0x65, 0x72, 0x5f, 0x67, 0x70, 0x2f, 0x69, 0x6e, 0x63, 0x72, 0x65, 0x61, 0x73, 0x65, 0x12, 0x95,
	0x01, 0x0a, 0x0b, 0x46, 0x72, 0x65, 0x65, 0x7a, 0x65, 0x41, 0x73, 0x73, 0x65, 0x74, 0x12, 0x2f,
	0x2e, 0x63, 0x6f, 0x6d, 0x70, 0x6f, 0x6e, 0x65, 0x6e, 0x74, 0x73, 0x2e, 0x73, 0x64, 0x6b, 0x73,
	0x2e, 0x67, 0x74, 0x73, 0x5f, 0x73, 0x68, 0x6f, 0x70, 0x2e, 0x70, 0x62, 0x2e, 0x46, 0x72, 0x65,
	0x65, 0x7a, 0x65, 0x41, 0x73, 0x73, 0x65, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a,
	0x30, 0x2e, 0x63, 0x6f, 0x6d, 0x70, 0x6f, 0x6e, 0x65, 0x6e, 0x74, 0x73, 0x2e, 0x73, 0x64, 0x6b,
	0x73, 0x2e, 0x67, 0x74, 0x73, 0x5f, 0x73, 0x68, 0x6f, 0x70, 0x2e, 0x70, 0x62, 0x2e, 0x46, 0x72,
	0x65, 0x65, 0x7a, 0x65, 0x41, 0x73, 0x73, 0x65, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x22, 0x23, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x1d, 0x3a, 0x01, 0x2a, 0x22, 0x18, 0x2f, 0x62,
	0x61, 0x63, 0x6b, 0x65, 0x6e, 0x64, 0x2f, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x67, 0x70, 0x2f, 0x68,
	0x6f, 0x6c, 0x64, 0x69, 0x6e, 0x67, 0x12, 0x94, 0x01, 0x0a, 0x08, 0x53, 0x75, 0x62, 0x41, 0x73,
	0x73, 0x65, 0x74, 0x12, 0x2c, 0x2e, 0x63, 0x6f, 0x6d, 0x70, 0x6f, 0x6e, 0x65, 0x6e, 0x74, 0x73,
	0x2e, 0x73, 0x64, 0x6b, 0x73, 0x2e, 0x67, 0x74, 0x73, 0x5f, 0x73, 0x68, 0x6f, 0x70, 0x2e, 0x70,
	0x62, 0x2e, 0x53, 0x75, 0x62, 0x41, 0x73, 0x73, 0x65, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x1a, 0x2d, 0x2e, 0x63, 0x6f, 0x6d, 0x70, 0x6f, 0x6e, 0x65, 0x6e, 0x74, 0x73, 0x2e, 0x73,
	0x64, 0x6b, 0x73, 0x2e, 0x67, 0x74, 0x73, 0x5f, 0x73, 0x68, 0x6f, 0x70, 0x2e, 0x70, 0x62, 0x2e,
	0x53, 0x75, 0x62, 0x41, 0x73, 0x73, 0x65, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x22, 0x2b, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x25, 0x3a, 0x01, 0x2a, 0x22, 0x20, 0x2f, 0x62, 0x61,
	0x63, 0x6b, 0x65, 0x6e, 0x64, 0x2f, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x67, 0x70, 0x2f, 0x63, 0x6f,
	0x6e, 0x66, 0x69, 0x72, 0x6d, 0x5f, 0x68, 0x6f, 0x6c, 0x64, 0x69, 0x6e, 0x67, 0x12, 0x9c, 0x01,
	0x0a, 0x0b, 0x52, 0x65, 0x74, 0x75, 0x72, 0x6e, 0x41, 0x73, 0x73, 0x65, 0x74, 0x12, 0x2f, 0x2e,
	0x63, 0x6f, 0x6d, 0x70, 0x6f, 0x6e, 0x65, 0x6e, 0x74, 0x73, 0x2e, 0x73, 0x64, 0x6b, 0x73, 0x2e,
	0x67, 0x74, 0x73, 0x5f, 0x73, 0x68, 0x6f, 0x70, 0x2e, 0x70, 0x62, 0x2e, 0x52, 0x65, 0x74, 0x75,
	0x72, 0x6e, 0x41, 0x73, 0x73, 0x65, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x30,
	0x2e, 0x63, 0x6f, 0x6d, 0x70, 0x6f, 0x6e, 0x65, 0x6e, 0x74, 0x73, 0x2e, 0x73, 0x64, 0x6b, 0x73,
	0x2e, 0x67, 0x74, 0x73, 0x5f, 0x73, 0x68, 0x6f, 0x70, 0x2e, 0x70, 0x62, 0x2e, 0x52, 0x65, 0x74,
	0x75, 0x72, 0x6e, 0x41, 0x73, 0x73, 0x65, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x22, 0x2a, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x24, 0x3a, 0x01, 0x2a, 0x22, 0x1f, 0x2f, 0x62, 0x61,
	0x63, 0x6b, 0x65, 0x6e, 0x64, 0x2f, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x67, 0x70, 0x2f, 0x72, 0x65,
	0x74, 0x75, 0x72, 0x6e, 0x5f, 0x68, 0x6f, 0x6c, 0x64, 0x69, 0x6e, 0x67, 0x42, 0x20, 0x5a, 0x1e,
	0x63, 0x6f, 0x6d, 0x70, 0x6f, 0x6e, 0x65, 0x6e, 0x74, 0x73, 0x2f, 0x73, 0x64, 0x6b, 0x73, 0x2f,
	0x67, 0x74, 0x73, 0x5f, 0x73, 0x68, 0x6f, 0x70, 0x2f, 0x70, 0x62, 0x3b, 0x70, 0x62, 0x62, 0x06,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var file_sdks_gts_shop_pb_gts_shop_proto_goTypes = []interface{}{
	(*vo.AuthorizationCodeRequest)(nil),     // 0: components.sdks.gts_shop.pb.AuthorizationCodeRequest
	(*vo.SendVerificationCodeRequest)(nil),  // 1: components.sdks.gts_shop.pb.SendVerificationCodeRequest
	(*vo.VerifyCodeRequest)(nil),            // 2: components.sdks.gts_shop.pb.VerifyCodeRequest
	(*vo.UserProfileRequest)(nil),           // 3: components.sdks.gts_shop.pb.UserProfileRequest
	(*vo.ListAssetRequest)(nil),             // 4: components.sdks.gts_shop.pb.ListAssetRequest
	(*vo.AddAssetRequest)(nil),              // 5: components.sdks.gts_shop.pb.AddAssetRequest
	(*vo.FreezeAssetRequest)(nil),           // 6: components.sdks.gts_shop.pb.FreezeAssetRequest
	(*vo.SubAssetRequest)(nil),              // 7: components.sdks.gts_shop.pb.SubAssetRequest
	(*vo.ReturnAssetRequest)(nil),           // 8: components.sdks.gts_shop.pb.ReturnAssetRequest
	(*vo.GetAccountResponse)(nil),           // 9: components.sdks.gts_shop.pb.GetAccountResponse
	(*vo.SendVerificationCodeResponse)(nil), // 10: components.sdks.gts_shop.pb.SendVerificationCodeResponse
	(*vo.VerifyCodeResponse)(nil),           // 11: components.sdks.gts_shop.pb.VerifyCodeResponse
	(*vo.ListAssetResponse)(nil),            // 12: components.sdks.gts_shop.pb.ListAssetResponse
	(*vo.AddAssetResponse)(nil),             // 13: components.sdks.gts_shop.pb.AddAssetResponse
	(*vo.FreezeAssetResponse)(nil),          // 14: components.sdks.gts_shop.pb.FreezeAssetResponse
	(*vo.SubAssetResponse)(nil),             // 15: components.sdks.gts_shop.pb.SubAssetResponse
	(*vo.ReturnAssetResponse)(nil),          // 16: components.sdks.gts_shop.pb.ReturnAssetResponse
}
var file_sdks_gts_shop_pb_gts_shop_proto_depIdxs = []int32{
	0,  // 0: components.sdks.gts_shop.pb.GtsShop.AuthCode:input_type -> components.sdks.gts_shop.pb.AuthorizationCodeRequest
	1,  // 1: components.sdks.gts_shop.pb.GtsShop.SendVerificationCode:input_type -> components.sdks.gts_shop.pb.SendVerificationCodeRequest
	2,  // 2: components.sdks.gts_shop.pb.GtsShop.VerifyCode:input_type -> components.sdks.gts_shop.pb.VerifyCodeRequest
	3,  // 3: components.sdks.gts_shop.pb.GtsShop.UserProfile:input_type -> components.sdks.gts_shop.pb.UserProfileRequest
	4,  // 4: components.sdks.gts_shop.pb.GtsShop.ListAsset:input_type -> components.sdks.gts_shop.pb.ListAssetRequest
	5,  // 5: components.sdks.gts_shop.pb.GtsShop.AddAsset:input_type -> components.sdks.gts_shop.pb.AddAssetRequest
	6,  // 6: components.sdks.gts_shop.pb.GtsShop.FreezeAsset:input_type -> components.sdks.gts_shop.pb.FreezeAssetRequest
	7,  // 7: components.sdks.gts_shop.pb.GtsShop.SubAsset:input_type -> components.sdks.gts_shop.pb.SubAssetRequest
	8,  // 8: components.sdks.gts_shop.pb.GtsShop.ReturnAsset:input_type -> components.sdks.gts_shop.pb.ReturnAssetRequest
	9,  // 9: components.sdks.gts_shop.pb.GtsShop.AuthCode:output_type -> components.sdks.gts_shop.pb.GetAccountResponse
	10, // 10: components.sdks.gts_shop.pb.GtsShop.SendVerificationCode:output_type -> components.sdks.gts_shop.pb.SendVerificationCodeResponse
	11, // 11: components.sdks.gts_shop.pb.GtsShop.VerifyCode:output_type -> components.sdks.gts_shop.pb.VerifyCodeResponse
	9,  // 12: components.sdks.gts_shop.pb.GtsShop.UserProfile:output_type -> components.sdks.gts_shop.pb.GetAccountResponse
	12, // 13: components.sdks.gts_shop.pb.GtsShop.ListAsset:output_type -> components.sdks.gts_shop.pb.ListAssetResponse
	13, // 14: components.sdks.gts_shop.pb.GtsShop.AddAsset:output_type -> components.sdks.gts_shop.pb.AddAssetResponse
	14, // 15: components.sdks.gts_shop.pb.GtsShop.FreezeAsset:output_type -> components.sdks.gts_shop.pb.FreezeAssetResponse
	15, // 16: components.sdks.gts_shop.pb.GtsShop.SubAsset:output_type -> components.sdks.gts_shop.pb.SubAssetResponse
	16, // 17: components.sdks.gts_shop.pb.GtsShop.ReturnAsset:output_type -> components.sdks.gts_shop.pb.ReturnAssetResponse
	9,  // [9:18] is the sub-list for method output_type
	0,  // [0:9] is the sub-list for method input_type
	0,  // [0:0] is the sub-list for extension type_name
	0,  // [0:0] is the sub-list for extension extendee
	0,  // [0:0] is the sub-list for field type_name
}

func init() { file_sdks_gts_shop_pb_gts_shop_proto_init() }
func file_sdks_gts_shop_pb_gts_shop_proto_init() {
	if File_sdks_gts_shop_pb_gts_shop_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_sdks_gts_shop_pb_gts_shop_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   0,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_sdks_gts_shop_pb_gts_shop_proto_goTypes,
		DependencyIndexes: file_sdks_gts_shop_pb_gts_shop_proto_depIdxs,
	}.Build()
	File_sdks_gts_shop_pb_gts_shop_proto = out.File
	file_sdks_gts_shop_pb_gts_shop_proto_rawDesc = nil
	file_sdks_gts_shop_pb_gts_shop_proto_goTypes = nil
	file_sdks_gts_shop_pb_gts_shop_proto_depIdxs = nil
}
