/*
129. Sum Root to Leaf Numbers

You are given the root of a binary tree containing digits from 0 to 9 only.
Each root-to-leaf path in the tree represents a number.

For example, the root-to-leaf path 1 -> 2 -> 3 represents the number 123.
Return the total sum of all root-to-leaf numbers. Test cases are generated so that the answer 
will fit in a 32-bit integer.
A leaf node is a node with no children.

Example 1:
Input: root = [1,2,3]
Output: 25
Explanation:
The root-to-leaf path 1->2 represents the number 12.
The root-to-leaf path 1->3 represents the number 13.
Therefore, sum = 12 + 13 = 25.

 */

public class SumRootToLeafNumbers {
    public static int sumNumbers(TreeNode root) {
        int sum = 0;
        int num = 0;
        return sumNumbersWrapper(root, sum, num);
    }


    public static int sumNumbersWrapper(TreeNode root, int sum, int num){
        if (root.left == null && root.right == null) {
            sum = sum + ((num*10) + root.val);
            return sum;
        }
        num = num*10 + root.val;
        if (root.left != null) {
            sum = sumNumbersWrapper(root.left, sum, num);
        }
        if (root.right != null) {
            sum = sumNumbersWrapper(root.right, sum, num);
        }
        return sum;
    }

    public static void main(String[] args) {
        TreeNode n5 = new TreeNode(1, null, null);
        TreeNode n4 = new TreeNode(5, null, null);
        TreeNode n3 = new TreeNode(0, null, null);
        TreeNode n2 = new TreeNode(9, n4, n5);
        TreeNode n1 = new TreeNode(4, n2, n3);
        System.out.println(sumNumbers(n1));
    }

}
/*
class TreeNode {
    int val;
    TreeNode left;
    TreeNode right;
    TreeNode() {}
    TreeNode(int val) { this.val = val; }
    TreeNode(int val, TreeNode left, TreeNode right) {
        this.val = val;
        this.left = left;
        this.right = right;
    }
}
*/