import java.math.BigDecimal;
import java.math.RoundingMode;
import java.util.Arrays;

public class TrimMean {
    public static double trimMean(int[] arr) {
        Arrays.sort(arr);
        double len = arr.length;
        double smallest5Per = Math.ceil((5.0/100.0) * (int)len);
        double largest5Per = (int)len-1-smallest5Per;
        double sum = 0;
        for (int i = (int)smallest5Per; i <= largest5Per; i++) {
            sum += arr[i];
        }
        return BigDecimal.valueOf(sum/(len*0.9)).setScale(5, RoundingMode.HALF_UP).doubleValue();
    }

    public static void main(String[] args) {
        int[] arr = new int[]{0, 1, 2, 3, 4, 5, 6, 7, 8, 9};
        System.out.println(trimMean(arr));
    }
}
