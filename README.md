# zerocoin


### trade-engine   
trade-engine模块实现了高效的撮合交易系统，专注于订单的撮合以及买卖盘的管理。
外部通过grpc的方式调用，只需要传入撮合必需的字段，完成撮合后，将trade-tickets发送到kafka，下游服务订阅相关的topic，进行后续处理。

特点：
- 支持grpc调用
- 使用跳表，实现bid和ask队列
- 高内聚低耦合，内部专注于撮合交易和买卖盘管理