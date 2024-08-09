# Kratos Project Template

## Api
```
http://127.0.0.1:31000/q/swagger-ui
```

## 编译Protobuf协议
```
make api
```
```
make config
```

## 部署服务
```
step1: cd gamefi_equipment && make build
step2: cp -rf configs ../gamefi_config bin/
step3: bin目录gamefi_config放到/app/gamefi_config目录, configs与gamefi_equipment上传到对应机器或打包到docker运行
step4: 启动 ./gamefi_equipment -conf qa
```
