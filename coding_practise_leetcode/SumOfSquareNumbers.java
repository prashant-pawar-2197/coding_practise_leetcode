/*
633. Sum of Square Numbers
Given a non-negative integer c, decide whether there're two integers a and b such that a2 + b2 = c.

*/

public class SumOfSquareNumbers {
    public static boolean judgeSquareSum(int c) {
        long l = 0;
        long r = (long) Math.sqrt(c);

        while (l <= r) {
            long sum = l * l + r * r;
            if (sum == c)
                return true;
            if (sum < c)
                ++l;
            else
                --r;
        }

        return false;
    }

    public static void main(String[] args) {
        System.out.println(judgeSquareSum(3));
    }
}
