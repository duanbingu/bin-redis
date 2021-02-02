package bredis

type StringResult struct {
	Result string
	Err error
}

func NewStringResult(result string, err error) *StringResult {
	return &StringResult{Result: result, Err: err}
}

/**
	返回key值 !nil panic抛出
 */
func (this *StringResult) Unwrap() string{
	if this.Err!=nil {
		panic(this.Err)
	}
	return this.Result
}
/**
	返回key值 !nil 抛出自定义value
*/
func (this *StringResult) Unwrap_Or(str string) string {
	if this.Err!=nil {
		return str
	}
	return this.Result
}
/**
	判断是否获取成功，未获取成功 则执行DB。true则返回redis数据
 */
func (this *StringResult) Unwrap_Or_Else(f func() string) string {
	if this.Err!=nil {
		return f()
	}
	return this.Result
}