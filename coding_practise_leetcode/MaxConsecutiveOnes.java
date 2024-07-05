public class MaxConsecutiveOnes {
    public static int findMaxConsecutiveOnes(int[] nums) {
        int res = 0;
        int currentMax = 0;
        for(int i=0; i< nums.length; i++){  
            if(nums[i] == 1) {
                currentMax++;
            } else {
                if(currentMax > res){
                    res = currentMax;
                }
                currentMax = 0;
            }
        }
        if(currentMax > res){
            res = currentMax;
        }
        return res;
    }

    public static void main(String[] args) {
        int[] arr = new int[]{1,0,1,1,0,1};
        System.out.println(findMaxConsecutiveOnes(arr));
    }
}
