package hw04lrucache

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

func (l *list) Len() int {
	return l.len
}

func (l *list) Front() *ListItem {
	return l.front
}

func (l *list) Back() *ListItem {
	return l.back
}

func (l *list) PushFront(v interface{}) *ListItem {
	temp := l.front

	switch i := v.(type) {
	case *ListItem:
		l.front = i
		l.front.Next = temp
		l.front.Prev = nil
	default:
		l.front = &ListItem{
			Value: v,
			Next:  temp,
			Prev:  nil,
		}
	}

	if temp != nil {
		temp.Prev = l.front
	}

	if l.len == 0 {
		l.back = l.front
	}

	l.len++
	return l.front
}

func (l *list) PushBack(v interface{}) *ListItem {
	temp := l.back
	l.back = &ListItem{
		Value: v,
		Next:  nil,
		Prev:  temp,
	}
	if temp != nil {
		temp.Next = l.back
	}

	if l.len == 0 {
		l.front = l.back
	}

	l.len++
	return l.back
}

func (l *list) Remove(i *ListItem) {
	prev := i.Prev
	next := i.Next

	if prev != nil {
		prev.Next = next
	} else {
		l.front = next
	}
	if next != nil {
		next.Prev = prev
	} else {
		l.back = prev
	}

	l.len--

	if l.len == 0 {
		l.back = nil
		l.front = nil
	}
}

func (l *list) MoveToFront(i *ListItem) {
	if l.len == 1 || i == l.front {
		return
	}
	l.Remove(i)
	l.PushFront(i)
}

func NewList() List {
	return new(list)
}
