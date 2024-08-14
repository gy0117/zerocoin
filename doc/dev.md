创建zerocoin-log topic   
`./kafka-topics.sh --create --topic zerocoin-log --bootstrap-server localhost:9092 --partitions 1 --replication-factor 1`

查看topic    
`./kafka-topics.sh --list --bootstrap-server localhost:9092`
