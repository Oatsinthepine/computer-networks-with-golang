# PA4 evidence

- 日期：2026-09-21
- 使用的 PA3 client：`assignments/pa03-upload-client/PA3.go`。
- 输入/输出大小：`upload.txt` 原始 64 字节；编号后的 `whatever.txt` 为 72 字节。
- 回复：`original: 64 bytes, new: 72 bytes`。
- 错误处理观察：大小行以 `TrimSpace` 去除换行后再转为整数；正文以有限字节流逐行读取，最后一行没有换行也可处理。
- 完成：本地端到端验证通过；未验证学校服务器或指定团队端口。
