public class MergeNodes {
    public ListNode mergeNodes(ListNode head) {
        while(head != null && head.val == 0){
            head = head.next;
        }
        ListNode res = head; 
        while(head != null){
            ListNode curr = head;
            while(curr.val != 0){
                head.val =  head.val + curr.next.val;
                curr = curr.next;
            }
            head.next = curr.next;
            head = head.next;
        }
        return res;
    }
}
