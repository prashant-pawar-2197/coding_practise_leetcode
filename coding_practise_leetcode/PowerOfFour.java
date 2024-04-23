public class PowerOfFour {
    public static boolean isPowerOfFour(int n) {
        boolean res = false;
        for (int i = 0; i <= n/4; i++) {
            int sum = (int)Math.pow(4, i); 
            if (sum == n) {
                res = true;
            } else if (sum > n) {
                break;
            }
        }
        return res;
    }

    // second approach. this is faster
    public static boolean isPowerOfFourV2(int n){
        if(n == 0){
            return false;
        }
        boolean res = true;
        while (n != 1) {
            if(n%4 != 0){
                res = false;
                break;
            }
            n = n/4;
        }
        return res;
    }

    public static void main(String[] args) {
        System.out.println(isPowerOfFour(8));
    }
}
