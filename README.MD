FileTailer
==========

A tool of monitoring error log and sending notify by dingtalk.

Usage
------

```
./filetailer
```

The configuration file named config.yml must place in the same directory.

Basic Configuration
-------------------

```yaml
watchers:
  - paths:
      - "./debug.log"
      - "./debug2.log"
    # 匹配一行开始的模式
    pattern: "^\\d{4}\\-\\d{2}\\-\\d{2} \\d{2}\\:\\d{2}\\:\\d{2},\\d{3} "
    rules:
      - cause: any     # 任意内容均允许展示
  - paths:
      - "./nginx.access.log"
    pattern: "\\[\\d{2}\\/[a-zA-Z]{3}\\/\\d{4}\\:\\d{2}\\:\\d{2}\\:\\d{2} +\\d{4}\\]"
    rules:
      - seperator: " "            # 匹配内容分隔符
        expression: "{print $9}"  # awk表达式
        cause: eq                 # awk表达式的内容和values的比较条件，包括eq、gt等
        # 对awk表达式出来的内容进行比较的具体值，符合其中任意一项则认为匹配上了，会允许通知
        values:
          - "500"
          #- "502"
          - "504"
        msg: "服务端错误:\n> %s"
notifiers:
  default: dingding    # 默认的notifier
  list:
    - key: dingding    # notifier的名称
      type: dingtalk   # 类型
      option:
        tokens:        # 钉钉机器人token，可以配置多个
        - xxxxxxxxxxxxxxxxxxxxxxxxxxxxx
      keyword: ERROR   # 关键词
      maxline: 4       # 最多显示的内容行数
      maxchar: 200     # 每行显示的最大字符数（第1行和最后1行不在此限制）
```
