import java.util.Arrays;

public class UniqNumOfOccurences {
    public static boolean uniqueOccurrences(int[] arr) {
        boolean res = true;
        Arrays.sort(arr);
        int cur = arr[0];
        int prevcount = 0;
        int currcount = 0;
        for (int i = 0; i < arr.length; i++) {
            int num = arr[i];
            if (num == cur) {
                currcount++;
            } else {
                cur = num;
                if (prevcount == currcount) {
                    return false;
                }
                prevcount = currcount;
                currcount = 0;
            }
            if (i == 0) {
                prevcount = currcount;
            }

        }
        return res;
    }

    public static void main(String[] args) {
        int[] arr = new int[] { 1, 2 };
        System.out.println(uniqueOccurrences(arr));
    }
}
