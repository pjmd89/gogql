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
	http.WsLocker.Lock()
	defer http.WsLocker.Unlock()
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
func (o *SourceEvents) listenExcecuteEvent(operationID OperationID, eventID EventID) (eventValue any) {
	eventValue = &SubscriptionClose{}
	http.WsLocker.Lock()
	operationEvents, ok := o.operationEvents[operationID]
	if !ok {
		return
	}

	subscriptionData, ok := operationEvents[eventID]
	if !ok {
		return
	}

	websocketChannel := http.WsIds[string(subscriptionData.socketID)]
	http.WsLocker.Unlock()

	select {
	case <-subscriptionData.channel:
		eventValue = subscriptionData.value
	case <-websocketChannel:
		http.WsLocker.Lock()
		defer http.WsLocker.Unlock()
		close(subscriptionData.channel)
		delete(o.operationEvents[operationID], eventID)
	}
	return
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
	http.WsLocker.Lock()
	defer http.WsLocker.Unlock()
	o.subscriptionEvents[operationID] <- value
}
