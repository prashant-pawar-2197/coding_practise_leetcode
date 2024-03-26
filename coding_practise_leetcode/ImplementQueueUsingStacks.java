import java.util.Stack;

class ImplementQueueUsingStacks {
    Stack<Integer> s1;
    Stack<Integer> s2;

    public ImplementQueueUsingStacks() {
        s1 = new Stack<>();
        s2 = new Stack<>();
    }
    
    public void push(int x) {
        s1.push(x);
    }
    
    public int pop() {
        s2.clear();
        int size = s1.size();
        for (int i = 0; i < size; i++) {
            s2.push(s1.pop());
        }
        int result = s2.pop();
        size = s2.size();
        for (int i = 0; i < size; i++) {
            s1.push(s2.pop());
        }
        return result;
    }
    
    public int peek() {
        s2.clear();
        int size = s1.size();
        for (int i = 0; i < size; i++) {
            s2.push(s1.pop());
        }
        int result = s2.peek();
        size = s2.size();
        for (int i = 0; i < size; i++) {
            s1.push(s2.pop());
        }
        return result;
    }
    
    public boolean empty() {
        return s1.isEmpty();
    }


    public static void main(String[] args) {
        ImplementQueueUsingStacks q1 = new ImplementQueueUsingStacks();
        q1.push(1);
        q1.push(2);
        System.out.println(q1.peek());
        System.out.println(q1.pop());
        System.out.println(q1.empty());
    }
}

/**
 * Your MyQueue object will be instantiated and called as such:
 * MyQueue obj = new MyQueue();
 * obj.push(x);
 * int param_2 = obj.pop();
 * int param_3 = obj.peek();
 * boolean param_4 = obj.empty();
 */