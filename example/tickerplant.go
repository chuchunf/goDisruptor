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
** receive data and update accordingly
 */
func receiveData(pooled *TickData, updated TickData) {
	pooled.seq = updated.seq
}

/*
** write the data to sequential log for recovery
 */
func writeLog(pooled *TickData) {
}

/*
** process the data
 */
func process(pooled *TickData) {
}

func TickerPlant() (func(data TickData), func(), func()) {
	instance := disruptor.NewDisruptor[TickData](1024)
	producer := instance.AddProducer(receiveData)
	logger := instance.AddConsumer(writeLog)
	processor := instance.AddConsumer(process)
	return producer, logger, processor
}
