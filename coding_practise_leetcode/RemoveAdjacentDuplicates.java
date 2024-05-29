/*
1047. Remove All Adjacent Duplicates In String
You are given a string s consisting of lowercase English letters. A duplicate removal consists 
of choosing two adjacent and equal letters and removing them.

We repeatedly make duplicate removals on s until we no longer can.
Return the final string after all such duplicate removals have been made. It can be proven that 
the answer is unique.


Example 1:

Input: s = "abbaca"
Output: "ca"
Explanation: 
For example, in "abbaca" we could remove "bb" since the letters are adjacent and equal, and this
 is the only possible move.  The result of this move is that the string is "aaca", of which only
"aa" is possible, so the final string is "ca".
 */

public class RemoveAdjacentDuplicates {
    public String removeDuplicates(String s) {
        Stack<Character> s1 = new Stack<>();
        s1.push(s.charAt(0));
        for(int i=1; i<s.length(); i++){
            Character ch = s.charAt(i);
            if(!s1.isEmpty() && s1.peek()== ch){
                s1.pop();
                continue;
            } else {
                s1.push(ch);
            }
        }
        StringBuilder str = new StringBuilder();
        while(!s1.isEmpty()){
            str.append(s1.pop());
        }
        return str.reverse().toString();
    }
}
