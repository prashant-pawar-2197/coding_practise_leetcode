class ListNode {
    int val;
    ListNode next;
    ListNode() {};
    ListNode(int val) { this.val = val; }
    ListNode(int val, ListNode next) { this.val = val; this.next = next; }
}

public class AddTwoNumbers {
    // approach till long only
    public static ListNode addTwoNumbers(ListNode l1, ListNode l2) {
        int multiplier = 0;
        long num1 = 0; 
        long num2 = 0;
        while(l1 != null){
            num1 = num1 + (long)(Math.pow(10, multiplier) * l1.val);
            multiplier++;
            l1 = l1.next;
        }
        multiplier = 0;
        while(l2 != null){
            num2 = num2 + (int)(Math.pow(10, multiplier) * l2.val);
            multiplier++;
            l2 = l2.next;
        }
        long finalSum = num1+num2;
        ListNode l3 = new ListNode();
        ListNode finalList = l3;
        while(finalSum > 0){
            long div = finalSum%10;
            l3.val = (int)div;
            finalSum /= 10;
            if (finalSum != 0){
                ListNode node = new ListNode();
                l3.next = node;
                l3 = l3.next;
            }
        }
        return finalList;
    }

    // universal approach
    public static ListNode addTwoNumberss(ListNode l1, ListNode l2) {
        ListNode l3 = new ListNode();
        ListNode finalList = l3;

        while(l1 != null && l2 != null){

        }
        return finalList;
    }

    public static void main(String[] args) {
        ListNode l34 = new ListNode(4, null);
        ListNode l33 = new ListNode(6, l34);
        ListNode l32 = new ListNode(5, l33);
        
        
        ListNode l31 = new ListNode(1, null);
        ListNode l30 = new ListNode(0, l31);
        ListNode l29 = new ListNode(0, l30);
        ListNode l28 = new ListNode(0, l29);
        ListNode l27 = new ListNode(0,l28);
        ListNode l26 = new ListNode(0, l27);
        ListNode l25 = new ListNode(0, l26);
        ListNode l24 = new ListNode(0, l25);
        ListNode l23 = new ListNode(0,l24);
        ListNode l22 = new ListNode(0,l23);
        ListNode l21 = new ListNode(0, l22);
        ListNode l20 = new ListNode(0, l21);
        ListNode l19 = new ListNode(0, l20);
        ListNode l18 = new ListNode(0, l19);
        ListNode l17 = new ListNode(0,l18);
        ListNode l16 = new ListNode(0, l17);
        ListNode l15 = new ListNode(0, l16);
        ListNode l14 = new ListNode(0, l15);
        ListNode l13 = new ListNode(0,l14);
        ListNode l12 = new ListNode(0,l13);
        ListNode l11 = new ListNode(0, l12);
        ListNode l10 = new ListNode(0, l11);
        ListNode l9 = new ListNode(0, l10);
        ListNode l8 = new ListNode(0, l9);
        ListNode l7 = new ListNode(0,l8);
        ListNode l6 = new ListNode(0, l7);
        ListNode l5 = new ListNode(0, l6);
        ListNode l4 = new ListNode(0, l5);
        ListNode l3 = new ListNode(0,l4);
        ListNode l2 = new ListNode(0,l3);
        ListNode l1 = new ListNode(1, l2);
        ListNode neww = addTwoNumbers(l1, l32);
        while (neww != null){
            System.out.println(neww.val);
            neww = neww.next;
        }
    }
}