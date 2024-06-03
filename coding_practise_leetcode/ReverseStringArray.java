public class ReverseStringArray {
    public void reverseString(char[] s) {
        int low = 0;
        int high = s.length-1;
        while(low <= high){
            char tmp = s[low];
            s[low] = s[high];
            s[high] = tmp;
            low++;
            high--;
        }
    }
}
