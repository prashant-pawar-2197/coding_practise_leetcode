/*
693. Binary Number with Alternating Bits

Given a positive integer, check whether it has alternating bits: namely, if two adjacent bits 
will always have different values.
Example 1:

Input: n = 5
Output: true
Explanation: The binary representation of 5 is: 101
Example 2:

Input: n = 7
Output: false
Explanation: The binary representation of 7 is: 111.

*/

public class HasAlternatingBits {
    public static boolean hasAlternatingBits(int n) {
        int prevBit = n&1;
        boolean res = true;
        while(n > 0){
            n = n >> 1;
            int curBit = n&1;
            if (curBit == prevBit) {
                res = false;
                break;
            }
            prevBit = curBit;
        } 
        return res;
    }

    public static void main(String[] args) {
        System.out.println(hasAlternatingBits(5));
    }
}
