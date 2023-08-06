package internal

//go:generate mockgen -package=mock -destination=./mock/sarama_mock.go github.com/IBM/sarama SyncProducer
