import java.util.Arrays;
import java.util.HashMap;

// 7 2 4 5 1
public class BuyAndSell {
    public static int maxProfit(int[] prices) {
        int buyPrice, currProfit, maxProfit;
        buyPrice = prices[0];
        currProfit = 0;
        maxProfit = 0;
        int arrLen = prices.length;
        for (int i = 1; i < arrLen; i++) {
            if (prices[i] < buyPrice) {
                buyPrice = prices[i];
            } else {
                currProfit = prices[i] - buyPrice;
                maxProfit = currProfit > maxProfit ? currProfit : maxProfit;
            }
        }
        return maxProfit;
    }

    public static void main(String[] args) {
        int[] arr = new int[]{3,2,6,5,0,3};
        System.out.println(maxProfit(arr));
    }
}
