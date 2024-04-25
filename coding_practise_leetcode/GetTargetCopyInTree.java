/*
1379. Find a Corresponding Node of a Binary Tree in a Clone of That Tree

Given two binary trees original and cloned and given a reference to a node target in the 
original tree.
The cloned tree is a copy of the original tree.
Return a reference to the same node in the cloned tree.
Note that you are not allowed to change any of the two trees or the target node and the 
answer must be a reference to a node in the cloned tree.

 */

public class GetTargetCopyInTree {
    public final TreeNode getTargetCopy(final TreeNode original, final TreeNode cloned, final TreeNode target) {
        TreeNode curr = cloned;
        return getTargetCopyWrapper(curr, target);
    }

    public static TreeNode getTargetCopyWrapper(TreeNode node, TreeNode target){
        if(node == null){
            return null;
        }
        if(node.val == target.val){
            return node;
        }
        TreeNode res = null;
        res = getTargetCopyWrapper(node.right, target);
        if(res != null){
            return res;
        }
        res = getTargetCopyWrapper(node.left, target);
        if(res != null){
            return res;
        }
        return null;
    }
}
