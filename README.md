# zerocoin      


### latest-update  
1. Using distributed transaction to processing match-orders
2. Optimize order queues using skip lists
3. Support token revocation and token renewal
4. Optimize order data with reading and writing separation



### trade-engine
The trade-engine module implements an efficient matching system, focused on order matching and order book management.

The external system calls via gRPC, only needing to pass in the required fields for matching. Once matching is complete, trade tickets are sent to Kafka, where downstream services subscribe to relevant topics for further processing.


Feature: 
- High cohesion and low coupling, with a focus on internal trade matching and order book management
- Supports gRPC calls
- Use a skip list to implement bid and ask queues

Improvements that can be added: 
- [x] Ability to cancel an order
- [ ] Monitoring
- [ ] Back up the order book in persistent storage



### todo
- [ ] Add payment module
- [ ] Add front-end page
