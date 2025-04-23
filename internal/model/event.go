package model

type Event struct {
	IngestionTime int64  `msgpack:"ingestion_time"`
	Timestamp     int64  `msgpack:"timestamp"`
	Message       string `msgpack:"message"`
}
