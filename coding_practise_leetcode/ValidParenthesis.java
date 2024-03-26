import java.util.HashMap;
import java.util.Stack;

public class ValidParenthesis {

    public static void main(String[] args) {
        ValidParenthesis.isValid("()");
    }

    public static boolean isValid(String s) {
        if (s.length() % 2 == 1 ) {
            return false;
        }
        Stack<Character> s1 = new Stack<>();
        HashMap<Character,Integer> hm = new HashMap<>();
        boolean result = true;
        hm.put('(', 1);
        hm.put('[', 1);
        hm.put('{', 1);
        for (int i = 0; i < s.length(); i++) {
            if (hm.containsKey(s.charAt(i))) {
                s1.add(s.charAt(i));
                continue;
            }
            if (!hm.containsKey(s.charAt(i))) {
                if(!s1.empty()){
                char c = s1.pop();
                switch (c) {
                    case '(':
                        if (s.charAt(i) != ')') {
                            result = false;
                            break;
                        }
                        break;
                    case '{':
                        if (s.charAt(i) != '}') {
                            result = false;
                            break;
                        }
                        break;
                    case '[':
                        if (s.charAt(i) != ']') {
                            result = false;
                            break;
                        }
                        break;
                    default:
                        break;
                }
                } else {
                    return false;
                }
            }
        }
        if(!s1.empty()){
            return false;
        }
        return result;
    }
}
