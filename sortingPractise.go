package main

import (
	"fmt"
	"sort"
)

type Item struct{
	name string
	quantitySold int
	price int
}

func sortByPopularityAndAmount(item []Item)  {
	
	sort.Slice(item, func(i, j int) bool {
		var sortedBySoldQuantity, sortedByLowerPrice bool
		sortedBySoldQuantity = item[i].quantitySold < item[j].quantitySold


        // sort by lowest sold price
        if item[i].quantitySold == item[j].quantitySold {
            sortedByLowerPrice = item[i].price < item[j].price
            return sortedByLowerPrice
        }
        return sortedBySoldQuantity
	})
}
func main()  {
	/*
	Sorting for basic data types
	numSlice := []int{2,1,8,4,6,0,8,6}
	nameSlice := []string{"Ram", "Sham","Lakhan","Sita","Gita"}
	sort.Ints(numSlice)
	sort.Strings(nameSlice)
	fmt.Println(numSlice)
	fmt.Println(nameSlice)
	*/
	sliceOfItems := []Item{{"Bat",23,200},{"Ball",46,100},{"Grip",46,50}}
	sortByPopularityAndAmount(sliceOfItems)
	fmt.Println(sliceOfItems)

}