// package main

// import (
// 	"fmt"
// )

// func panicHandler(funName string)  {
// 	if err := recover(); err != nil {
// 		fmt.Println("Some panic occurred in function --", funName, " error is ", err)
// 	}
// }
// func DivideNumbers(a , b int){
// 	num := a/b
// 	fmt.Println(num)
// }
// func standByFunc(){
// 	defer panicHandler("standBy")
// 	fmt.Println("dummyFunction entered")
// 	DivideNumbers(2,0)
// }

// func main() {
	
// 		standByFunc()
// 		fmt.Println("Came Back to normalflow")
// }