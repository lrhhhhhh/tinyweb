## intro
`tinyweb` 是个在 `net/http` 上构建的 和 gin 类似的 web 框架

## 特性和功能 
- 支持 Context，将 request 和 response 中封装到 Context 中 
- 使用 Trie 实现路由
   - 支持动态路由`/u/:city/:username`
- 分组控制
   - 支持 `group().group()`
- Middleware
   - handler 执行前
   - handler 执行后
   - handler 前后
   - 粒度：只支持到group
- 静态文件和模板渲染
   - 静态文件：`http.StripPrefix()`
   - 模板渲染：`html/template` 
   - 模板函数：`template.FuncMap`

## example
```shell
cd tinyweb/example/get_started
go run main.go 
```


