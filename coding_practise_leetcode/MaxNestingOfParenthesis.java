/*
1614. Maximum Nesting Depth of the Parentheses
Given a valid parentheses string s, return the nesting depth of s. The nesting depth is 
the maximum number of nested parentheses.

Example 1:
Input: s = "(1+(2*3)+((8)/4))+1"
Output: 3
Explanation:
Digit 8 is inside of 3 nested parentheses in the string.
Example 2:
Input: s = "(1)+((2))+(((3)))"
Output: 3
Explanation:
Digit 3 is inside of 3 nested parentheses in the string.
Example 3:
Input: s = "()(())((()()))"
Output: 3


LOGIC
 keep a counter variable and increment the counter variable when you see '('
 and when you see a ')' then compare the counter with a max variable.
 
 */
public class MaxNestingOfParenthesis {
    public static int maxDepth(String s) {
        int counter = 0;
        int max = 0;
        for(int i=0; i<s.length(); i++){
            char ch = s.charAt(i);
            if(ch == '('){
                counter++;
            } else if (ch == ')'){
                max = Math.max(max, counter);
                counter--;
            }
        }
        return max;
    }

    public static void main(String[] args) {
        System.out.println(maxDepth("(1+(2*3)+((8)/4))+1"));
    }   
}
