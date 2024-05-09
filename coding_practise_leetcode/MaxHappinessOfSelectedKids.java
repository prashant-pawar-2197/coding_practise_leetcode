import java.util.Arrays;


/*
3075. Maximize Happiness of Selected Children

You are given an array happiness of length n, and a positive integer k.

There are n children standing in a queue, where the ith child has happiness value happiness[i]. 
You want to select k children from these n children in k turns.

In each turn, when you select a child, the happiness value of all the children that have not 
been selected till now decreases by 1. Note that the happiness value cannot become negative and 
gets decremented only if it is positive.

Return the maximum sum of the happiness values of the selected children you can achieve by 
selecting k children.
 
Example 1:

Input: happiness = [1,2,3], k = 2
Output: 4
Explanation: We can pick 2 children in the following way:
- Pick the child with the happiness value == 3. The happiness value of the remaining children becomes [0,1].
- Pick the child with the happiness value == 1. The happiness value of the remaining child becomes [0]. Note that the happiness value cannot become less than 0.
The sum of the happiness values of the selected children is 3 + 1 = 4.
Example 2:

Input: happiness = [1,1,1,1], k = 2
Output: 1
Explanation: We can pick 2 children in the following way:
- Pick any child with the happiness value == 1. The happiness value of the remaining children becomes [0,0,0].
- Pick the child with the happiness value == 0. The happiness value of the remaining child becomes [0,0].
The sum of the happiness values of the selected children is 1 + 0 = 1.
Example 3:

Input: happiness = [2,3,4,5], k = 1
Output: 5
Explanation: We can pick 1 child in the following way:
- Pick the child with the happiness value == 5. The happiness value of the remaining children becomes [1,2,3].
The sum of the happiness values of the selected children is 5.
 */

/*
Intuition - Sort the given array and select K highest elements and store them in an array
Now iterate this array.
Add the 0th index element as it is, and for the elements after that reduce the number by the index
of for loop so to subtract the number of previous selections. If the number after reduction is
less than or equal to 0, then continue, else add into res
Finally return res after for loop completes
 */
public class MaxHappinessOfSelectedKids {
    public static long maximumHappinessSum(int[] happiness, int k) {
        long res = 0;
        int lenOfArr = happiness.length;
        Arrays.sort(happiness);
        int[] arr = new int[k];
        int i = lenOfArr-1;
        int j = 0;
        while (k > 0) {
            arr[j] = happiness[i];
            j++;
            i--;
            k--;          
        }
        for (int j2 = 0; j2 < arr.length; j2++) {
            int num = arr[j2];
            if (j2 == 0) {
                res += num;
                continue;
            } else {
                if(num-j2 <= 0){
                    continue;
                } else {
                    num = num-j2;
                }
                res += num;
            }
        }
        return res; 
    }

    public static void main(String[] args) {
        int[] arr = new int[]{2,83,62};
        System.out.println(maximumHappinessSum(arr, 3));
    }
}
