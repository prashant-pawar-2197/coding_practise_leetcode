/*
917. Reverse Only Letters

Given a string s, reverse the string according to the following rules:

All the characters that are not English letters remain in the same position.
All the English letters (lowercase or uppercase) should be reversed.
Return s after reversing it.

Example 1:
Input: s = "ab-cd"
Output: "dc-ba"
Example 2:
Input: s = "a-bC-dEf-ghIj"
Output: "j-Ih-gfE-dCba"
Example 3:
Input: s = "Test1ng-Leet=code-Q!"
Output: "Qedo1ct-eeLg=ntse-T!"
 
 */
public class ReverseOnlyLetters {
    public static String reverseOnlyLetters(String s) {
        int low = 0;
        int high = s.length()-1;
        char arr[] = s.toCharArray();
        while(low <= high){
            char ch = arr[low];
            if((ch >= 'A' && ch <= 'Z') || (ch >= 'a' && ch <= 'z')){
                while(!checkIfLetter(arr[high])){
                    high--;
                }
                arr[low] = arr[high];
                arr[high] = ch;
                low++;
                high--;
                continue;
            }
            low++;
        }
        StringBuilder str = new StringBuilder();
        for(int i=0; i<arr.length; i++){
            str.append(arr[i]);
        }
        return str.toString();
    }

    public static boolean checkIfLetter(char ch) {
        if(ch >= 'A' && ch <= 'Z'){
            return true;
        }
        if(ch >= 'a' && ch <= 'z'){
            return true;
        }
        return false;
    }

    public static void main(String[] args) {
        System.out.println(reverseOnlyLetters("Test1ng-Leet=code-Q!"));
    }
}
