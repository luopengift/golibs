package types


/*
[cpp] view plain copy
#string到int  
int,err:=strconv.Atoi(string)  
#string到int64  
int64, err := strconv.ParseInt(string, 10, 64)  
#int到string  
string:=strconv.Itoa(int)  
#int64到string  
string:=strconv.FormatInt(int64,10)  
*/

import (
	"strconv"
)

/*Type Int*/
type Int int

func (i Int) Int() int {return int(i)}
func (i Int) Int64() int64 {return int64(i)}
func (i Int) Float64() float64 {return float64(i)}
func (i Int) String() string {return strconv.Itoa(i.Int())}
func (i Int) Bytes() []byte {return []byte(i.String())}

/*Type Int64*/
type Int64 int64

func (i Int64) Int() int {return int(i)}
func (i Int64) Int64() int64 {return int64(i)}
func (i Int64) Float64() float64 {return float64(i)}
func (i Int64) String() string {return strconv.FormatInt(i.Int64(),10)}
func (i Int64) Bytes() []byte {return []byte(i.String())}

/*Type String*/
type String string

func (s String) Int() int {
	i, err := strconv.Atoi(s.String())
	if err != nil {
		return 0
	}
	return i
}

func (s String) Int64() int64 {
	i, err := strconv.ParseInt(s.String(), 10, 64)
	if err != nil {
		return 0
	}
	return i
}
func (s String) Float64() float64 {
	f, err := strconv.ParseFloat(s.String(), 64)
	if err != nil {
		return 0
	}
	return f
}
func (s String) String() string {return string(s)}
func (s String) Bytes() []byte {return []byte(s)}

/*Type Bytes*/
type Bytes []byte

func (b Bytes) Int() int {
        i, err := strconv.Atoi(b.String())
        if err != nil {
                return 0
        }
        return i

}

func (b Bytes) Int64() int64 {
        i, err := strconv.ParseInt(b.String(), 10, 64)
        if err != nil {
                return 0
        }
        return i

}
func (b Bytes) String() string {return string(b)}
func (b Bytes) Bytes() []byte {return []byte(b)}
