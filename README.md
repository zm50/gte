[![License](https://img.shields.io/badge/License-MIT-black.svg)](LICENSE)

GTE is a lightweight concurrent server framework based on Go language and epoll.

GTE 是一个基于Go语言和epoll实现的轻量级并发服务器框架。

## 关键特性
- 轻量级：基于epoll实现，占用资源少，适合高并发场景。
- 长连接：支持长连接，支持客户端主动断开连接。
- 并发：支持多路复用，数据解析与分发和业务逻辑采用多副本并发处理。
- 高性能：基于epoll实现，支持高并发场景下的高性能。
- 易用性：API简单易用。
- 插件支持：支持注册任务处理流的插件。
- 路由分组：支持路由分组，可以方便的管理和拓展路由。
- 扩展性：支持插件注册，支持路由分组，可以方便的扩展功能。

## 架构

消息处理流程：  
- 接收消息：gte框架通过epoll监听连接，将待读取数据的连接放入分发队列中，请求分发模块从分发队列中取出连接，并进行读取数据、数据拆包、请求封装等操作，然后将请求提交给任务处理队列。
- 处理消息：gte框架的任务处理模块从任务处理队列中取出请求，并基于请求中的消息ID执行对应服务端注册的任务处理流，来消费消息，完成业务处理。
- 发送消息：任务执行流完成消息的处理中，可以选择性的基于底层的长连接进行数据的回写。

![arch](/docs/gte_arch.png)
