# 构建grpc测试server与client

- 首先编写 `echo.proto`
- 运行IDL生成命令
`protoc -I ./proto --go_out=./proto/echo --go-grpc_out=./proto/echo ./proto/
  echo/*.proto`

- 使用生成的IDL单独构建 server 与 client 即可