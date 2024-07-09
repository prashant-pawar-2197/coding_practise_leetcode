/*

50. Pow(x, n)
Implement pow(x, n), which calculates x raised to the power n (i.e., xn)

 */
public class PowOfXRaiseToN {
    public double myPow(double x, int n) {
        if (x == 0) {
            return 0.0;
        }
        boolean flag = false;
        if (n < 0) {
            n *= -1;
            flag = true;
        }
        double sum = 0;
        sum = powWrapper(x, n);
        if (flag) {
            sum = 1.0 / sum;
        }
        return sum;
    }

    public static double powWrapper(double x, int n) {
        if (n == 0) {
            return 1.0;
        }
        double res = powWrapper(x, n / 2);
        res *= res;
        if (n % 2 != 0)
            res *= x;
        return res;
    }
}
