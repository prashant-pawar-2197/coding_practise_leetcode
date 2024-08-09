/*
643. Maximum Average Subarray I
You are given an integer array nums consisting of n elements, and an integer k.
Find a contiguous subarray whose length is equal to k that has the maximum average value and 
return this value. Any answer with a calculation error less than 10-5 will be accepted.

Example 1:
Input: nums = [1,12,-5,-6,50,3], k = 4
Output: 12.75000
Explanation: Maximum average is (12 - 5 - 6 + 50) / 4 = 51 / 4 = 12.75

Example 2:
Input: nums = [5], k = 1
Output: 5.00000

*/
public class MaxAvgSubArray {
    public static double findMaxAverage(int[] nums, int k) {
        double max = Integer.MIN_VALUE;
        for (int i = 1; i < nums.length; i++) {
            nums[i] = nums[i] + nums[i-1];
        }
        for (int i = 0; i <= nums.length - k; i++) {
            double sum = 0;
            if (i != 0) {
                sum = nums[i + k - 1] - nums[i - 1];
            } else {
                sum = nums[i + k - 1] - 0;
            }
            double div = sum / (k * 1.0);
            max = div > max ? div : max;
        }
        return max;
    }

    public static void main(String[] args) {
        int[] nums = new int[]{1,12,-5,-6,50,3};
        System.out.println(findMaxAverage(nums, 4));
    }
}
