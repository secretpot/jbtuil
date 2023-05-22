package types

type Set map[interface{}]bool

func (this *Set) Get() interface{} {
	for elem := range *this {
		return elem
	}
	return nil
}
func (this *Set) Add(elem interface{}) *Set {
	if _, ok := (*this)[elem]; !ok {
		(*this)[elem] = true
	}
	return this
}
func (this *Set) Del(elem interface{}) *Set {
	if _, ok := (*this)[elem]; ok {
		delete(*this, elem)
	}
	return this
}
func (this Set) Copy() Set {
	copySet := Set{}
	for k, v := range this {
		copySet[k] = v
	}
	return copySet
}
func SetToList(s Set) []interface{} {
	l := []interface{}{}
	for elem := range s {
		l = append(l, elem)
	}
	return l
}
