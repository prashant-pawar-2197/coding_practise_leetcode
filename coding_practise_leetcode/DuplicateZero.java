/*
1089. Duplicate Zeros
Given a fixed-length integer array arr, duplicate each occurrence of zero, shifting the 
remaining elements to the right.

Note that elements beyond the length of the original array are not written. Do the above 
modifications to the input array in place and do not return anything.

Example 1:

Input: arr = [1,0,2,3,0,4,5,0]
Output: [1,0,0,2,3,0,0,4]
Explanation: After calling your function, the input array is modified to: [1,0,0,2,3,0,0,4]
 */

public class DuplicateZero {
    public static void duplicateZeros(int[] arr) {
        int len = arr.length;
        int[] nums = new int[len];
        int counter = 0;
        for(int i=0; i<len; i++){
            nums[i] = arr[i];
        }
        for(int i=0; i<len; i++){
            if (counter >= len){
                break;
            }
            if(nums[i]==0){
                arr[counter] = 0;
                if(counter < len-1){
                    arr[counter+1] = 0;
                }
                counter = counter+2;
                continue;
            }
            arr[counter] = nums[i];
            counter++;
        }
    }
    public static void main(String[] args) {
        int[] arr = new int[]{1,0,2,3,0,4,5,0};
        duplicateZeros(arr);
        for (int i : arr) {
            System.out.print(i);
        }
    }
}
