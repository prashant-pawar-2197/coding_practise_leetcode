/*
Can Place Flowers
You have a long flowerbed in which some of the plots are planted, and some are not. 
However, flowers cannot be planted in adjacent plots.

Given an integer array flowerbed containing 0's and 1's, where 0 means empty and 1 
means not empty, and an integer n, return true if n new flowers can be planted in the 
flowerbed without violating the no-adjacent-flowers rule and false otherwise.

Example 1:

Input: flowerbed = [1,0,0,0,1], n = 1
Output: true
Example 2:

Input: flowerbed = [1,0,0,0,1], n = 2
Output: false
 */
public class CanPlaceFlowers {
    public static boolean canPlaceFlowers(int[] flowerbed, int n) {
        if(n == 0){
            return true;
        }
        if(flowerbed.length == 1){
            if(n > 1){
                return false;
            } else {
                if(flowerbed[0] == 0){
                    return true;
                } else return false;
            }
        }
        int res = 0;
        for(int i=0; i<flowerbed.length; ){
            if(flowerbed[i] == 0 ){
                if((i == flowerbed.length-1 && flowerbed[i-1] != 1) || (i == 0 && flowerbed[i+1] != 1)){
                    flowerbed[i] = 1;
                    res++;
                    i=i+1;
                    continue;
                }
                if(i < flowerbed.length-1 && flowerbed[i+1] != 1 && flowerbed[i-1] != 1){
                    res++;
                    flowerbed[i] = 1;
                }
            }
            i = i+1;
        }
        if(res >= n){
            return true;
        }
        return false;
    }
    public static void main(String[] args) {
        int[] arr = new int[]{0,0,1,0,0,0,0,1,0,1,0,0,0,1,0,0,1,0,1,0,1,0,0,0,1,0,1,0,1,0,0,1,0,0,0,0,0,1,0,1,0,0,0,1,0,0,1,0,0,0,1,0,0,1,0,0,1,0,0,0,1,0,0,0,0,1,0,0,1,0,0,0,0,1,0,0,0,1,0,1,0,0,0,0,0,0};
        System.out.println(canPlaceFlowers(arr, 17));
    }
}
