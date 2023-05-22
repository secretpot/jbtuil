package types

type Map map[interface{}]interface{}

func (this *Map) Get(key interface{}, or interface{}) interface{} {
	if value, ok := (*this)[key]; ok {
		return value
	} else {
		return or
	}
}
func (this *Map) Pop(key interface{}) interface{} {
	if value, ok := (*this)[key]; ok {
		delete(*this, key)
		return value
	} else {
		return value
	}
}
func (this *Map) Keys() []interface{} {
	keys := make([]interface{}, 0, len((*this)))
	for k := range *this {
		keys = append(keys, k)
	}
	return keys
}
func (this *Map) Values() []interface{} {
	vals := make([]interface{}, 0, len((*this)))
	for _, v := range *this {
		vals = append(vals, v)
	}
	return vals
}
func (this *Map) Update(other Map) {
	for k, v := range other {
		(*this)[k] = v
	}
}
func (this *Map) Merge(other Map) {
	for k, v := range other {
		if _, ok := (*this)[k]; !ok {
			(*this)[k] = v
		}
	}
}
func (this *Map) Clear() {
	for k := range *this {
		delete(*this, k)
	}
}
func (this *Map) Copy() Map {
	new := Map{}
	new.Update(*this)
	return new
}
