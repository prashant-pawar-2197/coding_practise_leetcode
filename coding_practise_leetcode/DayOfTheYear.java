/*
1154. Day of the Year
Given a string date representing a Gregorian calendar date formatted as YYYY-MM-DD, 
return the day number of the year.

Example 1:

Input: date = "2019-01-09"
Output: 9
Explanation: Given date is the 9th day of the year in 2019.
Example 2:

Input: date = "2019-02-10"
Output: 41

*/

public class DayOfTheYear {
    public static boolean isLeapYear(int year) {
        if (year % 4 == 0) {
            if (year % 100 == 0) {
                if (year % 400 == 0) {
                    return true;
                } else {
                    return false;
                }
            } else {
                return true;
            }
        } else {
            return false;
        }
    }

    public int dayOfYear(String date) {
        HashMap<Integer, Integer> hm = new HashMap<>();
        hm.put(1, 31);
        hm.put(2, 28);
        hm.put(3, 31);
        hm.put(4, 30);
        hm.put(5, 31);
        hm.put(6, 30);
        hm.put(7, 31);
        hm.put(8, 31);
        hm.put(9, 30);
        hm.put(10, 31);
        hm.put(11, 30);
        hm.put(12, 31);
        String[] arr = date.split("-");
        int res = 0;
        for (int i = 0; i < arr.length; i++) {
            int num = Integer.parseInt(arr[i]);
            if (i == 0 && Integer.parseInt(arr[i + 1]) > 2) {
                if (isLeapYear(num)) {
                    res++;
                }
            }
            if (i == 1) {
                for (int j = 1; j < num; j++) {
                    res = res + hm.get(j);
                }
            }
            if (i == 2) {
                res = res + num;
            }
        }
        return res;
    }
}
