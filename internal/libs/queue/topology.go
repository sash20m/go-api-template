package queue

// OwnedQueues returns the queues this service wants to exist in RabbitMQ.
func OwnedQueues() []QueueSpec {
	return []QueueSpec{
		{
			Name: UsersCreatedQueueName,
			Bindings: []Binding{
				{RoutingKey: UsersCreatedRoutingKey},
			},
		},
	}
}
