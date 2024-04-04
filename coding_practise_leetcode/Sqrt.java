public class Sqrt {
    public static int mySqrt(int x) {
        int result = x;
        int low = 0;
        int high = x/2;
        while(low <= high){
            int mid;
            mid = low + (high-low)/2;
            int mul = mid*mid;
            if (mul == x) {
                result = mid;    
                break;
            } else if (mul < x){
                low = mid+1;
                result = mid;
            } else {
                high = mid-1;
            }
        }
        return result;
    }

    public static void main(String[] args) {
        System.out.println(mySqrt(2147395599));
    }
}
