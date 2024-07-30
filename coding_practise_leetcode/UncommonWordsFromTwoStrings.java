/*
884. Uncommon Words from Two Sentences
A sentence is a string of single-space separated words where each word consists only of 
lowercase letters.
A word is uncommon if it appears exactly once in one of the sentences, and does not 
appear in the other sentence.
Given two sentences s1 and s2, return a list of all the uncommon words. You may return 
the answer in any order.

Example 1:

Input: s1 = "this apple is sweet", s2 = "this apple is sour"
Output: ["sweet","sour"]
Example 2:

Input: s1 = "apple apple", s2 = "banana"
Output: ["banana"]
 
Constraints:

1 <= s1.length, s2.length <= 200
s1 and s2 consist of lowercase English letters and spaces.
s1 and s2 do not have leading or trailing spaces.
All the words in s1 and s2 are separated by a single space.
 */
import java.util.*;

public class UncommonWordsFromTwoStrings {
    public static String[] uncommonFromSentences(String s1, String s2) {
        String[] arr1 = s1.split(" ");
        String[] arr2 = s2.split(" ");
        HashMap<String, Integer> hm = new HashMap<>();
        HashMap<String, Integer> hm2 = new HashMap<>();
        for (int i = 0; i < arr1.length; i++) {
            if (hm.get(arr1[i]) != null) {
                hm.put(arr1[i], hm.get(arr1[i]) + 1);
            } else {
                hm.put(arr1[i], 1);
            }
        }
        for (int i = 0; i < arr2.length; i++) {
            if (hm2.get(arr2[i]) != null) {
                hm2.put(arr2[i], hm2.get(arr2[i]) + 1);
            } else {
                hm2.put(arr2[i], 1);
            }
        }
        List<String> list = new ArrayList<>();
        for (int i = 0; i < arr1.length; i++) {
            if (!hm2.containsKey(arr1[i])) {
                if (hm.get(arr1[i]) == 1) {
                    list.add(arr1[i]);
                }
            }
        }
        for (int i = 0; i < arr2.length; i++) {
            if (!hm.containsKey(arr2[i])) {
                if (hm2.get(arr2[i]) == 1) {
                    list.add(arr2[i]);
                }
            }
        }
        String[] arr = new String[list.size()];
        for (int i = 0; i < list.size(); i++) {
            arr[i] = list.get(i);
        }
        return arr;
    }
    public static void main(String[] args) {
        for (String val : uncommonFromSentences("apple apple", "banana")) {
            System.out.print(val);
        }
    }
}
