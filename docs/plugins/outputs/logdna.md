# LogDNA output plugin for Fluentd
## Overview
This plugin has been designed to output logs to LogDNA.
More info at https://github.com/logdna/fluent-plugin-logdna
>Example Deployment: [Transport Nginx Access Logs into LogDNA with Logging Operator](../../../docs/example-logdna-nginx.md)

 #### Example output configurations
 ```
 spec:
  logdna:
    api_key: xxxxxxxxxxxxxxxxxxxxxxxxxxx
    hostname: logging-operator
    app: my-app
 ```

## Configuration
### LogDNA
#### Send your logs to LogDNA

| Variable Name | Type | Required | Default | Description |
|---|---|---|---|---|
| api_key | string | Yes | - | LogDNA Api key<br> |
| hostname | string | Yes | - | Hostname<br> |
| app | string | No | - | Application name<br> |
| buffer_chunk_limit | string | No | - | Do not increase past 8m (8MB) or your logs will be rejected by LogDNA server.<br> |
