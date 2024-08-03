/*
1460. Make Two Arrays Equal by Reversing Subarrays

You are given two integer arrays of equal length target and arr. In one step, you can 
select any non-empty subarray of arr and reverse it. You are allowed to make any number 
of steps.

Return true if you can make arr equal to target or false otherwise.
Example 1:

Input: target = [1,2,3,4], arr = [2,4,1,3]
Output: true
Explanation: You can follow the next steps to convert arr to target:
1- Reverse subarray [2,4,1], arr becomes [1,4,2,3]
2- Reverse subarray [4,2], arr becomes [1,2,4,3]
3- Reverse subarray [4,3], arr becomes [1,2,3,4]
There are multiple ways to convert arr to target, this is not the only way to do so.
Example 2:

Input: target = [7], arr = [7]
Output: true
Explanation: arr is equal to target without any reverses.
Example 3:

Input: target = [3,7,9], arr = [3,7,11]
Output: false
Explanation: arr does not have value 9 and it can never be converted to target.
 */

import java.util.HashMap;

public class CanBeEqual {
    public boolean canBeEqual(int[] target, int[] arr) {
        HashMap<Integer, Integer> hs = new HashMap<>();
        for(int i=0; i<target.length; i++){
            if(hs.containsKey(target[i])){
                hs.put(target[i], hs.get(target[i])+1);
            } else {
                hs.put(target[i], 1);
            }
        }
        for(int i=0; i<arr.length; i++){
            if(!hs.containsKey(arr[i]) || hs.get(arr[i]) == 0){
                return false;
            } else {
                hs.put(arr[i], hs.get(arr[i])-1);
            }
        }
        return true;
    }

    // idea behind this was capture the freq of target, then reduce the occurence when it 
    // occurs in arr... And then if we get any frequency lesser than 0 means that element
    // is either not present or the frequencies don;t match
    public static boolean canBeEqualV2(int[] target, int[] arr) {
        int freArr[] = new int[1001];
        for(int i=0; i<target.length; i++){
            freArr[target[i]]++;
        }
        for(int i=0; i<arr.length; i++){
            freArr[arr[i]]--;
        }
        for (int i = 0; i < freArr.length; i++) {
            if (freArr[i] < 0) {
                return false;
            }
        }
        return true;
    }
}
