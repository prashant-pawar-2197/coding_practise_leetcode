/*
977. Squares of a Sorted Array
Given an integer array nums sorted in non-decreasing order, return an array of the squares of 
each number sorted in non-decreasing order.

Example 1:

Input: nums = [-4,-1,0,3,10]
Output: [0,1,9,16,100]
Explanation: After squaring, the array becomes [16,1,0,9,100].
After sorting, it becomes [0,1,9,16,100].
 */
public class SortedSquares {
    public static int[] sortedSquares(int[] nums) {
        int[] res = new int[nums.length];
        int lo = 0;
        int hi = nums.length - 1;
        for (int i = nums.length - 1; i >= 0; i--) {
            if (Math.abs(nums[lo]) >= Math.abs(nums[hi])) {
                res[i] = nums[lo] * nums[lo];
                lo++;
            } else {
                res[i] = nums[hi] * nums[hi];
                hi--;
            }
        }
        return res;
    }

    public static void main(String[] args) {
        int[] nums = new int[] { -4, -1, 0, 3, 10 };
        for (Integer val : sortedSquares(nums)) {
            System.out.print(val);
        }
    }
}
