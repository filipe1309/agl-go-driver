package queue

type MockQueue struct {
	q []*QueueDTO
}

func (mq *MockQueue) Publish(msg []byte) error {
	dto := new(QueueDTO)
	dto.Unmarshal(msg)
	mq.q = append(mq.q, dto)

	return nil
}

func (mq *MockQueue) Consume(chanDTO chan<- QueueDTO) error {
	return nil
}
