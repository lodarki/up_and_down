package entity

type Map struct {
	m map[interface{}]interface{}
}

func NewMap() Map {
	return Map{
		m: make(map[interface{}]interface{}),
	}
}

func (m *Map) Put(key, value interface{}) {
	m.m[key] = value
}

func (m *Map) Get(key interface{}) interface{} {
	return m.m[key]
}

func (m *Map) Keys() []interface{} {
	var keySlice []interface{}
	for key := range m.m {
		keySlice = append(keySlice, key)
	}
	return keySlice
}

func (m *Map) KeysForInt64() []int64 {
	var keySlice []int64
	for key := range m.m {
		keySlice = append(keySlice, key.(int64))
	}
	return keySlice
}

func (m *Map) KeysForString() []string {
	var keySlice []string
	for key := range m.m {
		keySlice = append(keySlice, key.(string))
	}
	return keySlice
}

func (m *Map) Values() []interface{} {
	var valueSlice []interface{}
	for _, v := range m.m {
		valueSlice = append(valueSlice, v)
	}
	return valueSlice
}

func (m *Map) Entries() map[interface{}]interface{} {
	return m.m
}
