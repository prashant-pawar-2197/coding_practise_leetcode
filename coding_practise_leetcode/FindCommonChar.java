import java.util.ArrayList;
import java.util.List;

public class FindCommonChar {
    public static List<String> commonChars(String[] words) {
        int[] commonCnt = new int[26];
        Arrays.fill(commonCnt, Integer.MAX_VALUE);
        List<String> list = new ArrayList<>();
        for (String word : words) {
            int[] cnt = new int[26];
            for (Character ch : word.toCharArray()) {
                cnt[ch - 'a']++;
            }
            for (int i = 0; i < 26; i++) {
                commonCnt[i] = Math.min(commonCnt[i], cnt[i]);
            }
        }
        for (char ch = 'a'; ch <= 'z'; ch++) {
            for (int i = 0; i < commonCnt[ch - 'a']; i++) {
                list.add(String.valueOf(ch));
            }
        }
        return list;
    }

    public static void main(String[] args) {
        String[] arr = new String[] { "bella", "label", "roller" };
        System.out.println(commonChars(arr));
    }
}
