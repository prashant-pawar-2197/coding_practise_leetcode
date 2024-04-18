public class OneBitTwoBit {
    public static boolean isOneBitCharacter(int[] bits) {
        boolean res = false;
        boolean oneBitFound = false;
        int len = bits.length;
        for(int i=0; i<len;){
            int curr = bits[i];
            int nxt = 0;
            if (i+1 != len) {
                nxt = bits[i+1];
            }
            if (curr == 1 && (nxt == 1 || nxt == 0)) {
                i = i+2;
                continue;
            }
            if (bits[i] == 0 && i == len-1) {
                oneBitFound = true;
                break;
            }
            i++;
        }
        if (oneBitFound) {
            res = true;
        }
        return res;
    }
    
    /*
        Efficient approach     
        public boolean isOneBitCharacter(int[] bits) {
        boolean res = false;
        int len = bits.length;
        int i = 0;
        while(i<len-1){
            int curr = bits[i];
            if (curr == 1) {
                i = i+2;
                continue;
            }
            if (curr == 0) {
                i++;
            }
        }
        if (i == len-1) {
            res = true;
        }
        return res;
    }
    */

    public static void main(String[] args) {
        int[] arr = new int[]{0,0};
        System.out.println(isOneBitCharacter(arr));
    }
}
