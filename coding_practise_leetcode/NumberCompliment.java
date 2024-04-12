public class NumberCompliment {
    public static int findComplement(int num) {
        int res = 0;
        for (int i = 0; i <= 31; i++) {
            int bit = num >> i;
            if (bit == 0) {
                break;
            }
            if ((1 & bit) == 0) {
                res = res | (1 << i);
            }
        }
        return res;
    }
    public static void main(String[] args) {
        System.out.println(findComplement(1111));
    }
}
