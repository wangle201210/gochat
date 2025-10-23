# GoChat - AI 对话应用

基于 Fyne + Eino SDK 实现的 Golang AI 对话应用。

## 项目结构

```
gochat/
├── cmd/
│   └── gochat/           # 主程序入口
│       └── main.go
├── internal/             # 内部模块
│   ├── config/          # 配置管理
│   │   └── config.go
│   ├── models/          # 数据模型
│   │   └── message.go
│   ├── service/         # 服务层
│   │   └── ai/         # AI 服务
│   │       └── service.go
│   └── ui/             # 用户界面
│       └── chat.go
├── config.example.json  # 配置文件示例
├── go.mod
└── README.md
```

## 技术栈

- **UI 框架**: [Fyne v2](https://fyne.io/) - 跨平台 GUI 框架
- **AI SDK**: [Eino SDK](https://github.com/cloudwego/eino) - CloudWeGo AI 开发框架
- **语言**: Go 1.23.7+

## 功能特性

- ✅ 简洁的图形界面
- ✅ 流式对话响应
- ✅ 对话历史管理
- ✅ 支持多种 AI 模型（通过 Eino SDK）
- ✅ 配置文件管理

## 安装依赖

```bash
go mod tidy
```

## 配置

### 1. 创建配置文件

首次运行时，程序会在 `~/.gochat/config.json` 创建默认配置文件。

也可以手动复制示例配置：

```bash
cp config.example.json ~/.gochat/config.json
```

### 2. 编辑配置文件

```json
{
  "ai": {
    "provider": "openai",
    "model": "gpt-3.5-turbo",
    "api_key": "your-api-key-here",
    "base_url": "https://api.openai.com/v1"
  },
  "ui": {
    "window_width": 800,
    "window_height": 600,
    "theme": "light"
  }
}
```

**配置说明：**

- `ai.provider`: AI 服务提供商（目前支持 "openai"）
- `ai.model`: 模型名称（如 "gpt-3.5-turbo", "gpt-4" 等）
- `ai.api_key`: API 密钥
- `ai.base_url`: API 基础地址
- `ui.window_width`: 窗口宽度
- `ui.window_height`: 窗口高度
- `ui.theme`: 主题（"light" 或 "dark"）

### 获取 API Key

#### OpenAI

1. 访问 [OpenAI Platform](https://platform.openai.com/)
2. 注册并登录账号
3. 进入 API Keys 页面创建新的 API Key
4. 将 API Key 配置到 `config.json` 中

#### 兼容的 OpenAI API 服务

本应用支持所有兼容 OpenAI API 格式的服务，例如：
- **Azure OpenAI**: 修改 `base_url` 为你的 Azure endpoint
- **本地部署的模型** (如 LocalAI, Ollama): 修改 `base_url` 为本地地址
- **其他第三方服务**: 任何提供 OpenAI 兼容接口的服务

## 运行

```bash
go run cmd/gochat/main.go
```

## 构建

### 构建可执行文件

```bash
go build -o gochat cmd/gochat/main.go
```

### 跨平台构建

**Windows:**
```bash
GOOS=windows GOARCH=amd64 go build -o gochat.exe cmd/gochat/main.go
```

**macOS:**
```bash
GOOS=darwin GOARCH=amd64 go build -o gochat cmd/gochat/main.go
```

**Linux:**
```bash
GOOS=linux GOARCH=amd64 go build -o gochat cmd/gochat/main.go
```

## 使用说明

1. 启动应用后，会看到聊天界面
2. 在底部输入框中输入消息
3. 点击"发送"按钮或按 Enter 发送消息
4. AI 会以流式方式返回回复
5. 点击"清空历史"可以清除所有对话记录

## 代码模块说明

### config 模块

负责配置文件的加载、保存和管理。支持 JSON 格式配置。

**主要功能：**
- 加载配置文件
- 保存配置到文件
- 提供默认配置
- 获取配置文件路径

### models 模块

定义数据模型，主要是消息模型。

**Message 结构：**
- `ID`: 消息唯一标识
- `Role`: 消息角色（user/assistant/system）
- `Content`: 消息内容
- `Timestamp`: 消息时间戳

### service/ai 模块

封装 AI 服务逻辑，使用 Eino SDK 进行 AI 对话。

**主要功能：**
- 初始化 AI 模型
- 发送消息并获取回复
- 流式对话
- 管理对话历史
- 消息格式转换

### ui 模块

使用 Fyne 实现图形用户界面。

**主要组件：**
- 消息列表（List）
- 输入框（MultiLineEntry）
- 发送按钮
- 清空历史按钮

## 扩展支持

### 添加新的 AI 提供商

在 `internal/service/ai/service.go` 中的 `NewService` 函数添加新的 provider 分支：

```go
case "your-provider":
    chatModel, err = yourprovider.NewChatModel(context.Background(), &yourprovider.ChatModelConfig{
        BaseURL: cfg.BaseURL,
        Model:   cfg.Model,
        APIKey:  cfg.APIKey,
    })
```

## 许可证

MIT License

## 参考资源

- [Fyne 文档](https://docs.fyne.io/)
- [Eino SDK 文档](https://github.com/cloudwego/eino)
- [OpenAI Platform](https://platform.openai.com/)
