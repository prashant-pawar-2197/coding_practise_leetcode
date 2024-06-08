public class ReplaceElements {
    public static int[] replaceElements(int[] arr) {
        int len = arr.length;
        int res[] = new int[len];
        res[len - 1] = -1;
        if (len > 1) {
            res[len - 2] = arr[len - 1];
            int counter = len - 3;
            int curMax = arr[len - 1];
            for (int i = len - 2; i > 0; i--) {
                int val = arr[i];
                if (val > curMax) {
                    res[counter] = val;
                    curMax = val;
                    counter--;
                } else {
                    res[counter] = curMax;
                    counter--;
                }
            }
        }
        return res;
    }

    public static void main(String[] args) {
        int[] arr = new int[] { 17, 18, 5, 4, 6, 1 };
        System.out.println(replaceElements(arr));
    }
}
