import java.util.HashMap;

public class GetIntersectionOfTwoLinkedList {
    public static int getLengthOfList(ListNode headA){
        int length = 0;
        while (headA != null) {
            length++;
            headA = headA.next;
        }
        return length;
    }
    public static ListNode getIntersectionNode(ListNode headA, ListNode headB) {
        int lenA = getLengthOfList(headA);
        int lenB = getLengthOfList(headB);
        int diff = Math.abs(lenA - lenB);
        if (lenA > lenB) {
            while (diff > 0) {
                headA = headA.next;
                diff--;
            }
        } else {
            while (diff > 0) {
                headB = headB.next;
                diff--;
            }
        }

        while (headA != null && headB != null) {
            if (headA == headB) {
                return headA;
            }
            headA = headA.next;
            headB = headB.next;
        }
        return null;
    }
    public static void main(String[] args) {
            //[4,1,8,4,5]
            ListNode l11 = new ListNode(5);
            ListNode l10 = new ListNode(4,l11);
            ListNode l9 = new ListNode(8,l10);
            ListNode l8 = new ListNode(1,l9);
            ListNode l7 = new ListNode(4,l8);
            //[5,6,1,8,4,5]
            ListNode l6 = new ListNode(5);
            ListNode l5 = new ListNode(4,l6);
            ListNode l4 = new ListNode(8,l5);
            ListNode l3 = new ListNode(1,l4);
            ListNode l2 = new ListNode(6,l3);
            ListNode l1 = new ListNode(5,l2);
    
    
            ListNode lh = getIntersectionNode(l1, l7);
            while (lh != null) {
                System.out.println(lh.val);
            }
    }
}

class ListNode {
    int val;
    ListNode next;
    ListNode(int x) {
        val = x;
        next = null;
    }
    ListNode(int x, ListNode lx) {
        val = x;
        next = lx;
    }
}
