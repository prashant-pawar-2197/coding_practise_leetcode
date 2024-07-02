/*
350. Intersection of Two Arrays II

Given two integer arrays nums1 and nums2, return an array of their intersection. Each 
element in the result must appear as many times as it shows in both arrays and you may 
return the result in any order.


Example 1:

Input: nums1 = [1,2,2,1], nums2 = [2,2]
Output: [2,2]
Example 2:

Input: nums1 = [4,9,5], nums2 = [9,4,9,8,4]
Output: [4,9]
Explanation: [9,4] is also accepted.
 

Constraints:

1 <= nums1.length, nums2.length <= 1000
0 <= nums1[i], nums2[i] <= 1000
  
*/

import java.util.ArrayList;
import java.util.List;

public class IntersectionOfTwoArrayII {
    public int[] intersect(int[] nums1, int[] nums2) {
        int[] arr = new int[1001];
        List<Integer> list = new ArrayList<>();
        int arrr = nums1.length > nums2.length ? 1 : 0;
        switch(arrr){
            case 0 :
                for(int i=0; i<nums2.length; i++){
                    arr[nums2[i]]++;
                }
                for(int j=0; j<nums1.length; j++){
                    if(arr[nums1[j]] != 0){
                        list.add(nums1[j]);
                        arr[nums1[j]]--;
                    }
                }
            break; 
            case 1 :
                for(int i=0; i<nums1.length; i++){
                    arr[nums1[i]]++;
                }
                for(int j=0; j<nums2.length; j++){
                    if(arr[nums2[j]] != 0){
                        list.add(nums2[j]);
                        arr[nums2[j]]--;
                    }
                }
            break;
        }
        int size = list.size();
        int[] out = new int[size];
        for(int i = 0; i < size; i++) out[i] = list.get(i);
        return out;
    }
}
