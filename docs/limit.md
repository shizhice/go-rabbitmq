# Limits

### virtual host limit

> 说明：

`max-connections:` 当前 `virtual host` 允许的最大 `connection` 数。 `0`: 不允许客户端连接虚拟主机； `-1` : 不限制连接数

`max-queues:` 当前 `virtual host` 允许声明 `queue` 的最大数。`-1`: 不限制队列数

> 设置方法:
1. 在 [RabbitMQ Management](http://localhost:15672/#/limits)  中设置（路径：`Admin` > `Limits` > `Virtual host limits`）
2. 终端设置

```shell
# 设置虚拟主机的最大连接数
$ rabbitmqctl set_vhost_limits -p vhost_name '{"max-connections": 256}'

# 不限制连接数 
$ rabbitmqctl set_vhost_limits -p vhost_name '{"max-connections": 0}'

# 不允许客户端连接虚拟主机 
$ rabbitmqctl set_vhost_limits -p vhost_name '{"max-connections": -1}'

# 限制虚拟主机里最大的队列数
$ rabbitmqctl set_vhost_limits -p vhost_name '{"max-queues": 1024}'

# 不限制队列数
$ rabbitmqctl set_vhost_limits -p vhost_name '{"max-queues": -1}'

```

### user limit

> 说明：

`max-connections:` 当前用户允许创建的最大 `connection` 数。

`max-channels:` 当前用户，一个 `connection` 允许创建的最大 `channel` 数

ps: `user_limits` 默认是没有开启的，需终端开启 `rabbitmqctl enable_feature_flag user_limits`

> 设置方法:
1. 在 [RabbitMQ Management](http://localhost:15672/#/limits)  中设置（路径：`Admin` > `Limits` > `User limits`）
2. 终端设置

```shell
# 设置用户允许创建的最大连接数
$ rabbitmqctl set_user_limits golang '{"max-channels": 10}'

# 设置用户允许创建的最大channel数
$ rabbitmqctl set_user_limits golang '{"max-channels": 100}'

```

ps: 当用户创建的 `channel` 超出 `max-channels` 时，会出发异常，当前connection将被关闭














