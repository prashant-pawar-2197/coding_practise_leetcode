/*
2215. Find the Difference of Two Arrays

Given two 0-indexed integer arrays nums1 and nums2, return a list answer of size 2 where:

answer[0] is a list of all distinct integers in nums1 which are not present in nums2.
answer[1] is a list of all distinct integers in nums2 which are not present in nums1.
Note that the integers in the lists may be returned in any order.

Example 1:

Input: nums1 = [1,2,3], nums2 = [2,4,6]
Output: [[1,3],[4,6]]
 */

import java.util.ArrayList;
import java.util.HashSet;
import java.util.List;

public class FindDiffBetwArr {
    public static List<List<Integer>> findDifference(int[] nums1, int[] nums2) {
        List<List<Integer>> res = new ArrayList<>();
        HashSet<Integer> num1 = new HashSet<>();
        HashSet<Integer> num2 = new HashSet<>();
        for (int i = 0; i < nums1.length; i++) {
            num1.add(nums1[i]); 
        }
        for (int i = 0; i < nums2.length; i++) {
            num2.add(nums2[i]); 
        }
        List<Integer> l1 = new ArrayList<>();
        List<Integer> l2 = new ArrayList<>();
        for (int i : num1) {
            if (!num2.contains(i)) {
                l1.add(i);
                num2.remove(i);
            }
        }
        for (int i : num2) {
            if (!num1.contains(i)) {
                l2.add(i);
                num1.remove(i);
            }
        }
        res.add(l1);
        res.add(l2);
        return res;
    }

    public static void main(String[] args) {
        int arr1[] = new int[]{1,2,3,3};
        int arr2[] = new int[]{1,1,2,2};
        for (List<Integer> l1 : findDifference(arr1, arr2)) {
            for (int i : l1) {
                System.out.println(i);
            }
        }
    }
}
