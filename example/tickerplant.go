package example

import (
	disruptor "goDisruptor/pkg"
)

type TickData struct {
	seq     int32
	tick    int8
	dir     int8
	price   int8
	quality int8
	time    int64
}

/*
** receive data and update accoridngly
 */
func receiveData(pooled *TickData, updated TickData) {
	pooled.seq = updated.seq
}

/*
** write the data to sequencial log for recovery
 */
func writeLog(pooled *TickData) {
}

/*
** process the data
 */
func process(pooled *TickData) {
}

func TickerPlant() (func(data TickData), func(), func()) {
	disruptor := disruptor.NewDisruptor[TickData](1024)
	producer := disruptor.AddProducer(receiveData)
	logger := disruptor.AddConsumer(writeLog)
	processor := disruptor.AddConsumer(process)
	return producer, logger, processor
}
