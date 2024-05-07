/*
2816. Double a Number Represented as a Linked List

You are given the head of a non-empty linked list representing a non-negative integer without 
leading zeroes.

Return the head of the linked list after doubling it.
 */

public class DoubleNumInLL {
    public static ListNode doubleIt(ListNode head) {
        int num = doubleItWrapper(head, 0); 
        if (num != 0){
            ListNode ll = new ListNode(num);
            ll.next = head;
            head = ll;
        }
        return head;
    }
    public static int doubleItWrapper(ListNode head, int carry){
        if (head == null) {
            return 0;
        }
        if (head.next != null) {
            carry = doubleItWrapper(head.next, carry);
        }
        int sum = head.val*2 + carry;
        if (sum >= 10) {
            carry = sum/10;
            head.val = sum % 10;
        } else {
            head.val = sum;
            carry = 0;
        }
        return carry;
    }
}
