/*
2309. Greatest English Letter in Upper and Lower Case

Given a string of English letters s, return the greatest English letter which 
occurs as both a lowercase and uppercase letter in s. The returned letter should be 
in uppercase. If no such letter exists, return an empty string.

An English letter b is greater than another letter a if b appears after a in the 
English alphabet.

Example 1:

Input: s = "lEeTcOdE"
Output: "E"
Explanation:
The letter 'E' is the only letter to appear in both lower and upper case.

Example 2:

Input: s = "arRAzFif"
Output: "R"
Explanation:
The letter 'R' is the greatest letter to appear in both lower and upper case.
Note that 'A' and 'F' also appear in both lower and upper case, but 'R' is greater than 'F' or 'A'.

*/

import java.util.HashMap;

public class GreatestEnglishLetter {
    public static String greatestLetter(String s) {
        HashMap<Character, Boolean> hm = new HashMap<>();
        int max = -1;
        String res = "";
        for (int i = 0; i < s.length(); i++) {
            int val = (int) s.charAt(i);
            if (val < 91) {
                if (hm.get((char) (val + 32)) != null) {
                    if (val + 32 > max) {
                        max = val + 32;
                        System.out.println((char) max);
                    }
                } else {
                    System.out.println(val);
                    hm.put((char) val, true);
                }
            } else {
                if (hm.get((char) (val - 32)) != null) {
                    if (val > max) {
                        max = val;
                        System.out.println(max);
                    }
                } else {
                    hm.put((char) val, true);
                }
            }
        }
        if (max == -1) {
            return "";
        }
        res = Character.toString(Character.toUpperCase(max));
        return res;
    }

    // 2nd approach
    public String greatestLetterV2(String s) {
        boolean[] arr = new boolean[128];
        for (char c : s.toCharArray()) {
            arr[c] = true;
        }
        for (int i = 25; i >= 0; i--) {
            if (arr['a' + i] && arr['A' + i]) {
                return String.valueOf((char) ('A' + i));
            }
        }
        return "";
    }

    public static void main(String[] args) {
        System.out.println(greatestLetter("nzmguNAEtJHkQaWDVSKxRCUivXpGLBcsjeobYPFwTZqrhlyOIfdM"));
    }
}
