/*
1221. Split a String in Balanced Strings
Balanced strings are those that have an equal quantity of 'L' and 'R' characters.

Given a balanced string s, split it into some number of substrings such that:

Each substring is balanced.
Return the maximum number of balanced strings you can obtain.

Example 1:

Input: s = "RLRRLLRLRL"
Output: 4
Explanation: s can be split into "RL", "RRLL", "RL", "RL", each substring contains same 
number of 'L' and 'R'.

logic is increment counter is it is R and decrement if its L
if it becomes 0 then you have found balanced substring, then increment the counter
*/

public class BalancedString {
    public int balancedStringSplit(String s) {
        int res = 0;
        int Counter = 0;
        for(int i=0; i<s.length(); i++){
            if(s.charAt(i)=='L'){
                Counter++;
            } else if(s.charAt(i)=='R'){
                Counter--;
            }
            if(Counter == 0){
                res++;
                Counter = 0;
            }
        }
        return res;
    }    
}
