![logo](/docs/gte_logo.png)

[![License](https://img.shields.io/badge/License-MIT-black.svg)](LICENSE)

GTE is a lightweight concurrent server framework based on Go language and epoll.

GTE 是一个基于Go语言和epoll实现的轻量级并发服务器框架。

## 关键特性

- 轻量级：基于epoll实现，占用资源少，适合高并发场景。
- 长连接：支持长连接，支持客户端主动断开连接，支持开发者自定义连接属性。
- 并发：支持多路复用，数据解析与分发和业务逻辑采用多副本并发处理。
- 事件驱动：基于事件驱动，支持异步非阻塞IO。
- 多协议：支持使用TCP、Websocket协议。
- 高性能：基于epoll实现，支持高并发场景下的高性能。
- 易用性：API简单易用。
- 插件支持：支持注册任务处理流的插件。
- 流式处理：支持请求的流式处理，支持有状态函数。
- 路由分组：支持路由分组，可以方便的管理和拓展路由。
- 连接状态回调：支持连接状态变化时回调自定义的钩子函数，可以方便的进行连接状态的维护。
- 连接保活：通过客户端续租的方式实现连接保活，可以对于异常的连接进行清理。
- 扩展性：支持插件注册，支持路由分组，支持连接状态变化时回调，可以方便的扩展功能。

## 设计

对客户端连接进行管理与维护，基于epoll完成连接事件的监听，然后进行连接事件的处理（读取连接数据，解析连接数据，数据的业务处理）

提供网关模块、连接管理模块、请求分发模块，任务处理模块

- 网关模块接受客户端连接，并交给连接管理模块进行连接后续的维护与管理
- 连接管理模块进行客户端连接的管理，并通过epoll监听连接的事件，将事件发送到队列中提交给请求分发模块进行处理。监听连接的状态变化，并触发自定义的连接信号回调函数
- 请求分发模块接收连接队列中的连接，并进行读取数据、数据拆包、请求封装等操作，然后将请求提交到请求队列
- 任务处理模块接收请求队列中的请求，并基于请求中的消息ID执行对应服务端注册的任务处理流，来消费请求，完成业务处理

## 消息处理流程
- 接收消息：gte框架通过epoll监听连接，将待读取数据的连接放入分发队列中，请求分发模块从分发队列中取出连接，并进行读取数据、数据拆包、请求封装等操作，然后将请求提交给请求队列。
- 处理消息：gte框架的任务处理模块从任务处理队列中取出请求，并基于请求中的消息ID执行对应服务端注册的任务处理流，来消费请求，完成业务处理。
- 发送消息：任务执行流在请求的处理过程中，可以选择性的基于底层的长连接进行数据的回写。

## 架构

![arch](/docs/gte_arch.png)
