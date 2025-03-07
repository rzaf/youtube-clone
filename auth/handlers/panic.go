package handlers

import (
	"fmt"
	"github.com/rzaf/youtube-clone/database/pbs/helper"
)

// func PanicError(e error) {
// 	if e != nil {
// 		panic(e)
// 	}
// }

func PanicIfIsError(a any) {
	if a == nil {
		// PanicIfIsError is called with nil
		return
	}
	fmt.Printf("type(a):%T,value(a):%v\n", a, a)
	// checking underlying value of interface a and panicing if its non nil HttpError or error
	switch a2 := a.(type) {
	case *helper.HttpError:
		if a2 != nil {
			// fmt.Printf("PanicIfIsError got *helper.HttpError\n")
			panic(a2)
		}
	case error:
		if a2 != nil {
			fmt.Printf("PanicIfIsError got error\n")
			panic(a2)
		}
	}
	fmt.Printf("PanicIfIsError got no error\n")
}
