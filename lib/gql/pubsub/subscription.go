package pubsub

type Subscription struct{
	Id int
	Value interface{}
}
type pubSub struct{
	subscription map[int]chan *Subscription
}
var Storage pubSub = pubSub{subscription: make(map[int]chan *Subscription)}

func(o *pubSub) Subscribe(id int) *Subscription{
	if o.subscription[id] == nil{
		o.subscription[id] = make(chan *Subscription)
	}
	return <- o.subscription[id]
}
func(o *pubSub) Publish(id int, value interface{}){
	o.subscription[id] <- &Subscription{
		Id:id,
		Value: value,
	}
}