public class NumberOfOneBits {
    public static int hammingWeight(int n) {
        int result = 0;
        for(int i=0; i<32; i++){
            int bit = 1 & (n >> i);
            int res = bit & 1;
            if (res == 1) {
                result++;
            }
        }
        return result;
    }

    public static void main(String[] args) {
        System.out.println(hammingWeight(2147483645));
    }
}
