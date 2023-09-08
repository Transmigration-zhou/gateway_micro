技术栈：gin+gorm

- 运行管理面板（配合前端项目）
```bash
go run main.go -config=./conf/dev/ -endpoint dashboard
```

- 运行代理服务
```bash
go run main.go -config=./conf/dev/ -endpoint server
```

前端：https://github.com/Transmigration-zhou/gateway_vue
