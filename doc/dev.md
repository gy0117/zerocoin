创建zerocoin-log topic   
`./kafka-topics.sh --create --topic zerocoin-log --bootstrap-server localhost:9092 --partitions 1 --replication-factor 1`

查看topic    
`./kafka-topics.sh --list --bootstrap-server localhost:9092`


mysql
sudo mysql.server start
sudo mysql.server stop

redis
brew services start redis
brew services stop redis

kafka
zookeeper-server-start /opt/homebrew/etc/kafka/zookeeper.properties
kafka-server-start /opt/homebrew/etc/kafka/server.properties

mongodb

etcd
etcd
etcdctl --endpoints=localhost:2379 endpoint health  检查是否健康

cd到dtm目录下
./dtm -c ./conf.yml

账号密码
15996230001  xingzhi123
15996230002  xingzhi123


## TODO   
- [ ] 订单服务 dtm改造，下单：创建订单 + 冻结钱包里的钱，对应的是订单服务中插入一条订单记录，账户服务中冻结钱