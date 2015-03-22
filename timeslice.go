package timeshift

import "sync"

type Adder interface {
	Add(int64, interface{}) error
}

type Ranger interface {
	Range(int64, int64) []interface{}
}

type TimeSlice struct {
	gap   int64
	items map[int64]interface{}
	new   func() interface{}
	m     sync.RWMutex
}

func NewTimeSlice(gap int64, factory func() interface{}) *TimeSlice {
	return &TimeSlice{
		gap:   gap,
		items: make(map[int64]interface{}),
		new:   factory,
	}
}

func (m *TimeSlice) Add(t int64, v interface{}) error {
	m.m.Lock()
	defer m.m.Unlock()
	bt := m.boundary(t)
	value, ok := m.items[bt]
	if !ok {
		value = m.new()
		m.items[bt] = value
	}
	if adder, ok := value.(Adder); ok {
		return adder.Add(t, v)
	}
	m.items[bt] = v
	return nil
}

func (m *TimeSlice) Range(since int64, until int64) []interface{} {
	m.m.RLock()
	defer m.m.RUnlock()
	var r []interface{}
	for bt := m.boundary(since); bt < until; bt += m.gap {
		v, ok := m.items[bt]
		if !ok {
			continue
		}
		if ranger, ok := v.(Ranger); ok {
			r = append(r, ranger.Range(max(bt, since), min(bt+m.gap, until))...)
		} else {
			r = append(r, v)
		}
	}
	return r
}

func (m *TimeSlice) boundary(t int64) int64 {
	return t / m.gap * m.gap
}

func min(a int64, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

func max(a int64, b int64) int64 {
	if a > b {
		return a
	}
	return b
}
