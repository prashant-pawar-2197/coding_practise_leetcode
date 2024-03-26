package main

func main() {
  mymap := make(map[string]string)
  
  if(len(mymap) == 0){
    println("map is empty")
  } else {
    println("map is not empty")
  }

  mymap["Prashant"] = "Pawar"
  if(len(mymap) == 0){
    println("map is empty")
  } else {
    println("map is not empty")
  }

}