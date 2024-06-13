import java.util.*;

public class NumberOfLines {
    public static int[] numberOfLines(int[] widths, String s) {
        HashMap<Integer, Integer> hm = new HashMap<>();
        int[] res = new int[2];
        for (int i = 0; i < widths.length; i++) {
            hm.put((int)'a'+i, widths[i]);
        }
        int maxWidth = 100;
        int numLines = 0;
        for (int i = 0; i < s.length(); i++) {
            int currCharWidth = hm.get((int)s.charAt(i));
            if (currCharWidth > maxWidth) {
                maxWidth = 100;
                numLines++;
            }
            if (maxWidth > 0) {
                maxWidth = maxWidth - currCharWidth;
            } else {
                maxWidth = 100;
                numLines++;
            }
        }
        res[0] = numLines;
        res[1] = Math.abs(100-maxWidth);
        return res;    
    }

    public static void main(String[] args) {
        int[] widths = new int[]{4,10,10,10,10,10,10,10,10,10,10,10,10,10,10,10,10,10,10,10,10,10,10,10,10,10};
        String s = "bbbcccdddaaa";
        for (Integer val : numberOfLines(widths, s)) {
            System.out.println(val);
        }
    }
}
