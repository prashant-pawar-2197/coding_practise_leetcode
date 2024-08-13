/*
1189. Maximum Number of Balloons

Given a string text, you want to use the characters of text to form as many instances 
of the word "balloon" as possible.

You can use each character in text at most once. Return the maximum number of 
instances that can be formed.

Example 1:
Input: text = "nlaebolko"
Output: 1

Example 2:
Input: text = "loonbalxballpoon"
Output: 2

Example 3:
Input: text = "leetcode"
Output: 0

*/
public class MaxNumOfBallon {
    public static int maxNumberOfBalloons(String text) {
        int[] arr = new int[26];
        for(int i=0; i<text.length(); i++){
            switch (text.charAt(i)){
                case 'b':
                arr['b'-'a']++;
                break;
                case 'a':
                arr['a'-'a']++;
                break;
                case 'l':
                arr['l'-'a']++;
                break;
                case 'o':
                arr['o'-'a']++;
                break;
                case 'n':
                arr['n'-'a']++;
                break;
            }
        }
        String bal = "balloon";
        int res = 0;
        boolean flag = true;
        while(flag){
            int counter = 0;    
            for(int i=0; i<bal.length(); i++){
                if(arr[bal.charAt(i)-'a'] <= 0){
                    flag = false;
                    break;
                } else {
                    arr[bal.charAt(i)-'a']--;
                    counter++;
                }   
            }
            if (counter == 7) {
                res++;
            }
        }
        return res;
    }

    public static int maxNumberOfBalloonsV2(String text) {
        int[] arr = new int[26];
        for(int i=0; i<text.length(); i++){
            switch (text.charAt(i)){
                case 'b':
                arr['b'-'a']++;
                break;
                case 'a':
                arr['a'-'a']++;
                break;
                case 'l':
                arr['l'-'a']++;
                break;
                case 'o':
                arr['o'-'a']++;
                break;
                case 'n':
                arr['n'-'a']++;
                break;
            }
        }
        int multiplier = arr['b'-'a'];
        if (multiplier == 0) {
            return 0;
        }
        int counter = 0;
        for (int i = 0; i < arr.length; i++) {
            switch (i){
                case 'a'-'a', 'l'-'a', 'o'-'a', 'n'-'a':
                if(arr[i] % multiplier != 0){
                    return 0;
                } else {
                    counter++;
                }
            }
        }
        if (counter == 6) {
            return multiplier;
        }
        return 0;
    }

    public static void main(String[] args) {
        System.out.println(maxNumberOfBalloons("leetcode"));
    }
}
