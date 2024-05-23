import java.util.ArrayList;
/*

872. Leaf-Similar Trees
Consider all the leaves of a binary tree, from left to right order, the values of those leaves 
form a leaf value sequence.
For example, in the given tree above, the leaf value sequence is (6, 7, 4, 9, 8).

Two binary trees are considered leaf-similar if their leaf value sequence is the same.

Return true if and only if the two given trees with head nodes root1 and root2 are leaf-similar.

*/
class TreeNode {
    int val;
    TreeNode left;
    TreeNode right;

    TreeNode() {
    }

    TreeNode(int val) {
        this.val = val;
    }

    TreeNode(int val, TreeNode left, TreeNode right) {
        this.val = val;
        this.left = left;
        this.right = right;
    }
}

public class LeafSameTree {
    public static void inOrderTrav(TreeNode node, ArrayList<Integer> l1) {
        if (node == null)
            return;
        inOrderTrav(node.left, l1);
        if (node.left == null && node.right == null) {
            l1.add(node.val);
        }
        inOrderTrav(node.right, l1);
    }

    public static boolean leafSimilar(TreeNode root1, TreeNode root2) {
        boolean res = true;
        ArrayList<Integer> l1 = new ArrayList<>();
        ArrayList<Integer> l2 = new ArrayList<>();
        inOrderTrav(root1, l1);
        inOrderTrav(root2, l2);
        int l1Size = l1.size();
        int l2Size = l2.size();
        if (l1Size != l2Size) {
            res = false;
        } else if (l1Size == l2Size) {
            for (int i = 0; i < l1Size; i++) {
                int v1 = l1.get(i);
                int v2 = l2.get(i);
                if (v1 != v2) {
                    res = false;
                    break;
                }
            }
        }
        return res;
    }

    public static void main(String[] args) {
        TreeNode t9 = new TreeNode(200, null, null);
        TreeNode t8 = new TreeNode(2, null, null);
        TreeNode t7 = new TreeNode(1, t8, t9);
        TreeNode t6 = new TreeNode(200, null, null);
        TreeNode t5 = new TreeNode(2, null, null);
        TreeNode t4 = new TreeNode(1, t5, t6);
        // TreeNode t3 = new TreeNode(5, t9, t4);
        // TreeNode t2 = new TreeNode(1, t6, t5);
        // TreeNode t1 = new TreeNode(3, t3, t2);

        System.out.println(leafSimilar(t4, t7));

    }
}
