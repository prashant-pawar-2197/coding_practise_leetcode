import java.util.*;

/*
442. Find All Duplicates in an Array
    Given an integer array nums of length n where all the integers of nums are in the range [1, n] and each integer 
    appears once or twice, return an array of all the integers that appears twice.
    You must write an algorithm that runs in O(n) time and uses only constant extra space.
 
 */
public class FindingDupWithConstSpace {
    public static List<Integer> findDuplicates(int[] nums) {
        int lenOfArr = nums.length;
        List<Integer> l1 = new ArrayList<>();
        HashMap<Integer,Integer> l2 = new HashMap<>();
        for(int i=0; i < lenOfArr; i++){
            if (l2.containsKey(nums[i])){
                l1.add(nums[i]);
            } else {
                l2.put(nums[i],1);
            }
        }
        return l1;
    }

    public static void main(String[] args) {
        int []arr = new int[]{4,3,2,7,8,2,3,1};
        System.out.println(findDuplicates(arr));
    }
}

