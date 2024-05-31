/*
260. Single Number III
Given an integer array nums, in which exactly two elements appear only once and all the other 
elements appear exactly twice. Find the two elements that appear only once. You can return the 
answer in any order.

You must write an algorithm that runs in linear runtime complexity and uses only constant extra 
space.

Example 1:

Input: nums = [1,2,1,3,2,5]
Output: [3,5]
Explanation:  [5, 3] is also a valid answer.
 */
public class SingleNumber3 {
    public int[] singleNumber(int[] nums) {
        HashMap<Integer, Integer> hm = new HashMap<>();
        for(int i=0; i<nums.length; i++){
            int val = nums[i];
            if(hm.containsKey(val)){
                hm.remove(val);
            } else {
                hm.put(val, 1);
            }
        }
        int[] arr = new int[hm.size()];
        int counter = 0;
        for(Integer val : hm.keySet()){
   

import java.util.HashMap;

            counter++;
        }
        return arr;
    }
}
