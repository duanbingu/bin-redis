package bredis

type Iterator struct {
	data[] interface{}
	index int
}

func NewIterator(data []interface{}) *Iterator {
	return &Iterator{data: data}
}

/**
	判断是否有值
 */
func (this *Iterator) HasNext() bool {
	//判断数据是否为空
	if this.data == nil||len(this.data) ==0 {
		return false
	}
	return this.index<len(this.data)
}

//取值
func (this *Iterator) Next() (ret interface{}){
	ret = this.data[this.index]
	this.index = this.index+1
	return
}