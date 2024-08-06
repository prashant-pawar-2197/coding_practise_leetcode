/*
190. Reverse Bits

Reverse bits of a given 32 bits unsigned integer.

Note:

Note that in some languages, such as Java, there is no unsigned integer type. In this 
case, both input and output will be given as a signed integer type. They should not 
affect your implementation, as the integer's internal binary representation is the same, 
whether it is signed or unsigned.
In Java, the compiler represents the signed integers using 2's complement notation. 
Therefore, in Example 2 above, the input represents the signed integer -3 and the output 
represents the signed integer -1073741825.

*/

public class ReverseBits {
    
    // you need treat n as an unsigned value
    public static int reverseBits(int n) {
        int result = 0;
        for(int i=0; i<32; i++){
            int bit = 1 & (n >> i);
            result = result | (bit << 31-i);
        }
        return result;
    }

}

