public class ThreeConsecutiveOddNums {
    public boolean threeConsecutiveOdds(int[] arr) {
        int res = 0;
        boolean flag = false;
        for(int i=0; i<arr.length; i++){
            if(arr[i]%2 != 0){
                res++;
                if(res >= 3){
                    flag = true;
                    break;
                }
            } else {
                res = 0;
            }
        }
        return flag;
    }
}
