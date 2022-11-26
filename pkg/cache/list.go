package cache

type List interface {
	Len() int
	Front() *ListItem
	Back() *ListItem
	PushFront(v interface{}) *ListItem
	PushBack(v interface{}) *ListItem
	Remove(i *ListItem)
	MoveToFront(i *ListItem)
}

type ListItem struct {
	Value interface{}
	Next  *ListItem
	Prev  *ListItem
}

type list struct {
	front *ListItem
	back  *ListItem
	len   int
}

func (l list) Len() int {
	return l.len
}

func (l list) Front() *ListItem {
	return l.front
}

func (l list) Back() *ListItem {
	return l.back
}

func (l *list) PushFront(v interface{}) *ListItem {
	newItem := &ListItem{
		Value: v,
		Next:  nil,
		Prev:  nil,
	}

	if l.len == 0 {
		l.back = newItem
	} else {
		l.front.Prev = newItem
		newItem.Next = l.front
	}

	l.front = newItem
	l.len++

	return l.front
}

func (l *list) PushBack(v interface{}) *ListItem {
	newItem := &ListItem{
		Value: v,
		Next:  nil,
		Prev:  nil,
	}

	if l.len == 0 {
		l.front = newItem
	} else {
		l.back.Next = newItem
		newItem.Prev = l.back
	}

	l.back = newItem
	l.len++

	return l.back
}

func (l *list) Remove(i *ListItem) {
	if i.Prev == nil {
		if i.Next != nil {
			i.Next.Prev = nil
		}
		l.front = i.Next
	} else {
		i.Prev.Next = i.Next
	}

	if i.Next == nil {
		if i.Prev != nil {
			i.Prev.Next = nil
		}
		l.back = i.Prev
	} else {
		i.Next.Prev = i.Prev
	}

	l.len--
}

func (l *list) MoveToFront(i *ListItem) {
	if l.front == i {
		return
	}

	i.Prev.Next = i.Next
	if i.Next == nil {
		l.back = i.Prev
	} else {
		i.Next.Prev = i.Prev
	}

	i.Prev = nil
	i.Next = l.front
	l.front.Prev = i
	l.front = i
}

func NewList() List {
	return new(list)
}
