package state

import (
	"math"
	"reflect"
	"strconv"
)

func Invoke_Text(any interface{}, name string, args ...interface{}) string {
	//TODO: exception handling for invalid transform action

	inputs := make([]reflect.Value, len(args))
	for i := range args {
		inputs[i] = reflect.ValueOf(args[i])
	}

	ret := reflect.ValueOf(any).MethodByName(name).Call(inputs)

	return reflect.Value.String(ret[0])
}

func Invoke_Float(any interface{}, name string, args ...interface{}) float64 {
	//TODO: exception handling for invalid transform action

	inputs := make([]reflect.Value, len(args))
	for i := range args {
		inputs[i] = reflect.ValueOf(args[i])
	}

	ret := reflect.ValueOf(any).MethodByName(name).Call(inputs)

	return reflect.Value.Float(ret[0])
}

func Invoke_Int(any interface{}, name string, args ...interface{}) int64 {
	//TODO: exception handling for invalid transform action

	inputs := make([]reflect.Value, len(args))
	for i := range args {
		inputs[i] = reflect.ValueOf(args[i])
	}

	ret := reflect.ValueOf(any).MethodByName(name).Call(inputs)

	return reflect.Value.Int(ret[0])
}

type Functions struct{}

func (t_mC Functions) Transform_deziCelsius(mC int64) (C float32) {
	mC_float := float32(mC)
	C = mC_float / 10
	return C
}

func (t_int Functions) Transform_integerString(value_string string) (value_int int64) {
	value_int, err := strconv.ParseInt(value_string, 10, 0)
	if err != nil {
		value_int = 0
	}
	return value_int
}

func roundFloat(val float64, precision uint) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Round(val*ratio) / ratio
}

func (t_sc Functions) Transform_Statetext(old_state string) (new_state string) {
	switch old_state {
	case "St?rung":
		new_state = "Störung"
	case "Z?ndung":
		new_state = "Zündung"
	default:
		new_state = old_state
	}

	return new_state
}
