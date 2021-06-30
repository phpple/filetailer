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
file:
  paths:
    - "./debug.log"
    - "./debug2.log"
  pattern: "^\\d{4}\\-\\d{2}\\-\\d{2} \\d{2}\\:\\d{2}\\:\\d{2},\\d{3} "
notifer:
  name: dingtalk
  option:
    tokens:
      - xxxxxxxxxxxxxxxxxxxxxxxxxxxxx
    keyword: ERROR
```
