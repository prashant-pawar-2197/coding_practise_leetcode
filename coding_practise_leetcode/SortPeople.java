import java.util.Arrays;
import java.util.HashMap;

public class SortPeople {
    public static String[] sortPeople(String[] names, int[] heights) {
        HashMap<Integer, String> hm = new HashMap<>();
        for (int i = 0; i < names.length; i++) {
            hm.put(heights[i], names[i]);
        }
        Arrays.sort(heights);
        int counter = 0;
        String[] arr = new String[heights.length];
        for (int i = heights.length - 1; i >= 0; i--) {
            arr[counter] = hm.get(heights[i]);
            counter++;
        }
        return arr;
    }

    public static void main(String[] args) {
        String[] arr = new String[]{"Mary","John","Emma"};
        int[] arr1 = new int[]{180,165,170};
        sortPeople(arr, arr1);
    }
}
