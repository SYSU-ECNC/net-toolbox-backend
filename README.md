# toolBox backend

## 已完成模块：
- 用户飞书登录
- 将用户信息写入数据库
- agent注册管理
- agent脚本(在另一个项目中)

## 关于agent-client的接口
# Agent Client接口文档



## 1. POST	agent/new

- 说明

> 这是agent的注册接口，同时也是重置token后“重注册”的接口

- 请求头

| 参数名       | 类型   | 描述                             |
| :----------- | :----- | :------------------------------- |
| Content-Type | string | 定义请求参数类型application/json |

- 请求参数

| 参数名 | 类型   | 描述               |
| :----- | :----- | :----------------- |
| name   | string | agent的名字        |
| token  | string | agent注册用的token |

- 响应参数

| 参数名 | 类型 | 描述                                    |
| ------ | ---- | --------------------------------------- |
| status | bool | 注册状态，如注册成功为true，否则为false |

- 返回示例

```javascript
{
    "status": true
}
```





## 2. GET agent/task?agent_name=xxx

- 说明

> 这个接口是agent client获取被分配给它的任务，需要注意的是，这个获取task的接口并不是获取该agent被分配到的所有任务，而是每次响应一个任务，执行完一个任务后再请求下一个任务

- 请求头

| 参数名        | 类型   | 描述               |
| :------------ | :----- | :----------------- |
| Authorization | string | agent注册用的token |

- 请求参数

| 参数名 | 类型 | 描述 |
| :----- | :--- | :--- |
|        |      |      |

- 响应参数

| 参数名  | 类型   | 描述                                                |
| ------- | ------ | --------------------------------------------------- |
| task_id | int    | task_id，用于agent-client执行完该任务后上报结果使用 |
| command | string | 任务的具体内容，如"ping www.baidu.com"              |

- 返回示例

```javascript
{
    "task_id": 12345,
    "command": "ping www.baidu.com"
}
```



## 3. POST agent/task

- 说明

> 这个接口是用于agent client完成任务后提交任务执行结果的接口

- 请求头

| 参数名        | 类型   | 描述                             |
| :------------ | :----- | :------------------------------- |
| Authorization | string | agent注册用的token               |
| Content-Type  | string | 定义请求参数类型application/json |

- 请求参数

| 参数名  | 类型   | 描述         |
| :------ | :----- | :----------- |
| name    | string | agent的名字  |
| task_id | string | 提交任务的id |
| result  | string | 任务执行结果 |

- 响应参数

| 参数名 | 类型 | 描述 |
| ------ | ---- | ---- |
|        |      |      |

- 返回示例

```javascript

```

