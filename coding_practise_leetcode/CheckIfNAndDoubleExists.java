/*
1346. Check If N and Its Double Exist
Given an array arr of integers, check if there exist two indices i and j such that :

i != j
0 <= i, j < arr.length
arr[i] == 2 * arr[j]
 
Example 1:

Input: arr = [10,2,5,3]
Output: true
Explanation: For i = 0 and j = 2, arr[i] == 10 == 2 * 5 == 2 * arr[j]
Example 2:

Input: arr = [3,1,7,11]
Output: false
Explanation: There is no i and j that satisfy the conditions.
 */

public class CheckIfNAndDoubleExists {
    public boolean checkIfExist(int[] arr) {
        HashMap<Integer, Integer> hm = new HashMap<>();
        boolean flag = false;
        for(int i=0; i<arr.length; i++){
            int val = arr[i];
            if(hm.containsKey(val*2) && hm.get(val*2) != i){
                flag = true;
                break;
            } else if(val%2 ==0 && hm.containsKey(val/2) && hm.get(val/2) != i){
                flag = true;
                break;
            } else {
                hm.put(val, i);
            }
        }
        return flag;
    }
}


import java.util.HashMap;

