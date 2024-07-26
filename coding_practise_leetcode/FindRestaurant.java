/*
    599. Minimum Index Sum of Two Lists
Solved
Easy
Topics
Companies
Given two arrays of strings list1 and list2, find the common strings with the least index sum.

A common string is a string that appeared in both list1 and list2.

A common string with the least index sum is a common string such that if it appeared at list1[i] and list2[j] then i + j should be the minimum value among all the other common strings.

Return all the common strings with the least index sum. Return the answer in any order.

 

Example 1:

Input: list1 = ["Shogun","Tapioca Express","Burger King","KFC"], list2 = ["Piatti","The Grill at Torrey Pines","Hungry Hunter Steakhouse","Shogun"]
Output: ["Shogun"]
Explanation: The only common string is "Shogun".
Example 2:

Input: list1 = ["Shogun","Tapioca Express","Burger King","KFC"], list2 = ["KFC","Shogun","Burger King"]
Output: ["Shogun"]
Explanation: The common string with the least index sum is "Shogun" with index sum = (0 + 1) = 1.
Example 3:

Input: list1 = ["happy","sad","good"], list2 = ["sad","happy","good"]
Output: ["sad","happy"]
Explanation: There are three common strings:
"happy" with index sum = (0 + 1) = 1.
"sad" with index sum = (1 + 0) = 1.
"good" with index sum = (2 + 2) = 4.
The strings with the least index sum are "sad" and "happy".
 */

import java.util.HashMap;
import java.util.Stack;

public class FindRestaurant {
    public static String[] findRestaurant(String[] list1, String[] list2) {
        HashMap<String, Integer> hm = new HashMap<>();
        Stack<String> st = new Stack<>();
        for(int i=0; i<list1.length; i++){
            hm.put(list1[i],i);
        }
        int minIndex = Integer.MAX_VALUE;
        for(int i=0; i<list2.length; i++){
            if (hm.containsKey(list2[i])) {
                int sum = hm.get(list2[i])+i;
                if (sum < minIndex) {
                    minIndex = sum;
                    st.removeAllElements();
                    st.add(list2[i]);
                } else if (sum == minIndex) {
                    st.add(list2[i]);
                }
            }
        }
        String[] arr = new String[st.size()];
        int counter = 0;
        while (!st.isEmpty()) {
            arr[counter] = st.pop();
            counter++;
        }
        return arr;
    }

    public static void main(String[] args) {
        String[] arr1 = new String[]{"happy","sad","good"};
        String[] arr2 = new String[]{"sad","happy","good"};
        for (String val : findRestaurant(arr1, arr2)) {
            System.out.print(val);
        }
    }
}
