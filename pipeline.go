package pipeline

import (
	"fmt"
	"reflect"
	"sort"
	"unsafe"
)

type entry struct {
	p unsafe.Pointer
}

func (e *entry) Load() interface{} {
	return *(*interface{})(e.p)
}

type Stream struct {
	Items []*entry
}

func NewStream(e ...interface{}) *Stream {
	s := &Stream{Items: make([]*entry, 0, len(e))}
	// s.FLat(1, e)
	s.Add(e)
	return s
}

// Add 只展开一级
func (s *Stream) Add(e ...interface{}) *Stream {
	for _, item := range e {
		s.Items = append(s.Items, &entry{p: unsafe.Pointer(&item)})
	}
	return s
}

// Flat 如果interface e是集合类型，则逐级展开
func (s *Stream) FLatAll() *Stream {
	items := make([]interface{}, 0, len(s.Items))
	for _, item := range s.Items {
		items = append(items, item.Load())
	}
	s.Items = []*entry{}
	s.flatAll(items)
	return s
}

func (s *Stream) flatAll(e ...interface{}) {
	for _, item := range e {
		value := reflect.ValueOf(item)
		kind := value.Kind()
		if kind == reflect.Slice || kind == reflect.Array {
			for i := 0; i < value.Len(); i++ {
				// fmt.Println("FLat index of ", i, value.Len())
				s.flatAll(value.Index(i).Interface())
			}
		} else {
			s.Items = append(s.Items, &entry{p: unsafe.Pointer(&item)})
		}
	}
}

func (s *Stream) Flat() *Stream {
	items := make([]interface{}, 0, len(s.Items))
	for _, item := range s.Items {
		items = append(items, item.Load())
	}
	s.Items = []*entry{}
	s.flat(2, items)
	return s
}

func (s *Stream) flat(depth int, e ...interface{}) {
	for _, item := range e {
		value := reflect.ValueOf(item)
		kind := value.Kind()
		if (kind == reflect.Slice || kind == reflect.Array) && depth >= 0 {
			for i := 0; i < value.Len(); i++ {
				// fmt.Println("FLat index of ", i, value.Len(), depth)
				s.flat(depth-1, value.Index(i).Interface())
			}
		} else {
			s.Items = append(s.Items, &entry{p: unsafe.Pointer(&item)})
		}
	}
}

/* filter slice and return new slice*/
func (s *Stream) Filter(filter func(interface{}) bool) *Stream {
	list := make([]*entry, 0, len(s.Items))
	for _, entry := range s.Items {
		// fmt.Println(entry.Load())
		if !filter(entry.Load()) {
			list = append(list, entry)
		}
	}
	s.Items = list
	return s
}

/* iterate slice */
func (s *Stream) ForEach(act func(interface{})) {
	for _, entry := range s.Items {
		act(entry.Load())
	}
}

/* map slice and return new slice*/
func (s *Stream) Map(m func(interface{}) interface{}) *Stream {
	list := make([]*entry, 0, len(s.Items))
	for _, e := range s.Items {
		item := m(e.Load())
		list = append(list, &entry{p: unsafe.Pointer(&item)})
	}
	s.Items = list
	return s
}

// func FlatMap(func(interface{}) interface{}) {

//}

func (s *Stream) Limit(n int) *Stream {
	if n <= len(s.Items) {
		s.Items = s.Items[:n]
	}
	return s
}

func (s *Stream) Skip(n int) *Stream {
	if n >= len(s.Items) {
		s.Items = []*entry{}
	}
	s.Items = s.Items[n:]
	return s
}

func (s *Stream) Concat(s2 *Stream) *Stream {
	if s2 != nil {
		s.Items = append(s.Items, s2.Items...)
	}
	return s
}

func (s *Stream) Distinct() *Stream {
	m := make(map[interface{}]struct{}, len(s.Items))
	list := make([]*entry, 0, len(s.Items))
	for _, e := range s.Items {
		if _, ok := m[e.Load()]; !ok {
			m[e.Load()] = struct{}{}
			list = append(list, e)
		}
	}
	s.Items = list
	return s
}

func (s *Stream) Sort(order func(i, j int) bool) *Stream {
	sort.Slice(s.Items, order)
	return s
}

func (s *Stream) Index(i int) interface{} {
	return s.Items[i].Load()
}

func (s *Stream) Print(tag string) *Stream {
	fmt.Print(tag, ": ")
	for _, e := range s.Items {
		fmt.Print(e.Load(), " ")
	}
	fmt.Println()
	return s
}

//func (s *Stream) ToSlice() []string {
//	return s.Items
//}
