package yamlreader

type Scalar interface {
	Bool() bool
	Int() int64
	Float() float64
	Str() string
}

type DefaultScalar struct{}

func (v DefaultScalar) Bool() (r bool) {
	return
}

func (v DefaultScalar) Int() (r int64) {
	return
}

func (v DefaultScalar) Float() (r float64) {
	return
}

func (v DefaultScalar) Str() (r string) {
	return
}

type BoolScalar struct {
	DefaultScalar
	Value bool
}

func (v BoolScalar) Bool() bool {
	return v.Value
}

type IntScalar struct {
	DefaultScalar
	Value int64
}

func (v IntScalar) Int() int64 {
	return v.Value
}

type FloatScalar struct {
	DefaultScalar
	Value float64
}

func (v FloatScalar) Float() float64 {
	return v.Value
}

type StringScalar struct {
	DefaultScalar
	Value string
}

func (v StringScalar) Str() string {
	return v.Value
}
