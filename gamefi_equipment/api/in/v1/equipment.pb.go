// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        v4.25.0
// source: api/in/v1/equipment.proto

package v1

import (
	vo "gamefi_equipment/api/in/v1/vo"
	_ "github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2/options"
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

var File_api_in_v1_equipment_proto protoreflect.FileDescriptor

var file_api_in_v1_equipment_proto_rawDesc = []byte{
	0x0a, 0x19, 0x61, 0x70, 0x69, 0x2f, 0x69, 0x6e, 0x2f, 0x76, 0x31, 0x2f, 0x65, 0x71, 0x75, 0x69,
	0x70, 0x6d, 0x65, 0x6e, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x09, 0x61, 0x70, 0x69,
	0x2e, 0x69, 0x6e, 0x2e, 0x76, 0x31, 0x1a, 0x1c, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x61,
	0x70, 0x69, 0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1c, 0x61, 0x70, 0x69, 0x2f, 0x69, 0x6e, 0x2f, 0x76, 0x31, 0x2f,
	0x76, 0x6f, 0x2f, 0x65, 0x71, 0x75, 0x69, 0x70, 0x6d, 0x65, 0x6e, 0x74, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x1a, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x2d, 0x67, 0x65, 0x6e, 0x2d, 0x6f,
	0x70, 0x65, 0x6e, 0x61, 0x70, 0x69, 0x76, 0x32, 0x2f, 0x6f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73,
	0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x32, 0xa3, 0x0a, 0x0a, 0x10, 0x45, 0x71, 0x75, 0x69, 0x70, 0x6d, 0x65, 0x6e, 0x74,
	0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x94, 0x01, 0x0a, 0x0c, 0x41, 0x64, 0x64, 0x45,
	0x71, 0x75, 0x69, 0x70, 0x6d, 0x65, 0x6e, 0x74, 0x12, 0x1e, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x69,
	0x6e, 0x2e, 0x76, 0x31, 0x2e, 0x41, 0x64, 0x64, 0x45, 0x71, 0x75, 0x69, 0x70, 0x6d, 0x65, 0x6e,
	0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1f, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x69,
	0x6e, 0x2e, 0x76, 0x31, 0x2e, 0x41, 0x64, 0x64, 0x45, 0x71, 0x75, 0x69, 0x70, 0x6d, 0x65, 0x6e,
	0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x43, 0x92, 0x41, 0x1c, 0x12, 0x1a,
	0xe6, 0xb7, 0xbb, 0xe5, 0x8a, 0xa0, 0xe8, 0xa3, 0x85, 0xe5, 0xa4, 0x87, 0x28, 0xe6, 0x8e, 0x89,
	0xe8, 0x90, 0xbd, 0xe8, 0xa3, 0x85, 0xe5, 0xa4, 0x87, 0x29, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x1e,
	0x3a, 0x01, 0x2a, 0x22, 0x19, 0x2f, 0x76, 0x31, 0x2f, 0x69, 0x6e, 0x2f, 0x75, 0x73, 0x65, 0x72,
	0x2f, 0x65, 0x71, 0x75, 0x69, 0x70, 0x6d, 0x65, 0x6e, 0x74, 0x2d, 0x61, 0x64, 0x64, 0x12, 0x96,
	0x01, 0x0a, 0x10, 0x55, 0x70, 0x67, 0x72, 0x61, 0x64, 0x65, 0x45, 0x71, 0x75, 0x69, 0x70, 0x6d,
	0x65, 0x6e, 0x74, 0x12, 0x22, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x69, 0x6e, 0x2e, 0x76, 0x31, 0x2e,
	0x55, 0x70, 0x67, 0x72, 0x61, 0x64, 0x65, 0x45, 0x71, 0x75, 0x69, 0x70, 0x6d, 0x65, 0x6e, 0x74,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x23, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x69, 0x6e,
	0x2e, 0x76, 0x31, 0x2e, 0x55, 0x70, 0x67, 0x72, 0x61, 0x64, 0x65, 0x45, 0x71, 0x75, 0x69, 0x70,
	0x6d, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x39, 0x92, 0x41,
	0x0e, 0x12, 0x0c, 0xe5, 0xbc, 0xba, 0xe5, 0x8c, 0x96, 0xe8, 0xa3, 0x85, 0xe5, 0xa4, 0x87, 0x82,
	0xd3, 0xe4, 0x93, 0x02, 0x22, 0x3a, 0x01, 0x2a, 0x22, 0x1d, 0x2f, 0x76, 0x31, 0x2f, 0x69, 0x6e,
	0x2f, 0x75, 0x73, 0x65, 0x72, 0x2f, 0x65, 0x71, 0x75, 0x69, 0x70, 0x6d, 0x65, 0x6e, 0x74, 0x2d,
	0x75, 0x70, 0x67, 0x72, 0x61, 0x64, 0x65, 0x12, 0x97, 0x01, 0x0a, 0x11, 0x41, 0x64, 0x64, 0x46,
	0x69, 0x67, 0x68, 0x74, 0x45, 0x71, 0x75, 0x69, 0x70, 0x6d, 0x65, 0x6e, 0x74, 0x12, 0x23, 0x2e,
	0x61, 0x70, 0x69, 0x2e, 0x69, 0x6e, 0x2e, 0x76, 0x31, 0x2e, 0x41, 0x64, 0x64, 0x46, 0x69, 0x67,
	0x68, 0x74, 0x45, 0x71, 0x75, 0x69, 0x70, 0x6d, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x1a, 0x24, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x69, 0x6e, 0x2e, 0x76, 0x31, 0x2e, 0x41,
	0x64, 0x64, 0x46, 0x69, 0x67, 0x68, 0x74, 0x45, 0x71, 0x75, 0x69, 0x70, 0x6d, 0x65, 0x6e, 0x74,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x37, 0x92, 0x41, 0x0b, 0x12, 0x09, 0xe7,
	0xa9, 0xbf, 0xe8, 0xa3, 0x85, 0xe5, 0xa4, 0x87, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x23, 0x3a, 0x01,
	0x2a, 0x22, 0x1e, 0x2f, 0x76, 0x31, 0x2f, 0x69, 0x6e, 0x2f, 0x75, 0x73, 0x65, 0x72, 0x2f, 0x65,
	0x71, 0x75, 0x69, 0x70, 0x6d, 0x65, 0x6e, 0x74, 0x2d, 0x61, 0x64, 0x64, 0x66, 0x69, 0x67, 0x68,
	0x74, 0x12, 0x9f, 0x01, 0x0a, 0x13, 0x43, 0x6c, 0x65, 0x61, 0x72, 0x46, 0x69, 0x67, 0x68, 0x74,
	0x45, 0x71, 0x75, 0x69, 0x70, 0x6d, 0x65, 0x6e, 0x74, 0x12, 0x25, 0x2e, 0x61, 0x70, 0x69, 0x2e,
	0x69, 0x6e, 0x2e, 0x76, 0x31, 0x2e, 0x43, 0x6c, 0x65, 0x61, 0x72, 0x46, 0x69, 0x67, 0x68, 0x74,
	0x45, 0x71, 0x75, 0x69, 0x70, 0x6d, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x1a, 0x26, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x69, 0x6e, 0x2e, 0x76, 0x31, 0x2e, 0x43, 0x6c, 0x65,
	0x61, 0x72, 0x46, 0x69, 0x67, 0x68, 0x74, 0x45, 0x71, 0x75, 0x69, 0x70, 0x6d, 0x65, 0x6e, 0x74,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x39, 0x92, 0x41, 0x0b, 0x12, 0x09, 0xe8,
	0x84, 0xb1, 0xe8, 0xa3, 0x85, 0xe5, 0xa4, 0x87, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x25, 0x3a, 0x01,
	0x2a, 0x22, 0x20, 0x2f, 0x76, 0x31, 0x2f, 0x69, 0x6e, 0x2f, 0x75, 0x73, 0x65, 0x72, 0x2f, 0x65,
	0x71, 0x75, 0x69, 0x70, 0x6d, 0x65, 0x6e, 0x74, 0x2d, 0x63, 0x6c, 0x65, 0x61, 0x72, 0x66, 0x69,
	0x67, 0x68, 0x74, 0x12, 0x8a, 0x01, 0x0a, 0x0d, 0x4c, 0x69, 0x73, 0x74, 0x45, 0x71, 0x75, 0x69,
	0x70, 0x6d, 0x65, 0x6e, 0x74, 0x12, 0x1f, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x69, 0x6e, 0x2e, 0x76,
	0x31, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x45, 0x71, 0x75, 0x69, 0x70, 0x6d, 0x65, 0x6e, 0x74, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x20, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x69, 0x6e, 0x2e,
	0x76, 0x31, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x45, 0x71, 0x75, 0x69, 0x70, 0x6d, 0x65, 0x6e, 0x74,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x36, 0x92, 0x41, 0x0e, 0x12, 0x0c, 0xe8,
	0xa3, 0x85, 0xe5, 0xa4, 0x87, 0xe5, 0x88, 0x97, 0xe8, 0xa1, 0xa8, 0x82, 0xd3, 0xe4, 0x93, 0x02,
	0x1f, 0x3a, 0x01, 0x2a, 0x22, 0x1a, 0x2f, 0x76, 0x31, 0x2f, 0x69, 0x6e, 0x2f, 0x75, 0x73, 0x65,
	0x72, 0x2f, 0x65, 0x71, 0x75, 0x69, 0x70, 0x6d, 0x65, 0x6e, 0x74, 0x2d, 0x6c, 0x69, 0x73, 0x74,
	0x12, 0xa7, 0x01, 0x0a, 0x11, 0x4c, 0x69, 0x73, 0x74, 0x48, 0x65, 0x72, 0x6f, 0x45, 0x71, 0x75,
	0x69, 0x70, 0x6d, 0x65, 0x6e, 0x74, 0x12, 0x23, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x69, 0x6e, 0x2e,
	0x76, 0x31, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x48, 0x65, 0x72, 0x6f, 0x45, 0x71, 0x75, 0x69, 0x70,
	0x6d, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x24, 0x2e, 0x61, 0x70,
	0x69, 0x2e, 0x69, 0x6e, 0x2e, 0x76, 0x31, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x48, 0x65, 0x72, 0x6f,
	0x45, 0x71, 0x75, 0x69, 0x70, 0x6d, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x22, 0x47, 0x92, 0x41, 0x11, 0x12, 0x0f, 0xe8, 0x8b, 0xb1, 0xe9, 0x9b, 0x84, 0xe6, 0x80,
	0xbb, 0xe5, 0xb1, 0x9e, 0xe6, 0x80, 0xa7, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x2d, 0x3a, 0x01, 0x2a,
	0x22, 0x28, 0x2f, 0x76, 0x31, 0x2f, 0x69, 0x6e, 0x2f, 0x75, 0x73, 0x65, 0x72, 0x2f, 0x65, 0x71,
	0x75, 0x69, 0x70, 0x6d, 0x65, 0x6e, 0x74, 0x2d, 0x68, 0x65, 0x72, 0x6f, 0x65, 0x71, 0x75, 0x69,
	0x70, 0x6d, 0x65, 0x6e, 0x74, 0x2d, 0x6c, 0x69, 0x73, 0x74, 0x12, 0xb7, 0x01, 0x0a, 0x12, 0x42,
	0x61, 0x74, 0x63, 0x68, 0x48, 0x65, 0x72, 0x6f, 0x45, 0x71, 0x75, 0x69, 0x70, 0x6d, 0x65, 0x6e,
	0x74, 0x12, 0x24, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x69, 0x6e, 0x2e, 0x76, 0x31, 0x2e, 0x42, 0x61,
	0x74, 0x63, 0x68, 0x48, 0x65, 0x72, 0x6f, 0x45, 0x71, 0x75, 0x69, 0x70, 0x6d, 0x65, 0x6e, 0x74,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x25, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x69, 0x6e,
	0x2e, 0x76, 0x31, 0x2e, 0x42, 0x61, 0x74, 0x63, 0x68, 0x48, 0x65, 0x72, 0x6f, 0x45, 0x71, 0x75,
	0x69, 0x70, 0x6d, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x54,
	0x92, 0x41, 0x1d, 0x12, 0x1b, 0xe6, 0x89, 0xb9, 0xe9, 0x87, 0x8f, 0xe8, 0x8e, 0xb7, 0xe5, 0x8f,
	0x96, 0xe8, 0x8b, 0xb1, 0xe9, 0x9b, 0x84, 0xe6, 0x80, 0xbb, 0xe5, 0xb1, 0x9e, 0xe6, 0x80, 0xa7,
	0x82, 0xd3, 0xe4, 0x93, 0x02, 0x2e, 0x3a, 0x01, 0x2a, 0x22, 0x29, 0x2f, 0x76, 0x31, 0x2f, 0x69,
	0x6e, 0x2f, 0x75, 0x73, 0x65, 0x72, 0x2f, 0x65, 0x71, 0x75, 0x69, 0x70, 0x6d, 0x65, 0x6e, 0x74,
	0x2d, 0x68, 0x65, 0x72, 0x6f, 0x65, 0x71, 0x75, 0x69, 0x70, 0x6d, 0x65, 0x6e, 0x74, 0x2d, 0x62,
	0x61, 0x74, 0x63, 0x68, 0x12, 0x9e, 0x01, 0x0a, 0x12, 0x42, 0x72, 0x65, 0x61, 0x6b, 0x44, 0x6f,
	0x77, 0x6e, 0x45, 0x71, 0x75, 0x69, 0x70, 0x6d, 0x65, 0x6e, 0x74, 0x12, 0x24, 0x2e, 0x61, 0x70,
	0x69, 0x2e, 0x69, 0x6e, 0x2e, 0x76, 0x31, 0x2e, 0x42, 0x72, 0x65, 0x61, 0x6b, 0x44, 0x6f, 0x77,
	0x6e, 0x45, 0x71, 0x75, 0x69, 0x70, 0x6d, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x1a, 0x25, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x69, 0x6e, 0x2e, 0x76, 0x31, 0x2e, 0x42, 0x72,
	0x65, 0x61, 0x6b, 0x44, 0x6f, 0x77, 0x6e, 0x45, 0x71, 0x75, 0x69, 0x70, 0x6d, 0x65, 0x6e, 0x74,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x3b, 0x92, 0x41, 0x0e, 0x12, 0x0c, 0xe5,
	0x88, 0x86, 0xe8, 0xa7, 0xa3, 0xe8, 0xa3, 0x85, 0xe5, 0xa4, 0x87, 0x82, 0xd3, 0xe4, 0x93, 0x02,
	0x24, 0x3a, 0x01, 0x2a, 0x22, 0x1f, 0x2f, 0x76, 0x31, 0x2f, 0x69, 0x6e, 0x2f, 0x75, 0x73, 0x65,
	0x72, 0x2f, 0x65, 0x71, 0x75, 0x69, 0x70, 0x6d, 0x65, 0x6e, 0x74, 0x2d, 0x62, 0x72, 0x65, 0x61,
	0x6b, 0x64, 0x6f, 0x77, 0x6e, 0x1a, 0x11, 0x92, 0x41, 0x0e, 0x12, 0x0c, 0xe8, 0xa3, 0x85, 0xe5,
	0xa4, 0x87, 0xe7, 0xb3, 0xbb, 0xe7, 0xbb, 0x9f, 0x42, 0x1f, 0x5a, 0x1d, 0x67, 0x61, 0x6d, 0x65,
	0x66, 0x69, 0x5f, 0x65, 0x71, 0x75, 0x69, 0x70, 0x6d, 0x65, 0x6e, 0x74, 0x2f, 0x61, 0x70, 0x69,
	0x2f, 0x69, 0x6e, 0x2f, 0x76, 0x31, 0x3b, 0x76, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x33,
}

var file_api_in_v1_equipment_proto_goTypes = []interface{}{
	(*vo.AddEquipmentRequest)(nil),         // 0: api.in.v1.AddEquipmentRequest
	(*vo.UpgradeEquipmentRequest)(nil),     // 1: api.in.v1.UpgradeEquipmentRequest
	(*vo.AddFightEquipmentRequest)(nil),    // 2: api.in.v1.AddFightEquipmentRequest
	(*vo.ClearFightEquipmentRequest)(nil),  // 3: api.in.v1.ClearFightEquipmentRequest
	(*vo.ListEquipmentRequest)(nil),        // 4: api.in.v1.ListEquipmentRequest
	(*vo.ListHeroEquipmentRequest)(nil),    // 5: api.in.v1.ListHeroEquipmentRequest
	(*vo.BatchHeroEquipmentRequest)(nil),   // 6: api.in.v1.BatchHeroEquipmentRequest
	(*vo.BreakDownEquipmentRequest)(nil),   // 7: api.in.v1.BreakDownEquipmentRequest
	(*vo.AddEquipmentResponse)(nil),        // 8: api.in.v1.AddEquipmentResponse
	(*vo.UpgradeEquipmentResponse)(nil),    // 9: api.in.v1.UpgradeEquipmentResponse
	(*vo.AddFightEquipmentResponse)(nil),   // 10: api.in.v1.AddFightEquipmentResponse
	(*vo.ClearFightEquipmentResponse)(nil), // 11: api.in.v1.ClearFightEquipmentResponse
	(*vo.ListEquipmentResponse)(nil),       // 12: api.in.v1.ListEquipmentResponse
	(*vo.ListHeroEquipmentResponse)(nil),   // 13: api.in.v1.ListHeroEquipmentResponse
	(*vo.BatchHeroEquipmentResponse)(nil),  // 14: api.in.v1.BatchHeroEquipmentResponse
	(*vo.BreakDownEquipmentResponse)(nil),  // 15: api.in.v1.BreakDownEquipmentResponse
}
var file_api_in_v1_equipment_proto_depIdxs = []int32{
	0,  // 0: api.in.v1.EquipmentService.AddEquipment:input_type -> api.in.v1.AddEquipmentRequest
	1,  // 1: api.in.v1.EquipmentService.UpgradeEquipment:input_type -> api.in.v1.UpgradeEquipmentRequest
	2,  // 2: api.in.v1.EquipmentService.AddFightEquipment:input_type -> api.in.v1.AddFightEquipmentRequest
	3,  // 3: api.in.v1.EquipmentService.ClearFightEquipment:input_type -> api.in.v1.ClearFightEquipmentRequest
	4,  // 4: api.in.v1.EquipmentService.ListEquipment:input_type -> api.in.v1.ListEquipmentRequest
	5,  // 5: api.in.v1.EquipmentService.ListHeroEquipment:input_type -> api.in.v1.ListHeroEquipmentRequest
	6,  // 6: api.in.v1.EquipmentService.BatchHeroEquipment:input_type -> api.in.v1.BatchHeroEquipmentRequest
	7,  // 7: api.in.v1.EquipmentService.BreakDownEquipment:input_type -> api.in.v1.BreakDownEquipmentRequest
	8,  // 8: api.in.v1.EquipmentService.AddEquipment:output_type -> api.in.v1.AddEquipmentResponse
	9,  // 9: api.in.v1.EquipmentService.UpgradeEquipment:output_type -> api.in.v1.UpgradeEquipmentResponse
	10, // 10: api.in.v1.EquipmentService.AddFightEquipment:output_type -> api.in.v1.AddFightEquipmentResponse
	11, // 11: api.in.v1.EquipmentService.ClearFightEquipment:output_type -> api.in.v1.ClearFightEquipmentResponse
	12, // 12: api.in.v1.EquipmentService.ListEquipment:output_type -> api.in.v1.ListEquipmentResponse
	13, // 13: api.in.v1.EquipmentService.ListHeroEquipment:output_type -> api.in.v1.ListHeroEquipmentResponse
	14, // 14: api.in.v1.EquipmentService.BatchHeroEquipment:output_type -> api.in.v1.BatchHeroEquipmentResponse
	15, // 15: api.in.v1.EquipmentService.BreakDownEquipment:output_type -> api.in.v1.BreakDownEquipmentResponse
	8,  // [8:16] is the sub-list for method output_type
	0,  // [0:8] is the sub-list for method input_type
	0,  // [0:0] is the sub-list for extension type_name
	0,  // [0:0] is the sub-list for extension extendee
	0,  // [0:0] is the sub-list for field type_name
}

func init() { file_api_in_v1_equipment_proto_init() }
func file_api_in_v1_equipment_proto_init() {
	if File_api_in_v1_equipment_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_api_in_v1_equipment_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   0,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_api_in_v1_equipment_proto_goTypes,
		DependencyIndexes: file_api_in_v1_equipment_proto_depIdxs,
	}.Build()
	File_api_in_v1_equipment_proto = out.File
	file_api_in_v1_equipment_proto_rawDesc = nil
	file_api_in_v1_equipment_proto_goTypes = nil
	file_api_in_v1_equipment_proto_depIdxs = nil
}
