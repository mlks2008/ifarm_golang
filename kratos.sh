### 创建项目
#kratos new gamefi_equipment
#cd gamefi_equipment && go mod tidy
### 添加协议
#kratos proto add    api/in/v1/equipment.proto
#kratos proto client api/in/v1/equipment.proto
#kratos proto server api/in/v1/equipment.proto -t internal/service


cd components
### 公共协议
kratos proto client common/proto/http.proto
### sdk协议(调用外部应用服务)
kratos proto client sdks/gts_shop/pb/gts_shop.proto
kratos proto client sdks/gts_shop/pb/vo/gts_shop.proto
kratos proto client sdks/gamefi_platform/pb/gamefi_platform.proto
kratos proto client sdks/gamefi_platform/pb/vo/gamefi_platform.proto
kratos proto client sdks/gamefi_account/pb/gamefi_account.proto
kratos proto client sdks/gamefi_account/pb/vo/gamefi_account.proto
kratos proto client sdks/gtsportal/pb/gtsportal.proto
###### addtag
#protoc-go-inject-tag -input=sdks/gtscenter/pb/vo/gts_shop.pb.go


cd ../gamefi_equipment
### 服务的配制
kratos proto client internal/conf/conf.proto

### api协议
kratos proto client api/error.proto
kratos proto client api/in/v1/equipment.proto
kratos proto client api/in/v1/vo/equipment.proto