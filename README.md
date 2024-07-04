# VertexAI2OpenAIChat
VertexAI2OpenAIChat


配置 config.yaml 然后启动

```yaml
location: "us-central1" # GCP地区
keyfile: "xx-xx-xxx-xxxx.json" # Google Cloud Service Key文件名称，放在运行目录下
projectid: "xx-xxx-xxxxx" # 项目ID
authkey: "xxxxxxxxxxxxx" # 转换后的接口请求头部的Bearer Key,对接聊天面板的Key
listenaddr: ":8881" # 监听地址


```

### 服务密钥下载方法：

https://cloud.google.com/iam/docs/service-account-overview?hl=zh-cn
https://developers.google.com/workspace/guides/create-credentials?hl=zh-cn

如需获取服务帐号的凭据，请执行以下操作：

在 Google Cloud 控制台中，依次点击“菜单”图标 menu > IAM 和管理 > 服务帐号。
进入“服务账号”

选择您的服务帐号。
依次点击密钥 > 添加密钥 > 创建新密钥。
选择 JSON，然后点击创建。
系统会生成新的公钥/私钥对，并将其作为新文件下载到您的计算机。将下载的 JSON 文件另存为 credentials.json 到您的工作目录中。此文件是此密钥的唯一副本。如需了解如何安全地存储密钥，请参阅管理服务帐号密钥。

点击关闭。

### 关于模型

全区支持模型

google/gemini-1.5-flash
google/gemini-1.5-pro


claude 模型

publishers/anthropic/models/claude-3-5-sonnet

publishers/anthropic/models/claude-3-sonnet

publishers/anthropic/models/claude-3-opus

publishers/anthropic/models/claude-3-haiku



参考 https://cloud.google.com/vertex-ai/generative-ai/docs/partner-models/use-claude?hl=zh-cn#anthropic_claude_region_availability

Claude 3 Opus 可在以下区域使用：
us-east5 (Ohio)

Claude 3 Sonnet 可在以下区域使用：
us-central1 (Iowa)
asia-southeast1 (Singapore)

Claude 3 Haiku 可在以下区域使用：
us-central1 (Iowa)
europe-west4 (Netherlands)



