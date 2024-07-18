/*
1110. Delete Nodes And Return Forest

Given the root of a binary tree, each node in the tree has a distinct value.
After deleting all nodes with a value in to_delete, we are left with a forest 
(a disjoint union of trees).

Return the roots of the trees in the remaining forest. You may return the result in any 
order.

Example 1:
Input: root = [1,2,3,4,5,6,7], to_delete = [3,5]
Output: [[1,2,null,4],[6],[7]]

Example 2:

Input: root = [1,2,4,null,3], to_delete = [3]
Output: [[1,2,4]]
 
Constraints:

The number of nodes in the given tree is at most 1000.
Each node has a distinct value between 1 and 1000.
to_delete.length <= 1000
to_delete contains distinct values between 1 and 1000.

*/

import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;

public class DeleteNodeReturnForest {
    public static List<TreeNode> delNodes(TreeNode root, int[] to_delete) {
        HashMap<Integer, Integer> hm = new HashMap<>();
        for (int i = 0; i < to_delete.length; i++) {
            hm.put(to_delete[i], 1);
        }
        if (root == null) {
            return null;
        }
        List<TreeNode> l1 = new ArrayList<>();
        if (!hm.containsKey(root.val)) {
            l1.add(root);
        }
        if (to_delete.length == 0) {
            return l1;
        }
        delNodesWrapper(root, hm, l1);
        return l1;
    }

    public static boolean delNodesWrapper(TreeNode root, HashMap<Integer, Integer> hm, List<TreeNode> l1) {
        if (root == null) {
            return false;
        }
        if (hm.containsKey(root.val)) {
            if (root.left != null) {
                delNodesWrapper(root.left, hm, l1);
                if (!hm.containsKey(root.left.val)) {
                    l1.add(root.left);
                }
            }
            if (root.right != null) {
                delNodesWrapper(root.right, hm, l1);
                if (!hm.containsKey(root.right.val)) {
                    l1.add(root.right);
                }
            }
            return true;
        }
        if (root.left != null) {
            boolean val = delNodesWrapper(root.left, hm, l1);
            if (val) {
                root.left = null;
            }
        }
        if (root.right != null) {
            boolean val = delNodesWrapper(root.right, hm, l1);
            if (val) {
                root.right = null;
            }
        }
        return false;
    }

    public static void main(String[] args) {
        // TreeNode n7 = new TreeNode(7, null, null);
        // TreeNode n6 = new TreeNode(6, null, null);
        TreeNode n5 = new TreeNode(5, null, null);
        TreeNode n4 = new TreeNode(4, n5, null);
        TreeNode n3 = new TreeNode(3, null, null);
        TreeNode n2 = new TreeNode(2, n3, null);
        TreeNode n1 = new TreeNode(1, n2, n4);

        int[] del = new int[] { 4, 1 };
        for (TreeNode val : delNodes(n1, del)) {
            System.out.println(val.val);
        }
    }
}
