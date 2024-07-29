/*
860. Lemonade Change

At a lemonade stand, each lemonade costs $5. Customers are standing in a queue to buy from you and order one at a time (in the order specified by bills). Each customer will only buy one lemonade and pay with either a $5, $10, or $20 bill. You must provide the correct change to each customer so that the net transaction is that the customer pays $5.

Note that you do not have any change in hand at first.

Given an integer array bills where bills[i] is the bill the ith customer pays, return true if you can provide every customer with the correct change, or false otherwise.

 

Example 1:

Input: bills = [5,5,5,10,20]
Output: true
Explanation: 
From the first 3 customers, we collect three $5 bills in order.
From the fourth customer, we collect a $10 bill and give back a $5.
From the fifth customer, we give a $10 bill and a $5 bill.
Since all customers got correct change, we output true.
Example 2:

Input: bills = [5,5,10,10,20]
Output: false
Explanation: 
From the first two customers in order, we collect two $5 bills.
For the next two customers in order, we collect a $10 bill and give back a $5 bill.
For the last customer, we can not give the change of $15 back because we only have two $10 bills.
Since not every customer received the correct change, the answer is false.
 */

import java.util.HashMap;

public class LemonadeChange {
    public static boolean lemonadeChange(int[] bills) {
        if (bills[0] == 10 || bills[1] == 20) {
            return false;
        }
        HashMap<Integer, Integer> hm = new HashMap<>();
        hm.put(bills[0], 1);
        for (int i = 1; i < bills.length; i++) {
            switch (bills[i]) {
                case 5:
                    hm.put(5, (hm.get(5) == null ? 0 : hm.get(5)) + 1);
                    break;
                case 10:
                    int val = hm.get(5) == null ? 0 : hm.get(5);
                    if (val == 0) {
                        return false;
                    } else {
                        hm.put(5, val - 1);
                        hm.put(10, (hm.get(10) == null ? 0 : hm.get(10)) + 1);
                    }
                    break;
                case 20:
                    int val5 = hm.get(5);
                    int val10 = hm.get(10) == null ? 0 : hm.get(10);
                    if (val10 == 0 && val5 < 3) {
                        return false;
                    }
                    if (val5 >= 1 && val10 >= 1) {
                        hm.put(10, val10 - 1);
                        hm.put(5, val5 - 1);
                    } else if (val5 >= 3) {
                        hm.put(5, val5 - 3);
                    } else {
                        return false;
                    }
                    break;
            }
        }
        return true;
    }
    
    public static void main(String[] args) {
        int[] arr = new int[]{5,5,10,20,5,5,5,5,5,5,5,5,5,10,5,5,20,5,20,5};
        System.out.println(lemonadeChange(arr));
    }
}
