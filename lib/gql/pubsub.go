package gql

import "github.com/pjmd89/gogql/lib/http"

var PubSub = &SourceEvents{
	subscriptionEvents: make(map[OperationID]chan interface{}, 0),
	operationEvents:    make(map[OperationID]map[EventID]*Subscription, 0),
}

// se ejecuta en el momento levantar el esquema, 1 por cada subscripcion existente en el esquema
func (o *SourceEvents) createSubscriptionEvent(operationID OperationID) {
	o.subscriptionEvents[operationID] = make(chan interface{})
	go o.listenSubscriptionEvent(operationID)
}

// se ejecuta en el momento de crear una subscripcion cuando llega la peticion desde websocket
func (o *SourceEvents) createExcecuteEvent(eventId EventID, operationId OperationID, socketId SocketID, requestID RequestID, mt int) {

	if o.operationEvents[operationId] == nil {
		o.operationEvents[operationId] = make(map[EventID]*Subscription, 0)
	}
	o.operationEvents[operationId][eventId] = &Subscription{
		eventID:     eventId,
		operationID: operationId,
		socketID:    socketId,
		requestID:   requestID,
		messageType: mt,
		channel:     make(chan bool),
	}
}
func (o *SourceEvents) listenExcecuteEvent(operationID OperationID, eventID EventID) interface{} {
	select {
	case <-o.operationEvents[operationID][eventID].channel:
		return o.operationEvents[operationID][eventID].value
	case <-http.WsIds[string(o.operationEvents[operationID][eventID].socketID)]:
		close(o.operationEvents[operationID][eventID].channel)
		delete(o.operationEvents[operationID], eventID)
		return &SubscriptionClose{}
	}
}
func (o *SourceEvents) listenSubscriptionEvent(operationID OperationID) {
	for {
		r := <-o.subscriptionEvents[operationID]
		if len(o.operationEvents) > 0 {
			for _, v := range o.operationEvents[operationID] {
				v.value = r
				v.channel <- true
			}
		}
	}
}
func (o *SourceEvents) Publish(operationID OperationID, value interface{}) {
	o.subscriptionEvents[operationID] <- value
}
