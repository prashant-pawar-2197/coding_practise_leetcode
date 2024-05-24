/*
819. Most Common Word
Given a string paragraph and a string array of the banned words banned, return the most 
frequent word that is not banned. It is guaranteed there is at least one word that is not 
banned, and that the answer is unique.

The words in paragraph are case-insensitive and the answer should be returned in lowercase.

Example 1:
Input: paragraph = "Bob hit a ball, the hit BALL flew far after it was hit.", banned = ["hit"]
Output: "ball"
"!?',;."
 */

import java.util.HashMap;
import java.util.HashSet;

public class MostCommonWord {
    public String mostCommonWord(String paragraph, String[] banned) {
        HashMap<String, Integer> hm = new HashMap<>();
        HashSet<String> hs = new HashSet<>();
        for (int i = 0; i < banned.length; i++) {
            hs.add(banned[i]);
        }
        String[] arr = paragraph.toLowerCase().split("\\W+");

        for (int i = 0; i < arr.length; i++) {
            if (hs.contains(arr[i])) {
                continue;
            } else {
                if (!hm.containsKey(arr[i])) {
                    hm.put(arr[i], 1);
                } else {
                    hm.put(arr[i], hm.get(arr[i])+1);
                }
            }
        }
        
        int maxCount = 0;
        String res = "";
        for (String word : hm.keySet()) {
            int cnt = hm.get(word);
            if (cnt > maxCount) {
                maxCount = cnt;
                res = word;
            }
        }
        return res;
    }
}
