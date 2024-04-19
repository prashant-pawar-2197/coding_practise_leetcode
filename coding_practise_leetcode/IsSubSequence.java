public class IsSubSequence {
    public static boolean isSubsequence(String s, String t) {
        int lenOfS = s.length();
        int lenOfT = t.length();
        int j = 0;
        int i = 0;
        while(i<lenOfS && j<lenOfT){
            if(s.charAt(i) == t.charAt(j)){
                j++;
                i++;
            } else {
                j++;
            }
        }
        boolean res = false;
        if(i == lenOfS){
            res = true;
        }
        return res;
    }

    public static void main(String[] args) {
        System.out.println(isSubsequence("abc","ahbgdc"));
    }
}
